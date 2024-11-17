package handlers

import (
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type SetInfoHandler struct {
	Handler
	nextCommand bottypes.Command
}

func NewSetInfoHandler(gs GlobalStater) *SetInfoHandler {

	sh := &SetInfoHandler{}
	sh.gs = gs

	sh.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.SetInfoStartCommand: {
			sh.SetInfoStartHandler,
		},

		// Не работает модификатор CommandBackable
		cmd.SetNameCommand:    {sh.ModifyHandler(sh.SetNameHandler, []int{StateBackable})},
		cmd.SetSurnameCommand: {sh.ModifyHandler(sh.SetSurnameHandler, []int{StateBackable})},
		cmd.SetAgeCommand:     {sh.ModifyHandler(sh.SetAgeHandler, []int{StateBackable})},
		cmd.SetInfoEndCommand: {sh.SetInfoEndHandler},
		cmd.AnyCommand:        {},
	}

	return sh
}

func (handler *SetInfoHandler) Handle(params HandlerParams) ([]HandlerResponse, error) {
	var res []HandlerResponse

	var commandToHandle bottypes.Command

	if params.message.Command.IsCommand() {
		commandToHandle = params.command
	} else {
		commandToHandle = handler.nextCommand
	}

	handleFuncs, ok := handler.GetCommandFromMap(commandToHandle)
	if !ok {
		panic("wrong handler")
	}

	for _, handleFunc := range handleFuncs {
		response, err := handleFunc(params)
		if err != nil {
			return []HandlerResponse{}, err
		}
		res = append(res, response)
	}

	return res, nil
}

func (handler *SetInfoHandler) SetInfoStartHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Let me know you a bit closer"}
	res.messages = append(res.messages, retMessage)
	handler.nextCommand = cmd.SetNameCommand
	res.postCommandsHandle = append(res.postCommandsHandle, cmd.SetNameCommand)
	res.nextState = "info-state"

	return res, nil
}

func (handler *SetInfoHandler) SetNameHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	handler.nextCommand = cmd.SetSurnameCommand

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your name"}
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, cmd.AnyCommand)

	return res, nil
}

func (handler *SetInfoHandler) SetSurnameHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	handler.gs.SetName(params.message.Info.Text)

	handler.nextCommand = cmd.SetAgeCommand

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your surname"}
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, cmd.AnyCommand)

	return res, nil
}

func (handler *SetInfoHandler) SetAgeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	handler.gs.SetSurname(params.message.Info.Text)

	handler.nextCommand = cmd.SetInfoEndCommand

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your age"}
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, cmd.AnyCommand)

	return res, nil
}

func (handler *SetInfoHandler) SetInfoEndHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	age, err := strconv.Atoi(params.message.Info.Text)

	if err != nil {
		retMessage := bottypes.Message{ChatID: chatID, Text: "Wrong age, try again!"}
		res.messages = append(res.messages, retMessage)
		res.nextCommands = append(res.nextCommands, cmd.AnyCommand, cmd.BackStateCommand)
		return res, nil
	}

	handler.gs.SetAge(age)

	retMessage := bottypes.Message{ChatID: chatID, Text: "My gratitude"}
	res.messages = append(res.messages, retMessage)

	res.postCommandsHandle = append(res.postCommandsHandle, cmd.ShowCommandsCommand)
	res.nextState = "start-state"

	return res, nil
}
