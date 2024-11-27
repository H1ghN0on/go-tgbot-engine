package handlers_example

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type StartHandler struct {
	Handler
}

func NewStartHandler(gs ExampleGlobalStater) *StartHandler {

	h := &StartHandler{}
	h.gs = gs

	h.Commands = map[bottypes.Command][]func(params handlers.HandlerParams) (handlers.HandlerResponse, error){
		cmd.StartCommand:        {h.ModifyHandler(h.StartHandler, []int{handlers.RemovableByTrigger})},
		cmd.ShowCommandsCommand: {h.ModifyHandler(h.ShowCommandsHandler, []int{handlers.RemovableByTrigger})},
		cmd.LevelOneCommand:     {h.ModifyHandler(h.LevelOneHandler, []int{handlers.RemoveTriggerer})},
		cmd.LevelTwoCommand:     {h.ModifyHandler(h.LevelTwoHandler, []int{handlers.RemoveTriggerer})},
		cmd.LevelThreeCommand:   {h.ModifyHandler(h.LevelThreeHandler, []int{handlers.RemoveTriggerer})},
		cmd.BigMessagesCommand:  {h.ModifyHandler(h.BigMessagesHandler, []int{handlers.RemoveTriggerer})},
	}

	return h
}

func (handler *StartHandler) Handle(params handlers.HandlerParams) ([]handlers.HandlerResponse, error) {
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

func (handler *StartHandler) StartHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	return res, nil
}

func (handler *StartHandler) LevelOneHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "YOUR FINAL SIGHT!"})
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	return handlers.HandlerResponse{Messages: res.Messages, NextState: "start-state", PostCommandsHandle: res.PostCommandsHandle}, nil
}

func (handler *StartHandler) LevelTwoHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "SEE AS I SEE!"})
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	return handlers.HandlerResponse{Messages: res.Messages, NextState: "start-state", PostCommandsHandle: res.PostCommandsHandle}, nil
}

func (handler *StartHandler) LevelThreeHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "FEEL WITH ME!"})
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	return handlers.HandlerResponse{Messages: res.Messages, NextState: "start-state", PostCommandsHandle: res.PostCommandsHandle}, nil
}

func (handler *StartHandler) ShowCommandsHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	if handler.gs.GetName() != "" && handler.gs.GetSurname() != "" {
		retMessage := bottypes.Message{ChatID: chatID, Text: "Hello, " + handler.gs.GetName() + " " + handler.gs.GetSurname() + "!"}
		res.Messages = append(res.Messages, retMessage)
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

	buttonRow8 := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "Calendar", Command: cmd.CalendarStartCommand},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow1, buttonRow2, buttonRow3, buttonRow4, buttonRow5, buttonRow6, buttonRow7, buttonRow8)

	res.Messages = append(res.Messages, retMessage)

	return handlers.HandlerResponse{Messages: res.Messages}, nil
}

func (handler *StartHandler) BigMessagesHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	chatID := params.Message.Info.ChatID
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.Messages = append(res.Messages, bottypes.Message{ChatID: chatID, Text: "Давно выяснено, что при оценке дизайна и композиции читаемый текст мешает сосредоточиться. Lorem Ipsum используют потому, что тот обеспечивает более или менее стандартное заполнение шаблона, а также реальное распределение букв и пробелов в абзацах, которое не получается при простой дубликации 'Здесь ваш текст.. Здесь ваш текст.. Здесь ваш текст..' Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию, так что поиск по ключевым словам 'lorem ipsum' сразу показывает, как много веб-страниц всё ещё дожидаются своего настоящего рождения. За прошедшие годы текст Lorem Ipsum получил много версий. Некоторые версии появились по ошибке, некоторые - намеренно (например, юмористические варианты)."})
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	return handlers.HandlerResponse{Messages: res.Messages, NextState: "start-state", PostCommandsHandle: res.PostCommandsHandle}, nil
}
