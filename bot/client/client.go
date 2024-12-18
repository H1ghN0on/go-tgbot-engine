package client

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ClientError struct {
	message string
}

func (err ClientError) Error() string {
	return err.message
}

type HandlerResponser interface {
	GetMessages() []bottypes.Message
	ContainsTrigger(bottypes.Trigger) bool
	GetNextCommands() []bottypes.Command
	GetNextCommandToParse() bottypes.ParseableCommand
}

type CommandHandlerRequester interface {
	GetMessage() bottypes.ParsedMessage
}

type CommandHandlerResponser interface {
	GetResponses() []HandlerResponser
}

type CommandHandler interface {
	NewCommandHandlerRequest(msg bottypes.ParsedMessage) CommandHandlerRequester
	Handle(req CommandHandlerRequester) (CommandHandlerResponser, error)
}

type Client struct {
	cmdhandler         CommandHandler
	api                *tgbotapi.BotAPI
	chatID             int64
	lastMessage        bottypes.Message
	messagesToRemove   []bottypes.Message
	nextCommandToParse bottypes.ParseableCommand
}

func (client Client) GetUserID() int64 {
	return client.chatID
}

func (client Client) parseCommand(message bottypes.Message) bottypes.Command {
	command := bottypes.Command{
		Command: message.Text,
		Data:    "",
	}

	parseCommand := client.nextCommandToParse

	if parseCommand.ParseType == bottypes.NoParse ||
		parseCommand.Command.Command == "" {
		return command
	}

	for _, exception := range parseCommand.Exceptions {
		if exception.Command == message.Text {
			return command
		}
	}

	if parseCommand.ParseType == bottypes.AnyTextParse {
		command.Command = parseCommand.Command.Command
		command.Data = message.Text
		return command
	}

	if parseCommand.ParseType == bottypes.DynamicButtonParse {
		if !strings.HasPrefix(message.Text, parseCommand.Command.Command) {
			return command
		}

		data, _ := strings.CutPrefix(message.Text, parseCommand.Command.Command)
		command.Command = parseCommand.Command.Command
		command.Data = data
		return command
	}

	return command
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
	
	if client.lastMessage.ID == 0 {
		return fmt.Errorf("keyboard has no message to attach")
	}

	if hasText && message.Text != client.lastMessage.Text {
		if hasButtons {
			req := tgbotapi.NewEditMessageTextAndMarkup(client.lastMessage.ChatID, client.lastMessage.ID, message.Text, keyboard)
			req.ParseMode = client.lastMessage.ParseMode.СonvertToAPI()
			_, err := client.api.Request(req)
			if err != nil {
				return err
			}
			return nil
		} else {
			req := tgbotapi.NewEditMessageText(client.lastMessage.ChatID, client.lastMessage.ID, message.Text)
			req.ParseMode = client.lastMessage.ParseMode.СonvertToAPI()
			_, err := client.api.Request(req)
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		if hasButtons {
			_, err := client.api.Request(tgbotapi.NewEditMessageReplyMarkup(client.lastMessage.ChatID, client.lastMessage.ID, keyboard))
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
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Text, string(button.Command.Command+button.Command.Data)))
			}
			for _, button := range buttonRow.CheckboxButtons {
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Text, string(button.Command.Command)))
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
		return ClientError{message: "Send keyboard error: no keyboard"}
	}

	if client.lastMessage.ID == 0 {
		err := client.SendText(bottypes.Message{ChatID: message.ChatID, Text: message.Text})
		if err != nil {
			return ClientError{message: "Send keyboard error: " + err.Error()}
		}
	}

	err := client.SetupKeyboard(message, keyboard)
	if err != nil {
		return ClientError{message: "Send keyboard error: " + err.Error()}
	}

	client.lastMessage = bottypes.Message{
		ID:         client.lastMessage.ID,
		ChatID:     client.lastMessage.ChatID,
		Text:       message.Text,
		ButtonRows: message.ButtonRows,
	}

	return nil
}

func (client *Client) SendText(message bottypes.Message) error {

	msg := tgbotapi.NewMessage(client.chatID, message.Text)
	msg.ParseMode = message.ParseMode.СonvertToAPI()
	sent, err := client.api.Send(msg)
	if err != nil {
		return ClientError{message: "Send message error: " + err.Error()}
	}

	client.lastMessage = bottypes.Message{
		ID:         sent.MessageID,
		ChatID:     sent.Chat.ID,
		Text:       sent.Text,
		ButtonRows: message.ButtonRows,
		ParseMode:  message.ParseMode,
	}

	return nil
}

func (client *Client) SendMessage(message bottypes.Message) error {

	msg := tgbotapi.NewMessage(client.chatID, message.Text)
	msg.ParseMode = message.ParseMode.СonvertToAPI()
	keyboard, exists := client.PrepareKeyboard(message)
	if exists {
		msg.ReplyMarkup = keyboard
	}

	sent, err := client.api.Send(msg)
	if err != nil {
		return ClientError{message: "Send message error: " + err.Error()}
	}

	client.lastMessage = bottypes.Message{
		ID:         sent.MessageID,
		ChatID:     sent.Chat.ID,
		Text:       sent.Text,
		ButtonRows: message.ButtonRows,
		ParseMode:  client.lastMessage.ParseMode,
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
		logger.Client().Critical(fmt.Errorf("error text: %q, error: %w", responseMessage.Text, sendErr).Error())
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
			MessageID: v.ID,
			ChatID:    v.ChatID,
		}
		_, err := client.api.Request(msgToDelete)
		if err != nil {
			return fmt.Errorf("remove error: %w", err)
		}
	}

	return nil
}

func (client *Client) setMyCommands(chatID int64, res []HandlerResponser) error {
	var commands []tgbotapi.BotCommand

	if len(res) == 0 {
		_, err := client.api.Request(tgbotapi.NewDeleteMyCommandsWithScope(tgbotapi.NewBotCommandScopeChat(chatID)))
		return err
	}

	lastRes := res[len(res)-1]
	nextCommands := lastRes.GetNextCommands()

	if len(nextCommands) == 0 {
		_, err := client.api.Request(tgbotapi.NewDeleteMyCommandsWithScope(tgbotapi.NewBotCommandScopeChat(chatID)))
		return err
	}

	for _, command := range nextCommands {
		commands = append(commands, tgbotapi.BotCommand{Command: command.Command, Description: command.Description})
	}

	_, err := client.api.Request(tgbotapi.NewSetMyCommandsWithScope(
		tgbotapi.NewBotCommandScopeChat(chatID),
		commands...))
	return err
}

func (client *Client) setNextCommandToParse(command bottypes.ParseableCommand) {
	if command.Command.Command == "" || command.ParseType == bottypes.NoParse {
		client.nextCommandToParse = bottypes.ParseableCommand{}
	} else {
		logger.Client().Info("command", command.Command.Command, "will be parsed")
		client.nextCommandToParse = command
	}
}

func (client *Client) HandleNewMessage(receivedMessage bottypes.Message) {

	command := client.parseCommand(receivedMessage)

	parsedMessage := bottypes.ParsedMessage{
		Info:    receivedMessage,
		Command: command,
	}

	req := client.cmdhandler.NewCommandHandlerRequest(parsedMessage)
	handlerResult, err := client.cmdhandler.Handle(req)
	if err != nil {
		client.sendErrorMessage(parsedMessage.Info.ChatID, fmt.Errorf("bot error: %w", err))
		return
	}

	for _, response := range handlerResult.GetResponses() {
		for _, message := range response.GetMessages() {

			if message.ChatID == 0 {
				logger.Client().Critical("receiver of sending message is unknown")
				panic("Chat ID = 0")
			}

			if len(message.ButtonRows) == 0 && message.Text == "" {
				logger.Client().Warning("trying to send empty message, skipped")
				continue
			}

			if message.Text != "" && len(message.ButtonRows) == 0 {
				err := client.SendText(message)
				if err != nil {
					logger.Client().Critical(fmt.Errorf("error text:: %q, error: %w", message.Text, err).Error())
					panic(err)
				}
			} else if len(message.ButtonRows) != 0 && response.ContainsTrigger(bottypes.StartKeyboardTrigger) {
				err := client.SendKeyboard(message)
				if err != nil {
					panic(err)
				}
			} else if len(message.ButtonRows) != 0 {
				err := client.SendMessage(message)
				if err != nil {
					panic(err)
				}
			}

			if response.ContainsTrigger(bottypes.AddToNextRemoveTrigger) {
				logger.Client().Info("message", strconv.Itoa(client.lastMessage.ID), "marked to remove")
				client.addToRemoveMessagesQueue(client.lastMessage)
			}
		}
		if response.ContainsTrigger(bottypes.RemoveTrigger) {
			logger.Client().Info("removing all marked messages")
			err := client.removeMessagesByTrigger()
			if err != nil {
				panic(err)
			}
		}
		client.setNextCommandToParse(response.GetNextCommandToParse())
	}

	err = client.setMyCommands(parsedMessage.Info.ChatID, handlerResult.GetResponses())
	if err != nil {
		panic(err)
	}
}

func NewClient(api *tgbotapi.BotAPI, ch CommandHandler, chatID int64) *Client {
	return &Client{
		api:        api,
		cmdhandler: ch,
		chatID:     chatID,
	}
}
