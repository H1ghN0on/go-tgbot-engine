package main

import (
	"fmt"
	"os"

	"github.com/H1ghN0on/go-tgbot-engine/bot"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	"github.com/H1ghN0on/go-tgbot-engine/statemachine"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	var sm statemachine.StateMachine

	startState := statemachine.State{
		Name: "start-state",
		AvailableCommands: []string{
			"/level_one",
			"/level_two",
			"/level_three",
			"/show_commands",
			"/keyboard_start",
			"/create_error",
			"/level_four_start",
			"/big_messages",
		},
	}

	levelFourState := statemachine.State{
		Name: "level-four-state",
		AvailableCommands: []string{
			"/level_four_one",
			"/level_four_two",
			"/level_four_three",
			"/level_four_four",
		},
	}

	keyboardState := statemachine.State{
		Name: "keyboard-state",
		AvailableCommands: []string{
			"/keyboard_one",
			"/keyboard_two",
			"/keyboard_three",
		},
	}

	startState.AvailableStates = append(startState.AvailableStates, levelFourState, keyboardState)
	levelFourState.AvailableStates = append(levelFourState.AvailableStates, startState)
	keyboardState.AvailableStates = append(keyboardState.AvailableStates, startState)

	sm.AddStates(startState, levelFourState, keyboardState)

	err := sm.SetStateByName("start-state")
	if err != nil {
		panic(err.Error())
	}

	tgBotKey, exists := os.LookupEnv("TELEGRAM_BOT_API")
	if !exists {
		panic(".env does not contain TELEGRAM_BOT_API")
	}
	botAPI, err := tgbotapi.NewBotAPI(tgBotKey)

	if err != nil {
		fmt.Println("pizdec")
	}

	client := bot.Client{C: botAPI, Sm: &sm}
	bot := &bot.Bot{
		Cmdhandler: handlers.NewCommandHandler(&sm),
		Client:     &client,
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := client.C.GetUpdatesChan(u)

	for update := range updates {
		bot.ListenMessages(update)
	}
}

// Реализация кнопки назад через запоминание всех постов, что привязано к командному хэндлеру
