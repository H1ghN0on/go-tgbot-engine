package handlers

import (
	"slices"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands/example"
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

type HandlerParams struct {
	Command bottypes.Command
	Message bottypes.ParsedMessage
}

func (hr HandlerResponse) GetMessages() []bottypes.Message {
	return hr.Messages
}

func (hr HandlerResponse) ContainsTrigger(trigger bottypes.Trigger) bool {
	return slices.Contains(hr.Triggers, trigger)
}

func (hr HandlerResponse) GetNextCommands() []bottypes.Command {
	return hr.NextCommands
}

func (hr HandlerResponse) GetNextCommandToParse() bottypes.ParseableCommand {
	return hr.NextCommandToParse
}

type HandlerResponse struct {
	Messages           []bottypes.Message
	Triggers           []bottypes.Trigger
	NextState          string
	PostCommandsHandle PostCommands
	NextCommands       []bottypes.Command
	NextCommandToParse bottypes.ParseableCommand
}

type PostCommands struct {
	Commands      []bottypes.Command
	IsBackCommand bool
}

type Handler struct {
	Commands map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error)
}

func (handler *Handler) HandleBackCommand(params HandlerParams) ([]HandlerResponse, error) {
	return []HandlerResponse{}, nil
}

func (handler *Handler) ModifyHandler(handlerFoo func(HandlerParams) (HandlerResponse, error), modifiers []int) func(HandlerParams) (HandlerResponse, error) {
	return func(params HandlerParams) (HandlerResponse, error) {
		if slices.Contains(modifiers, Nothingness) {
			return HandlerResponse{
				Triggers: []bottypes.Trigger{bottypes.NothingTrigger},
			}, nil
		}

		response, err := handlerFoo(params)
		if err != nil {
			return HandlerResponse{}, err
		}

		for idx, message := range response.Messages {

			if slices.Contains(modifiers, StateBackable) {
				response.Messages[idx].ButtonRows = append(response.Messages[idx].ButtonRows, bottypes.ButtonRows{
					Buttons: []bottypes.Button{
						{ChatID: message.ChatID, Text: "Back", Command: cmd.BackStateCommand},
					},
				})
			}

			if slices.Contains(modifiers, CommandBackable) {
				response.Messages[idx].ButtonRows = append(response.Messages[idx].ButtonRows, bottypes.ButtonRows{
					Buttons: []bottypes.Button{
						{ChatID: message.ChatID, Text: "Back", Command: cmd.BackCommandCommand},
					},
				})
			}

			if slices.Contains(modifiers, CheckboxableOne) {
				checkboxEmoji := "✅"
				emptyEmoji := "⬜"

				for rowIdx := range response.Messages[idx].ButtonRows {
					for checkboxIdx, checkboxButton := range response.Messages[idx].ButtonRows[rowIdx].CheckboxButtons {
						if checkboxButton.Active {
							response.Messages[idx].ButtonRows[rowIdx].CheckboxButtons[checkboxIdx].Text = checkboxEmoji + " " + checkboxButton.Text
						} else {
							response.Messages[idx].ButtonRows[rowIdx].CheckboxButtons[checkboxIdx].Text = emptyEmoji + " " + checkboxButton.Text
						}
					}
				}
			}

			if slices.Contains(modifiers, CheckboxableTwo) {
				checkboxEmoji := "✅"
				emptyEmoji := "⬜"
				for rowIdx := range response.Messages[idx].ButtonRows {
					for checkboxIdx, checkboxButton := range response.Messages[idx].ButtonRows[rowIdx].CheckboxButtons {
						newCheckboxButton := checkboxButton
						if checkboxButton.Active {
							newCheckboxButton.Text = checkboxEmoji
						} else {
							newCheckboxButton.Text = emptyEmoji
						}
						response.Messages[idx].ButtonRows[rowIdx].CheckboxButtons = append([]bottypes.CheckboxButton{newCheckboxButton}, response.Messages[idx].ButtonRows[rowIdx].CheckboxButtons...)
						response.Messages[idx].ButtonRows[rowIdx].CheckboxButtons[checkboxIdx+1].Command = cmd.NothingnessCommand
					}
				}
			}
		}

		if slices.Contains(modifiers, StateBackable) {
			response.NextCommands = append(response.NextCommands, cmd.BackStateCommand)
			response.NextCommandToParse.Exceptions = append(response.NextCommandToParse.Exceptions, cmd.BackStateCommand)
		}

		if slices.Contains(modifiers, CommandBackable) {
			response.NextCommands = append(response.NextCommands, cmd.BackCommandCommand)
			response.NextCommandToParse.Exceptions = append(response.NextCommandToParse.Exceptions, cmd.BackCommandCommand)
		}

		if slices.Contains(modifiers, RemovableByTrigger) {
			response.Triggers = append(response.Triggers, bottypes.AddToNextRemoveTrigger)
		}

		if slices.Contains(modifiers, KeyboardStarter) {
			response.Triggers = append(response.Triggers, bottypes.StartKeyboardTrigger)
		}

		if slices.Contains(modifiers, KeyboardStopper) {
			response.Triggers = append(response.Triggers, bottypes.StopKeyboardTrigger)
		}

		if slices.Contains(modifiers, RemoveTriggerer) {
			response.Triggers = append(response.Triggers, bottypes.RemoveTrigger)
		}

		return response, nil
	}
}

func (handler Handler) GetCommands() []bottypes.Command {
	commands := make([]bottypes.Command, len(handler.Commands))

	i := 0
	for k := range handler.Commands {
		commands[i] = k
		i++
	}

	return commands
}

func (handler Handler) GetCommandFromMap(command bottypes.Command) ([]func(HandlerParams) (HandlerResponse, error), bool) {
	for cmd, funcs := range handler.Commands {
		if cmd.Command == command.Command {
			return funcs, true
		}
	}

	return nil, false
}

func (handler *Handler) EmptyHandler(params HandlerParams) (HandlerResponse, error) {
	return HandlerResponse{}, nil
}

type HandlerResponseError struct {
	Message string
}

func (res HandlerResponseError) Error() string {
	return res.Message
}
