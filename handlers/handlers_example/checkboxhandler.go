package handlers_example

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type CheckboxHandler struct {
	Handler
	firstCheckbox  bool
	secondCheckbox bool
	thirdCheckbox  bool
	fourthCheckbox bool
}

func NewCheckboxHandler(gs ExampleGlobalStater) *CheckboxHandler {

	sh := &CheckboxHandler{}
	sh.gs = gs

	sh.Commands = map[bottypes.Command][]func(params handlers.HandlerParams) (handlers.HandlerResponse, error){
		cmd.CheckboxStartCommand: {sh.SetCheckboxesStartHandler,
			sh.ModifyHandler(sh.InitCheckboxesHandler, []int{handlers.CheckboxableOne, handlers.KeyboardStarter, handlers.StateBackable})},
		cmd.CheckboxFirstCommand:  {sh.ModifyHandler(sh.FirstCheckboxHandler, []int{handlers.CheckboxableOne, handlers.KeyboardStarter, handlers.StateBackable})},
		cmd.CheckboxSecondCommand: {sh.ModifyHandler(sh.SecondCheckboxHandler, []int{handlers.CheckboxableOne, handlers.KeyboardStarter, handlers.StateBackable})},
		cmd.CheckboxThirdCommand:  {sh.ModifyHandler(sh.ThirdCheckboxHandler, []int{handlers.CheckboxableOne, handlers.KeyboardStarter, handlers.StateBackable})},
		cmd.CheckboxFourthCommand: {sh.ModifyHandler(sh.FourthCheckboxHandler, []int{handlers.CheckboxableOne, handlers.KeyboardStarter, handlers.StateBackable})},
		cmd.CheckboxAcceptCommand: {sh.ModifyHandler(sh.AcceptCheckboxHandler, []int{handlers.KeyboardStopper})},
		cmd.NothingnessCommand:    {sh.ModifyHandler(sh.EmptyHandler, []int{handlers.Nothingness})},
	}

	return sh
}

func (handler *CheckboxHandler) Handle(params handlers.HandlerParams) ([]handlers.HandlerResponse, error) {
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

func (handler *CheckboxHandler) SetCheckboxesStartHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Let's try some checkboxes"}
	res.Messages = append(res.Messages, retMessage)
	res.NextState = "checkbox-state"

	return res, nil
}

func (handler *CheckboxHandler) InitCheckboxesHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID
	retMessage := bottypes.Message{ChatID: chatID, Text: ""}

	handler.firstCheckbox = false
	handler.secondCheckbox = false
	handler.thirdCheckbox = false
	handler.fourthCheckbox = false

	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = handler.GetCommands()

	return res, nil
}

func (handler *CheckboxHandler) FirstCheckboxHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	handler.firstCheckbox = !handler.firstCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = handler.GetCommands()

	return res, nil
}

func (handler *CheckboxHandler) SecondCheckboxHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	handler.secondCheckbox = !handler.secondCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = handler.GetCommands()

	return res, nil
}

func (handler *CheckboxHandler) ThirdCheckboxHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	handler.thirdCheckbox = !handler.thirdCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = handler.GetCommands()

	return res, nil
}

func (handler *CheckboxHandler) FourthCheckboxHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	handler.fourthCheckbox = !handler.fourthCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = handler.GetCommands()

	return res, nil
}

func (handler *CheckboxHandler) AcceptCheckboxHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "You selected"})

	if handler.firstCheckbox {
		res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "First"})
	}

	if handler.secondCheckbox {
		res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "Second"})
	}

	if handler.thirdCheckbox {
		res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "Third"})
	}

	if handler.fourthCheckbox {
		res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "Fourth"})
	}

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	res.NextState = "start-state"

	return res, nil
}

func (handler *CheckboxHandler) gatherAllCheckboxes(chatID int64) []bottypes.ButtonRows {

	return []bottypes.ButtonRows{
		{
			CheckboxButtons: []bottypes.CheckboxButton{
				{ChatID: chatID, Text: "Checkbox one", Command: cmd.CheckboxFirstCommand, Active: handler.firstCheckbox},
			},
		},
		{
			CheckboxButtons: []bottypes.CheckboxButton{
				{ChatID: chatID, Text: "Checkbox two", Command: cmd.CheckboxSecondCommand, Active: handler.secondCheckbox},
			},
		},
		{
			CheckboxButtons: []bottypes.CheckboxButton{
				{ChatID: chatID, Text: "Checkbox three", Command: cmd.CheckboxThirdCommand, Active: handler.thirdCheckbox},
			},
		},
		{
			CheckboxButtons: []bottypes.CheckboxButton{
				{ChatID: chatID, Text: "Checkbox four", Command: cmd.CheckboxFourthCommand, Active: handler.fourthCheckbox},
			},
		},
		{
			Buttons: []bottypes.Button{
				{ChatID: chatID, Text: "Accept", Command: cmd.CheckboxAcceptCommand},
			},
		},
	}
}
