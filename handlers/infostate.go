package handlers

import (
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type SetInfoHandler struct {
	Handler

	name    string
	surname string
	age     int
}

func NewSetInfoHandler(gs GlobalStater) *SetInfoHandler {

	sh := &SetInfoHandler{}
	sh.gs = gs

	sh.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.SetInfoStartCommand: {sh.SetInfoStartHandler},

		cmd.SetNameCommand:    {sh.ModifyHandler(sh.SetNameHandler, []int{StateBackable})},
		cmd.SetSurnameCommand: {sh.ModifyHandler(sh.SetSurnameHandler, []int{CommandBackable})},
		cmd.SetAgeCommand:     {sh.ModifyHandler(sh.SetAgeHandler, []int{CommandBackable})},
		cmd.SetInfoEndCommand: {sh.SetInfoEndHandler},
	}

	return sh
}

func (handler *SetInfoHandler) Handle(params HandlerParams) ([]HandlerResponse, error) {
	var res []HandlerResponse

	handleFuncs, ok := handler.GetCommandFromMap(params.command)

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

	handler.name = ""
	handler.surname = ""
	handler.age = 0

	retMessage := bottypes.Message{ChatID: chatID, Text: "Let me know you a bit closer"}
	res.messages = append(res.messages, retMessage)
	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.SetNameCommand)
	res.nextState = "info-state"

	return res, nil
}

func (handler *SetInfoHandler) SetNameHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your name"}
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, cmd.SetSurnameCommand)
	res.nextCommandToParse = bottypes.ParseableCommand{
		Command:   cmd.SetSurnameCommand,
		ParseType: bottypes.AnyTextParse,
	}

	return res, nil
}

func (handler *SetInfoHandler) SetSurnameHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	handler.name = params.command.Data

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your surname"}
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, cmd.SetAgeCommand)
	res.nextCommandToParse = bottypes.ParseableCommand{
		Command:   cmd.SetAgeCommand,
		ParseType: bottypes.AnyTextParse,
	}

	return res, nil
}

func (handler *SetInfoHandler) SetAgeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	handler.surname = params.command.Data

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your age"}
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, cmd.SetInfoEndCommand)
	res.nextCommandToParse = bottypes.ParseableCommand{
		Command:   cmd.SetInfoEndCommand,
		ParseType: bottypes.AnyTextParse,
	}

	return res, nil
}

func (handler *SetInfoHandler) SetInfoEndHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	age, err := strconv.Atoi(params.command.Data)

	if err != nil || age < 0 || age > 100 {
		retMessage := bottypes.Message{ChatID: chatID, Text: "Wrong age, try again!"}
		res.messages = append(res.messages, retMessage)
		res.nextCommands = append(res.nextCommands, cmd.SetInfoEndCommand)
		res.nextCommandToParse = bottypes.ParseableCommand{
			Command:   cmd.SetInfoEndCommand,
			ParseType: bottypes.AnyTextParse,
		}

		return res, nil
	}

	handler.age = age

	retMessage := bottypes.Message{ChatID: chatID, Text: "My gratitude"}
	res.messages = append(res.messages, retMessage)

	handler.gs.SetName(handler.name)
	handler.gs.SetSurname(handler.surname)
	handler.gs.SetAge(handler.age)

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.ShowCommandsCommand)
	res.nextState = "start-state"

	return res, nil
}
