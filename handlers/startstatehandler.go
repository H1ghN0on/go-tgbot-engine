package handlers

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

type StartHandler struct {
	Handler
}

func NewStartHandler(gs GlobalStater) *StartHandler {

	h := &StartHandler{}
	h.gs = gs

	h.commands = map[string][]func(params HandlerParams) HandlerResponse{
		"/show_commands": {h.ShowCommandsHandler},
		"/level_one":     {h.LevelOneHandler},
		"/level_two":     {h.LevelTwoHandler},
		"/level_three":   {h.LevelThreeHandler},
		"/big_messages":  {h.BigMessagesHandler},
	}

	return h
}

func (handler *StartHandler) InitHandler() {

}

func (handler *StartHandler) Handle(command string, params HandlerParams) ([]HandlerResponse, bool) {
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

func (handler *StartHandler) DeinitHandler() {

}

func (handler *Handler) LevelOneHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SIGHT!"})

	return HandlerResponse{messages: res.messages, nextState: "start-state"}
}

func (handler *Handler) LevelTwoHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "SEE AS I SEE!"})

	return HandlerResponse{messages: res.messages, nextState: "start-state"}
}

func (handler *Handler) LevelThreeHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "FEEL WITH ME!"})

	return HandlerResponse{messages: res.messages, nextState: "start-state"}
}

func (handler *Handler) ShowCommandsHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	if handler.gs.GetName() != "" && handler.gs.GetSurname() != "" {
		retMessage := bottypes.Message{ChatID: chatID, Text: "Hello, " + handler.gs.GetName() + " " + handler.gs.GetSurname() + "!"}
		res.messages = append(res.messages, retMessage)
	}

	retMessage := bottypes.Message{ChatID: chatID, Text: "Here you go"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Level 1", Command: bottypes.Command{Text: "/level_one"}},
			{ChatID: chatID, Text: "Level 2", Command: bottypes.Command{Text: "/level_two"}},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Level 3", Command: bottypes.Command{Text: "/level_three"}},
			{ChatID: chatID, Text: "Create error", Command: bottypes.Command{Text: "/create_error"}},
		},
	}

	buttonRow3 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Curtain Call", Command: bottypes.Command{Text: "/level_four_start"}},
		},
	}

	buttonRow4 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Big messages", Command: bottypes.Command{Text: "/big_messages"}},
		},
	}

	buttonRow5 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Keyboard", Command: bottypes.Command{Text: "/keyboard_start"}},
		},
	}

	buttonRow6 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Set Info", Command: bottypes.Command{Text: "/set_info_start"}},
		},
	}
	buttonRow7 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Checkboxes", Command: bottypes.Command{Text: "/checkboxes_start"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2, buttonRow3, buttonRow4, buttonRow5, buttonRow6, buttonRow7)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages, nextState: "start-state"}
}

func (handler *Handler) BigMessagesHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})

	return HandlerResponse{messages: res.messages, nextState: "start-state"}
}
