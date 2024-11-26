package handlers_example

import (
	"fmt"
	"slices"
	"time"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"

	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands/example"
)

type ExampleGlobalStater interface {
	GetName() string
	GetSurname() string
	GetAge() int

	SetName(name string)
	SetSurname(surname string)
	SetAge(age int)

	GetScheduleFirst() []time.Time
	GetScheduleSecond() []time.Time

	GetDataForDynamicKeyboard() map[string][]string
}

type Handler struct {
	handlers.Handler
	gs ExampleGlobalStater
}

func (handler Handler) FindCommandInTheList(command bottypes.Command) (bottypes.Command, error) {
	if command.Command == "" || !command.IsCommand() {
		return bottypes.Command{}, fmt.Errorf("not command received")
	}

	index := slices.IndexFunc(cmd.Commands, func(com bottypes.Command) bool { return command.Equal(com) })
	if index == -1 {
		return bottypes.Command{}, fmt.Errorf("unknown command received")
	}

	trueCommand := cmd.Commands[index]
	trueCommand.Data = command.Data

	return trueCommand, nil
}
