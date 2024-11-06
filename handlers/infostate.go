package handlers

import (
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
)

type SetInfoHandler struct {
	Handler
	nextCommand bottypes.Command
}

func NewSetInfoHandler(gs GlobalStater) *SetInfoHandler {

	sh := &SetInfoHandler{}
	sh.gs = gs

	sh.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		"/set_info_start": {
			sh.SetInfoStartHandler,
		},

		// Не работает модификатор CommandBackable
		"/set_name":     {sh.ModifyHandler(sh.SetNameHandler, []int{StateBackable})},
		"/set_surname":  {sh.ModifyHandler(sh.SetSurnameHandler, []int{StateBackable})},
		"/set_age":      {sh.ModifyHandler(sh.SetAgeHandler, []int{StateBackable})},
		"/set_info_end": {sh.SetInfoEndHandler},
		"*":             {},
	}

	return sh
}

func (handler *SetInfoHandler) Handle(command bottypes.Command, params HandlerParams) ([]HandlerResponse, error) {
	var res []HandlerResponse

	var commandToHandle bottypes.Command

	if command.IsCommand() {
		commandToHandle = command
	} else {
		commandToHandle = handler.nextCommand
	}

	handleFuncs, ok := handler.Handler.commands[commandToHandle]
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
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Let me know you a bit closer"}
	res.messages = append(res.messages, retMessage)
	handler.nextCommand = "/set_name"
	res.postCommandsHandle = append(res.postCommandsHandle, "/set_name")
	res.nextState = "info-state"

	return res, nil
}

func (handler *SetInfoHandler) SetNameHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	handler.nextCommand = "/set_surname"

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your name"}
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, "*")

	return res, nil
}

func (handler *SetInfoHandler) SetSurnameHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	handler.gs.SetName(params.message.Text)

	handler.nextCommand = "/set_age"

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your surname"}
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, "*")

	return res, nil
}

func (handler *SetInfoHandler) SetAgeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	handler.gs.SetSurname(params.message.Text)

	handler.nextCommand = "/set_info_end"

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your age"}
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, "*")

	return res, nil
}

func (handler *SetInfoHandler) SetInfoEndHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	age, err := strconv.Atoi(params.message.Text)

	if err != nil {
		retMessage := bottypes.Message{ChatID: chatID, Text: "Wrong age, try again!"}
		res.messages = append(res.messages, retMessage)
		res.nextCommands = append(res.nextCommands, "*", "/back_state")
		return res, nil
	}

	handler.gs.SetAge(age)

	retMessage := bottypes.Message{ChatID: chatID, Text: "My gratitude"}
	res.messages = append(res.messages, retMessage)

	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")
	res.nextState = "start-state"

	return res, nil
}
