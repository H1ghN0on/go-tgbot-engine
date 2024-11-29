package handlers_example

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands/example"
)

var next_month_symbol = ">"
var prev_month_symbol = "<"
var next_year_symbol = ">>"
var prev_year_symbol = "<<"
var empty_symbol = " "
var date_time_format = "2006-01-02 15:04"
var date_format = "2006-01-02"
var time_format = "15:04"

type CalendarHandler struct {
	Handler

	currentTime time.Time
	chosenDate  time.Time

	availableTime []time.Time
}

func NewCalendarHandler(gs ExampleGlobalStater) *CalendarHandler {

	h := &CalendarHandler{}
	h.gs = gs

	h.Commands = map[bottypes.Command][]func(params handlers.HandlerParams) (handlers.HandlerResponse, error){
		cmd.CalendarStartCommand:        {h.ModifyHandler(h.CalendarStartHandler, []int{})},
		cmd.CalendarChooseCommand:       {h.ModifyHandler(h.CalendarChooseHandler, []int{handlers.KeyboardStarter, handlers.StateBackable})},
		cmd.CalendarChooseFirstCommand:  {h.ModifyHandler(h.CalendarChooseFirstHandler, []int{})},
		cmd.CalendarChooseSecondCommand: {h.ModifyHandler(h.CalendarChooseSecondHandler, []int{})},
		cmd.CalendarLaunchCommand:       {h.ModifyHandler(h.CalendarLaunchHandler, []int{handlers.KeyboardStarter, handlers.CommandBackable, handlers.RemovableByTrigger})},
		cmd.CalendarNextMonthCommand:    {h.ModifyHandler(h.CalendarNextMonthHandler, []int{})},
		cmd.CalendarPrevMonthCommand:    {h.ModifyHandler(h.CalendarPrevMonthHandler, []int{})},
		cmd.CalendarNextYearCommand:     {h.ModifyHandler(h.CalendarNextYearHandler, []int{})},
		cmd.CalendarPrevYearCommand:     {h.ModifyHandler(h.CalendarPrevYearHandler, []int{})},
		cmd.CalendarSetDayCommand:       {h.ModifyHandler(h.CalendarSetDayHandler, []int{handlers.KeyboardStarter, handlers.CommandBackable, handlers.RemovableByTrigger})},
		cmd.CalendarSetTimeCommand:      {h.ModifyHandler(h.CalendarSetTimeHandler, []int{handlers.KeyboardStarter})},
		cmd.CalendarFinishCommand:       {h.ModifyHandler(h.CalendarFinishHandler, []int{handlers.KeyboardStopper, handlers.RemoveTriggerer})},
	}

	return h
}

func (handler *CalendarHandler) Handle(params handlers.HandlerParams) ([]handlers.HandlerResponse, error) {
	var res []handlers.HandlerResponse

	handleFuncs, ok := handler.GetCommandFromMap(params.Command)
	if !ok {
		panic("wrong handler")
	}

	for _, handleFunc := range handleFuncs {
		response, err := handleFunc(params)
		if err != nil {
			return []handlers.HandlerResponse{}, err
		}
		res = append(res, response)
	}

	return res, nil
}

func (handler *CalendarHandler) HandleBackCommand(params handlers.HandlerParams) ([]handlers.HandlerResponse, error) {
	var response []handlers.HandlerResponse

	var res handlers.HandlerResponse

	if params.Command.Equal(cmd.CalendarChooseCommand) {
		handler.availableTime = nil
	}

	if params.Command.Equal(cmd.CalendarLaunchCommand) {
		handler.chosenDate = time.Time{}
	}

	response = append(response, res)
	return response, nil
}

func (handler *CalendarHandler) CalendarStartHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.CalendarChooseCommand)
	res.NextState = "calendar-state"

	handler.currentTime = time.Now()

	return res, nil
}

func (handler *CalendarHandler) CalendarChooseHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	calendar1 := bottypes.Button{
		ChatID:  params.Message.Info.ChatID,
		Text:    "Schedule 1",
		Command: cmd.CalendarChooseFirstCommand,
	}

	calendar2 := bottypes.Button{
		ChatID:  params.Message.Info.ChatID,
		Text:    "Schedule 2",
		Command: cmd.CalendarChooseSecondCommand,
	}

	res.Messages = append(res.Messages, bottypes.Message{
		ChatID: params.Message.Info.ChatID,
		Text:   "Choose schedule",
		ButtonRows: []bottypes.ButtonRows{
			{Buttons: []bottypes.Button{calendar1, calendar2}},
		},
	})

	res.NextCommands = append(res.NextCommands, cmd.CalendarChooseFirstCommand, cmd.CalendarChooseSecondCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarChooseFirstHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	handler.availableTime = handler.gs.GetScheduleFirst()
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarChooseSecondHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	handler.availableTime = handler.gs.GetScheduleSecond()
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarLaunchHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	calendar := handler.buildCalendar()

	var buttonRows []bottypes.ButtonRows

	for _, calendarRow := range calendar {
		var buttonRow bottypes.ButtonRows
		for _, calendarData := range calendarRow {

			button := bottypes.Button{
				ChatID:  params.Message.Info.ChatID,
				Text:    calendarData,
				Command: cmd.NothingnessCommand,
			}

			if button.Text == "" {
				button.Text = " "
			}

			if button.Text == next_year_symbol {
				button.Command = cmd.CalendarNextYearCommand
			} else if button.Text == prev_year_symbol {
				button.Command = cmd.CalendarPrevYearCommand
			} else if button.Text == next_month_symbol {
				button.Command = cmd.CalendarNextMonthCommand
			} else if button.Text == prev_month_symbol {
				button.Command = cmd.CalendarPrevMonthCommand
			}

			currentDay, err := strconv.Atoi(button.Text)

			if err == nil {
				button.Command = cmd.CalendarSetDayCommand
				t := time.Date(handler.currentTime.Year(), handler.currentTime.Month(), currentDay, 0, 0, 0, 0, time.UTC)
				button.Command.Data = t.Format(date_format)
			}

			buttonRow.Buttons = append(buttonRow.Buttons, button)
		}
		buttonRows = append(buttonRows, buttonRow)
	}

	res.Messages = append(res.Messages, bottypes.Message{
		Text:       "Launch the calendar",
		ChatID:     params.Message.Info.ChatID,
		ButtonRows: buttonRows,
	})

	res.NextCommandToParse.ParseType = bottypes.DynamicButtonParse
	res.NextCommandToParse.Command = cmd.CalendarSetDayCommand
	res.NextCommandToParse.Exceptions = append(res.NextCommandToParse.Exceptions,
		cmd.NothingnessCommand,
		cmd.CalendarPrevMonthCommand,
		cmd.CalendarNextMonthCommand,
		cmd.CalendarNextYearCommand,
		cmd.CalendarPrevYearCommand,
	)

	res.NextCommands = append(res.NextCommands,
		cmd.NothingnessCommand,
		cmd.CalendarPrevMonthCommand,
		cmd.CalendarNextMonthCommand,
		cmd.CalendarNextYearCommand,
		cmd.CalendarPrevYearCommand,
		cmd.CalendarSetDayCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarNextMonthHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	handler.currentTime = handler.currentTime.AddDate(0, 1, 0)
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarPrevMonthHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	handler.currentTime = handler.currentTime.AddDate(0, -1, 0)
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarNextYearHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	handler.currentTime = handler.currentTime.AddDate(1, 0, 0)
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarPrevYearHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	handler.currentTime = handler.currentTime.AddDate(-1, 0, 0)
	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.CalendarLaunchCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarSetDayHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	date, err := time.Parse(date_format, params.Command.Data)
	if err != nil {
		panic(err)
	}

	sort.Slice(handler.availableTime, func(i, j int) bool {
		return handler.availableTime[i].Before(handler.availableTime[j])
	})

	handler.chosenDate = date

	var buttons []bottypes.ButtonRows
	for _, time := range handler.availableTime {
		if dateEqualByDay(time, date) {
			dateString := time.Format(time_format)
			buttons = append(buttons, bottypes.ButtonRows{
				Buttons: []bottypes.Button{
					{
						ChatID: params.Message.Info.ChatID,
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

	res.Messages = append(res.Messages, bottypes.Message{
		Text:       "Select time",
		ButtonRows: buttons,
		ChatID:     params.Message.Info.ChatID,
	})

	res.NextCommandToParse.ParseType = bottypes.DynamicButtonParse
	res.NextCommandToParse.Command = cmd.CalendarSetTimeCommand
	res.NextCommands = append(res.NextCommands, cmd.CalendarSetTimeCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarSetTimeHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	date, err := time.Parse(time_format, params.Command.Data)
	if err != nil {
		panic(err)
	}

	handler.chosenDate = handler.chosenDate.Add(
		time.Hour*time.Duration(date.Hour()) + time.Minute*time.Duration(date.Minute()))

	res.Messages = append(res.Messages, bottypes.Message{
		Text:   "You have chosen " + handler.chosenDate.Format(date_time_format),
		ChatID: params.Message.Info.ChatID,
	})

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.CalendarFinishCommand)

	return res, nil
}

func (handler *CalendarHandler) CalendarFinishHandler(params handlers.HandlerParams) (handlers.HandlerResponse, error) {
	var res handlers.HandlerResponse

	res.PostCommandsHandle.Commands = append(res.PostCommandsHandle.Commands, cmd.ShowCommandsCommand)
	res.NextState = "start-state"

	return res, nil
}

var months []string = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
var daysOfWeek []string = []string{"MO", "TU", "WE", "TH", "FR", "SA", "SU"}

func (handler CalendarHandler) buildCalendar() [][]string {
	var data [][]string

	var date = handler.currentTime

	data = append(data, buildYear(date))
	data = append(data, buildDaysOfWeek())
	data = append(data, buildWeeks(date, handler.availableTime)...)
	data = append(data, buildFooter())

	return data
}

func buildYear(t time.Time) []string {
	if t.Month() < time.January || t.Month() > time.December {
		panic("wrong month")
	}
	return []string{prev_year_symbol, months[t.Month()-1] + empty_symbol + fmt.Sprint(t.Year()), next_year_symbol}
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
		return dateEqualByDay(timeOfDay, t)
	})
}

func dateEqualByDay(t time.Time, other time.Time) bool {
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

		if hasTimeDayInSlice(date, availableTime) {
			weeks[currentWeek][getCorrectWeekday(weekday)] = fmt.Sprint(date.Day())
		}

		if weekday == time.Sunday {
			weeks = append(weeks, make([]string, 7))
			currentWeek++
		}
	}

	return weeks
}

func buildFooter() []string {
	return []string{prev_month_symbol, empty_symbol, next_month_symbol}
}

func getCorrectWeekday(weekday time.Weekday) time.Weekday {
	if weekday == time.Sunday {
		return time.Saturday
	}

	return weekday - 1
}
