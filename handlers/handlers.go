package handlers

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

type HandlerParams struct {
	message bottypes.Message
}

type HandlerResponse struct {
	Messages          []bottypes.Message
	ShouldSwitchState string
	ShouldRemoveLast  bool
	SetKeyboard       bool
}

func LevelOneHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SIGHT!"})

	return HandlerResponse{Messages: res.Messages}
}

func LevelTwoHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "SEE AS I SEE!"})

	return HandlerResponse{Messages: res.Messages}
}

func LevelThreeHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "FEEL WITH ME!"})

	return HandlerResponse{Messages: res.Messages}
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
			{ChatID: chatID, Text: "Big Messages", Command: bottypes.Command{Text: "/big_messages"}},
		},
	}

	buttonRow5 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Keyboard", Command: bottypes.Command{Text: "/keyboard_start"}},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2, buttonRow3, buttonRow4, buttonRow5)
	res.Messages = append(res.Messages, retMessage)

	return HandlerResponse{Messages: res.Messages, ShouldSwitchState: "start-state", ShouldRemoveLast: false}
}

func KeyboardStartHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Here you go"}

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
	res.Messages = append(res.Messages, retMessage)

	return HandlerResponse{Messages: res.Messages, ShouldSwitchState: "keyboard-state", SetKeyboard: true}
}

func KeyboardOneHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Here you go"}

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
	res.Messages = append(res.Messages, retMessage)

	return HandlerResponse{Messages: res.Messages, SetKeyboard: true}
}

func KeyboardTwoHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Here you go"}

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
	res.Messages = append(res.Messages, retMessage)

	return HandlerResponse{Messages: res.Messages, SetKeyboard: true}
}

func KeyboardThreeHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse
	chatID := params.message.ChatID

	retMessage := bottypes.Message{ChatID: chatID, Text: "Alabama certified moment"}
	res.Messages = append(res.Messages, retMessage)

	return HandlerResponse{Messages: res.Messages, ShouldSwitchState: "start-state", SetKeyboard: false}
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

	res.Messages = append(res.Messages, response)

	return HandlerResponse{Messages: res.Messages, ShouldSwitchState: "level-four-state"}
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

	res.Messages = append(res.Messages, response)

	return HandlerResponse{Messages: res.Messages}
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

	res.Messages = append(res.Messages, response)

	return HandlerResponse{Messages: res.Messages}
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

	res.Messages = append(res.Messages, response)

	return HandlerResponse{Messages: res.Messages}
}

func LevelFourFourHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "FOUR!"})

	return HandlerResponse{Messages: res.Messages, ShouldSwitchState: "start-state"}
}

func BigMessagesHandler(params HandlerParams) HandlerResponse {
	var res HandlerResponse

	chatID := params.message.ChatID
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})

	return HandlerResponse{Messages: res.Messages, ShouldSwitchState: "start-state"}
}
