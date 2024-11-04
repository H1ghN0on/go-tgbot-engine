package handlers

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

type LevelFourHandler struct {
	Handler
}

func NewLevelFourHandler(gs GlobalStater) *LevelFourHandler {

	h := &LevelFourHandler{}
	h.gs = gs

	h.commands = map[string][]func(params HandlerParams) HandlerResponse{
		"/level_four_start": {h.ModifyHandler(h.LevelFourStartHandler, []int{StateBackable, RemovableByTrigger})},
		"/level_four_one":   {h.ModifyHandler(h.LevelFourOneHandler, []int{StateBackable, RemovableByTrigger})},
		"/level_four_two":   {h.ModifyHandler(h.LevelFourTwoHandler, []int{StateBackable, RemovableByTrigger})},
		"/level_four_three": {h.ModifyHandler(h.LevelFourThreeHandler, []int{StateBackable, RemovableByTrigger})},
		"/level_four_four":  {h.ModifyHandler(h.LevelFourFourHandler, []int{RemovableByTrigger, RemoveTriggerer})},
	}

	return h
}

func (handler *LevelFourHandler) InitHandler() {

}

func (handler *LevelFourHandler) Handle(command string, params HandlerParams) ([]HandlerResponse, bool) {
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

func (handler *LevelFourHandler) DeinitHandler() {

}

func (handler *Handler) LevelFourStartHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SCENE BEGINS!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "ONE", Command: bottypes.Command{Text: "/level_four_one"}},
		},
	})

	res.messages = append(res.messages, response)

	return HandlerResponse{messages: res.messages, nextState: "level-four-state"}
}

func (handler *Handler) LevelFourOneHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "ONE!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "TWO", Command: bottypes.Command{Text: "/level_four_two"}},
		},
	})

	res.messages = append(res.messages, response)

	return HandlerResponse{messages: res.messages}
}

func (handler *Handler) LevelFourTwoHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "TWO!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "THREE", Command: bottypes.Command{Text: "/level_four_three"}},
		},
	})

	res.messages = append(res.messages, response)

	return HandlerResponse{messages: res.messages}
}

func (handler *Handler) LevelFourThreeHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "THREE!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "FOUR", Command: bottypes.Command{Text: "/level_four_four"}},
		},
	})

	res.messages = append(res.messages, response)

	return HandlerResponse{messages: res.messages}
}

func (handler *Handler) LevelFourFourHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "FOUR!"})
	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")

	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}
}
