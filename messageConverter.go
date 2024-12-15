package messageConverter

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageConverter struct {
	// api    *tgbotapi.BotAPI
	// chatID int64
}

func (msgCnv MessageConverter) PrepareKeyboard(message bottypes.Message) (tgbotapi.InlineKeyboardMarkup, bool) {
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

// func (msgCnv MessageConverter) msgBaseInfo(message bottypes.Message) bottypes.Message {

// 	msgBase := bottypes.Message{
// 		ID:         message.ID,
// 		ChatID:     message.ChatID,
// 		UserName:   message.UserName,
// 		Text:       message.Text,
// 		ParseMode:  message.ParseMode,
// 		ButtonRows: message.ButtonRows,
// 	}

// 	return msgBase
// }

func (msgCnv MessageConverter) ConvertTypeNewMessage(message bottypes.Message) tgbotapi.MessageConfig {

	// zxc := msgCnv.msgBaseInfo(message)
	zxc := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           message.ChatID,
			ReplyToMessageID: 0,
		},
		Text:                  message.Text,
		DisableWebPagePreview: false,
	}

	return zxc
}

func (msgCnv MessageConverter) ConvertTypeNewEditMessageText(message bottypes.Message) tgbotapi.EditMessageTextConfig {

	zxc := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    message.ChatID,
			MessageID: message.ID,
		},
		Text: message.Text,
	}

	return zxc
}

func (msgCnv MessageConverter) ConvertTypeNewEditMessageTextAndMarkup(message bottypes.Message) tgbotapi.EditMessageTextConfig {

	keyboard, exists := msgCnv.PrepareKeyboard(message)

	if !exists {
		logger.GlobalLogger.Warning("|from: messageConverter|send keyboard error: no keyboard")
	}

	zxc := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      message.ChatID,
			MessageID:   message.ID,
			ReplyMarkup: &keyboard,
		},
		Text: message.Text,
	}
	return zxc
}

func (msgCnv MessageConverter) ConvertTypeNewInlineKeyboardButtonData(message bottypes.Message, options bottypes.Button) tgbotapi.InlineKeyboardButton {

	// var commandData bottypes.Button
	commandData := options

	data := string(commandData.Command.Command + commandData.Command.Data)
	zxc := tgbotapi.InlineKeyboardButton{
		Text:         message.Text,
		CallbackData: &data,
	}
	return zxc
}
