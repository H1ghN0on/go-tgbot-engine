package handlers_example

import (
	"slices"
	"strings"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands/example"
)

type DynamicKeyboardHandler struct {
	Handler
	selectedItems []string
}

func NewDynamicKeyboardHandler(gs ExampleGlobalStater) *DynamicKeyboardHandler {

	h := &DynamicKeyboardHandler{}
	h.gs = gs

	h.Commands = map[bottypes.Command][]func(params handlers.HandlerParams) (handlers.HandlerResponse, error){
		cmd.DynamicKeyboardStartCommand:       {h.ModifyHandler(h.DynamicKeyboardStartHandler, []int{handlers.RemovableByTrigger})},
		cmd.DynamicKeyboardFirstStageCommand:  {h.ModifyHandler(h.DynamicKeyboardFirstHandler, []int{handlers.KeyboardStarter, handlers.RemovableByTrigger, handlers.StateBackable})},
		cmd.DynamicKeyboardSecondStageCommand: {h.ModifyHandler(h.DynamicKeyboardSecondHandler, []int{handlers.KeyboardStarter, handlers.RemovableByTrigger, handlers.CommandBackable})},
		cmd.DynamicKeyboardFinishCommand:      {h.ModifyHandler(h.DynamicKeyboardFinishHandler, []int{handlers.RemoveTriggerer})},
	}

	return h
}

func (handler *DynamicKeyboardHandler) HandleBackCommand(params handlers.HandlerParams) ([]handlers.HandlerResponse, error) {
	var response []handlers.HandlerResponse

	backHandler := handler.ModifyHandler(handler.HandleBackCommandImpl, []int{handlers.RemovableByTrigger})

	res, err := backHandler(params)

	if err != nil {
		return response, err
	}

	response = append(response, res)
	return response, nil
}

func (handler *DynamicKeyboardHandler) HandleBackCommandImpl(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	handler.selectedItems = handler.selectedItems[:len(handler.selectedItems)-1]

	var response handlers.HandlerResponse
	retMessage := bottypes.Message{ChatID: params.Message.Info.ChatID, Text: "You clicked back"}
	retMessage2 := bottypes.Message{ChatID: params.Message.Info.ChatID, Text: "ㅤ"}
	response.Messages = append(response.Messages, retMessage, retMessage2)

	return response, nil
}

func (handler *DynamicKeyboardHandler) Handle(params handlers.HandlerParams) ([]handlers.HandlerResponse, error) {
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

func (handler *DynamicKeyboardHandler) DynamicKeyboardStartHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	handler.selectedItems = nil
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.DynamicKeyboardFirstStageCommand)
	res.NextState = "dynamic-keyboard-state"

	return res, nil
}

func (handler *DynamicKeyboardHandler) DynamicKeyboardFirstHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Stack songs one"}

	data := handler.gs.GetDataForDynamicKeyboard()

	for _, text := range data["first_stage"] {
		retMessage.ButtonRows = append(retMessage.ButtonRows, bottypes.ButtonRows{
			Buttons: []bottypes.Button{
				{ChatID: chatID, Text: text, Command: bottypes.Command{
					Command:     cmd.DynamicKeyboardSecondStageCommand.Command,
					Description: cmd.DynamicKeyboardSecondStageCommand.Description,
					Data:        text,
				}},
			},
		})
	}

	res.Messages = append(res.Messages, retMessage)
	res.NextCommandToParse = bottypes.ParseableCommand{
		Command:   cmd.DynamicKeyboardSecondStageCommand,
		ParseType: bottypes.DynamicButtonParse,
	}

	res.NextCommands = append(res.NextCommands, cmd.DynamicKeyboardSecondStageCommand)

	return res, nil
}

func (handler *DynamicKeyboardHandler) DynamicKeyboardSecondHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Stack songs two"}

	data := handler.gs.GetDataForDynamicKeyboard()

	if !slices.Contains(data["first_stage"], params.Command.Data) {
		res.NextCommandToParse = bottypes.ParseableCommand{
			Command: cmd.DynamicKeyboardSecondStageCommand,
		}
		res.NextCommands = append(res.NextCommands, cmd.DynamicKeyboardSecondStageCommand)
		return res, nil
	}

	handler.selectedItems = append(handler.selectedItems, params.Command.Data)

	for _, text := range data["second_stage"] {
		retMessage.ButtonRows = append(retMessage.ButtonRows, bottypes.ButtonRows{
			Buttons: []bottypes.Button{
				{ChatID: chatID, Text: text, Command: bottypes.Command{
					Command:     cmd.DynamicKeyboardFinishCommand.Command,
					Description: cmd.DynamicKeyboardFinishCommand.Description,
					Data:        text,
				}},
			},
		})
	}

	res.Messages = append(res.Messages, retMessage)
	res.NextCommandToParse = bottypes.ParseableCommand{
		Command:   cmd.DynamicKeyboardFinishCommand,
		ParseType: bottypes.DynamicButtonParse,
	}
	res.NextCommands = append(res.NextCommands, cmd.DynamicKeyboardFinishCommand)

	return res, nil
}

func (handler *DynamicKeyboardHandler) DynamicKeyboardFinishHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID

	handler.selectedItems = append(handler.selectedItems, params.Command.Data)

	messageText := strings.Join(handler.selectedItems, ", ")

	retMessage := bottypes.Message{ChatID: chatID, Text: "You selected " + messageText}
	res.Messages = append(res.Messages, retMessage)

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	res.NextState = "start-state"

	return res, nil
}
