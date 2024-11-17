package handlers

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type DynamicKeyboardHandler struct {
	Handler
}

func NewDynamicKeyboardhandler(gs GlobalStater) *DynamicKeyboardHandler {

	h := &DynamicKeyboardHandler{}
	h.gs = gs

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.DynamicKeyboardStartCommand:       {h.ModifyHandler(h.DynamicKeyboardStartHandler, []int{RemovableByTrigger})},
		cmd.DynamicKeyboardFirstStageCommand:  {h.ModifyHandler(h.DynamicKeyboardFirstHandler, []int{KeyboardStarter, RemovableByTrigger})},
		cmd.DynamicKeyboardSecondStageCommand: {h.ModifyHandler(h.DynamicKeyboardSecondHandler, []int{KeyboardStarter, RemovableByTrigger})},
	}

	return h
}

func (handler *DynamicKeyboardHandler) Handle(params HandlerParams) ([]HandlerResponse, error) {
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

func (handler *DynamicKeyboardHandler) DynamicKeyboardStartHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	res.postCommandsHandle = append(res.postCommandsHandle, cmd.DynamicKeyboardFirstStageCommand)
	res.nextState = "dynamic-keyboard-state"

	return res, nil
}

func (handler *DynamicKeyboardHandler) DynamicKeyboardFirstHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Stack songs one"}

	data := handler.gs.GetDataForDynamicKeyboard()

	for _, text := range data["first_stage"] {
		retMessage.ButtonRows = append(retMessage.ButtonRows, bottypes.ButtonRows{
			Buttons: []bottypes.Button{
				{ChatID: chatID, Text: text, Command: cmd.KeyboardTwoCommand, Data: text},
			},
		})
	}

	res.messages = append(res.messages, retMessage)

	return res, nil
}

func (handler *DynamicKeyboardHandler) DynamicKeyboardSecondHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	res.postCommandsHandle = append(res.postCommandsHandle, cmd.KeyboardOneCommand)
	res.nextState = "dynamic-keyboard-state"

	return res, nil
}
