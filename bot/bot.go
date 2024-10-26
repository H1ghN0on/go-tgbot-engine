package bot

import (
	"errors"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	errs "github.com/H1ghN0on/go-tgbot-engine/errors"
	hdl "github.com/H1ghN0on/go-tgbot-engine/handlers"
	sm "github.com/H1ghN0on/go-tgbot-engine/statemachine"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotError struct {
	code    errs.ErrorType
	message string
}

func (err BotError) Error() string {
	return err.message
}

type Client struct {
	C                    *tgbotapi.BotAPI
	Sm                   *sm.StateMachine
	LastMessagesToRemove []bottypes.Message
	keyboardMessage      bottypes.Message
}

type Bot struct {
	Cmdhandler hdl.CommandHandler
	Client     *Client
}

func (bot *Bot) SendMessage(message bottypes.Message, shouldRemoveLastMessage bool, setKeyboard bool) error {

	var keyboard tgbotapi.InlineKeyboardMarkup

	msg := tgbotapi.NewMessage(message.ChatID, message.Text)
	if len(message.ButtonRows) > 0 {
		var buttonRows [][]tgbotapi.InlineKeyboardButton
		for _, buttonRow := range message.ButtonRows {
			var buttons []tgbotapi.InlineKeyboardButton
			for _, button := range buttonRow.Buttons {
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Command.Text))
			}
			buttonRows = append(buttonRows, buttons)
		}
		keyboard = tgbotapi.NewInlineKeyboardMarkup(buttonRows...)
		msg.ReplyMarkup = keyboard
	}

	if setKeyboard && bot.Client.keyboardMessage.ID != 0 {
		_, err := bot.Client.C.Request(tgbotapi.NewEditMessageReplyMarkup(bot.Client.keyboardMessage.ChatID, bot.Client.keyboardMessage.ID, keyboard))
		if err != nil {
			panic(err.Error())
		}
		return nil
	}

	if !setKeyboard && bot.Client.keyboardMessage.ID != 0 {
		bot.Client.keyboardMessage = bottypes.Message{}
	}

	sent, err := bot.Client.C.Send(msg)
	if err != nil {
		return BotError{code: errs.SendMessageError, message: "Send message error: " + err.Error()}
	}

	if shouldRemoveLastMessage {
		bot.Client.LastMessagesToRemove = append(bot.Client.LastMessagesToRemove, bottypes.Message{
			ID:         sent.MessageID,
			ChatID:     sent.Chat.ID,
			Text:       message.Text,
			ButtonRows: message.ButtonRows,
		})
	}

	if setKeyboard {
		bot.Client.keyboardMessage = bottypes.Message{
			ID:     sent.MessageID,
			ChatID: sent.Chat.ID,
			Text:   sent.Text,
		}
	}

	return nil
}

func (bot *Bot) ListenMessages(update tgbotapi.Update) {

	var receivedMessage bottypes.Message

	if update.Message != nil {

		receivedMessage = bottypes.Message{
			ID:     update.Message.MessageID,
			ChatID: update.Message.Chat.ID,
			Text:   update.Message.Text}

	} else if update.CallbackQuery != nil {

		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := bot.Client.C.Request(callback); err != nil {
			panic(err)
		}

		receivedMessage = bottypes.Message{
			ID:     update.CallbackQuery.Message.MessageID,
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   update.CallbackQuery.Data}

	} else {
		panic("Unknown message received")
	}

	handlerResult, err := bot.Cmdhandler.Handle(receivedMessage)

	if err != nil {
		var responseMessage bottypes.Message

		responseMessage.ChatID = update.CallbackQuery.Message.Chat.ID
		if errors.As(err, &errs.CommandHandlerError{}) {
			responseMessage.Text = "Command handler error: " + err.Error()
		} else if errors.As(err, &errs.StateMachineError{}) {
			responseMessage.Text = "State machine error: " + err.Error()
		} else {
			responseMessage.Text = "Unknown error occured"
		}

		err := bot.SendMessage(responseMessage, false, false)
		if err != nil {
			panic(err)
		}
	} else {
		for _, messageToRemove := range bot.Client.LastMessagesToRemove {
			bot.Client.C.Request(tgbotapi.NewDeleteMessage(messageToRemove.ChatID, messageToRemove.ID))
		}
		bot.Client.LastMessagesToRemove = nil

		for _, response := range handlerResult.Responses {
			for _, v := range response.Messages {
				err := bot.SendMessage(v, response.ShouldRemoveLast, response.SetKeyboard)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
