package handlers

import (
	"fmt"
)

var command_queue_max_size int = 30

type BackHandler struct {
	Handler
	sm           StateMachiner
	commandQueue []string
}

func NewBackHandler(gs GlobalStater, sm StateMachiner) *BackHandler {

	h := &BackHandler{}
	h.gs = gs
	h.sm = sm

	h.commands = map[string][]func(params HandlerParams) (HandlerResponse, error){
		"/back_state":   {h.ModifyHandler(h.BackStateHandler, []int{RemoveTriggerer, KeyboardStopper})},
		"/back_command": {h.BackCommandHandler},
	}

	return h
}

func (handler *BackHandler) UpdateLastCommand(command string) {
	if command != "/back_state" && command != "/back_command" {
		if len(handler.commandQueue) == 0 || command != handler.commandQueue[len(handler.commandQueue)-1] {
			if len(handler.commandQueue) > command_queue_max_size {
				handler.commandQueue = handler.commandQueue[1:]
			}
			handler.commandQueue = append(handler.commandQueue, command)
		}
	}

	for _, command := range handler.commandQueue {
		fmt.Println(command)
	}
	fmt.Println("------------------")
}

func (handler *BackHandler) ClearCommandQueue() {
	handler.commandQueue = nil
}

func (handler *BackHandler) InitHandler() {

}

func (handler *BackHandler) Handle(command string, params HandlerParams) ([]HandlerResponse, bool, error) {
	var res []HandlerResponse

	handleFuncs, ok := handler.Handler.commands[command]
	if !ok {
		panic("wrong handler")
	}

	for _, handleFunc := range handleFuncs {
		response, err := handleFunc(params)
		if err != nil {
			return []HandlerResponse{}, true, err
		}
		res = append(res, response)
	}

	return res, true, nil
}

func (handler *BackHandler) DeinitHandler() {

}

func (handler *BackHandler) BackStateHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	prevState := handler.sm.GetPreviousState()
	newActiveCommand := prevState.GetStartCommand()

	return HandlerResponse{messages: res.messages, nextState: prevState.GetName(), postCommandsHandle: []string{newActiveCommand}}, nil
}

func (handler *BackHandler) BackCommandHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	if len(handler.commandQueue) < 1 {
		return HandlerResponse{}, HandlerResponseError{message: "cannot return to previous command"}
	}

	currentCommand := handler.commandQueue[len(handler.commandQueue)-2]
	handler.commandQueue = handler.commandQueue[:len(handler.commandQueue)-1]

	return HandlerResponse{messages: res.messages, postCommandsHandle: []string{currentCommand}}, nil
}
