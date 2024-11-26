package globalstate_example

import (
	"math/rand/v2"
	"time"
)

type ExampleGlobalState struct {
	name    string
	surname string
	age     int
}

func (gs ExampleGlobalState) GetName() string {
	return gs.name
}

func (gs ExampleGlobalState) GetSurname() string {
	return gs.surname
}

func (gs ExampleGlobalState) GetAge() int {
	return gs.age
}

func (gs *ExampleGlobalState) SetName(name string) {
	gs.name = name
}

func (gs *ExampleGlobalState) SetSurname(surname string) {
	gs.surname = surname
}

func (gs *ExampleGlobalState) SetAge(age int) {
	gs.age = age
}

func (gs *ExampleGlobalState) GetDataForDynamicKeyboard() map[string][]string {
	return map[string][]string{
		"first_stage":  {"Necromantic", "Gotcha gotcha", "Hanipaganda"},
		"second_stage": {"Bad magus", "Third eye", "Midnight parade"},
	}
}

func (gs *ExampleGlobalState) GetScheduleFirst() (res []time.Time) {
	today := time.Now()

	for i := 0; i < 5; i++ {
		hour := rand.IntN(24)
		minute := rand.IntN(60)
		second := rand.IntN(60)

		res = append(res, time.Date(today.Year(), today.Month(), today.Day(), hour, minute, second, 0, time.UTC))
	}

	return res
}

func (gs *ExampleGlobalState) GetScheduleSecond() (res []time.Time) {
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
