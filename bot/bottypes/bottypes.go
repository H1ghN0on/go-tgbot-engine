package bottypes

import (
	"slices"
	"strings"
)

type ParseCommandType int

type User struct {
	UserID int64
}

const (
	AnyTextParse       ParseCommandType = iota
	DynamicButtonParse ParseCommandType = iota
	NoParse            ParseCommandType = iota
)

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

type ParsedMessage struct {
	Info    Message
	Command Command
}

type Trigger int

type Command struct {
	Command     string
	Description string
	Data        string
	SkipOnBack  bool
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

func (cmd Command) InSlice(commands []Command) bool {
	return slices.ContainsFunc(commands, func(command Command) bool {
		return cmd.Equal(command)
	})
}

type ParseableCommand struct {
	Command    Command
	Exceptions []Command
	ParseType  ParseCommandType
}

const (
	RemoveTrigger          Trigger = iota
	AddToNextRemoveTrigger Trigger = iota
	StartKeyboardTrigger   Trigger = iota
	StopKeyboardTrigger    Trigger = iota
	NothingTrigger         Trigger = iota
)
