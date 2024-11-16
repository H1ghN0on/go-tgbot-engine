package handlers

import (
	"slices"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
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

type HandlerParams struct {
	message bottypes.Message
}

type HandlerResponse struct {
	messages           []bottypes.Message
	triggers           []bottypes.Trigger
	nextState          string
	postCommandsHandle []bottypes.Command
	nextCommands       []bottypes.Command
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
						{ChatID: message.ChatID, Text: "Back", Command: "/back_state"},
					},
				})
			}

			if slices.Contains(modifiers, CommandBackable) {
				response.messages[idx].ButtonRows = append(response.messages[idx].ButtonRows, bottypes.ButtonRows{
					Buttons: []bottypes.Button{
						{ChatID: message.ChatID, Text: "Back", Command: "/back_command"},
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
						response.messages[idx].ButtonRows[rowIdx].CheckboxButtons[checkboxIdx+1].Command = "/nothingness"
					}
				}
			}
		}

		if slices.Contains(modifiers, StateBackable) {
			response.nextCommands = append(response.nextCommands, "/back_state")
		}

		if slices.Contains(modifiers, CommandBackable) {
			response.nextCommands = append(response.nextCommands, "/back_command")
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
