package handlers

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type StartHandler struct {
	Handler
}

func NewStartHandler(gs GlobalStater) *StartHandler {

	h := &StartHandler{}
	h.gs = gs

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.ShowCommandsCommand: {h.ShowCommandsHandler},
		cmd.LevelOneCommand:     {h.LevelOneHandler},
		cmd.LevelTwoCommand:     {h.LevelTwoHandler},
		cmd.LevelThreeCommand:   {h.LevelThreeHandler},
		cmd.BigMessagesCommand:  {h.BigMessagesHandler},
	}

	return h
}

func (handler *StartHandler) Handle(params HandlerParams) ([]HandlerResponse, error) {
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

func (handler *StartHandler) LevelOneHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SIGHT!"})
	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.ShowCommandsCommand)
	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}, nil
}

func (handler *StartHandler) LevelTwoHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "SEE AS I SEE!"})
	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.ShowCommandsCommand)
	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}, nil
}

func (handler *StartHandler) LevelThreeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "FEEL WITH ME!"})
	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.ShowCommandsCommand)
	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}, nil
}

func (handler *StartHandler) ShowCommandsHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse
	chatID := params.message.Info.ChatID

	if handler.gs.GetName() != "" && handler.gs.GetSurname() != "" {
		retMessage := bottypes.Message{ChatID: chatID, Text: "Hello, " + handler.gs.GetName() + " " + handler.gs.GetSurname() + "!"}
		res.messages = append(res.messages, retMessage)
	}

	retMessage := bottypes.Message{ChatID: chatID, Text: "Here you go"}

	buttonRow1 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Level 1", Command: cmd.LevelOneCommand},
			{ChatID: chatID, Text: "Level 2", Command: cmd.LevelTwoCommand},
		},
	}

	buttonRow2 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Level 3", Command: cmd.LevelThreeCommand},
			// {ChatID: chatID, Text: "Create error", Command: "/create_error"},
		},
	}

	buttonRow3 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Curtain Call", Command: cmd.LevelFourStartCommand},
		},
	}

	buttonRow4 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Big messages", Command: cmd.BigMessagesCommand},
		},
	}

	buttonRow5 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Keyboard", Command: cmd.KeyboardStartCommand},
			{ChatID: chatID, Text: "Dynamic Keyboard", Command: cmd.DynamicKeyboardStartCommand},
		},
	}

	buttonRow6 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Set Info", Command: cmd.SetInfoStartCommand},
		},
	}
	buttonRow7 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Checkboxes", Command: cmd.CheckboxStartCommand},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2, buttonRow3, buttonRow4, buttonRow5, buttonRow6, buttonRow7)

	res.messages = append(res.messages, retMessage)

	return HandlerResponse{messages: res.messages}, nil
}

func (handler *StartHandler) BigMessagesHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	chatID := params.message.Info.ChatID
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.messages = append(res.messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.ShowCommandsCommand)
	return HandlerResponse{messages: res.messages, nextState: "start-state", postCommandsHandle: res.postCommandsHandle}, nil
}
