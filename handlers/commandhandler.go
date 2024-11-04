package handlers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/H1ghN0on/go-tgbot-engine/bot"
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
)

type Stater interface {
	GetName() string
	GetStartCommand() string
	GetAvailableCommands() []string
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
}

type Handlerable interface {
	GetCommands() []string
	InitHandler()
	Handle(command string, params HandlerParams) ([]HandlerResponse, bool, error)
	DeinitHandler()
}

type BackHandlerable interface {
	Handlerable
	UpdateLastCommand(command string)
	ClearCommandQueue()
}

type CommandHandlerError struct {
	message string
}

func (err CommandHandlerError) Error() string {
	return err.message
}

type CommandHandlerRequest struct {
	receivedMessage   bottypes.Message
	shouldUpdateQueue bool
}

func (req CommandHandlerRequest) GetMessage() bottypes.Message {
	return req.receivedMessage
}

func (req CommandHandlerRequest) ShouldUpdateQueue() bool {
	return req.shouldUpdateQueue
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
	sm            StateMachiner
	gs            GlobalStater
	handlers      []Handlerable
	activeHandler Handlerable
	backHandler   BackHandlerable
}

func (ch *CommandHandler) NewCommandHandlerRequest(msg bottypes.Message, shouldUpdateQueue bool) bot.CommandHandlerRequester {
	return &CommandHandlerRequest{
		receivedMessage:   msg,
		shouldUpdateQueue: shouldUpdateQueue,
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

func (ch *CommandHandler) checkCommandInState(command string) bool {
	return slices.Contains(ch.sm.GetActiveState().GetAvailableCommands(), command) ||
		(slices.Contains(ch.sm.GetActiveState().GetAvailableCommands(), "*") && !strings.HasPrefix(command, "/"))
}

func (ch *CommandHandler) handlePostCommands(message bottypes.Message, responses []HandlerResponse) ([]bot.HandlerResponser, error) {

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

func (ch *CommandHandler) handleCommand(command string, receivedMessage bottypes.Message) (bot.CommandHandlerResponser, error) {
	var res CommandHandlerResponse

	for _, handler := range ch.handlers {
		if slices.Contains(handler.GetCommands(), command) || ch.activeHandler != nil {

			if ch.activeHandler != nil {
				if slices.Contains(ch.backHandler.GetCommands(), command) {
					responses, isFinished, err := ch.backHandler.Handle(command, HandlerParams{message: receivedMessage})

					if err != nil {
						return CommandHandlerResponse{}, CommandHandlerError{message: err.Error()}
					}

					if isFinished {
						ch.activeHandler.DeinitHandler()
						ch.activeHandler = nil
					}
					ch.backHandler.UpdateLastCommand(command)
					res.responses = responses
					break
				}

			} else {
				ch.activeHandler = handler
				ch.activeHandler.InitHandler()
			}

			responses, isFinished, err := ch.activeHandler.Handle(command, HandlerParams{message: receivedMessage})

			if err != nil {
				return CommandHandlerResponse{}, CommandHandlerError{message: err.Error()}
			}

			if isFinished {
				ch.activeHandler.DeinitHandler()
				ch.activeHandler = nil
			}
			ch.backHandler.UpdateLastCommand(command)
			res.responses = responses
			break
		}
	}

	if len(res.responses) == 0 {
		return CommandHandlerResponse{}, CommandHandlerError{message: "this command is unknown"}
	}

	err := ch.updateState(res)
	if err != nil {
		return CommandHandlerResponse{}, CommandHandlerError{message: "command handler error: " + err.Error()}
	}

	postCommandResponses, err := ch.handlePostCommands(receivedMessage, res.responses)
	if err != nil {
		return CommandHandlerResponse{}, err
	}

	for _, postCommandResponse := range postCommandResponses {
		res.responses = append(res.responses, postCommandResponse.(HandlerResponse))
	}

	return res, nil
}

func (ch *CommandHandler) Handle(req bot.CommandHandlerRequester) (bot.CommandHandlerResponser, error) {

	receivedMessage := req.GetMessage()
	command := receivedMessage.Text

	if !ch.checkCommandInState(command) {
		return CommandHandlerResponse{}, CommandHandlerError{message: "this command is not available "}
	}

	chRes, err := ch.handleCommand(command, receivedMessage)
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
	backHandler := NewBackHandler(gs, sm)

	ch.handlers = append(ch.handlers,
		setInfoHandler,
		keyboardHandler,
		levelFourHandler,
		startHandler,
		checkboxHandler,
		backHandler)

	ch.backHandler = backHandler

	return ch
}
