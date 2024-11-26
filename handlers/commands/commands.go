package commands

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

var NothingnessCommand bottypes.Command = bottypes.Command{Command: "/nothingness", Description: "Пустышка", SkipOnBack: true}

// Back Handler
var BackCommandCommand bottypes.Command = bottypes.Command{Command: "/back_command", Description: "Вернуться к предыдущей команде"}
var BackStateCommand bottypes.Command = bottypes.Command{Command: "/back_state", Description: "Вернуться к предыдущему состоянию"}

var Commands = []bottypes.Command{
	NothingnessCommand,
	BackCommandCommand,
	BackStateCommand,
}
