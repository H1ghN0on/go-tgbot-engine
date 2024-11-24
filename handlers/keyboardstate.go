package handlers

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type KeyboardHandler struct {
	Handler
}

func NewKeyboardHandler(gs GlobalStater) *KeyboardHandler {

	h := &KeyboardHandler{}
	h.gs = gs

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.KeyboardStartCommand:  {h.ModifyHandler(h.KeyboardStartHandler, []int{RemovableByTrigger})},
		cmd.KeyboardOneCommand:    {h.ModifyHandler(h.KeyboardOneHandler, []int{StateBackable, KeyboardStarter, RemovableByTrigger})},
		cmd.KeyboardTwoCommand:    {h.ModifyHandler(h.KeyboardTwoHandler, []int{CommandBackable, KeyboardStarter, RemovableByTrigger})},
		cmd.KeyboardThreeCommand:  {h.ModifyHandler(h.KeyboardThreeHandler, []int{CommandBackable, RemovableByTrigger, KeyboardStarter})},
		cmd.KeyboardFinishCommand: {h.ModifyHandler(h.KeyboardFinishHandler, []int{RemoveTriggerer, KeyboardStopper})},
	}

	return h
}

func (handler *KeyboardHandler) Handle(params HandlerParams) ([]HandlerResponse, error) {
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

func (handler *KeyboardHandler) KeyboardStartHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.KeyboardOneCommand)
	res.nextState = "keyboard-state"

	return res, nil
}

func (handler *KeyboardHandler) KeyboardOneHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 1"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Let me", Command: cmd.KeyboardTwoCommand},
			{ChatID: chatID, Text: "No, Let me!", Command: cmd.KeyboardTwoCommand},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "LET ME!", Command: cmd.KeyboardTwoCommand},
			{ChatID: chatID, Text: "l-ll-let me *blushes*", Command: cmd.KeyboardTwoCommand},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, cmd.KeyboardTwoCommand)

	return res, nil
}

func (handler *KeyboardHandler) KeyboardTwoHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 2"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Push me", Command: cmd.KeyboardThreeCommand},
			{ChatID: chatID, Text: "No, push me!", Command: cmd.KeyboardThreeCommand},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "PUSH ME!", Command: cmd.KeyboardThreeCommand},
			{ChatID: chatID, Text: "p-pp-push me *blushes*", Command: cmd.KeyboardThreeCommand},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, cmd.KeyboardThreeCommand)

	return res, nil
}

func (handler *KeyboardHandler) KeyboardThreeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 3"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Approach me", Command: cmd.KeyboardFinishCommand},
			{ChatID: chatID, Text: "No, approach me!", Command: cmd.KeyboardFinishCommand},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "APPROACH ME!", Command: cmd.KeyboardFinishCommand},
			{ChatID: chatID, Text: "a-aa-approach me *blushes*", Command: cmd.KeyboardFinishCommand},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)

	res.messages = append(res.messages, retMessage)
	res.nextCommands = append(res.nextCommands, cmd.KeyboardFinishCommand)

	return res, nil
}

func (handler *KeyboardHandler) KeyboardFinishHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Alabama certified moment"}
	res.messages = append(res.messages, retMessage)

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.ShowCommandsCommand)
	res.nextState = "start-state"

	return res, nil
}
