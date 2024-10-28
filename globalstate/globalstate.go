package globalstate

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
