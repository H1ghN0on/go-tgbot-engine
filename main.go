package main

import (
	"os"

	"github.com/H1ghN0on/go-tgbot-engine/bot"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	"github.com/H1ghN0on/go-tgbot-engine/statemachine"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/joho/godotenv"
)

func configurateStateMachine(sm *statemachine.StateMachine) {
	startState := statemachine.NewState(
		"start-state",

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
		"/level_four_two",
		"/level_four_three",
		"/level_four_four",
	)

	keyboardState := statemachine.NewState(
		"keyboard-state",

		"/keyboard_one",
		"/keyboard_two",
		"/keyboard_three",
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

// Реализация кнопки назад через запоминание всех постов, что привязано к командному хэндлеру
