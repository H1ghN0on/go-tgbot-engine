package handlers

import (
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
)

type SetInfoHandler SequentHandler

func NewSetInfoHandler(gs GlobalStater) *SetInfoHandler {

	sh := &SetInfoHandler{}
	sh.gs = gs

	sh.commands = map[string][]func(params HandlerParams) (HandlerResponse, error){
		"/set_info_start": {
			sh.SetInfoStartHandler,
		},
		"/set_name":     {sh.ModifyHandler(sh.SetNameHandler, []int{StateBackable})},
		"/set_surname":  {sh.ModifyHandler(sh.SetSurnameHandler, []int{CommandBackable})},
		"/set_age":      {sh.ModifyHandler(sh.SetAgeHandler, []int{CommandBackable})},
		"/set_info_end": {sh.SetInfoEndHandler},
	}

	sh.commandSequence = map[int]string{
		0: "/set_info_start",
		1: "/set_name",
		2: "/set_surname",
		3: "/set_age",
		4: "/set_info_end",
	}

	sh.active = 0

	return sh
}

func (handler *SetInfoHandler) InitHandler() {
	handler.active = 0
}

func (handler *SetInfoHandler) Handle(command string, params HandlerParams) ([]HandlerResponse, bool, error) {
	var res []HandlerResponse
	currentCommand, ok := handler.commandSequence[handler.active]
	if !ok {
		panic("wrong handler")
	}

	handleFuncs, ok := handler.Handler.commands[currentCommand]
	if !ok {
		panic("wrong handler")
	}

	for _, handleFunc := range handleFuncs {
		response, err := handleFunc(params)
		if err != nil {
			isFinished := handler.active == len(handler.commandSequence)
			return []HandlerResponse{}, isFinished, err
		}
		res = append(res, response)
	}

	handler.active++

	isFinished := handler.active == len(handler.commandSequence)

	return res, isFinished, nil
}

func (handler *SetInfoHandler) DeinitHandler() {
	handler.active = 0
}

func (handler *SetInfoHandler) SetInfoStartHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Let me know you a bit closer"}
	res.messages = append(res.messages, retMessage)
	res.postCommandsHandle = append(res.postCommandsHandle, "/set_name")

	return HandlerResponse{messages: res.messages, nextState: "info-state", postCommandsHandle: res.postCommandsHandle}, nil
}

func (handler *Handler) SetNameHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your name"}
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}, nil
}

func (handler *Handler) SetSurnameHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	handler.gs.SetName(params.message.Text)

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your surname"}
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}, nil
}

func (handler *Handler) SetAgeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	handler.gs.SetSurname(params.message.Text)

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your age"}
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}, nil
}

func (handler *Handler) SetInfoEndHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	age, err := strconv.Atoi(params.message.Text)

	if err != nil {
		return HandlerResponse{}, HandlerResponseError{message: "Wrong age"}
	}

	handler.gs.SetAge(age)

	retMessage := bottypes.Message{ChatID: chatID, Text: "My gratitude"}
	res.messages = append(res.messages, retMessage)

	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")

	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}, nil
}
