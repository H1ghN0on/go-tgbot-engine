package handlers

import (
	"fmt"
	"slices"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

var command_queue_max_size int = 30

type BackHandler struct {
	Handler
	sm           StateMachiner
	commandQueue []bottypes.Command
}

func (handler Handler) FindCommandInTheList(command bottypes.Command) (bottypes.Command, error) {
	if command.Command == "" || !command.IsCommand() {
		return bottypes.Command{}, fmt.Errorf("not command received")
	}

	index := slices.IndexFunc(cmd.Commands, func(com bottypes.Command) bool { return command.Equal(com) })
	if index == -1 {
		return bottypes.Command{}, fmt.Errorf("unknown command received")
	}

	trueCommand := cmd.Commands[index]
	trueCommand.Data = command.Data

	return trueCommand, nil
}

func NewBackHandler(sm StateMachiner) *BackHandler {

	h := &BackHandler{}
	h.sm = sm

	h.Commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.BackStateCommand:   {h.ModifyHandler(h.BackStateHandler, []int{RemoveTriggerer, KeyboardStopper})},
		cmd.BackCommandCommand: {h.BackCommandHandler},
	}

	return h
}

func (handler *BackHandler) UpdateLastCommand(command bottypes.Command) {
	if !command.Equal(cmd.BackStateCommand) && !command.Equal(cmd.BackCommandCommand) && !command.Equal(cmd.NothingnessCommand) {
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

	handleFuncs, ok := handler.GetCommandFromMap(params.Command)
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

	res.NextState = prevState.GetName()
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, newActiveCommand)

	return res, nil
}

func (handler *BackHandler) BackCommandHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	if len(handler.commandQueue) < 1 {
		return HandlerResponse{}, HandlerResponseError{Message: "cannot return to previous command"}
	}

	lastCommand := handler.commandQueue[len(handler.commandQueue)-1]
	currentCommand := lastCommand

	for currentCommand.SkipOnBack || currentCommand.Equal(lastCommand) {
		handler.commandQueue = handler.commandQueue[:len(handler.commandQueue)-1]

		if len(handler.commandQueue) < 1 {
			return HandlerResponse{}, HandlerResponseError{Message: "cannot return to previous command"}
		}

		currentCommand = handler.commandQueue[len(handler.commandQueue)-1]
	}

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, currentCommand)
	res.PostCommandsHandle.IsBackCommand = true

	return res, nil
}
