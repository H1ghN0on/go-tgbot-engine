package handlers

import (
	"fmt"
	"slices"

	"github.com/H1ghN0on/go-tgbot-engine/bot"
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/logger"

	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
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
		if v.NextState() != "" {
			statesSet[v.NextState()] = true
		}
	}

	return len(statesSet) > 1
}

type GlobalStater interface {
	GetName() string
	GetSurname() string
	GetAge() int

	SetName(name string)
	SetSurname(surname string)
	SetAge(age int)

	GetDataForDynamicKeyboard() map[string][]string
}

type Handlerable interface {
	GetCommands() []bottypes.Command
	Handle(params HandlerParams) ([]HandlerResponse, error)
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

func (chr CommandHandlerResponse) GetResponses() []bot.HandlerResponser {
	var convertedResponses []bot.HandlerResponser

	for _, response := range chr.responses {
		convertedResponses = append(convertedResponses, response)
	}
	return convertedResponses
}

type CommandHandler struct {
	sm           StateMachiner
	gs           GlobalStater
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

func (ch *CommandHandler) NewCommandHandlerRequest(msg bottypes.ParsedMessage) bot.CommandHandlerRequester {
	return &CommandHandlerRequest{
		receivedMessage: msg,
	}
}

func (ch *CommandHandler) updateState(res CommandHandlerResponse) error {

	if hasMultipleStatesInCommand(res) {
		return CommandHandlerError{message: "multiple states in commands are forbidden"}
	}

	for _, response := range res.responses {
		if response.NextState() == "" {
			continue
		}

		err := ch.sm.SetStateByName(response.NextState())
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
		ch.nextCommands = append(ch.nextCommands, response.nextCommands...)
	}
}

func isCommandInSlice(commands []bottypes.Command, command bottypes.Command) bool {
	return slices.ContainsFunc(commands, func(cmd bottypes.Command) bool { return command.Equal(cmd) })
}

func (ch *CommandHandler) checkCommandInState(command bottypes.Command) bool {
	return isCommandInSlice(ch.sm.GetActiveState().GetAvailableCommands(), command) ||
		(isCommandInSlice(ch.sm.GetActiveState().GetAvailableCommands(), cmd.AnyCommand) && !command.IsCommand())
}

func (ch *CommandHandler) checkCommandInNextCommands(command bottypes.Command) bool {
	return len(ch.nextCommands) == 0 || isCommandInSlice(ch.nextCommands, command) ||
		(isCommandInSlice(ch.nextCommands, cmd.AnyCommand) && !command.IsCommand())
}

func (ch *CommandHandler) checkCommandInHandler(command bottypes.Command, handler Handlerable) bool {
	commandsToCheck := handler.GetCommands()
	return isCommandInSlice(commandsToCheck, command) ||
		(isCommandInSlice(commandsToCheck, cmd.AnyCommand) && !command.IsCommand())
}

func (ch *CommandHandler) hasCommandInHandler(commands []bottypes.Command, handler Handlerable) bool {

	commandsToCheck := handler.GetCommands()
	for _, command := range commands {
		if isCommandInSlice(commandsToCheck, command) || (isCommandInSlice(commandsToCheck, cmd.AnyCommand) && !command.IsCommand()) {
			return true
		}
	}
	return false
}

func (ch *CommandHandler) handlePostCommands(message bottypes.ParsedMessage, responses []HandlerResponse) ([]bot.HandlerResponser, error) {

	var res []bot.HandlerResponser

	for idx, response := range responses {
		for _, commandToHandle := range response.postCommandsHandle {
			handleRes, err := ch.handleCommand(commandToHandle, message)
			if err != nil {
				return []bot.HandlerResponser{}, CommandHandlerError{message: "handle post commands: " + err.Error()}
			}
			res = append(res, handleRes.GetResponses()...)
		}
		responses[idx].postCommandsHandle = nil
	}

	return res, nil
}

func (ch *CommandHandler) handleCommand(command bottypes.Command, receivedMessage bottypes.ParsedMessage) (bot.CommandHandlerResponser, error) {
	var res CommandHandlerResponse

	logger.CommandHandler().Info("trying to handle", command.String())

	for _, handler := range ch.handlers {
		commandToCheck := command
		if !command.IsCommand() && ch.hasCommandInHandler(ch.nextCommands, handler) {
			commandToCheck = cmd.AnyCommand
		}

		if ch.checkCommandInHandler(commandToCheck, handler) {

			responses, err := handler.Handle(HandlerParams{command: commandToCheck, message: receivedMessage})

			if err != nil {
				logger.CommandHandler().Critical("error while handling command", command.String(), err.Error())
				return CommandHandlerResponse{}, CommandHandlerError{message: err.Error()}
			}

			ch.backHandler.UpdateLastCommand(command)
			res.responses = responses
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

func (ch *CommandHandler) Handle(req bot.CommandHandlerRequester) (bot.CommandHandlerResponser, error) {

	receivedMessage := req.GetMessage()

	if !ch.checkCommandInNextCommands(receivedMessage.Command) {
		logger.CommandHandler().Critical(receivedMessage.Command.String(), "is not in next commands (", convertCommandsToString(ch.nextCommands), ")")
		return CommandHandlerResponse{}, CommandHandlerError{message: "this command is not available (not in next commands)"}
	}

	if !ch.checkCommandInState(receivedMessage.Command) {
		logger.CommandHandler().Critical(receivedMessage.Command.String(), "is not in state commands (", convertCommandsToString(ch.sm.GetActiveState().GetAvailableCommands()), ")")
		return CommandHandlerResponse{}, CommandHandlerError{message: "this command is not available (not in state)"}
	}

	chRes, err := ch.handleCommand(receivedMessage.Command, receivedMessage)
	if err != nil {
		return CommandHandlerResponse{}, err
	}

	return chRes, nil
}

func NewCommandHandler(sm StateMachiner, gs GlobalStater) *CommandHandler {
	ch := &CommandHandler{
		sm: sm,
		gs: gs,
	}

	setInfoHandler := NewSetInfoHandler(gs)
	keyboardHandler := NewKeyboardhandler(gs)
	levelFourHandler := NewLevelFourHandler(gs)
	startHandler := NewStartHandler(gs)
	checkboxHandler := NewCheckboxHandler(gs)
	dynamicKeyboardHandler := NewDynamicKeyboardhandler(gs)
	backHandler := NewBackHandler(gs, sm)

	ch.handlers = append(ch.handlers,
		setInfoHandler,
		keyboardHandler,
		levelFourHandler,
		startHandler,
		checkboxHandler,
		dynamicKeyboardHandler,
		backHandler)

	ch.backHandler = backHandler

	return ch
}
