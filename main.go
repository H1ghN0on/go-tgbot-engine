package main

import (
	"log"
	"os"
	"strconv"

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

	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	loggerStatus, exists := os.LookupEnv("LOGGER_LEVEL")
	if !exists {
		log.Println(".env does not contain LOGGER_LEVEL")
	}
	levelInt, err := strconv.Atoi(loggerStatus)
	if err != nil {
		log.Println("Error:", err, "\nChoosed default LogLevel `INFO`")
		levelInt = 0
	}

	mainLogger := logger.NewLogger(levelInt)

	mainLogger.Info("This is an info message")
	mainLogger.Warning("This is a warning message")

	mainLogger.StateMachine.Warning("StateMachine warning")
	mainLogger.CommandHandler.Critical("CommandHandler error")

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

	commandHandler := handlers.NewCommandHandler(&sm)

	client := bot.NewClient(botAPI, commandHandler)
	client.ListenMessages()
}
