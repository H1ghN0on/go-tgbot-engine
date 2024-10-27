package handlers

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

type HandlerParams struct {
	message bottypes.Message
}

type HandlerResponse struct {
	messages   []bottypes.Message
	nextState  string
	isKeyboard bool
}

func (hr HandlerResponse) GetMessages() []bottypes.Message {
	return hr.messages
}

func (hr HandlerResponse) NextState() string {
	return hr.nextState
}

func (hr HandlerResponse) IsKeyboard() bool {
	return hr.isKeyboard
}

func LevelOneHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SIGHT!"})

	return HandlerResponse{messages: res.messages}
}

func LevelTwoHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "SEE AS I SEE!"})

	return HandlerResponse{messages: res.messages}
}

func LevelThreeHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "FEEL WITH ME!"})

	return HandlerResponse{messages: res.messages}
}

func ShowCommandsHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

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

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2, buttonRow3, buttonRow4, buttonRow5)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages, nextState: "start-state"}
}

func KeyboardStartHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 1"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Let me", Command: bottypes.Command{Text: "/keyboard_one"}},
			{ChatID: chatID, Text: "No, Let me!", Command: bottypes.Command{Text: "/keyboard_one"}},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "LET ME!", Command: bottypes.Command{Text: "/keyboard_one"}},
			{ChatID: chatID, Text: "l-ll-let me *blushes*", Command: bottypes.Command{Text: "/keyboard_one"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages, nextState: "keyboard-state", isKeyboard: true}
}

func KeyboardOneHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: ""}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Push me", Command: bottypes.Command{Text: "/keyboard_two"}},
			{ChatID: chatID, Text: "No, push me!", Command: bottypes.Command{Text: "/keyboard_two"}},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "PUSH ME!", Command: bottypes.Command{Text: "/keyboard_two"}},
			{ChatID: chatID, Text: "p-pp-push me *blushes*", Command: bottypes.Command{Text: "/keyboard_two"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages, isKeyboard: true}
}

func KeyboardTwoHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Option 3"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Approach me", Command: bottypes.Command{Text: "/keyboard_three"}},
			{ChatID: chatID, Text: "No, approach me!", Command: bottypes.Command{Text: "/keyboard_three"}},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "APPROACH ME!", Command: bottypes.Command{Text: "/keyboard_three"}},
			{ChatID: chatID, Text: "a-aa-approach me *blushes*", Command: bottypes.Command{Text: "/keyboard_three"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2)
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages, isKeyboard: true}
}

func KeyboardThreeHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Alabama certified moment"}
	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages, nextState: "start-state"}
}

func LevelFourStartHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SCENE BEGINS!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "ONE", Command: bottypes.Command{Text: "/level_four_one"}},
		},
	})

	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Back", Command: bottypes.Command{Text: "/state_back"}},
		},
	})

	res.messages = append(res.messages, response)

	return HandlerResponse{messages: res.messages, nextState: "level-four-state"}
}

func LevelFourOneHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "ONE!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "TWO", Command: bottypes.Command{Text: "/level_four_two"}},
		},
	})

	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Back", Command: bottypes.Command{Text: "/state_back"}},
		},
	})

	res.messages = append(res.messages, response)

	return HandlerResponse{messages: res.messages}
}

func LevelFourTwoHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "TWO!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "THREE", Command: bottypes.Command{Text: "/level_four_three"}},
		},
	})

	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Back", Command: bottypes.Command{Text: "/state_back"}},
		},
	})

	res.messages = append(res.messages, response)

	return HandlerResponse{messages: res.messages}
}

func LevelFourThreeHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID

	response := bottypes.Message{ChatID: chatID, Text: "THREE!"}
	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "FOUR", Command: bottypes.Command{Text: "/level_four_four"}},
		},
	})

	response.ButtonRows = append(response.ButtonRows, bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Back", Command: bottypes.Command{Text: "/state_back"}},
		},
	})

	res.messages = append(res.messages, response)

	return HandlerResponse{messages: res.messages}
}

func LevelFourFourHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "FOUR!"})

	return HandlerResponse{messages: res.messages, nextState: "start-state"}
}

func BigMessagesHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})

	return HandlerResponse{messages: res.messages, nextState: "start-state"}
}

func BackHandler(params HandlerParams) HandlerResponse {
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
