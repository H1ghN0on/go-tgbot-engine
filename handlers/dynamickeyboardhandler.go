package handlers

import (
	"slices"
	"strings"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type DynamicKeyboardHandler struct {
	Handler

	selectedItems []string
}

func NewDynamicKeyboardHandler(gs GlobalStater) *DynamicKeyboardHandler {

	h := &DynamicKeyboardHandler{}
	h.gs = gs

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.DynamicKeyboardStartCommand:       {h.ModifyHandler(h.DynamicKeyboardStartHandler, []int{RemovableByTrigger})},
		cmd.DynamicKeyboardFirstStageCommand:  {h.ModifyHandler(h.DynamicKeyboardFirstHandler, []int{KeyboardStarter, RemovableByTrigger, StateBackable})},
		cmd.DynamicKeyboardSecondStageCommand: {h.ModifyHandler(h.DynamicKeyboardSecondHandler, []int{KeyboardStarter, RemovableByTrigger, CommandBackable})},
		cmd.DynamicKeyboardFinishCommand:      {h.ModifyHandler(h.DynamicKeyboardFinishHandler, []int{RemoveTriggerer})},
	}

	return h
}

func (handler *DynamicKeyboardHandler) HandleBackCommand(params HandlerParams) ([]HandlerResponse, error) {
	var response []HandlerResponse

	backHandler := handler.ModifyHandler(handler.HandleBackCommandImpl, []int{RemovableByTrigger})

	res, err := backHandler(params)

	if err != nil {
		return response, err
	}

	response = append(response, res)
	return response, nil
}

func (handler *DynamicKeyboardHandler) HandleBackCommandImpl(params HandlerParams) (HandlerResponse, error) {
	handler.selectedItems = handler.selectedItems[:len(handler.selectedItems)-1]

	var response HandlerResponse
	retMessage := bottypes.Message{ChatID: params.message.Info.ChatID, Text: "You clicked back"}
	retMessage2 := bottypes.Message{ChatID: params.message.Info.ChatID, Text: "ㅤ"}
	response.messages = append(response.messages, retMessage, retMessage2)

	return response, nil
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

	handler.selectedItems = nil
	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.DynamicKeyboardFirstStageCommand)
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
				{ChatID: chatID, Text: text, Command: bottypes.Command{
					Command:     cmd.DynamicKeyboardSecondStageCommand.Command,
					Description: cmd.DynamicKeyboardSecondStageCommand.Description,
					Data:        text,
				}},
			},
		})
	}

	res.messages = append(res.messages, retMessage)
	res.nextCommandToParse = bottypes.ParseableCommand{
		Command:   cmd.DynamicKeyboardSecondStageCommand,
		ParseType: bottypes.DynamicButtonParse,
	}

	res.nextCommands = append(res.nextCommands, cmd.DynamicKeyboardSecondStageCommand)

	return res, nil
}

func (handler *DynamicKeyboardHandler) DynamicKeyboardSecondHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Stack songs two"}

	data := handler.gs.GetDataForDynamicKeyboard()

	if !slices.Contains(data["first_stage"], params.command.Data) {
		res.nextCommandToParse = bottypes.ParseableCommand{
			Command: cmd.DynamicKeyboardSecondStageCommand,
		}
		res.nextCommands = append(res.nextCommands, cmd.DynamicKeyboardSecondStageCommand)
		return res, nil
	}

	handler.selectedItems = append(handler.selectedItems, params.command.Data)

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

	res.messages = append(res.messages, retMessage)
	res.nextCommandToParse = bottypes.ParseableCommand{
		Command:   cmd.DynamicKeyboardFinishCommand,
		ParseType: bottypes.DynamicButtonParse,
	}
	res.nextCommands = append(res.nextCommands, cmd.DynamicKeyboardFinishCommand)

	return res, nil
}

func (handler *DynamicKeyboardHandler) DynamicKeyboardFinishHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID

	handler.selectedItems = append(handler.selectedItems, params.command.Data)

	messageText := strings.Join(handler.selectedItems, ", ")

	retMessage := bottypes.Message{ChatID: chatID, Text: "You selected " + messageText}
	res.messages = append(res.messages, retMessage)

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.ShowCommandsCommand)
	res.nextState = "start-state"

	return res, nil
}
