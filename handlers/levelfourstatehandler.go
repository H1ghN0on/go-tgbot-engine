package handlers

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type LevelFourHandler struct {
	Handler
}

func NewLevelFourHandler(gs GlobalStater) *LevelFourHandler {

	h := &LevelFourHandler{}
	h.gs = gs

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.LevelFourStartCommand: {h.ModifyHandler(h.LevelFourStartHandler, []int{StateBackable, RemovableByTrigger})},
		cmd.LevelFourOneCommand:   {h.ModifyHandler(h.LevelFourOneHandler, []int{StateBackable, RemovableByTrigger})},
		cmd.LevelFourTwoCommand:   {h.ModifyHandler(h.LevelFourTwoHandler, []int{StateBackable, RemovableByTrigger})},
		cmd.LevelFourThreeCommand: {h.ModifyHandler(h.LevelFourThreeHandler, []int{StateBackable, RemovableByTrigger})},
		cmd.LevelFourFourCommand:  {h.ModifyHandler(h.LevelFourFourHandler, []int{RemovableByTrigger, RemoveTriggerer})},
	}

	return h
}

func (handler *LevelFourHandler) Handle(params HandlerParams) ([]HandlerResponse, error) {
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

func (handler *LevelFourHandler) LevelFourStartHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SCENE BEGINS!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "ONE", Command: cmd.LevelFourOneCommand},
		},
	})

	res.messages = append(res.messages, response)
	res.nextCommands = append(res.nextCommands, cmd.LevelFourOneCommand)
	res.nextState = "level-four-state"

	return res, nil
}

func (handler *LevelFourHandler) LevelFourOneHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "ONE!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "TWO", Command: cmd.LevelFourTwoCommand},
		},
	})

	res.messages = append(res.messages, response)
	res.nextCommands = append(res.nextCommands, cmd.LevelFourTwoCommand)

	return res, nil
}

func (handler *LevelFourHandler) LevelFourTwoHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "TWO!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "THREE", Command: cmd.LevelFourThreeCommand},
		},
	})

	res.messages = append(res.messages, response)
	res.nextCommands = append(res.nextCommands, cmd.LevelFourThreeCommand)

	return res, nil
}

func (handler *LevelFourHandler) LevelFourThreeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "THREE!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "FOUR", Command: cmd.LevelFourFourCommand},
		},
	})

	res.messages = append(res.messages, response)
	res.nextCommands = append(res.nextCommands, cmd.LevelFourFourCommand)

	return res, nil
}

func (handler *LevelFourHandler) LevelFourFourHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "FOUR!"})
	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.ShowCommandsCommand)
	res.nextState = "start-state"

	return res, nil
}
