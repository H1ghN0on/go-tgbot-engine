package handlers_example

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands/example"
)

type KeyboardHandler struct {
	Handler
}

func NewKeyboardHandler(gs ExampleGlobalStater) *KeyboardHandler {

	h := &KeyboardHandler{}
	h.gs = gs

	h.Commands = map[bottypes.Command][]func(params handlers.HandlerParams) (handlers.HandlerResponse, error){
		cmd.KeyboardStartCommand:  {h.ModifyHandler(h.KeyboardStartHandler, []int{handlers.RemovableByTrigger})},
		cmd.KeyboardOneCommand:    {h.ModifyHandler(h.KeyboardOneHandler, []int{handlers.StateBackable, handlers.KeyboardStarter, handlers.RemovableByTrigger})},
		cmd.KeyboardTwoCommand:    {h.ModifyHandler(h.KeyboardTwoHandler, []int{handlers.CommandBackable, handlers.KeyboardStarter, handlers.RemovableByTrigger})},
		cmd.KeyboardThreeCommand:  {h.ModifyHandler(h.KeyboardThreeHandler, []int{handlers.CommandBackable, handlers.RemovableByTrigger, handlers.KeyboardStarter})},
		cmd.KeyboardFinishCommand: {h.ModifyHandler(h.KeyboardFinishHandler, []int{handlers.RemoveTriggerer, handlers.KeyboardStopper})},
	}

	return h
}

func (handler *KeyboardHandler) Handle(params handlers.HandlerParams) ([]handlers.HandlerResponse, error) {
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

func (handler *KeyboardHandler) KeyboardStartHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.KeyboardOneCommand)
	res.NextState = "keyboard-state"

	return res, nil
}

func (handler *KeyboardHandler) KeyboardOneHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

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
	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = append(res.NextCommands, cmd.KeyboardTwoCommand)

	return res, nil
}

func (handler *KeyboardHandler) KeyboardTwoHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

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
	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = append(res.NextCommands, cmd.KeyboardThreeCommand)

	return res, nil
}

func (handler *KeyboardHandler) KeyboardThreeHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

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

	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = append(res.NextCommands, cmd.KeyboardFinishCommand)

	return res, nil
}

func (handler *KeyboardHandler) KeyboardFinishHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Alabama certified moment"}
	res.Messages = append(res.Messages, retMessage)

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	res.NextState = "start-state"

	return res, nil
}
