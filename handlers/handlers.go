package handlers

import (
	"slices"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

const (
	StateBackable = iota
	CommandBackable
	RemovableByTrigger
	KeyboardStarter
	KeyboardStopper
	RemoveTriggerer
	CheckboxableOne
	CheckboxableTwo
	Nothingness
)

type Handler struct {
	gs       GlobalStater
	commands map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error)
}

func (handler *Handler) HandleBackCommand(params HandlerParams) ([]HandlerResponse, error) {
	return []HandlerResponse{}, nil
}

type HandlerResponseError struct {
	message string
}

func (res HandlerResponseError) Error() string {
	return res.message
}

func (handler Handler) GetCommands() []bottypes.Command {
	commands := make([]bottypes.Command, len(handler.commands))

	i := 0
	for k := range handler.commands {
		commands[i] = k
		i++
	}

	return commands
}

func (handler Handler) GetCommandFromMap(command bottypes.Command) ([]func(HandlerParams) (HandlerResponse, error), bool) {
	for cmd, funcs := range handler.commands {
		if cmd.Command == command.Command {
			return funcs, true
		}
	}

	return nil, false
}

type HandlerParams struct {
	command bottypes.Command
	message bottypes.ParsedMessage
}

type HandlerResponse struct {
	messages           []bottypes.Message
	triggers           []bottypes.Trigger
	nextState          string
	postCommandsHandle PostCommands
	nextCommands       []bottypes.Command
	nextCommandToParse bottypes.ParseableCommand
}

type PostCommands struct {
	commands      []bottypes.Command
	isBackCommand bool
}

func (hr HandlerResponse) GetMessages() []bottypes.Message {
	return hr.messages
}

func (hr HandlerResponse) NextState() string {
	return hr.nextState
}

func (hr HandlerResponse) ContainsTrigger(trigger bottypes.Trigger) bool {
	return slices.Contains(hr.triggers, trigger)
}

func (hr HandlerResponse) GetNextCommands() []bottypes.Command {
	return hr.nextCommands
}

func (hr HandlerResponse) GetNextCommandToParse() bottypes.ParseableCommand {
	return hr.nextCommandToParse
}

func NewHandler(gs GlobalStater) *Handler {
	return &Handler{
		gs: gs,
	}
}

func (handler *Handler) ModifyHandler(handlerFoo func(HandlerParams) (HandlerResponse, error), modifiers []int) func(HandlerParams) (HandlerResponse, error) {
	return func(params HandlerParams) (HandlerResponse, error) {
		if slices.Contains(modifiers, Nothingness) {
			return HandlerResponse{
				triggers: []bottypes.Trigger{bottypes.NothingTrigger},
			}, nil
		}

		response, err := handlerFoo(params)
		if err != nil {
			return HandlerResponse{}, err
		}

		for idx, message := range response.messages {

			if slices.Contains(modifiers, StateBackable) {
				response.messages[idx].ButtonRows = append(response.messages[idx].ButtonRows, bottypes.ButtonRows{
					Buttons: []bottypes.Button{
						{ChatID: message.ChatID, Text: "Back", Command: cmd.BackStateCommand},
					},
				})
			}

			if slices.Contains(modifiers, CommandBackable) {
				response.messages[idx].ButtonRows = append(response.messages[idx].ButtonRows, bottypes.ButtonRows{
					Buttons: []bottypes.Button{
						{ChatID: message.ChatID, Text: "Back", Command: cmd.BackCommandCommand},
					},
				})
			}

			if slices.Contains(modifiers, CheckboxableOne) {
				checkboxEmoji := "✅"
				emptyEmoji := "⬜"

				for rowIdx := range response.messages[idx].ButtonRows {
					for checkboxIdx, checkboxButton := range response.messages[idx].ButtonRows[rowIdx].CheckboxButtons {
						if checkboxButton.Active {
							response.messages[idx].ButtonRows[rowIdx].CheckboxButtons[checkboxIdx].Text = checkboxEmoji + " " + checkboxButton.Text
						} else {
							response.messages[idx].ButtonRows[rowIdx].CheckboxButtons[checkboxIdx].Text = emptyEmoji + " " + checkboxButton.Text
						}
					}
				}
			}

			if slices.Contains(modifiers, CheckboxableTwo) {
				checkboxEmoji := "✅"
				emptyEmoji := "⬜"
				for rowIdx := range response.messages[idx].ButtonRows {
					for checkboxIdx, checkboxButton := range response.messages[idx].ButtonRows[rowIdx].CheckboxButtons {
						newCheckboxButton := checkboxButton
						if checkboxButton.Active {
							newCheckboxButton.Text = checkboxEmoji
						} else {
							newCheckboxButton.Text = emptyEmoji
						}
						response.messages[idx].ButtonRows[rowIdx].CheckboxButtons = append([]bottypes.CheckboxButton{newCheckboxButton}, response.messages[idx].ButtonRows[rowIdx].CheckboxButtons...)
						response.messages[idx].ButtonRows[rowIdx].CheckboxButtons[checkboxIdx+1].Command = cmd.NothingnessCommand
					}
				}
			}
		}

		if slices.Contains(modifiers, StateBackable) {
			response.nextCommands = append(response.nextCommands, cmd.BackStateCommand)
			response.nextCommandToParse.Exceptions = append(response.nextCommandToParse.Exceptions, cmd.BackStateCommand)
		}

		if slices.Contains(modifiers, CommandBackable) {
			response.nextCommands = append(response.nextCommands, cmd.BackCommandCommand)
			response.nextCommandToParse.Exceptions = append(response.nextCommandToParse.Exceptions, cmd.BackCommandCommand)
		}

		if slices.Contains(modifiers, RemovableByTrigger) {
			response.triggers = append(response.triggers, bottypes.AddToNextRemoveTrigger)
		}

		if slices.Contains(modifiers, KeyboardStarter) {
			response.triggers = append(response.triggers, bottypes.StartKeyboardTrigger)
		}

		if slices.Contains(modifiers, KeyboardStopper) {
			response.triggers = append(response.triggers, bottypes.StopKeyboardTrigger)
		}

		if slices.Contains(modifiers, RemoveTriggerer) {
			response.triggers = append(response.triggers, bottypes.RemoveTrigger)
		}

		return response, nil
	}
}

func (handler *Handler) EmptyHandler(params HandlerParams) (HandlerResponse, error) {
	return HandlerResponse{}, nil
}
