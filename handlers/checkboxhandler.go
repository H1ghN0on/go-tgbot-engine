package handlers

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
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

	sh.commands = map[string][]func(params HandlerParams) (HandlerResponse, error){
		"/checkboxes_start": {sh.SetCheckboxesStartHandler,
			sh.ModifyHandler(sh.InitCheckboxesHandler, []int{CheckboxableTwo, KeyboardStarter})},
		"/checkboxes_first":  {sh.ModifyHandler(sh.FirstCheckboxHandler, []int{CheckboxableTwo, KeyboardStarter})},
		"/checkboxes_second": {sh.ModifyHandler(sh.SecondCheckboxHandler, []int{CheckboxableTwo, KeyboardStarter})},
		"/checkboxes_third":  {sh.ModifyHandler(sh.ThirdCheckboxHandler, []int{CheckboxableTwo, KeyboardStarter})},
		"/checkboxes_fourth": {sh.ModifyHandler(sh.FourthCheckboxHandler, []int{CheckboxableTwo, KeyboardStarter})},
		"/checkboxes_accept": {sh.ModifyHandler(sh.AcceptCheckboxHandler, []int{KeyboardStopper})},
		"/nothingness":       {sh.ModifyHandler(sh.EmptyHandler, []int{Nothingness})},
	}

	return sh
}

func (handler *CheckboxHandler) InitHandler() {
	handler.firstCheckbox = false
	handler.secondCheckbox = false
	handler.thirdCheckbox = false
	handler.fourthCheckbox = false
}

func (handler *CheckboxHandler) Handle(command string, params HandlerParams) ([]HandlerResponse, bool, error) {
	var res []HandlerResponse

	handleFuncs, ok := handler.Handler.commands[command]
	if !ok {
		panic("wrong handler")
	}

	isFinished := command == "/checkboxes_accept"

	for _, handleFunc := range handleFuncs {
		response, err := handleFunc(params)
		if err != nil {
			return []HandlerResponse{}, isFinished, err
		}
		res = append(res, response)
	}

	return res, isFinished, nil
}

func (handler *CheckboxHandler) DeinitHandler() {

}

func (handler *CheckboxHandler) SetCheckboxesStartHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Let's try some checkboxes"}
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages, nextState: "checkbox-state"}, nil
}

func (handler *CheckboxHandler) InitCheckboxesHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID
	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.messages = append(res.messages, retMessage)
	return HandlerResponse{messages: res.messages}, nil
}

func (handler *CheckboxHandler) FirstCheckboxHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	handler.firstCheckbox = !handler.firstCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}, nil
}

func (handler *CheckboxHandler) SecondCheckboxHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	handler.secondCheckbox = !handler.secondCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}, nil
}

func (handler *CheckboxHandler) ThirdCheckboxHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	handler.thirdCheckbox = !handler.thirdCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}, nil
}

func (handler *CheckboxHandler) FourthCheckboxHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	handler.fourthCheckbox = !handler.fourthCheckbox

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}
	retMessage.ButtonRows = handler.gatherAllCheckboxes(chatID)

	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}, nil
}

func (handler *CheckboxHandler) AcceptCheckboxHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

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

	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")

	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}, nil
}

func (handler *CheckboxHandler) gatherAllCheckboxes(chatID int64) []bottypes.ButtonRows {

	return []bottypes.ButtonRows{
		{
			CheckboxButtons: []bottypes.CheckboxButton{
				{ChatID: chatID, Text: "Checkbox one", Command: bottypes.Command{Text: "/checkboxes_first"}, Active: handler.firstCheckbox},
			},
		},
		{
			CheckboxButtons: []bottypes.CheckboxButton{
				{ChatID: chatID, Text: "Checkbox two", Command: bottypes.Command{Text: "/checkboxes_second"}, Active: handler.secondCheckbox},
			},
		},
		{
			CheckboxButtons: []bottypes.CheckboxButton{
				{ChatID: chatID, Text: "Checkbox three", Command: bottypes.Command{Text: "/checkboxes_third"}, Active: handler.thirdCheckbox},
			},
		},
		{
			CheckboxButtons: []bottypes.CheckboxButton{
				{ChatID: chatID, Text: "Checkbox four", Command: bottypes.Command{Text: "/checkboxes_fourth"}, Active: handler.fourthCheckbox},
			},
		},
		{
			Buttons: []bottypes.Button{
				{ChatID: chatID, Text: "Accept", Command: bottypes.Command{Text: "/checkboxes_accept"}},
			},
		},
	}
}
