package handlers

import (
	"fmt"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
)

type BackHandler struct {
	Handler
	sm           StateMachiner
	commandQueue []string
}

func NewBackHandler(gs GlobalStater, sm StateMachiner) *BackHandler {

	h := &BackHandler{}
	h.gs = gs
	h.sm = sm

	h.commands = map[string][]func(params HandlerParams) HandlerResponse{
		"/back_state":   {h.ModifyHandler(h.BackStateHandler, []int{RemoveTriggerer, KeyboardStopper})},
		"/back_command": {h.BackCommandHandler},
	}

	return h
}

func (handler *BackHandler) UpdateLastCommand(command string) {
	if command != "/back_state" && command != "/back_command" {
		if len(handler.commandQueue) == 0 || command != handler.commandQueue[len(handler.commandQueue)-1] {
			handler.commandQueue = append(handler.commandQueue, command)
		}
	}

	for _, command := range handler.commandQueue {
		fmt.Println(command)
	}
	fmt.Println("------------------")
}

func (handler *BackHandler) InitHandler() {

}

func (handler *BackHandler) Handle(command string, params HandlerParams) ([]HandlerResponse, bool) {
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

func (handler *BackHandler) DeinitHandler() {

}

func (handler *BackHandler) BackStateHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	prevState := handler.sm.GetPreviousState()
	newActiveCommand := prevState.GetStartCommand()

	return HandlerResponse{messages: res.messages, nextState: prevState.GetName(), postCommandsHandle: []string{newActiveCommand}}
}

func (handler *BackHandler) BackCommandHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	if len(handler.commandQueue) < 1 {
		return HandlerResponse{messages: []bottypes.Message{{ChatID: params.message.ChatID, Text: "pizda"}}}
	}

	currentCommand := handler.commandQueue[len(handler.commandQueue)-2]
	handler.commandQueue = handler.commandQueue[:len(handler.commandQueue)-1]
	// logger.NewLogger(0).Info("prevPrevCommand: " + handler.prevPrevCommand + " prevCommand: " + handler.prevCommand + " current command: " + handler.currentCommand)

	return HandlerResponse{messages: res.messages, postCommandsHandle: []string{currentCommand}}
}
