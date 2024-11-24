package bot

import (
	"math/rand/v2"
	"time"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
)

/* Notification messages */

type UserNotificationType int

const (
	OnlyMe            UserNotificationType = iota
	AllConnectedUsers UserNotificationType = iota
)

func (bot Bot) ChooseUserNotificator(nfType UserNotificationType) func() []bottypes.User {
	switch nfType {
	case OnlyMe:
		return bot.GetAllConnectedUsers()
	case AllConnectedUsers:
		return bot.GetOnlyMe()
	}

	panic("unknown notificator")
}

func (bot Bot) TimeNotification() []bottypes.Message {
	var messages []bottypes.Message

	t := time.Now().Format(time.RFC850)
	messages = append(messages, bottypes.Message{
		Text: "The time is " + t,
	})

	return messages
}

func (bot Bot) RandomTrackNotification() []bottypes.Message {

	var messages []bottypes.Message

	var tracks = []string{"Wire", "Senior Grang Botanist", "Ehiztaria", "Inbred Basilisk", "The Abhorrence", "The Legionary", "Silent Scream"}

	randomNumber := rand.IntN(len(tracks))
	messages = append(messages, bottypes.Message{
		Text: tracks[randomNumber],
	})
	return messages

}

/* Users for notification */

func (bot Bot) GetAllConnectedUsers() func() []bottypes.User {

	return func() []bottypes.User {
		var users []bottypes.User

		for _, client := range bot.clients {
			users = append(users, bottypes.User{
				UserID: client.GetUserID(),
			})
		}

		return users
	}
}

func (bot Bot) GetOnlyMe() func() []bottypes.User {
	return func() []bottypes.User {
		var users []bottypes.User
		users = append(users, bottypes.User{
			UserID: 872451555,
		})
		return users
	}
}
