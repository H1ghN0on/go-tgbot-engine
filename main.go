package main

import (
	"log"
	"os"
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot"
	"github.com/H1ghN0on/go-tgbot-engine/globalstate"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	"github.com/H1ghN0on/go-tgbot-engine/logger"
	"github.com/H1ghN0on/go-tgbot-engine/statemachine"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/joho/godotenv"
)

func configurateStateMachine(sm *statemachine.StateMachine) {
	startState := statemachine.NewState(
		"start-state",

		"/show_commands",

		"/level_one",
		"/level_two",
		"/level_three",
		"/show_commands",
		"/keyboard_start",
		"/create_error",
		"/level_four_start",
		"/big_messages",
		"/set_info_start",
		"/checkboxes_start",
	)

	levelFourState := statemachine.NewState(
		"level-four-state",

		"/level_four_start",

		"/level_four_start",
		"/level_four_one",
		"/level_four_two",
		"/level_four_three",
		"/level_four_four",
		"/back_state",
	)

	keyboardState := statemachine.NewState(
		"keyboard-state",

		"/keyboard_start",

		"/keyboard_start",
		"/keyboard_one",
		"/keyboard_two",
		"/keyboard_three",
		"/keyboard_finish",
		"/back_state",
		"/back_command",
	)

	infoState := statemachine.NewState(
		"info-state",

		"/set_info_start",

		"/set_info_start",
		"/set_name",
		"/set_surname",
		"/set_age",
		"/set_info_end",
		"/back_state",
		"/set_info_end",
		"*",
	)

	checkboxState := statemachine.NewState(
		"checkbox-state",

		"/checkboxes_start",

		"/checkboxes_start",
		"/checkboxes_first",
		"/checkboxes_second",
		"/checkboxes_third",
		"/checkboxes_fourth",
		"/checkboxes_accept",
		"/nothingness",
	)

	startState.SetAvailableStates(*levelFourState, *keyboardState, *infoState, *checkboxState, *startState)
	levelFourState.SetAvailableStates(*startState)
	keyboardState.SetAvailableStates(*startState)
	infoState.SetAvailableStates((*startState))
	checkboxState.SetAvailableStates((*startState))

	sm.AddStates(*startState, *levelFourState, *keyboardState, *infoState, *checkboxState)

	err := sm.SetStateByName("start-state")
	if err != nil {
		panic(err.Error())
	}
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

	mainLogger := logger.NewLogger(levelInt)
	mainLogger.Info("This is an info message")

	tgBotKey, exists := os.LookupEnv("TELEGRAM_BOT_API")
	if !exists {
		panic(".env does not contain TELEGRAM_BOT_API")
	}
	botAPI, err := tgbotapi.NewBotAPI(tgBotKey)
	if err != nil {
		panic(err.Error())
	}

	var sm statemachine.StateMachine
	configurateStateMachine(&sm)

	var gs globalstate.GlobalState

	commandHandler := handlers.NewCommandHandler(&sm, &gs)

	client := bot.NewClient(botAPI, commandHandler)
	client.ListenMessages()
}
