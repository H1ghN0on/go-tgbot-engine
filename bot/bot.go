package bot

import (
	"errors"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	errs "github.com/H1ghN0on/go-tgbot-engine/errors"

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
	cmdhandler           CommandHandler
	api                  *tgbotapi.BotAPI
	lastMessagesToRemove []bottypes.Message
	keyboardMessage      bottypes.Message
}

type HandlerResponser interface {
	GetMessages() []bottypes.Message
	ShouldSwitchState() string
	ShouldRemoveLast() bool
	SetKeyboard() bool
}

type CommandHandlerResponser interface {
	GetResponses() []HandlerResponser
}

type CommandHandler interface {
	Handle(bottypes.Message) (CommandHandlerResponser, error)
}

func (bot *Client) SendMessage(message bottypes.Message, shouldRemoveLastMessage bool, setKeyboard bool) error {

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

	if setKeyboard && bot.keyboardMessage.ID != 0 {
		_, err := bot.api.Request(tgbotapi.NewEditMessageReplyMarkup(bot.keyboardMessage.ChatID, bot.keyboardMessage.ID, keyboard))
		if err != nil {
			panic(err.Error())
		}
		return nil
	}

	if !setKeyboard && bot.keyboardMessage.ID != 0 {
		bot.keyboardMessage = bottypes.Message{}
	}

	sent, err := bot.api.Send(msg)
	if err != nil {
		return BotError{code: errs.SendMessageError, message: "Send message error: " + err.Error()}
	}

	if shouldRemoveLastMessage {
		bot.lastMessagesToRemove = append(bot.lastMessagesToRemove, bottypes.Message{
			ID:         sent.MessageID,
			ChatID:     sent.Chat.ID,
			Text:       message.Text,
			ButtonRows: message.ButtonRows,
		})
	}

	if setKeyboard {
		bot.keyboardMessage = bottypes.Message{
			ID:     sent.MessageID,
			ChatID: sent.Chat.ID,
			Text:   sent.Text,
		}
	}

	return nil
}

func (client *Client) ListenMessages() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := client.api.GetUpdatesChan(u)

	for update := range updates {
		var receivedMessage bottypes.Message

		if update.Message != nil {

			receivedMessage = bottypes.Message{
				ID:     update.Message.MessageID,
				ChatID: update.Message.Chat.ID,
				Text:   update.Message.Text,
			}

		} else if update.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := client.api.Request(callback); err != nil {
				panic(err)
			}

			receivedMessage = bottypes.Message{
				ID:     update.CallbackQuery.Message.MessageID,
				ChatID: update.CallbackQuery.Message.Chat.ID,
				Text:   update.CallbackQuery.Data}

		} else {
			panic("Unknown message received")
		}

		handlerResult, err := client.cmdhandler.Handle(receivedMessage)

		if err != nil {
			var responseMessage bottypes.Message

			responseMessage.ChatID = receivedMessage.ChatID
			if errors.As(err, &errs.CommandHandlerError{}) {
				responseMessage.Text = "Command handler error: " + err.Error()
			} else if errors.As(err, &errs.StateMachineError{}) {
				responseMessage.Text = "State machine error: " + err.Error()
			} else {
				responseMessage.Text = "Unknown error occured"
			}

			err := client.SendMessage(responseMessage, false, false)
			if err != nil {
				panic(err)
			}
		} else {
			for _, messageToRemove := range client.lastMessagesToRemove {
				client.api.Request(tgbotapi.NewDeleteMessage(messageToRemove.ChatID, messageToRemove.ID))
			}
			client.lastMessagesToRemove = nil

			for _, response := range handlerResult.GetResponses() {
				for _, v := range response.GetMessages() {
					err := client.SendMessage(v, response.ShouldRemoveLast(), response.SetKeyboard())
					if err != nil {
						panic(err)
					}
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
