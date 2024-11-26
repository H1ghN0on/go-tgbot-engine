package handlers_example

import (
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands/example"
)

type SetInfoHandler struct {
	Handler

	name    string
	surname string
	age     int
}

func NewSetInfoHandler(gs ExampleGlobalStater) *SetInfoHandler {

	sh := &SetInfoHandler{}
	sh.gs = gs

	sh.Commands = map[bottypes.Command][]func(params handlers.HandlerParams) (handlers.HandlerResponse, error){
		cmd.SetInfoStartCommand: {sh.SetInfoStartHandler},

		cmd.SetNameCommand:    {sh.ModifyHandler(sh.SetNameHandler, []int{handlers.StateBackable})},
		cmd.SetSurnameCommand: {sh.ModifyHandler(sh.SetSurnameHandler, []int{handlers.CommandBackable})},
		cmd.SetAgeCommand:     {sh.ModifyHandler(sh.SetAgeHandler, []int{handlers.CommandBackable})},
		cmd.SetInfoEndCommand: {sh.SetInfoEndHandler},
	}

	return sh
}

func (handler *SetInfoHandler) Handle(params handlers.HandlerParams) ([]handlers.HandlerResponse, error) {
	var res []handlers.HandlerResponse

	handleFuncs, ok := handler.GetCommandFromMap(params.Command)

	if !ok {
		panic("wrong handler")
	}

	for _, handleFunc := range handleFuncs {
		response, err := handleFunc(params)
		if err != nil {
			return []handlers.HandlerResponse{}, err
		}
		res = append(res, response)
	}

	return res, nil
}

func (handler *SetInfoHandler) SetInfoStartHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	handler.name = ""
	handler.surname = ""
	handler.age = 0

	retMessage := bottypes.Message{ChatID: chatID, Text: "Let me know you a bit closer"}
	res.Messages = append(res.Messages, retMessage)
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.SetNameCommand)
	res.NextState = "info-state"

	return res, nil
}

func (handler *SetInfoHandler) SetNameHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your name"}
	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = append(res.NextCommands, cmd.SetSurnameCommand)
	res.NextCommandToParse = bottypes.ParseableCommand{
		Command:   cmd.SetSurnameCommand,
		ParseType: bottypes.AnyTextParse,
	}

	return res, nil
}

func (handler *SetInfoHandler) SetSurnameHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	handler.name = params.Command.Data

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your surname"}
	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = append(res.NextCommands, cmd.SetAgeCommand)
	res.NextCommandToParse = bottypes.ParseableCommand{
		Command:   cmd.SetAgeCommand,
		ParseType: bottypes.AnyTextParse,
	}

	return res, nil
}

func (handler *SetInfoHandler) SetAgeHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	handler.surname = params.Command.Data

	retMessage := bottypes.Message{ChatID: chatID, Text: "Enter your age"}
	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = append(res.NextCommands, cmd.SetInfoEndCommand)
	res.NextCommandToParse = bottypes.ParseableCommand{
		Command:   cmd.SetInfoEndCommand,
		ParseType: bottypes.AnyTextParse,
	}

	return res, nil
}

func (handler *SetInfoHandler) SetInfoEndHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	age, err := strconv.Atoi(params.Command.Data)

	if err != nil || age < 0 || age > 100 {
		retMessage := bottypes.Message{ChatID: chatID, Text: "Wrong age, try again!"}
		res.Messages = append(res.Messages, retMessage)
		res.NextCommands = append(res.NextCommands, cmd.SetInfoEndCommand)
		res.NextCommandToParse = bottypes.ParseableCommand{
			Command:   cmd.SetInfoEndCommand,
			ParseType: bottypes.AnyTextParse,
		}

		return res, nil
	}

	handler.age = age

	retMessage := bottypes.Message{ChatID: chatID, Text: "My gratitude"}
	res.Messages = append(res.Messages, retMessage)

	handler.gs.SetName(handler.name)
	handler.gs.SetSurname(handler.surname)
	handler.gs.SetAge(handler.age)

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	res.NextState = "start-state"

	return res, nil
}
