package handlers

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
)

type CalendarHandler struct {
	Handler

	currentMonth time.Month
	currentYear  int

	chosenDate time.Time

	availableTime []time.Time
}

func NewCalendarHandler(gs GlobalStater) *CalendarHandler {

	h := &CalendarHandler{}
	h.gs = gs

	h.commands = map[bottypes.Command][]func(params HandlerParams) (HandlerResponse, error){
		cmd.CalendarStartCommand:        {h.ModifyHandler(h.CalendarStartHandler, []int{})},
		cmd.CalendarChooseCommand:       {h.ModifyHandler(h.CalendarChooseHandler, []int{KeyboardStarter, StateBackable})},
		cmd.CalendarChooseFirstCommand:  {h.ModifyHandler(h.CalendarChooseFirstHandler, []int{})},
		cmd.CalendarChooseSecondCommand: {h.ModifyHandler(h.CalendarChooseSecondHandler, []int{})},
		cmd.CalendarLaunchCommand:       {h.ModifyHandler(h.CalendarLaunchHandler, []int{KeyboardStarter, CommandBackable, RemovableByTrigger})},
		cmd.CalendarNextMonthCommand:    {h.ModifyHandler(h.CalendarNextMonthHandler, []int{})},
		cmd.CalendarPrevMonthCommand:    {h.ModifyHandler(h.CalendarPrevMonthHandler, []int{})},
		cmd.CalendarNextYearCommand:     {h.ModifyHandler(h.CalendarNextYearHandler, []int{})},
		cmd.CalendarPrevYearCommand:     {h.ModifyHandler(h.CalendarPrevYearHandler, []int{})},
		cmd.CalendarSetDayCommand:       {h.ModifyHandler(h.CalendarSetDayHandler, []int{KeyboardStarter})},
		cmd.CalendarSetTimeCommand:      {h.ModifyHandler(h.CalendarSetTimeHandler, []int{KeyboardStarter})},
		cmd.CalendarFinishCommand:       {h.ModifyHandler(h.CalendarFinishHandler, []int{KeyboardStopper, RemoveTriggerer})},
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

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.CalendarChooseCommand)
	res.nextState = "calendar-state"

	handler.currentMonth = time.Now().Month()
	if handler.currentMonth == time.December {
		handler.currentMonth = time.January
		handler.currentYear++
	} else {
		handler.currentMonth++
	}
	handler.currentYear = time.Now().Year()

	return res, nil
}

func (handler *CalendarHandler) CalendarChooseHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	calendar1 := bottypes.Button{
		ChatID:  params.message.Info.ChatID,
		Text:    "Schedule 1",
		Command: cmd.CalendarChooseFirstCommand,
	}

	calendar2 := bottypes.Button{
		ChatID:  params.message.Info.ChatID,
		Text:    "Schedule 2",
		Command: cmd.CalendarChooseSecondCommand,
	}

	res.messages = append(res.messages, bottypes.Message{
		ChatID: params.message.Info.ChatID,
		Text:   "Choose schedule",
		ButtonRows: []bottypes.ButtonRows{
			{Buttons: []bottypes.Button{calendar1, calendar2}},
		},
	})

	res.nextCommands = append(res.nextCommands, cmd.CalendarChooseCommand, cmd.CalendarChooseSecondCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarChooseFirstHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	handler.availableTime = handler.gs.GetScheduleFirst()
	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarChooseSecondHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	handler.availableTime = handler.gs.GetScheduleSecond()
	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarLaunchHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	calendar := handler.buildCalendar()

	var buttonRows []bottypes.ButtonRows

	for _, calendarRow := range calendar {
		var buttonRow bottypes.ButtonRows
		for _, calendarData := range calendarRow {

			button := bottypes.Button{
				ChatID:  params.message.Info.ChatID,
				Text:    calendarData,
				Command: cmd.NothingnessCommand,
			}

			if button.Text == "" {
				button.Text = " "
			}

			currentDay, err := strconv.Atoi(button.Text)

			if err == nil {
				button.Command = cmd.CalendarSetDayCommand
				t := time.Date(handler.currentYear, handler.currentMonth-1, currentDay, 0, 0, 0, 0, time.UTC)
				button.Command.Data = t.Format("2006-01-02")
			}

			if button.Text == ">>" {
				button.Command = cmd.CalendarNextYearCommand
			}

			if button.Text == "<<" {
				button.Command = cmd.CalendarPrevYearCommand
			}

			if button.Text == ">" {
				button.Command = cmd.CalendarNextMonthCommand
			}

			if button.Text == "<" {
				button.Command = cmd.CalendarPrevMonthCommand
			}

			buttonRow.Buttons = append(buttonRow.Buttons, button)
		}
		buttonRows = append(buttonRows, buttonRow)
	}

	res.messages = append(res.messages, bottypes.Message{
		Text:       "Launch the calendar",
		ChatID:     params.message.Info.ChatID,
		ButtonRows: buttonRows,
	})

	res.nextCommandToParse.ParseType = bottypes.DynamicButtonParse
	res.nextCommandToParse.Command = cmd.CalendarSetDayCommand
	res.nextCommandToParse.Exceptions = append(res.nextCommandToParse.Exceptions,
		cmd.NothingnessCommand,
		cmd.CalendarPrevMonthCommand,
		cmd.CalendarNextMonthCommand,
		cmd.CalendarNextYearCommand,
		cmd.CalendarPrevYearCommand,
	)

	res.nextCommands = append(res.nextCommands,
		cmd.NothingnessCommand,
		cmd.CalendarPrevMonthCommand,
		cmd.CalendarNextMonthCommand,
		cmd.CalendarNextYearCommand,
		cmd.CalendarPrevYearCommand,
		cmd.CalendarSetDayCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarNextMonthHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	if handler.currentMonth == time.December {
		handler.currentMonth = time.January
		handler.currentYear++
	} else {
		handler.currentMonth++
	}

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarPrevMonthHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	if handler.currentMonth == time.January {
		handler.currentMonth = time.December
		handler.currentYear--
	} else {
		handler.currentMonth--
	}

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarNextYearHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	handler.currentYear++

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarPrevYearHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	handler.currentYear--

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarSetDayHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	date, err := time.Parse("2006-01-02", params.command.Data)
	if err != nil {
		panic(err)
	}

	sort.Slice(handler.availableTime, func(i, j int) bool {
		return handler.availableTime[i].Before(handler.availableTime[j])
	})

	handler.chosenDate = date

	var buttons []bottypes.ButtonRows
	for _, time := range handler.availableTime {
		if dateEqual(time, date) {
			dateString := time.Format("15:04")
			buttons = append(buttons, bottypes.ButtonRows{
				Buttons: []bottypes.Button{
					{
						ChatID: params.message.Info.ChatID,
						Text:   dateString,
						Command: bottypes.Command{
							Command:     cmd.CalendarSetTimeCommand.Command,
							Description: cmd.CalendarSetTimeCommand.Description,
							Data:        dateString,
						},
					},
				},
			})
		}
	}

	res.messages = append(res.messages, bottypes.Message{
		Text:       "Select time",
		ButtonRows: buttons,
		ChatID:     params.message.Info.ChatID,
	})

	res.nextCommandToParse.ParseType = bottypes.DynamicButtonParse
	res.nextCommandToParse.Command = cmd.CalendarSetTimeCommand
	res.nextCommands = append(res.nextCommands, cmd.CalendarSetTimeCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarSetTimeHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	date, err := time.Parse("15:04", params.command.Data)
	if err != nil {
		panic(err)
	}

	handler.chosenDate = handler.chosenDate.Add(
		time.Hour*time.Duration(date.Hour()) + time.Minute*time.Duration(date.Minute()))

	res.messages = append(res.messages, bottypes.Message{
		Text:   "You have chosen " + handler.chosenDate.Format("2006-01-02 15:04"),
		ChatID: params.message.Info.ChatID,
	})

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.CalendarFinishCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarFinishHandler(params HandlerParams) (HandlerResponse, error) {
	var res HandlerResponse

	res.postCommandsHandle.commands = append(res.postCommandsHandle.commands, cmd.ShowCommandsCommand)
	res.nextState = "start-state"

	return res, nil
}

var months []string = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
var daysOfWeek []string = []string{"MO", "TU", "WE", "TH", "FR", "SA", "SU"}

func (handler CalendarHandler) buildCalendar() [][]string {
	var data [][]string

	var date = time.Date(handler.currentYear, handler.currentMonth, 0, 0, 0, 0, 0, time.UTC)

	data = append(data, buildYear(date))
	data = append(data, buildDaysOfWeek())
	data = append(data, buildWeeks(date, handler.availableTime)...)
	data = append(data, buildFooter())

	return data
}

func buildYear(t time.Time) []string {
	if t.Month() < 1 || t.Month() > 12 {
		panic("wrong month")
	}
	return []string{"<<", months[t.Month()-1] + " " + fmt.Sprint(t.Year()), ">>"}
}

func buildDaysOfWeek() []string {
	return daysOfWeek
}

func daysInMonth(t time.Time) int {
	t = time.Date(t.Year(), t.Month(), 32, 0, 0, 0, 0, time.UTC)
	daysInMonth := 32 - t.Day()
	return daysInMonth
}

func rangeCurrentMonth(t time.Time) func() time.Time {

	startRange := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	endRange := time.Date(t.Year(), t.Month(), daysInMonth(startRange), 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if startRange.After(endRange) {
			return time.Time{}
		}
		date := startRange
		startRange = startRange.AddDate(0, 0, 1)
		return date
	}
}

func hasTimeDayInSlice(t time.Time, availableTime []time.Time) bool {
	return slices.ContainsFunc(availableTime, func(timeOfDay time.Time) bool {
		return dateEqual(timeOfDay, t)
	})
}

func dateEqual(t time.Time, other time.Time) bool {
	return other.Day() == t.Day() &&
		other.Month() == t.Month() &&
		other.Year() == t.Year()
}

func buildWeeks(t time.Time, availableTime []time.Time) [][]string {
	var weeks [][]string

	currentWeek := 0
	weeks = append(weeks, make([]string, 7))
	for day := rangeCurrentMonth(t); ; {
		date := day()
		if date.IsZero() {
			break
		}

		weekday := date.Weekday()

		if weekday == time.Sunday {
			if hasTimeDayInSlice(date, availableTime) {
				weeks[currentWeek][getCorrectWeekday(weekday)] = fmt.Sprint(date.Day())
			}

			weeks = append(weeks, make([]string, 7))
			currentWeek++
			continue
		}

		if hasTimeDayInSlice(date, availableTime) {
			weeks[currentWeek][getCorrectWeekday(weekday)] = fmt.Sprint(date.Day())
		}
	}

	return weeks
}

func getCorrectWeekday(weekday time.Weekday) time.Weekday {
	if weekday == time.Sunday {
		return time.Saturday
	}

	return weekday - 1
}

func buildFooter() []string {
	return []string{"<", "", ">"}
}
