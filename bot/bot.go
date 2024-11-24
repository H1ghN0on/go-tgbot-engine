package bot

import (
	"fmt"
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/bot/client"
	"github.com/H1ghN0on/go-tgbot-engine/bot/notificator"
	"github.com/H1ghN0on/go-tgbot-engine/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Notificationer interface {
	GetMessages() []bottypes.Message
	GetUsers() []bottypes.User
	GetTimeoutSec() int
}

type BotError struct {
	message string
}

func (err BotError) Error() string {
	return err.message
}

type Notificator interface {
	AddNotification(notificator.Notificationer)
	Start()
	Stop()
}

type Bot struct {
	api         *tgbotapi.BotAPI
	clients     map[int64]*client.Client
	notificator Notificator
	onNewClient func() client.CommandHandler
}

func (client *Bot) parseMessage(update tgbotapi.Update) (bottypes.Message, int64, error) {
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
			Text:   update.CallbackQuery.Data,
		}

	} else if update.EditedMessage != nil {
		chatID = update.EditedMessage.Chat.ID
		return bottypes.Message{}, chatID, BotError{message: "editing is forbidden"}
	} else {
		return bottypes.Message{}, 0, BotError{message: "unknown message received"}
	}

	return receivedMessage, chatID, nil
}

func (bot *Bot) ListenMessages() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.api.GetUpdatesChan(u)

	if bot.notificator != nil {
		bot.notificator.Start()
	}

	logger.Bot().Info("listening messsages")

	for update := range updates {
		var receivedMessage bottypes.Message

		receivedMessage, chatID, err := bot.parseMessage(update)

		if err != nil {
			logger.Bot().Critical(string(chatID), fmt.Errorf("parse error: %w", err).Error())
			continue
		}

		logger.Bot().Info("new message received from", strconv.Itoa(int(receivedMessage.ChatID)))

		var activeClient *client.Client

		activeClient, ok := bot.clients[receivedMessage.ChatID]

		if !ok {
			logger.Bot().Info("adding new client with id", string(receivedMessage.ChatID))

			ch := bot.onNewClient()

			activeClient = client.NewClient(bot.api, ch, receivedMessage.ChatID)
			bot.clients[receivedMessage.ChatID] = activeClient
		}

		go func() {
			activeClient.HandleNewMessage(receivedMessage)
		}()
	}
}

func (bot *Bot) notificationHandler(notification notificator.Notificationer) {
	logger.Bot().Info(
		"notification timeout, sending",
		fmt.Sprint(len(notification.GetMessages())),
		"messages to", fmt.Sprint(len(notification.GetUsers())),
		"users")

	if len(notification.GetMessages()) == 0 || len(notification.GetUsers()) == 0 {
		return
	}

	for _, user := range notification.GetUsers() {
		for _, message := range notification.GetMessages() {
			message := tgbotapi.NewMessage(user.UserID, message.Text)
			bot.api.Send(message)
		}
	}
}

func (bot *Bot) AddStaticNotification(messages []bottypes.Message, users UserNotificationType, timeout int) {
	if bot.notificator == nil {
		return
	}
	notification := notificator.NewStaticNotification(messages, bot.ChooseUserNotificator(users), timeout)
	bot.notificator.AddNotification(notification)
}

func (bot *Bot) AddDynamicNotification(messages func() []bottypes.Message, users UserNotificationType, timeout int) {
	if bot.notificator == nil {
		return
	}
	notification := notificator.NewDynamicNotification(messages, bot.ChooseUserNotificator(users), timeout)
	bot.notificator.AddNotification(notification)
}

func NewBot(api *tgbotapi.BotAPI, onNewClient func() client.CommandHandler, useNotificator bool) *Bot {
	bot := &Bot{
		api: api,
	}
	bot.clients = make(map[int64]*client.Client)
	bot.onNewClient = onNewClient

	if useNotificator {
		bot.notificator = notificator.NewNotificator(
			nil,
			bot.notificationHandler)
	}

	return bot
}
