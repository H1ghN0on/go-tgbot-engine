package handlers

import (
	"slices"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
)

const (
	StateBackable = iota
	CommandBackable
	RemovableByTrigger
	Keyboardable
	RemoveTriggerer
	CheckboxableOne
	CheckboxableTwo
	Nothingness
)

type Handler struct {
	gs       GlobalStater
	commands map[string][]func(params HandlerParams) HandlerResponse
}

func (handler Handler) GetCommands() []string {
	commands := make([]string, len(handler.commands))

	i := 0
	for k := range handler.commands {
		commands[i] = k
		i++
	}

	return commands
}

type SequentHandler struct {
	Handler
	commandSequence map[int]string
	active          int
}

type HandlerParams struct {
	message bottypes.Message
}

type HandlerResponse struct {
	messages             []bottypes.Message
	nextState            string
	isKeyboard           bool
	isRemovableByTrigger bool
	isRemoveTriggered    bool
	isNothing            bool
}

func (hr HandlerResponse) GetMessages() []bottypes.Message {
	return hr.messages
}

func (hr HandlerResponse) NextState() string {
	return hr.nextState
}

func (hr HandlerResponse) IsKeyboard() bool {
	return hr.isKeyboard
}

func (hr HandlerResponse) IsRemovableByTrigger() bool {
	return hr.isRemovableByTrigger
}

func NewHandler(gs GlobalStater) *Handler {
	return &Handler{
		gs: gs,
	}
}

func (handler *Handler) ModifyHandler(handlerFoo func(HandlerParams) HandlerResponse, modifiers []int) func(HandlerParams) HandlerResponse {
	return func(params HandlerParams) HandlerResponse {
		if slices.Contains(modifiers, Nothingness) {
			return HandlerResponse{
				isNothing: true,
			}
		}

		response := handlerFoo(params)
		for idx, message := range response.messages {

			if slices.Contains(modifiers, StateBackable) {
				response.messages[idx].ButtonRows = append(response.messages[idx].ButtonRows, bottypes.ButtonRows{
					Buttons: []bottypes.Button{
						{ChatID: message.ChatID, Text: "Back", Command: bottypes.Command{Text: "/back_state"}},
					},
				})
			}

			if slices.Contains(modifiers, CommandBackable) {
				response.messages[idx].ButtonRows = append(response.messages[idx].ButtonRows, bottypes.ButtonRows{
					Buttons: []bottypes.Button{
						{ChatID: message.ChatID, Text: "Back", Command: bottypes.Command{Text: "/back_command"}},
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
						response.messages[idx].ButtonRows[rowIdx].CheckboxButtons[checkboxIdx+1].Command.Text = "/nothingness"
					}
				}
			}
		}

		if slices.Contains(modifiers, RemovableByTrigger) {
			response.isRemovableByTrigger = true
		}

		if slices.Contains(modifiers, Keyboardable) {
			response.isKeyboard = true
		}

		if slices.Contains(modifiers, RemoveTriggerer) {
			response.isRemoveTriggered = true
		}

		return response
	}
}

func (handler *Handler) EmptyHandler(params HandlerParams) HandlerResponse {
	return HandlerResponse{}
}
