package handlers_example

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands/example"
)

type ParseModeHandler struct {
	Handler
}

func NewParseModeHandler(gs ExampleGlobalStater) *ParseModeHandler {

	h := &ParseModeHandler{}
	h.gs = gs

	h.Commands = map[bottypes.Command][]func(params handlers.HandlerParams) (handlers.HandlerResponse, error){
		cmd.ParseModeKeyboardStartCommand:  {h.ModifyHandler(h.ParseModeKeyboardStartHandler, []int{})},
		cmd.ParseModeStartCommand:          {h.ModifyHandler(h.ParseModeStartHandler, []int{handlers.StateBackable, handlers.RemovableByTrigger})},
		cmd.ParseModeMarkdownV2Command:     {h.ModifyHandler(h.ParseModeMarkdownV2Handler, []int{})},
		cmd.ParseModeHTMLCommand:           {h.ModifyHandler(h.ParseModeHTMLHandler, []int{})},
		cmd.ParseModeKeyboardFinishCommand: {h.ModifyHandler(h.ParseModeKeyboardFinishHandler, []int{handlers.RemoveTriggerer})},
	}
	return h
}

func (handler *ParseModeHandler) Handle(params handlers.HandlerParams) ([]handlers.HandlerResponse, error) {
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

func (handler *ParseModeHandler) ParseModeKeyboardStartHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ParseModeStartCommand)
	res.NextState = "parse-mode-keyboard-state"

	return res, nil
}

func (handler *ParseModeHandler) ParseModeStartHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID

	retMessage := bottypes.Message{ParseMode: bottypes.MarkdownV2, ChatID: chatID, Text: "*Choose parse mode*"}

	buttonRow := bottypes.ButtonRows{
		Buttons: []bottypes.Button{
			{ChatID: chatID, Text: "MarkdownV2", Command: cmd.ParseModeMarkdownV2Command},
			{ChatID: chatID, Text: "HTML", Command: cmd.ParseModeHTMLCommand},
		},
	}

	retMessage.ButtonRows = append(retMessage.ButtonRows, buttonRow)
	res.Messages = append(res.Messages, retMessage)
	res.NextCommands = append(res.NextCommands, cmd.ParseModeMarkdownV2Command, cmd.ParseModeHTMLCommand)
	return res, nil
}

func (handler *ParseModeHandler) ParseModeMarkdownV2Handler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID
	modeText := `MODES
				1 Crossed out: ~Example Text~
				2 Bold: *Example Text*
				3 Hidden: ||Example Text||
				4 Underlined: __Example Text__
				5 Italics: _Example Text_
				6 Link: [Example Text](https://avatars.mds.yandex.net/i?id=38e4a7500b85ae5d0ecc8f48c1128512f5511828-7942262-images-thumbs&n=13)
				7 References: [God](https://t.me/durov_russia)
				8 Emoji: ![🙀](tg://emoji?id=5458425656759032455)`
	modeText2 := "9 Monospaced: `Example Text`\n10 Code:\n```python\ndef hello(name: str = 'World'):\n    print(f'Hello, {name}!')\n\nhello()\nhello('Ivan')```"
	combinedText := modeText + "\n" + modeText2
	retMessage := bottypes.Message{ParseMode: bottypes.MarkdownV2, ChatID: chatID, Text: combinedText}
	res.Messages = append(res.Messages, retMessage)
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ParseModeKeyboardFinishCommand)
	return res, nil
}

func (handler *ParseModeHandler) ParseModeHTMLHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse
	chatID := params.Message.Info.ChatID
	boldText := "<b>Bold Text</b>\n<i>Italic Text</i>\n<s>Strikethrough Text</s>\n<code>Inline Monospaced Text</code>\n<pre>\nBlock of\nMonospaced Text\n</pre>\n<a href='https://example.com'>Clickable Link</a>\n<tg-spoiler>Spoiler Text</tg-spoiler>"

	retMessage := bottypes.Message{ParseMode: bottypes.HTML, ChatID: chatID, Text: boldText}

	res.Messages = append(res.Messages, retMessage)

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ParseModeKeyboardFinishCommand)
	return res, nil
}

func (handler *ParseModeHandler) ParseModeKeyboardFinishHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	res.NextState = "start-state"

	return res, nil
}
