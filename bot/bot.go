package bot

import (
	"fmt"
	"slices"

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
	cmdhandler          CommandHandler
	api                 *tgbotapi.BotAPI
	lastMessage         bottypes.Message
	messagesToRemove    []bottypes.Message
	lastKeyboardMessage bottypes.Message
}

type HandlerResponser interface {
	GetMessages() []bottypes.Message
	NextState() string
	ContainsTrigger(bottypes.Trigger) bool
}

type CommandHandlerRequester interface {
	GetMessage() bottypes.Message
	ShouldUpdateQueue() bool
}

type CommandHandlerResponser interface {
	GetResponses() []HandlerResponser
}

type CommandHandler interface {
	NewCommandHandlerRequest(msg bottypes.Message, shouldUpdateQueue bool) CommandHandlerRequester
	Handle(req CommandHandlerRequester) (CommandHandlerResponser, error)
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

	} else if update.EditedMessage != nil {
		chatID = update.EditedMessage.Chat.ID
		return bottypes.Message{}, chatID, BotError{message: "editing is forbidden"}
	} else {
		return bottypes.Message{}, 0, BotError{message: "unknown message received"}
	}

	return receivedMessage, chatID, nil
}

func (client Client) compareMessages(a bottypes.Message) func(bottypes.Message) bool {
	return func(b bottypes.Message) bool {
		return a.ChatID == b.ChatID && a.ID == b.ID
	}
}

func (client *Client) addToRemoveMessagesQueue(message bottypes.Message) {
	if !slices.ContainsFunc(client.messagesToRemove, client.compareMessages(message)) {
		client.messagesToRemove = append(client.messagesToRemove, message)
	}
}

func (client *Client) SetupKeyboard(message bottypes.Message, keyboard tgbotapi.InlineKeyboardMarkup) error {
	hasText := message.Text != ""
	hasButtons := len(message.ButtonRows) != 0

	var attachMessage bottypes.Message
	if client.lastKeyboardMessage.ID != 0 {
		attachMessage = client.lastKeyboardMessage
	} else if client.lastMessage.ID != 0 {
		attachMessage = client.lastMessage
		client.lastKeyboardMessage = client.lastMessage
	} else {
		return fmt.Errorf("keyboard is no message to attach")
	}

	if hasText {
		if hasButtons {
			_, err := client.api.Request(tgbotapi.NewEditMessageTextAndMarkup(attachMessage.ChatID, attachMessage.ID, message.Text, keyboard))
			if err != nil {
				return err
			}
			return nil
		} else {
			_, err := client.api.Request(tgbotapi.NewEditMessageText(attachMessage.ChatID, attachMessage.ID, message.Text))
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		if hasButtons {
			_, err := client.api.Request(tgbotapi.NewEditMessageReplyMarkup(attachMessage.ChatID, attachMessage.ID, keyboard))
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("keyboard setup error")
}

func (client *Client) PrepareKeyboard(message bottypes.Message) (tgbotapi.InlineKeyboardMarkup, bool) {
	var keyboard tgbotapi.InlineKeyboardMarkup

	if len(message.ButtonRows) == 0 {
		return tgbotapi.InlineKeyboardMarkup{}, false
	}

	if len(message.ButtonRows) > 0 {
		var buttonRows [][]tgbotapi.InlineKeyboardButton
		for _, buttonRow := range message.ButtonRows {
			var buttons []tgbotapi.InlineKeyboardButton
			for _, button := range buttonRow.Buttons {
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Command.Text))
			}
			for _, button := range buttonRow.CheckboxButtons {
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Command.Text))
			}
			buttonRows = append(buttonRows, buttons)
		}

		keyboard = tgbotapi.NewInlineKeyboardMarkup(buttonRows...)
	}

	return keyboard, true
}

func (client *Client) SendKeyboard(message bottypes.Message) error {
	keyboard, exists := client.PrepareKeyboard(message)

	if !exists {
		return BotError{message: "Send keyboard error: no keyboard"}
	}

	if client.lastMessage.ID != 0 || client.lastKeyboardMessage.ID != 0 {
		err := client.SetupKeyboard(message, keyboard)
		if err != nil {
			return BotError{message: "Send keyboard error: " + err.Error()}
		}
	}
	return nil
}

func (client *Client) SendMessage(message bottypes.Message) error {

	msg := tgbotapi.NewMessage(message.ChatID, message.Text)
	keyboard, exists := client.PrepareKeyboard(message)
	if exists {
		msg.ReplyMarkup = keyboard
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

	sendErr := client.SendMessage(responseMessage)
	if sendErr != nil {
		panic(sendErr)
	}
}

func (client *Client) removeMessagesByTrigger() error {
	defer func() {
		client.messagesToRemove = nil
	}()

	// Reverse or not to reverse...

	// for i := len(client.messagesToRemove) - 1; i >= 0; i-- {
	// 	msgToDelete := tgbotapi.DeleteMessageConfig{
	// 		ChatID:    client.messagesToRemove[i].ChatID,
	// 		MessageID: client.messagesToRemove[i].ID,
	// 	}
	// 	_, err := client.api.Request(msgToDelete)
	// 	if err != nil {
	// 		return fmt.Errorf("remove error: %w", err)
	// 	}
	// }

	for _, v := range client.messagesToRemove {
		msgToDelete := tgbotapi.DeleteMessageConfig{
			ChatID:    v.ChatID,
			MessageID: v.ID,
		}
		_, err := client.api.Request(msgToDelete)
		if err != nil {
			return fmt.Errorf("remove error: %w", err)
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

		receivedMessage, chatID, err := client.parseMessage(update)
		if err != nil {
			client.sendErrorMessage(chatID, fmt.Errorf("parse error: %w", err))
			continue
		}

		req := client.cmdhandler.NewCommandHandlerRequest(receivedMessage, true)
		handlerResult, err := client.cmdhandler.Handle(req)
		if err != nil {
			client.sendErrorMessage(chatID, fmt.Errorf("handle command error: %w", err))
			continue
		}

		for _, response := range handlerResult.GetResponses() {
			for _, message := range response.GetMessages() {
				if message.ChatID == 0 {
					panic("Chat ID = 0")
				}

				if response.ContainsTrigger(bottypes.NothingTrigger) {
					continue
				}

				if response.ContainsTrigger(bottypes.StartKeyboardTrigger) {
					err := client.SendKeyboard(message)

					if err != nil {
						panic(err)
					}

					if response.ContainsTrigger(bottypes.AddToNextRemoveTrigger) {
						if client.lastKeyboardMessage.ID != 0 {
							client.addToRemoveMessagesQueue(client.lastKeyboardMessage)
						} else {
							client.addToRemoveMessagesQueue(client.lastMessage)
						}
					}

				} else {
					err := client.SendMessage(message)

					if err != nil {
						panic(err)
					}

					if response.ContainsTrigger(bottypes.AddToNextRemoveTrigger) {
						client.addToRemoveMessagesQueue(client.lastMessage)
					}
				}

			}
			if response.ContainsTrigger(bottypes.RemoveTrigger) {
				err := client.removeMessagesByTrigger()
				if err != nil {
					panic(err)
				}

				if response.ContainsTrigger(bottypes.StopKeyboardTrigger) {
					client.lastKeyboardMessage = bottypes.Message{}
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
