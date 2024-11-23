package bot

import (
	"fmt"
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/bot/client"
	"github.com/H1ghN0on/go-tgbot-engine/globalstate"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	"github.com/H1ghN0on/go-tgbot-engine/logger"
	"github.com/H1ghN0on/go-tgbot-engine/statemachine"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotError struct {
	message string
}

func (err BotError) Error() string {
	return err.message
}

type Bot struct {
	api     *tgbotapi.BotAPI
	clients map[int64]*client.Client
	gs      *globalstate.GlobalState
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
			sm := statemachine.NewStateMachine()
			ch := handlers.NewCommandHandler(sm, bot.gs)
			activeClient = client.NewClient(bot.api, ch)
			bot.clients[receivedMessage.ChatID] = activeClient
		}

		activeClient.HandleNewMessage(receivedMessage)
	}
}

func NewBot(api *tgbotapi.BotAPI, gs *globalstate.GlobalState) *Bot {
	bot := &Bot{
		api: api,
		gs:  gs,
	}

	bot.clients = make(map[int64]*client.Client)
	return bot
}
