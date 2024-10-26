package bottypes

type Command struct {
	Text string
}

type Button struct {
	ChatID  int64
	Text    string
	Command Command
}

type ButtonRows struct {
	Buttons []Button
}

type Message struct {
	ID         int
	ChatID     int64
	Text       string
	ButtonRows []ButtonRows
}
