package handlers

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

type KeyboardHandler struct {
	Handler
}

func NewKeyboardhandler(gs GlobalStater) *KeyboardHandler {

	h := &KeyboardHandler{}
	h.gs = gs

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		"/keyboard_start":  {h.ModifyHandler(h.KeyboardStartHandler, []int{RemovableByTrigger})},
		"/keyboard_one":    {h.ModifyHandler(h.KeyboardOneHandler, []int{StateBackable, KeyboardStarter, RemovableByTrigger})},
		"/keyboard_two":    {h.ModifyHandler(h.KeyboardTwoHandler, []int{CommandBackable, KeyboardStarter, RemovableByTrigger})},
		"/keyboard_three":  {h.ModifyHandler(h.KeyboardThreeHandler, []int{CommandBackable, RemovableByTrigger, KeyboardStarter})},
		"/keyboard_finish": {h.ModifyHandler(h.KeyboardFinishHandler, []int{RemoveTriggerer, KeyboardStopper})},
	}

	return h
}

func (handler *KeyboardHandler) Handle(command bottypes.Command, params HandlerParams) ([]HandlerResponse, error) {
	var res []HandlerResponse

	handleFuncs, ok := handler.Handler.commands[command]
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

func (handler *KeyboardHandler) KeyboardStartHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	res.postCommandsHandle = append(res.postCommandsHandle, "/keyboard_one")
	res.nextState = "keyboard-state"

	return res, nil
}

func (handler *KeyboardHandler) KeyboardOneHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 1"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Let me", Command: "/keyboard_two"},
			{ChatID: chatID, Text: "No, Let me!", Command: "/keyboard_two"},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "LET ME!", Command: "/keyboard_two"},
			{ChatID: chatID, Text: "l-ll-let me *blushes*", Command: "/keyboard_two"},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, "/keyboard_two")

	return res, nil
}

func (handler *KeyboardHandler) KeyboardTwoHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 2"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Push me", Command: "/keyboard_three"},
			{ChatID: chatID, Text: "No, push me!", Command: "/keyboard_three"},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "PUSH ME!", Command: "/keyboard_three"},
			{ChatID: chatID, Text: "p-pp-push me *blushes*", Command: "/keyboard_three"},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, "/keyboard_three")

	return res, nil
}

func (handler *KeyboardHandler) KeyboardThreeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 3"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Approach me", Command: "/keyboard_finish"},
			{ChatID: chatID, Text: "No, approach me!", Command: "/keyboard_finish"},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "APPROACH ME!", Command: "/keyboard_finish"},
			{ChatID: chatID, Text: "a-aa-approach me *blushes*", Command: "/keyboard_finish"},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)

	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, "/keyboard_finish")

	return res, nil
}

func (handler *KeyboardHandler) KeyboardFinishHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Alabama certified moment"}
	res.messages = append(res.messages, retMessage)

	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")
	res.nextState = "start-state"

	return res, nil
}
