package messageСonverter

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PrepareKeyboard(message bottypes.Message) (tgbotapi.InlineKeyboardMarkup, bool) {
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

func NewMessage(message bottypes.Message) tgbotapi.MessageConfig {
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: message.ChatID,
		},
		Text:      message.Text,
		ParseMode: message.ParseMode.СonvertToAPI(),
	}
}

func NewEditMessageText(message bottypes.Message, lastMessage bottypes.Message) tgbotapi.EditMessageTextConfig {
	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    lastMessage.ChatID,
			MessageID: lastMessage.ID,
		},
		Text:      message.Text,
		ParseMode: message.ParseMode.СonvertToAPI(),
	}
}

func NewEditMessageTextAndMarkup(message bottypes.Message, lastMessage bottypes.Message) tgbotapi.EditMessageTextConfig {

	keyboard, exists := PrepareKeyboard(message)

	if !exists {
		logger.MessageConverter().Warning("|from: messageConverter|send keyboard error: no keyboard")
	}

	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      lastMessage.ChatID,
			MessageID:   lastMessage.ID,
			ReplyMarkup: &keyboard,
		},
		Text:      message.Text,
		ParseMode: message.ParseMode.СonvertToAPI(),
	}
}
