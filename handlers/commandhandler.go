package handlers

import (
	"fmt"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/bot/client"
	"github.com/H1ghN0on/go-tgbot-engine/logger"
)

type Stater interface {
	GetName() string
	GetStartCommand() bottypes.Command
	GetAvailableCommands() []bottypes.Command
	GetAvailableStates() []Stater
	CanRestart() bool
}

type StateMachiner interface {
	AddStates(states ...Stater)
	SetStateByName(stateName string) error
	SetState(state Stater) error
	GetActiveState() Stater
	GetPreviousState() Stater
}

func hasMultipleStatesInCommand(res CommandHandlerResponse) bool {
	statesSet := make(map[string]bool)
	for _, v := range res.responses {
		if v.NextState != "" {
			statesSet[v.NextState] = true
		}
	}

	return len(statesSet) > 1
}

type Handlerable interface {
	GetCommands() []bottypes.Command
	Handle(params HandlerParams) ([]HandlerResponse, error)
	HandleBackCommand(params HandlerParams) ([]HandlerResponse, error)
	FindCommandInTheList(bottypes.Command) (bottypes.Command, error)
}

type BackHandlerable interface {
	Handlerable
	UpdateLastCommand(command bottypes.Command)
	ClearCommandQueue()
}

type CommandHandlerError struct {
	message string
}

func (err CommandHandlerError) Error() string {
	return err.message
}

type CommandHandlerRequest struct {
	receivedMessage bottypes.ParsedMessage
}

func (req CommandHandlerRequest) GetMessage() bottypes.ParsedMessage {
	return req.receivedMessage
}

type CommandHandlerResponse struct {
	responses []HandlerResponse
}

func (chr CommandHandlerResponse) GetResponses() []client.HandlerResponser {
	var convertedResponses []client.HandlerResponser

	for _, response := range chr.responses {
		convertedResponses = append(convertedResponses, response)
	}
	return convertedResponses
}

type CommandHandler struct {
	sm           StateMachiner
	handlers     []Handlerable
	backHandler  BackHandlerable
	nextCommands []bottypes.Command
}

func convertCommandsToString(commands []bottypes.Command) string {
	var ret string
	for _, command := range commands {
		ret += command.String() + " "
	}
	if len(ret) != 0 {
		ret = ret[:len(ret)-1]
	}

	return ret
}

func (ch *CommandHandler) NewCommandHandlerRequest(msg bottypes.ParsedMessage) client.CommandHandlerRequester {
	return CommandHandlerRequest{
		receivedMessage: msg,
	}
}

func (ch *CommandHandler) updateState(res CommandHandlerResponse) error {

	if hasMultipleStatesInCommand(res) {
		return CommandHandlerError{message: "multiple states in commands are forbidden"}
	}

	for _, response := range res.responses {
		if response.NextState == "" {
			continue
		}

		err := ch.sm.SetStateByName(response.NextState)
		ch.backHandler.ClearCommandQueue()
		if err != nil {
			return CommandHandlerError{message: fmt.Errorf("handler error: %w", err).Error()}
		}
		break
	}
	return nil
}

func (ch *CommandHandler) updateNextCommands(responses []HandlerResponse) {
	ch.nextCommands = nil
	for _, response := range responses {
		ch.nextCommands = append(ch.nextCommands, response.NextCommands...)
	}
}

func (ch *CommandHandler) checkCommandInState(command bottypes.Command) bool {
	return command.InSlice(ch.sm.GetActiveState().GetAvailableCommands())
}

func (ch *CommandHandler) checkCommandInNextCommands(command bottypes.Command) bool {
	return len(ch.nextCommands) == 0 || command.InSlice(ch.nextCommands)
}

func (ch *CommandHandler) checkCommandInHandler(command bottypes.Command, handler Handlerable) bool {
	commandsToCheck := handler.GetCommands()
	return command.InSlice(commandsToCheck)
}

func (ch *CommandHandler) handlePostCommands(message bottypes.ParsedMessage, responses []HandlerResponse) ([]client.HandlerResponser, error) {

	var res []client.HandlerResponser

	for idx, response := range responses {
		for _, commandToHandle := range response.PostCommandsHandle.Commands {
			handleRes, err := ch.handleCommand(commandToHandle, message, response.PostCommandsHandle.IsBackCommand)
			if err != nil {
				return []client.HandlerResponser{}, CommandHandlerError{message: "handle post commands: " + err.Error()}
			}
			res = append(res, handleRes.GetResponses()...)
		}
		responses[idx].PostCommandsHandle.Commands = nil
	}

	return res, nil
}

func (ch *CommandHandler) handleCommand(
	command bottypes.Command,
	receivedMessage bottypes.ParsedMessage,
	shouldHandleBack bool) (client.CommandHandlerResponser, error) {

	var res CommandHandlerResponse

	logger.CommandHandler().Info("trying to handle", command.String())

	for _, handler := range ch.handlers {

		if ch.checkCommandInHandler(command, handler) {

			trueCommand, err := handler.FindCommandInTheList(command)
			if err != nil {
				logger.CommandHandler().Critical("error while finding true command in the list", command.String(), err.Error())
				return CommandHandlerResponse{}, CommandHandlerError{message: err.Error()}
			}

			if shouldHandleBack {
				responses, err := handler.HandleBackCommand(HandlerParams{Command: trueCommand, Message: receivedMessage})

				if err != nil {
					logger.CommandHandler().Critical("error while handling command", trueCommand.String(), err.Error())
					return CommandHandlerResponse{}, CommandHandlerError{message: err.Error()}
				}

				res.responses = append(res.responses, responses...)
			}

			responses, err := handler.Handle(HandlerParams{Command: trueCommand, Message: receivedMessage})

			if err != nil {
				logger.CommandHandler().Critical("error while handling command", trueCommand.String(), err.Error())
				return CommandHandlerResponse{}, CommandHandlerError{message: err.Error()}
			}

			ch.backHandler.UpdateLastCommand(trueCommand)
			res.responses = append(res.responses, responses...)
			break
		}
	}

	if len(res.responses) == 0 {
		logger.CommandHandler().Warning("command", command.String(), "is unknown")
		return CommandHandlerResponse{}, CommandHandlerError{message: "this command is unknown"}
	}

	err := ch.updateState(res)
	if err != nil {
		return CommandHandlerResponse{}, CommandHandlerError{message: "command handler error: " + err.Error()}
	}
	ch.updateNextCommands(res.responses)

	postCommandResponses, err := ch.handlePostCommands(receivedMessage, res.responses)
	if err != nil {
		return CommandHandlerResponse{}, err
	}

	for _, postCommandResponse := range postCommandResponses {
		res.responses = append(res.responses, postCommandResponse.(HandlerResponse))
	}

	ch.updateNextCommands(res.responses)

	logger.CommandHandler().Info("command", command.String(), "handled successfully")

	return res, nil
}

func (ch *CommandHandler) Handle(req client.CommandHandlerRequester) (client.CommandHandlerResponser, error) {

	receivedMessage := req.GetMessage()

	if !ch.checkCommandInNextCommands(receivedMessage.Command) {
		logger.CommandHandler().Critical(receivedMessage.Command.String(), "is not in next commands (", convertCommandsToString(ch.nextCommands), ")")
		return CommandHandlerResponse{}, CommandHandlerError{message: "this command is not available (not in next commands)"}
	}

	if !ch.checkCommandInState(receivedMessage.Command) {
		logger.CommandHandler().Critical(receivedMessage.Command.String(), "is not in state commands (", convertCommandsToString(ch.sm.GetActiveState().GetAvailableCommands()), ")")
		return CommandHandlerResponse{}, CommandHandlerError{message: "this command is not available (not in state)"}
	}

	chRes, err := ch.handleCommand(receivedMessage.Command, receivedMessage, false)
	if err != nil {
		return CommandHandlerResponse{}, err
	}

	return chRes, nil
}

func NewCommandHandler(handlers []Handlerable, sm StateMachiner) *CommandHandler {
	ch := &CommandHandler{
		sm: sm,
	}

	backHandler := NewBackHandler(sm)

	ch.handlers = append(ch.handlers, handlers...)
	ch.handlers = append(ch.handlers, backHandler)

	ch.backHandler = backHandler

	return ch
}
