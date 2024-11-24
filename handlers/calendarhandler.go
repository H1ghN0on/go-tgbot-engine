package handlers

import (
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type CalendarHandler struct {
	Handler
}

func NewCalendarHandler(gs GlobalStater) *CalendarHandler {

	h := &CalendarHandler{}
	h.gs = gs

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.CalendarStartCommand:  {h.ModifyHandler(h.CalendarStartHandler, []int{})},
		cmd.CalendarLaunchCommand: {h.ModifyHandler(h.CalendarLaunchHandler, []int{RemovableByTrigger, StateBackable})},
	}

	return h
}

func (handler *CalendarHandler) Handle(params HandlerParams) ([]HandlerResponse, error) {
	var res []HandlerResponse

	handleFuncs, ok := handler.GetCommandFromMap(params.command)
	if !ok {
		panic("wrong handler")
	}

	for _, handleFunc := range handleFuncs {
		response, err := handleFunc(params)
		if err != nil {
			return []HandlerResponse{}, err
		}
		res = append(res, response)
	}

	return res, nil
}

func (handler *CalendarHandler) CalendarStartHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.CalendarLaunchCommand)
	res.nextState = "calendar-state"

	return res, nil
}

func (handler *CalendarHandler) CalendarLaunchHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	res.messages = append(res.messages, bottypes.Message{
		Text:   "Launch the calendar",
		ChatID: params.message.Info.ChatID,
	})

	return res, nil
}
