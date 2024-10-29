package handlers

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

type KeyboardHandler struct {
	Handler
}

func NewKeyboardhandler(gs GlobalStater) *KeyboardHandler {

	h := &KeyboardHandler{}
	h.gs = gs

	h.commands = map[string][]func(params HandlerParams) HandlerResponse{
		"/keyboard_start": {h.ModifyHandler(h.KeyboardStartHandler, []int{StateBackable, KeyboardStarter, RemovableByTrigger})},
		"/keyboard_one":   {h.ModifyHandler(h.KeyboardOneHandler, []int{CommandBackable, KeyboardStarter, RemovableByTrigger})},
		"/keyboard_two":   {h.ModifyHandler(h.KeyboardTwoHandler, []int{CommandBackable, KeyboardStarter, RemovableByTrigger})},
		"/keyboard_three": {h.ModifyHandler(h.KeyboardThreeHandler, []int{RemoveTriggerer, KeyboardStopper})},
	}

	return h
}

func (handler *KeyboardHandler) InitHandler() {

}

func (handler *KeyboardHandler) Handle(command string, params HandlerParams) ([]HandlerResponse, bool) {
	var res []HandlerResponse

	handleFuncs, ok := handler.Handler.commands[command]
	if !ok {
		panic("wrong handler")
	}

	for _, handleFunc := range handleFuncs {
		response := handleFunc(params)
		res = append(res, response)
	}

	return res, true
}

func (handler *KeyboardHandler) DeinitHandler() {

}

func (handler *Handler) KeyboardStartHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 1"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Let me", Command: bottypes.Command{Text: "/keyboard_one"}},
			{ChatID: chatID, Text: "No, Let me!", Command: bottypes.Command{Text: "/keyboard_one"}},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "LET ME!", Command: bottypes.Command{Text: "/keyboard_one"}},
			{ChatID: chatID, Text: "l-ll-let me *blushes*", Command: bottypes.Command{Text: "/keyboard_one"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages, nextState: "keyboard-state"}
}

func (handler *Handler) KeyboardOneHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Push me", Command: bottypes.Command{Text: "/keyboard_two"}},
			{ChatID: chatID, Text: "No, push me!", Command: bottypes.Command{Text: "/keyboard_two"}},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "PUSH ME!", Command: bottypes.Command{Text: "/keyboard_two"}},
			{ChatID: chatID, Text: "p-pp-push me *blushes*", Command: bottypes.Command{Text: "/keyboard_two"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}
}

func (handler *Handler) KeyboardTwoHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 3"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Approach me", Command: bottypes.Command{Text: "/keyboard_three"}},
			{ChatID: chatID, Text: "No, approach me!", Command: bottypes.Command{Text: "/keyboard_three"}},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "APPROACH ME!", Command: bottypes.Command{Text: "/keyboard_three"}},
			{ChatID: chatID, Text: "a-aa-approach me *blushes*", Command: bottypes.Command{Text: "/keyboard_three"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}
}

func (handler *Handler) KeyboardThreeHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Alabama certified moment"}
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages, nextState: "start-state"}
}
