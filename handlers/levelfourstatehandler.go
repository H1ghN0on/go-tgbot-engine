package handlers

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

type LevelFourHandler struct {
	Handler
}

func NewLevelFourHandler(gs GlobalStater) *LevelFourHandler {

	h := &LevelFourHandler{}
	h.gs = gs

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		"/level_four_start": {h.ModifyHandler(h.LevelFourStartHandler, []int{StateBackable, RemovableByTrigger})},
		"/level_four_one":   {h.ModifyHandler(h.LevelFourOneHandler, []int{StateBackable, RemovableByTrigger})},
		"/level_four_two":   {h.ModifyHandler(h.LevelFourTwoHandler, []int{StateBackable, RemovableByTrigger})},
		"/level_four_three": {h.ModifyHandler(h.LevelFourThreeHandler, []int{StateBackable, RemovableByTrigger})},
		"/level_four_four":  {h.ModifyHandler(h.LevelFourFourHandler, []int{RemovableByTrigger, RemoveTriggerer})},
	}

	return h
}

func (handler *LevelFourHandler) Handle(command bottypes.Command, params HandlerParams) ([]HandlerResponse, error) {
	var res []HandlerResponse

	handleFuncs, ok := handler.Handler.commands[command]
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

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SCENE BEGINS!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "ONE", Command: "/level_four_one"},
		},
	})

	res.messages = append(res.messages, response)
	res.nextCommands = append(res.nextCommands, "/level_four_one")
	res.nextState = "level-four-state"

	return res, nil
}

func (handler *LevelFourHandler) LevelFourOneHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "ONE!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "TWO", Command: "/level_four_two"},
		},
	})

	res.messages = append(res.messages, response)
	res.nextCommands = append(res.nextCommands, "/level_four_two")

	return res, nil
}

func (handler *LevelFourHandler) LevelFourTwoHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "TWO!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "THREE", Command: "/level_four_three"},
		},
	})

	res.messages = append(res.messages, response)
	res.nextCommands = append(res.nextCommands, "/level_four_three")

	return res, nil
}

func (handler *LevelFourHandler) LevelFourThreeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "THREE!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "FOUR", Command: "/level_four_four"},
		},
	})

	res.messages = append(res.messages, response)
	res.nextCommands = append(res.nextCommands, "/level_four_four")

	return res, nil
}

func (handler *LevelFourHandler) LevelFourFourHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "FOUR!"})
	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")
	res.nextState = "start-state"

	return res, nil
}
