package handlers

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

type KeyboardHandler struct {
	Handler
}

func NewKeyboardhandler(gs GlobalStater) *KeyboardHandler {

	h := &KeyboardHandler{}
	h.gs = gs

	h.commands = map[string][]func(params HandlerParams) HandlerResponse{
		"/keyboard_start":  {h.ModifyHandler(h.KeyboardStartHandler, []int{RemovableByTrigger})},
		"/keyboard_one":    {h.ModifyHandler(h.KeyboardOneHandler, []int{StateBackable, KeyboardStarter, RemovableByTrigger})},
		"/keyboard_two":    {h.ModifyHandler(h.KeyboardTwoHandler, []int{CommandBackable, KeyboardStarter, RemovableByTrigger})},
		"/keyboard_three":  {h.ModifyHandler(h.KeyboardThreeHandler, []int{CommandBackable, RemovableByTrigger, KeyboardStarter})},
		"/keyboard_finish": {h.ModifyHandler(h.KeyboardFinishHandler, []int{RemoveTriggerer, KeyboardStopper})},
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

	res.postCommandsHandle = append(res.postCommandsHandle, "/keyboard_one")

	return HandlerResponse{nextState: "keyboard-state", postCommandsHandle: res.postCommandsHandle}
}

func (handler *Handler) KeyboardOneHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 1"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Let me", Command: bottypes.Command{Text: "/keyboard_two"}},
			{ChatID: chatID, Text: "No, Let me!", Command: bottypes.Command{Text: "/keyboard_two"}},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "LET ME!", Command: bottypes.Command{Text: "/keyboard_two"}},
			{ChatID: chatID, Text: "l-ll-let me *blushes*", Command: bottypes.Command{Text: "/keyboard_two"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}
}

func (handler *Handler) KeyboardTwoHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 2"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Push me", Command: bottypes.Command{Text: "/keyboard_three"}},
			{ChatID: chatID, Text: "No, push me!", Command: bottypes.Command{Text: "/keyboard_three"}},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "PUSH ME!", Command: bottypes.Command{Text: "/keyboard_three"}},
			{ChatID: chatID, Text: "p-pp-push me *blushes*", Command: bottypes.Command{Text: "/keyboard_three"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}
}

func (handler *Handler) KeyboardThreeHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 3"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Approach me", Command: bottypes.Command{Text: "/keyboard_finish"}},
			{ChatID: chatID, Text: "No, approach me!", Command: bottypes.Command{Text: "/keyboard_finish"}},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "APPROACH ME!", Command: bottypes.Command{Text: "/keyboard_finish"}},
			{ChatID: chatID, Text: "a-aa-approach me *blushes*", Command: bottypes.Command{Text: "/keyboard_finish"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}
}

func (handler *Handler) KeyboardFinishHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Alabama certified moment"}
	res.messages = append(res.messages, retMessage)

	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")
	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}
}
