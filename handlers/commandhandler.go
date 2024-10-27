package handlers

import (
	"fmt"
	"slices"

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
	commandsQueue []string
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

	if len(ch.commandsQueue) > 10 {
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

	previousState := ch.sm.GetPreviousState()
	if previousState.GetName() == "" {
		return CommandHandlerResponse{}, CommandHandlerError{message: "can not return to previous state"}
	}
	ch.sm.SetState(previousState)
	receivedMessageBack := req.GetMessage()
	receivedMessageBack.Text = previousState.GetStartCommand()
	res, err := ch.Handle(CommandHandlerRequest{
		receivedMessage:   receivedMessageBack,
		shouldUpdateQueue: false,
	})

	if err != nil {
		return CommandHandlerResponse{}, CommandHandlerError{message: fmt.Errorf("back error: %w", err).Error()}
	}

	if req.ShouldUpdateQueue() {
		ch.updateCommandsQueue(req.GetMessage().Text)
	}

	return res, nil
}

func (ch *CommandHandler) handleBackCommandRequest(req bot.CommandHandlerRequester) (bot.CommandHandlerResponser, error) {

	if len(ch.commandsQueue) < 2 {
		return CommandHandlerResponse{}, CommandHandlerError{message: "can not return to previous command"}
	}
	lastCommand := ch.commandsQueue[len(ch.commandsQueue)-2]
	if lastCommand == "" {
		return CommandHandlerResponse{}, CommandHandlerError{message: "can not return to previous command"}
	}

	receivedMessageBack := req.GetMessage()
	receivedMessageBack.Text = lastCommand

	res, err := ch.Handle(CommandHandlerRequest{
		receivedMessage:   receivedMessageBack,
		shouldUpdateQueue: false,
	})
	if err != nil {
		return CommandHandlerResponse{}, CommandHandlerError{message: fmt.Errorf("back command error: %w", err).Error()}
	}

	if req.ShouldUpdateQueue() {
		ch.updateCommandsQueue(req.GetMessage().Text)
	}

	return res, nil
}

func (ch *CommandHandler) Handle(req bot.CommandHandlerRequester) (bot.CommandHandlerResponser, error) {
	var res CommandHandlerResponse

	receivedMessage := req.GetMessage()

	if !slices.Contains(ch.sm.GetActiveState().GetAvailableCommands(), receivedMessage.Text) {
		return CommandHandlerResponse{}, CommandHandlerError{message: "this command is not available "}
	}

	command := receivedMessage.Text

	switch command {
	case "/level_one":
		handleRes := LevelOneHandler(HandlerParams{message: receivedMessage})
		showHandleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		lolRes := LevelFourStartHandler(HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, handleRes, showHandleRes, lolRes)
	case "/level_two":
		levelhandleRes := LevelTwoHandler(HandlerParams{message: receivedMessage})
		showHandleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, levelhandleRes, showHandleRes)
	case "/level_three":
		levelhandleRes := LevelThreeHandler(HandlerParams{message: receivedMessage})
		showHandleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, levelhandleRes, showHandleRes)
	case "/level_four_start":
		levelhandleRes := MakeStateBackable(LevelFourStartHandler, HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, levelhandleRes)
	case "/level_four_one":
		levelhandleRes := MakeStateBackable(LevelFourOneHandler, HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, levelhandleRes)
	case "/level_four_two":
		levelhandleRes := MakeStateBackable(LevelFourTwoHandler, HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, levelhandleRes)
	case "/level_four_three":
		levelhandleRes := MakeStateBackable(LevelFourThreeHandler, HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, levelhandleRes)
	case "/level_four_four":
		levelhandleRes := LevelFourFourHandler(HandlerParams{message: receivedMessage})
		showHandleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, levelhandleRes, showHandleRes)
	case "/big_messages":
		levelhandleRes := BigMessagesHandler(HandlerParams{message: receivedMessage})
		handleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, levelhandleRes, handleRes)
	case "/show_commands":
		handleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, handleRes)
	case "/keyboard_start":
		handleRes := MakeStateBackable(KeyboardStartHandler, HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, handleRes)
	case "/keyboard_one":
		handleRes := MakeCommandBackable(KeyboardOneHandler, HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, handleRes)
	case "/keyboard_two":
		handleRes := MakeCommandBackable(KeyboardTwoHandler, HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, handleRes)
	case "/keyboard_three":
		handleRes := KeyboardThreeHandler(HandlerParams{message: receivedMessage})
		showHandleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.responses = append(res.responses, handleRes, showHandleRes)
	case "/back_command":
		handleRes, err := ch.handleBackCommandRequest(req)
		if err != nil {
			return CommandHandlerResponse{}, CommandHandlerError{message: fmt.Errorf("back state error: %w", err).Error()}
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

	if hasMultipleStatesInCommand(res) {
		return CommandHandlerResponse{}, CommandHandlerError{message: "multiple states in commands are forbidden"}
	}

	for _, response := range res.responses {
		if response.NextState() == "" {
			continue
		}

		err := ch.sm.SetStateByName(response.NextState())
		if err != nil {
			return CommandHandlerResponse{}, CommandHandlerError{message: fmt.Errorf("handle error: %w", err).Error()}
		}
		break
	}

	if req.ShouldUpdateQueue() {
		ch.updateCommandsQueue(command)
	}

	return res, nil
}

func NewCommandHandler(sm StateMachiner) *CommandHandler {
	return &CommandHandler{
		sm: sm,
	}
}
