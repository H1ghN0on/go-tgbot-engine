package handlers

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

var command_queue_max_size int = 30

type BackHandler struct {
	Handler
	sm           StateMachiner
	commandQueue []bottypes.Command
}

func NewBackHandler(gs GlobalStater, sm StateMachiner) *BackHandler {

	h := &BackHandler{}
	h.gs = gs
	h.sm = sm

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.BackStateCommand:   {h.ModifyHandler(h.BackStateHandler, []int{RemoveTriggerer, KeyboardStopper})},
		cmd.BackCommandCommand: {h.BackCommandHandler},
	}

	return h
}

func (handler *BackHandler) UpdateLastCommand(command bottypes.Command) {
	if !command.Equal(cmd.BackStateCommand) && !command.Equal(cmd.BackCommandCommand) {
		if len(handler.commandQueue) == 0 || command != handler.commandQueue[len(handler.commandQueue)-1] {
			if len(handler.commandQueue) > command_queue_max_size {
				handler.commandQueue = handler.commandQueue[1:]
			}
			handler.commandQueue = append(handler.commandQueue, command)
		}
	}
}

func (handler *BackHandler) ClearCommandQueue() {
	handler.commandQueue = nil
}

func (handler *BackHandler) Handle(params HandlerParams) ([]HandlerResponse, error) {
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

func (handler *BackHandler) BackStateHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	prevState := handler.sm.GetPreviousState()
	newActiveCommand := prevState.GetStartCommand()

	return HandlerResponse{messages: res.messages, nextState: prevState.GetName(), postCommandsHandle: []bottypes.Command{newActiveCommand}}, nil
}

func (handler *BackHandler) BackCommandHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	if len(handler.commandQueue) < 1 {
		return HandlerResponse{}, HandlerResponseError{message: "cannot return to previous command"}
	}

	currentCommand := handler.commandQueue[len(handler.commandQueue)-2]
	handler.commandQueue = handler.commandQueue[:len(handler.commandQueue)-1]

	return HandlerResponse{messages: res.messages, postCommandsHandle: []bottypes.Command{currentCommand}}, nil
}
