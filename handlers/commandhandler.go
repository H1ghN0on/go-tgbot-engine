package handlers

import (
	"slices"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	errs "github.com/H1ghN0on/go-tgbot-engine/errors"
)

type StateMachiner interface {
	AddStates(states ...bottypes.State)
	SetStateByName(stateName string) error
	SetState(state bottypes.State) error
	GetActiveState() bottypes.State
}

func hasMultipleStatesInCommand(res CommandHandlerResponse) bool {
	statesSet := make(map[string]bool)
	for _, v := range res.Responses {
		if v.ShouldSwitchState != "" {
			statesSet[v.ShouldSwitchState] = true
		}
	}

	return len(statesSet) > 1
}

type CommandHandlerResponse struct {
	Responses []HandlerResponse
}

type CommandHandler struct {
	sm StateMachiner
}

func (ch *CommandHandler) Handle(receivedMessage bottypes.Message) (CommandHandlerResponse, error) {
	var res CommandHandlerResponse

	if !slices.Contains(ch.sm.GetActiveState().AvailableCommands, bottypes.Command{Text: receivedMessage.Text}) {
		return CommandHandlerResponse{}, errs.CommandHandlerError{Code: errs.UnknownCommand, Message: "This command is not available "}
	}

	// if !ch.sm.GetActiveState().ContainsCommand(bottypes.Command{Text: receivedMessage.Text}) {
	// 	return CommandHandlerResponse{}, errs.CommandHandlerError{Code: errs.UnknownCommand, Message: "This command is not available "}
	// }

	switch receivedMessage.Text {
	case "/level_one":
		handleRes := LevelOneHandler(HandlerParams{message: receivedMessage})
		showHandleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		lolRes := LevelFourStartHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, handleRes, showHandleRes, lolRes)
	case "/level_two":
		levelhandleRes := LevelTwoHandler(HandlerParams{message: receivedMessage})
		showHandleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, levelhandleRes, showHandleRes)
	case "/level_three":
		levelhandleRes := LevelThreeHandler(HandlerParams{message: receivedMessage})
		showHandleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, levelhandleRes, showHandleRes)
	case "/level_four_start":
		levelhandleRes := LevelFourStartHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, levelhandleRes)
	case "/level_four_one":
		levelhandleRes := LevelFourOneHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, levelhandleRes)
	case "/level_four_two":
		levelhandleRes := LevelFourTwoHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, levelhandleRes)
	case "/level_four_three":
		levelhandleRes := LevelFourThreeHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, levelhandleRes)
	case "/level_four_four":
		levelhandleRes := LevelFourFourHandler(HandlerParams{message: receivedMessage})
		showHandleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, levelhandleRes, showHandleRes)
	case "/big_messages":
		levelhandleRes := BigMessagesHandler(HandlerParams{message: receivedMessage})
		handleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, levelhandleRes, handleRes)
	case "/show_commands":
		handleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, handleRes)
	case "/keyboard_start":
		handleRes := KeyboardStartHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, handleRes)
	case "/keyboard_one":
		handleRes := KeyboardOneHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, handleRes)
	case "/keyboard_two":
		handleRes := KeyboardTwoHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, handleRes)
	case "/keyboard_three":
		handleRes := KeyboardThreeHandler(HandlerParams{message: receivedMessage})
		showHandleRes := ShowCommandsHandler(HandlerParams{message: receivedMessage})
		res.Responses = append(res.Responses, handleRes, showHandleRes)
	case "/create_error":
		return CommandHandlerResponse{}, errs.CommandHandlerError{Code: errs.UnknownCommand, Message: "This command is unknown "}
	default:
		return CommandHandlerResponse{}, errs.CommandHandlerError{Code: errs.UnknownCommand, Message: "This command is unknown "}
	}

	if hasMultipleStatesInCommand(res) {
		return CommandHandlerResponse{}, errs.StateMachineError{Code: errs.MultipleStatesFromCommand, Message: "Multiple states are forbidden"}
	}

	for _, response := range res.Responses {
		if response.ShouldSwitchState == "" {
			continue
		}

		err := ch.sm.SetStateByName(response.ShouldSwitchState)
		if err != nil {
			return CommandHandlerResponse{}, errs.StateMachineError{Code: errs.WrongState, Message: err.Error()}
		}
		break
	}

	return res, nil
}

func NewCommandHandler(sm StateMachiner) CommandHandler {
	return CommandHandler{
		sm: sm,
	}
}
