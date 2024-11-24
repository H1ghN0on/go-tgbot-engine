package globalstate

import (
	"math/rand/v2"
	"time"
)

type GlobalState struct {
	name    string
	surname string
	age     int
}

func (gs GlobalState) GetName() string {
	return gs.name
}

func (gs GlobalState) GetSurname() string {
	return gs.surname
}

func (gs GlobalState) GetAge() int {
	return gs.age
}

func (gs *GlobalState) SetName(name string) {
	gs.name = name
}

func (gs *GlobalState) SetSurname(surname string) {
	gs.surname = surname
}

func (gs *GlobalState) SetAge(age int) {
	gs.age = age
}

func (gs *GlobalState) GetDataForDynamicKeyboard() map[string][]string {
	return map[string][]string{
		"first_stage":  {"Necromantic", "Gotcha gotcha", "Hanipaganda"},
		"second_stage": {"Bad magus", "Third eye", "Midnight parade"},
	}
}

func (gs *GlobalState) GetScheduleFirst() (res []time.Time) {
	today := time.Now()

	for i := 0; i < 5; i++ {
		hour := rand.IntN(24)
		minute := rand.IntN(60)
		second := rand.IntN(60)

		res = append(res, time.Date(today.Year(), today.Month(), today.Day(), hour, minute, second, 0, time.UTC))
	}

	return res
}

func (gs *GlobalState) GetScheduleSecond() (res []time.Time) {
	today := time.Now()

	for i := 0; i < 10; i++ {
		today = today.AddDate(0, 0, 1)
		for i := 0; i < 10; i++ {
			hour := rand.IntN(24)
			minute := rand.IntN(60)
			second := rand.IntN(60)

			res = append(res, time.Date(today.Year(), today.Month(), today.Day(), hour, minute, second, 0, time.UTC))
		}
	}

	return res
}
