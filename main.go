package main

import (
	"os"

	"github.com/H1ghN0on/go-tgbot-engine/bot"
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
	)

	levelFourState := statemachine.NewState(
		"level-four-state",

		"/level_four_one",

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
		"/back_state",
		"/back_command",
	)

	startState.SetAvailableStates(*levelFourState, *keyboardState)
	levelFourState.SetAvailableStates(*startState)
	keyboardState.SetAvailableStates(*startState)

	sm.AddStates(*startState, *levelFourState, *keyboardState)

	err := sm.SetStateByName("start-state")
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	// tgBotKey, exists := os.LookupEnv("TELEGRAM_BOT_API")
	// if !exists {
	// 	panic(".env does not contain TELEGRAM_BOT_API")
	// }
	mainLogger := logger.NewLogger(logger.Critical)

	mainLogger.Info("This is an info message")
	mainLogger.Warning("This is a warning message")

	mainLogger.StateMachine.Warning("StateMachine warning")
	mainLogger.CommandHandler.Critical("CommandHandler error")

	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	tgBotKey, exists := os.LookupEnv("TELEGRAM_BOT_API")
	if !exists {
		panic(".env does not contain TELEGRAM_BOT_API")
	}
	botAPI, err := tgbotapi.NewBotAPI(tgBotKey)
	if err != nil {
		// panic(err.Error())
		mainLogger.Warning(err.Error())

	}

	var sm statemachine.StateMachine
	configurateStateMachine(&sm)

	commandHandler := handlers.NewCommandHandler(&sm)

	client := bot.NewClient(botAPI, commandHandler)
	client.ListenMessages()
}
