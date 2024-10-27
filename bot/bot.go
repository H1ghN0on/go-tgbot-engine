package bot

import (
	"fmt"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotError struct {
	message string
}

func (err BotError) Error() string {
	return err.message
}

type Client struct {
	cmdhandler  CommandHandler
	api         *tgbotapi.BotAPI
	lastMessage bottypes.Message
}

type HandlerResponser interface {
	GetMessages() []bottypes.Message
	NextState() string
	IsKeyboard() bool
}

type CommandHandlerResponser interface {
	GetResponses() []HandlerResponser
}

type CommandHandler interface {
	Handle(bottypes.Message) (CommandHandlerResponser, error)
}

func (client *Client) parseMessage(update tgbotapi.Update) (bottypes.Message, int64, error) {
	var receivedMessage bottypes.Message
	var chatID int64

	if update.Message != nil {

		chatID = update.Message.Chat.ID

		receivedMessage = bottypes.Message{
			ID:     update.Message.MessageID,
			ChatID: chatID,
			Text:   update.Message.Text,
		}

	} else if update.CallbackQuery != nil {

		chatID = update.CallbackQuery.Message.Chat.ID

		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := client.api.Request(callback); err != nil {
			return bottypes.Message{}, chatID, BotError{message: "callback request failed"}
		}

		receivedMessage = bottypes.Message{
			ID:     update.CallbackQuery.Message.MessageID,
			ChatID: chatID,
			Text:   update.CallbackQuery.Data}

	} else {
		return bottypes.Message{}, 0, BotError{message: "unknown message received"}
	}

	return receivedMessage, chatID, nil
}

func (client *Client) SendMessage(message bottypes.Message, isKeyboard bool) error {

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

	if isKeyboard && client.lastMessage.ID != 0 {
		if message.Text != "" {
			_, err := client.api.Request(tgbotapi.NewEditMessageTextAndMarkup(client.lastMessage.ChatID, client.lastMessage.ID, message.Text, keyboard))
			if err != nil {
				panic(err.Error())
			}
		} else {
			_, err := client.api.Request(tgbotapi.NewEditMessageReplyMarkup(client.lastMessage.ChatID, client.lastMessage.ID, keyboard))
			if err != nil {
				panic(err.Error())
			}
		}
		return nil
	}

	sent, err := client.api.Send(msg)
	if err != nil {
		return BotError{message: "Send message error: " + err.Error()}
	}

	client.lastMessage = bottypes.Message{
		ID:     sent.MessageID,
		ChatID: sent.Chat.ID,
		Text:   sent.Text,
	}

	return nil
}

func (client *Client) sendErrorMessage(chatID int64, err error) {
	if chatID == 0 {
		panic("unknown chat ID")
	}

	var responseMessage bottypes.Message
	responseMessage.ChatID = chatID
	responseMessage.Text = err.Error()

	sendErr := client.SendMessage(responseMessage, false)
	if sendErr != nil {
		panic(sendErr)
	}
}

func (client *Client) ListenMessages() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := client.api.GetUpdatesChan(u)

	for update := range updates {
		var receivedMessage bottypes.Message

		receivedMessage, chatID, err := client.parseMessage(update)
		if err != nil {
			client.sendErrorMessage(chatID, fmt.Errorf("parse error: %w", err))
			continue
		}

		handlerResult, err := client.cmdhandler.Handle(receivedMessage)
		if err != nil {
			client.sendErrorMessage(chatID, fmt.Errorf("handle command error: %w", err))
			continue
		}

		for _, response := range handlerResult.GetResponses() {
			for _, v := range response.GetMessages() {
				err := client.SendMessage(v, response.IsKeyboard())
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func NewClient(api *tgbotapi.BotAPI, ch CommandHandler) *Client {
	return &Client{
		cmdhandler: ch,
		api:        api,
	}
}
