package handlers_example

import (
	"time"

	"github.com/H1ghN0on/go-tgbot-engine/handlers"
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
