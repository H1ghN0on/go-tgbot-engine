package handlers_example

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands/example"
)

type LevelFourHandler struct {
	Handler
}

func NewLevelFourHandler(gs ExampleGlobalStater) *LevelFourHandler {

	h := &LevelFourHandler{}
	h.gs = gs

	h.Commands = map[bottypes.Command][]func(params handlers.HandlerParams) (handlers.HandlerResponse, error){
		cmd.LevelFourStartCommand: {h.ModifyHandler(h.LevelFourStartHandler, []int{handlers.StateBackable, handlers.RemovableByTrigger})},
		cmd.LevelFourOneCommand:   {h.ModifyHandler(h.LevelFourOneHandler, []int{handlers.StateBackable, handlers.RemovableByTrigger})},
		cmd.LevelFourTwoCommand:   {h.ModifyHandler(h.LevelFourTwoHandler, []int{handlers.StateBackable, handlers.RemovableByTrigger})},
		cmd.LevelFourThreeCommand: {h.ModifyHandler(h.LevelFourThreeHandler, []int{handlers.StateBackable, handlers.RemovableByTrigger})},
		cmd.LevelFourFourCommand:  {h.ModifyHandler(h.LevelFourFourHandler, []int{handlers.RemovableByTrigger, handlers.RemoveTriggerer})},
	}

	return h
}

func (handler *LevelFourHandler) Handle(params handlers.HandlerParams) ([]handlers.HandlerResponse, error) {
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

func (handler *LevelFourHandler) LevelFourStartHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SCENE BEGINS!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "ONE", Command: cmd.LevelFourOneCommand},
		},
	})

	res.Messages = append(res.Messages, response)
	res.NextCommands = append(res.NextCommands, cmd.LevelFourOneCommand)
	res.NextState = "level-four-state"

	return res, nil
}

func (handler *LevelFourHandler) LevelFourOneHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "ONE!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "TWO", Command: cmd.LevelFourTwoCommand},
		},
	})

	res.Messages = append(res.Messages, response)
	res.NextCommands = append(res.NextCommands, cmd.LevelFourTwoCommand)

	return res, nil
}

func (handler *LevelFourHandler) LevelFourTwoHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "TWO!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "THREE", Command: cmd.LevelFourThreeCommand},
		},
	})

	res.Messages = append(res.Messages, response)
	res.NextCommands = append(res.NextCommands, cmd.LevelFourThreeCommand)

	return res, nil
}

func (handler *LevelFourHandler) LevelFourThreeHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "THREE!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "FOUR", Command: cmd.LevelFourFourCommand},
		},
	})

	res.Messages = append(res.Messages, response)
	res.NextCommands = append(res.NextCommands, cmd.LevelFourFourCommand)

	return res, nil
}

func (handler *LevelFourHandler) LevelFourFourHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "FOUR!"})
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	res.NextState = "start-state"

	return res, nil
}
