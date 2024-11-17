package bottypes

import "strings"

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

type Trigger int

type Command struct {
	Command     string
	Description string
}

func (cmd Command) String() string {
	return cmd.Command
}

func (cmd Command) IsCommand() bool {
	return strings.HasPrefix(cmd.Command, "/")
}

func (cmd Command) Equal(other Command) bool {
	return cmd.Command == other.Command
}

const (
	RemoveTrigger          Trigger = iota
	AddToNextRemoveTrigger Trigger = iota
	StartKeyboardTrigger   Trigger = iota
	StopKeyboardTrigger    Trigger = iota
	NothingTrigger         Trigger = iota
)
