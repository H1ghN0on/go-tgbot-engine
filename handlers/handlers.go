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
