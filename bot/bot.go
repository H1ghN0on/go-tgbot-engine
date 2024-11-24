package bot

import (
	"fmt"
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/bot/client"
	"github.com/H1ghN0on/go-tgbot-engine/bot/notificator"
	"github.com/H1ghN0on/go-tgbot-engine/globalstate"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	"github.com/H1ghN0on/go-tgbot-engine/logger"
	"github.com/H1ghN0on/go-tgbot-engine/statemachine"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var shouldUseNotifications bool = false

type BotError struct {
	message string
}

func (err BotError) Error() string {
	return err.message
}

type Notificator interface {
	Start()
	Stop()
}

type Bot struct {
	api         *tgbotapi.BotAPI
	clients     map[int64]*client.Client
	gs          *globalstate.GlobalState
	notificator Notificator
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
			sm := statemachine.NewStateMachine()
			ch := handlers.NewCommandHandler(sm, bot.gs)
			activeClient = client.NewClient(bot.api, ch, receivedMessage.ChatID)
			bot.clients[receivedMessage.ChatID] = activeClient
		}

		activeClient.HandleNewMessage(receivedMessage)
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

func NewBot(api *tgbotapi.BotAPI, gs *globalstate.GlobalState) *Bot {
	bot := &Bot{
		api: api,
		gs:  gs,
	}

	bot.clients = make(map[int64]*client.Client)

	if shouldUseNotifications {
		notification := notificator.NewStaticNotification([]bottypes.Message{{Text: "Ravevenge"}, {Text: "Glass Cages"}}, bot.OnlyMe, 5)
		notification2 := notificator.NewStaticNotification([]bottypes.Message{{Text: "Crypteque"}}, bot.AllConnectedUsers, 10)

		dynamicNotification := notificator.NewDynamicNotification(bot.TimeNotification, bot.OnlyMe, 5)
		dynamicNotification2 := notificator.NewDynamicNotification(bot.RandomTrackNotification, bot.AllConnectedUsers, 10)

		bot.notificator = notificator.NewNotificator(
			[]notificator.Notificationer{*notification, *notification2, *dynamicNotification, *dynamicNotification2},
			bot.notificationHandler)
	}

	return bot
}
