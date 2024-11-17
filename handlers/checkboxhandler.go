package handlers

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type CheckboxHandler struct {
	Handler
	firstCheckbox  bool
	secondCheckbox bool
	thirdCheckbox  bool
	fourthCheckbox bool
}

func NewCheckboxHandler(gs GlobalStater) *CheckboxHandler {

	sh := &CheckboxHandler{}
	sh.gs = gs

	sh.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.CheckboxStartCommand: {sh.SetCheckboxesStartHandler,
			sh.ModifyHandler(sh.InitCheckboxesHandler, []int{CheckboxableOne, KeyboardStarter, StateBackable})},
		cmd.CheckboxFirstCommand:  {sh.ModifyHandler(sh.FirstCheckboxHandler, []int{CheckboxableOne, KeyboardStarter, StateBackable})},
		cmd.CheckboxSecondCommand: {sh.ModifyHandler(sh.SecondCheckboxHandler, []int{CheckboxableOne, KeyboardStarter, StateBackable})},
		cmd.CheckboxThirdCommand:  {sh.ModifyHandler(sh.ThirdCheckboxHandler, []int{CheckboxableOne, KeyboardStarter, StateBackable})},
		cmd.CheckboxFourthCommand: {sh.ModifyHandler(sh.FourthCheckboxHandler, []int{CheckboxableOne, KeyboardStarter, StateBackable})},
		cmd.CheckboxAcceptCommand: {sh.ModifyHandler(sh.AcceptCheckboxHandler, []int{KeyboardStopper})},
		cmd.NothingnessCommand:    {sh.ModifyHandler(sh.EmptyHandler, []int{Nothingness})},
	}

	return sh
}

func (handler *CheckboxHandler) Handle(params HandlerParams) ([]HandlerResponse, error) {
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

func (handler *CheckboxHandler) SetCheckboxesStartHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Let's try some checkboxes"}
	res.messages = append(res.messages, retMessage)
	res.nextState = "checkbox-state"

	return res, nil
}

func (handler *CheckboxHandler) InitCheckboxesHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID
	retMessage := bottypes.Message{ChatID: chatID, Text: ""}

	handler.firstCheckbox = false
	handler.secondCheckbox = false
	handler.thirdCheckbox = false
	handler.fourthCheckbox = false

	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.messages = append(res.messages, retMessage)
	res.nextCommands = handler.GetCommands()

	return res, nil
}

func (handler *CheckboxHandler) FirstCheckboxHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	handler.firstCheckbox = !handler.firstCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.messages = append(res.messages, retMessage)
	res.nextCommands = handler.GetCommands()

	return res, nil
}

func (handler *CheckboxHandler) SecondCheckboxHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	handler.secondCheckbox = !handler.secondCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.messages = append(res.messages, retMessage)
	res.nextCommands = handler.GetCommands()

	return res, nil
}

func (handler *CheckboxHandler) ThirdCheckboxHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	handler.thirdCheckbox = !handler.thirdCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.messages = append(res.messages, retMessage)
	res.nextCommands = handler.GetCommands()

	return res, nil
}

func (handler *CheckboxHandler) FourthCheckboxHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	handler.fourthCheckbox = !handler.fourthCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.messages = append(res.messages, retMessage)
	res.nextCommands = handler.GetCommands()

	return res, nil
}

func (handler *CheckboxHandler) AcceptCheckboxHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "You selected"})

	if handler.firstCheckbox {
		res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "First"})
	}

	if handler.secondCheckbox {
		res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Second"})
	}

	if handler.thirdCheckbox {
		res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Third"})
	}

	if handler.fourthCheckbox {
		res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Fourth"})
	}

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.ShowCommandsCommand)
	res.nextState = "start-state"

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
