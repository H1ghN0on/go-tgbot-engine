package bottypes

type Command struct {
	Text string
}

type Button struct {
	ChatID  int64
	Text    string
	Command Command
}

type CheckboxButton struct {
	ChatID  int64
	Text    string
	Command Command
	Active  bool
}

type ButtonRows struct {
	Buttons         []Button
	CheckboxButtons []CheckboxButton
}

type Message struct {
	ID         int
	ChatID     int64
	Text       string
	ButtonRows []ButtonRows
}
