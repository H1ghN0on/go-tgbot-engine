package main

import (
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/H1ghN0on/go-tgbot-engine/bot"
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/bot/client"
	"github.com/H1ghN0on/go-tgbot-engine/globalstate"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
	handlersExample "github.com/H1ghN0on/go-tgbot-engine/handlers/handlers_example"
	"github.com/H1ghN0on/go-tgbot-engine/logger"
	"github.com/H1ghN0on/go-tgbot-engine/statemachine"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var gs globalstate.ExampleGlobalState

func configureCommandHandler(sm *statemachine.StateMachine, gs *globalstate.ExampleGlobalState) *handlers.CommandHandler {
	var handlerables []handlers.Handlerable

	handlerables = append(handlerables,
		handlersExample.NewSetInfoHandler(gs),
		handlersExample.NewKeyboardHandler(gs),
		handlersExample.NewLevelFourHandler(gs),
		handlersExample.NewStartHandler(gs),
		handlersExample.NewCheckboxHandler(gs),
		handlersExample.NewDynamicKeyboardHandler(gs),
		handlersExample.NewCalendarHandler(gs))

	commandHandler := handlers.NewCommandHandler(handlerables, sm)
	return commandHandler
}

func configureStateMachine() *statemachine.StateMachine {

	startState := statemachine.NewState(
		"start-state",
		
		cmd.ShowCommandsCommand,
		
		cmd.StartCommand,
		cmd.LevelOneCommand,
		cmd.LevelTwoCommand,
		cmd.LevelThreeCommand,
		cmd.ShowCommandsCommand,
		cmd.KeyboardStartCommand,
		cmd.LevelFourStartCommand,
		cmd.BigMessagesCommand,
		cmd.SetInfoStartCommand,
		cmd.CheckboxStartCommand,
		cmd.DynamicKeyboardStartCommand,
		cmd.CalendarStartCommand,
	)

	levelFourState := statemachine.NewState(
		"level-four-state",

		cmd.LevelFourStartCommand,

		cmd.LevelFourStartCommand,
		cmd.LevelFourOneCommand,
		cmd.LevelFourTwoCommand,
		cmd.LevelFourThreeCommand,
		cmd.LevelFourFourCommand,
		cmd.BackStateCommand,
	)

	keyboardState := statemachine.NewState(
		"keyboard-state",

		cmd.KeyboardStartCommand,

		cmd.KeyboardStartCommand,
		cmd.KeyboardOneCommand,
		cmd.KeyboardTwoCommand,
		cmd.KeyboardThreeCommand,
		cmd.KeyboardFinishCommand,
		cmd.BackStateCommand,
		cmd.BackCommandCommand,
	)

	infoState := statemachine.NewState(
		"info-state",

		cmd.SetInfoStartCommand,

		cmd.SetInfoStartCommand,
		cmd.SetNameCommand,
		cmd.SetSurnameCommand,
		cmd.SetAgeCommand,
		cmd.SetInfoEndCommand,
		cmd.BackStateCommand,
		cmd.BackCommandCommand,
	)

	checkboxState := statemachine.NewState(
		"checkbox-state",

		cmd.CheckboxStartCommand,

		cmd.CheckboxStartCommand,
		cmd.CheckboxFirstCommand,
		cmd.CheckboxSecondCommand,
		cmd.CheckboxThirdCommand,
		cmd.CheckboxFourthCommand,
		cmd.CheckboxAcceptCommand,
		cmd.BackStateCommand,
		cmd.NothingnessCommand,
	)

	dynamicKeyboardState := statemachine.NewState(
		"dynamic-keyboard-state",

		cmd.DynamicKeyboardStartCommand,

		cmd.DynamicKeyboardFirstStageCommand,
		cmd.DynamicKeyboardSecondStageCommand,
		cmd.DynamicKeyboardFinishCommand,
		cmd.BackStateCommand,
		cmd.BackCommandCommand,
	)

	calendarState := statemachine.NewState(
		"calendar-state",

		cmd.CalendarStartCommand,

		cmd.CalendarLaunchCommand,
		cmd.CalendarChooseCommand,
		cmd.CalendarChooseFirstCommand,
		cmd.CalendarChooseSecondCommand,
		cmd.CalendarNextMonthCommand,
		cmd.CalendarPrevMonthCommand,
		cmd.CalendarNextYearCommand,
		cmd.CalendarPrevYearCommand,
		cmd.CalendarSetDayCommand,
		cmd.CalendarSetTimeCommand,
		cmd.CalendarFinishCommand,
		cmd.BackStateCommand,
		cmd.BackCommandCommand,
		cmd.NothingnessCommand,
	)
	
	startState.SetAvailableStates(*levelFourState, *keyboardState, *infoState, *checkboxState, *startState, *dynamicKeyboardState, *calendarState)
	levelFourState.SetAvailableStates(*startState)
	keyboardState.SetAvailableStates(*startState)
	infoState.SetAvailableStates(*startState)
	checkboxState.SetAvailableStates(*startState)
	dynamicKeyboardState.SetAvailableStates(*startState)
	calendarState.SetAvailableStates(*startState)

	sm := &statemachine.StateMachine{}

	sm.AddStates(*startState, *levelFourState, *keyboardState, *infoState, *checkboxState, *dynamicKeyboardState, *calendarState)

	err := sm.SetStateByName("start-state")
	if err != nil {
		panic(err.Error())
	}

	return sm
}

func onNewClient() client.CommandHandler {
	sm := configureStateMachine()
	ch := configureCommandHandler(sm, &gs)
	return ch
}

func timeNotification() []bottypes.Message {
	var messages []bottypes.Message

	t := time.Now().Format(time.RFC850)
	messages = append(messages, bottypes.Message{
		Text: "The time is " + t,
	})

	return messages
}

func randomTrackNotification() []bottypes.Message {

	var messages []bottypes.Message

	var tracks = []string{"Wire", "Senior Grang Botanist", "Ehiztaria", "Inbred Basilisk", "The Abhorrence", "The Legionary", "Silent Scream"}

	randomNumber := rand.IntN(len(tracks))
	messages = append(messages, bottypes.Message{
		Text: tracks[randomNumber],
	})
	return messages

}

func main() {

	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	loggerStatus, exists := os.LookupEnv("LOGGER_LEVEL")
	if !exists {
		log.Println(".env does not contain LOGGER_LEVEL")
	}
	levelInt, err := strconv.Atoi(loggerStatus)
	if err != nil {
		log.Println("strconv.Atoi error:", err)
		levelInt = 0
	}

	logger.InitGlobalLoggerSettings(logger.LoggerSettings{
		Level: logger.LogLevel(levelInt),
	})

	tgBotKey, exists := os.LookupEnv("TELEGRAM_BOT_API")
	if !exists {
		panic(".env does not contain TELEGRAM_BOT_API")
	}
	botAPI, err := tgbotapi.NewBotAPI(tgBotKey)
	if err != nil {
		panic(err.Error())
	}

	exampleBot := bot.NewBot(botAPI, onNewClient, false)
	exampleBot.AddStaticNotification([]bottypes.Message{{Text: "Ravevenge"}}, bot.OnlyStorm, 10)
	exampleBot.AddStaticNotification([]bottypes.Message{{Text: "Crypteque"}}, bot.AllConnectedUsers, 5)

	exampleBot.AddDynamicNotification(timeNotification, bot.OnlyStorm, 10)
	exampleBot.AddDynamicNotification(randomTrackNotification, bot.AllConnectedUsers, 5)
	exampleBot.ListenMessages()
}
