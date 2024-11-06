package handlers

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

type StartHandler struct {
	Handler
}

func NewStartHandler(gs GlobalStater) *StartHandler {

	h := &StartHandler{}
	h.gs = gs

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		"/show_commands": {h.ShowCommandsHandler},
		"/level_one":     {h.LevelOneHandler},
		"/level_two":     {h.LevelTwoHandler},
		"/level_three":   {h.LevelThreeHandler},
		"/big_messages":  {h.BigMessagesHandler},
	}

	return h
}

func (handler *StartHandler) Handle(command bottypes.Command, params HandlerParams) ([]HandlerResponse, error) {
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

func (handler *StartHandler) LevelOneHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SIGHT!"})
	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")
	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}, nil
}

func (handler *StartHandler) LevelTwoHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "SEE AS I SEE!"})
	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")
	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}, nil
}

func (handler *StartHandler) LevelThreeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "FEEL WITH ME!"})
	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")
	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}, nil
}

func (handler *StartHandler) ShowCommandsHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.ChatID

	if handler.gs.GetName() != "" && handler.gs.GetSurname() != "" {
		retMessage := bottypes.Message{ChatID: chatID, Text: "Hello, " + handler.gs.GetName() + " " + handler.gs.GetSurname() + "!"}
		res.messages = append(res.messages, retMessage)
	}

	retMessage := bottypes.Message{ChatID: chatID, Text: "Here you go"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Level 1", Command: "/level_one"},
			{ChatID: chatID, Text: "Level 2", Command: "/level_two"},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Level 3", Command: "/level_three"},
			{ChatID: chatID, Text: "Create error", Command: "/create_error"},
		},
	}

	buttonRow3 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Curtain Call", Command: "/level_four_start"},
		},
	}

	buttonRow4 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Big messages", Command: "/big_messages"},
		},
	}

	buttonRow5 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Keyboard", Command: "/keyboard_start"},
		},
	}

	buttonRow6 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Set Info", Command: "/set_info_start"},
		},
	}
	buttonRow7 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Checkboxes", Command: "/checkboxes_start"},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2, buttonRow3, buttonRow4, buttonRow5, buttonRow6, buttonRow7)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}, nil
}

func (handler *StartHandler) BigMessagesHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.postCommandsHandle = append(res.postCommandsHandle, "/show_commands")
	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}, nil
}
