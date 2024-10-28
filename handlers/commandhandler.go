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
	Handle(command string, params HandlerParams) ([]HandlerResponse, bool)
	DeinitHandler()
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
	responses     []HandlerResponse
	triggerRemove bool
	isNothing     bool
}

func (chr CommandHandlerResponse) GetResponses() []bot.HandlerResponser {
	var convertedResponses []bot.HandlerResponser

	for _, response := range chr.responses {
		convertedResponses = append(convertedResponses, response)
	}
	return convertedResponses
}

func (chr CommandHandlerResponse) TriggerRemove() bool {
	return chr.triggerRemove
}

func (chr CommandHandlerResponse) IsNothing() bool {
	return chr.isNothing
}

type CommandHandler struct {
	sm            StateMachiner
	gs            GlobalStater
	commandsQueue []string
	handlers      []Handlerable
	activeHandler Handlerable
}

func (ch *CommandHandler) updateCommandsQueue(command string) {
	switch command {
	case "/back_command":
		ch.commandsQueue = ch.commandsQueue[:len(ch.commandsQueue)-1]
	case "/back_state":
		break
	default:
		ch.commandsQueue = append(ch.commandsQueue, command)
	}

	if len(ch.commandsQueue) > 20 {
		ch.commandsQueue = ch.commandsQueue[1:len(ch.commandsQueue)]
	}
}

func (ch *CommandHandler) NewCommandHandlerRequest(msg bottypes.Message, shouldUpdateQueue bool) bot.CommandHandlerRequester {
	return &CommandHandlerRequest{
		receivedMessage:   msg,
		shouldUpdateQueue: shouldUpdateQueue,
	}
}

func (ch *CommandHandler) handleBackStateRequest(req bot.CommandHandlerRequester) (bot.CommandHandlerResponser, error) {

	var res CommandHandlerResponse

	previousState := ch.sm.GetPreviousState()
	if previousState.GetName() == "" {
		return CommandHandlerResponse{}, CommandHandlerError{message: "can not return to previous state"}
	}
	ch.sm.SetState(previousState)
	receivedMessageBack := req.GetMessage()
	receivedMessageBack.Text = previousState.GetStartCommand()

	handleRes, err := ch.Handle(CommandHandlerRequest{
		receivedMessage:   receivedMessageBack,
		shouldUpdateQueue: false,
	})

	if err != nil {
		return CommandHandlerResponse{}, CommandHandlerError{message: fmt.Errorf("back error: %w", err).Error()}
	}

	if req.ShouldUpdateQueue() {
		ch.updateCommandsQueue(req.GetMessage().Text)
	}

	res = handleRes.(CommandHandlerResponse)
	res.triggerRemove = true

	return res, nil
}

func (ch *CommandHandler) handleBackCommandRequest(req bot.CommandHandlerRequester) (bot.CommandHandlerResponser, error) {

	var res CommandHandlerResponse

	if len(ch.commandsQueue) < 2 {
		return CommandHandlerResponse{}, CommandHandlerError{message: "can not return to previous command"}
	}
	lastCommand := ch.commandsQueue[len(ch.commandsQueue)-2]
	if lastCommand == "" {
		return CommandHandlerResponse{}, CommandHandlerError{message: "can not return to previous command"}
	}

	receivedMessageBack := req.GetMessage()
	receivedMessageBack.Text = lastCommand

	handleRes, err := ch.Handle(CommandHandlerRequest{
		receivedMessage:   receivedMessageBack,
		shouldUpdateQueue: false,
	})

	if err != nil {
		return CommandHandlerResponse{}, CommandHandlerError{message: fmt.Errorf("back command error: %w", err).Error()}
	}

	if req.ShouldUpdateQueue() {
		ch.updateCommandsQueue(req.GetMessage().Text)
	}

	res = handleRes.(CommandHandlerResponse)

	return res, nil
}

func (ch *CommandHandler) checkCommandInState(command string) bool {
	return slices.Contains(ch.sm.GetActiveState().GetAvailableCommands(), command) ||
		(slices.Contains(ch.sm.GetActiveState().GetAvailableCommands(), "*") && !strings.HasPrefix(command, "/"))
}

func (ch *CommandHandler) moveToAnotherState(message bottypes.Message) []HandlerResponse {

	startStateCommand := ch.sm.GetActiveState().GetStartCommand()

	if message.Text == startStateCommand {
		return nil
	}

	for _, handler := range ch.handlers {
		if slices.Contains(handler.GetCommands(), startStateCommand) {
			handler.InitHandler()
			responses, _ := handler.Handle(startStateCommand, HandlerParams{message: message})
			handler.DeinitHandler()
			return responses
		}
	}

	return nil
}

func (ch *CommandHandler) Handle(req bot.CommandHandlerRequester) (bot.CommandHandlerResponser, error) {
	var res CommandHandlerResponse

	receivedMessage := req.GetMessage()
	command := receivedMessage.Text

	if !ch.checkCommandInState(command) {
		return CommandHandlerResponse{}, CommandHandlerError{message: "this command is not available "}
	}

	// Trying to handle in main handlers

	for _, handler := range ch.handlers {
		if slices.Contains(handler.GetCommands(), command) || ch.activeHandler != nil {
			if ch.activeHandler == nil {
				ch.activeHandler = handler
				ch.activeHandler.InitHandler()
			}
			responses, isFinished := ch.activeHandler.Handle(command, HandlerParams{message: receivedMessage})
			if isFinished {
				ch.activeHandler.DeinitHandler()
				ch.activeHandler = nil
			}
			if req.ShouldUpdateQueue() {
				ch.updateCommandsQueue(command)
			}
			res.responses = responses
			break
		}
	}

	// Trying to handle in additional handlers

	if len(res.responses) == 0 {
		switch command {
		case "/back_command":
			handleRes, err := ch.handleBackCommandRequest(req)
			if err != nil {
				return CommandHandlerResponse{}, CommandHandlerError{message: fmt.Errorf("back command error: %w", err).Error()}
			}
			return handleRes, nil

		case "/back_state":
			handleRes, err := ch.handleBackStateRequest(req)

			if err != nil {
				return CommandHandlerResponse{}, CommandHandlerError{message: fmt.Errorf("back state error: %w", err).Error()}
			}
			return handleRes, nil
		case "/create_error":
			return CommandHandlerResponse{}, CommandHandlerError{message: "this command is unknown"}
		default:
			return CommandHandlerResponse{}, CommandHandlerError{message: "this command is unknown"}
		}
	}

	if hasMultipleStatesInCommand(res) {
		return CommandHandlerResponse{}, CommandHandlerError{message: "multiple states in commands are forbidden"}
	}

	// Trying to move to another state (if so, handle start command of new state)

	for _, response := range res.responses {
		if response.NextState() == "" {
			continue
		}
		err := ch.sm.SetStateByName(response.NextState())
		if err != nil {
			return CommandHandlerResponse{}, CommandHandlerError{message: fmt.Errorf("handler error: %w", err).Error()}
		}
		stateResponses := ch.moveToAnotherState(receivedMessage)
		if stateResponses != nil {
			res.responses = append(res.responses, stateResponses...)
		}
		break
	}

	// Remove trigger marked messages if needed

	for _, response := range res.responses {
		if response.isNothing {
			res.isNothing = true
		}

		if response.isRemoveTriggered {
			res.triggerRemove = true
		}
	}

	return res, nil
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

	ch.handlers = append(ch.handlers, setInfoHandler, keyboardHandler, levelFourHandler, startHandler, checkboxHandler)

	return ch
}
