package main

import (
	"fmt"
	"os"

	"github.com/H1ghN0on/go-tgbot-engine/bot"
	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
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

	startState := bottypes.State{
		Name: "start-state",
		AvailableCommands: []bottypes.Command{
			{Text: "/level_one"},
			{Text: "/level_two"},
			{Text: "/level_three"},
			{Text: "/show_commands"},
			{Text: "/keyboard_start"},
			{Text: "/create_error"},
			{Text: "/level_four_start"},
			{Text: "/big_messages"},
		},
	}

	levelFourState := bottypes.State{
		Name: "level-four-state",
		AvailableCommands: []bottypes.Command{
			{Text: "/level_four_one"},
			{Text: "/level_four_two"},
			{Text: "/level_four_three"},
			{Text: "/level_four_four"},
		},
	}

	keyboardState := bottypes.State{
		Name: "keyboard-state",
		AvailableCommands: []bottypes.Command{
			{Text: "/keyboard_one"},
			{Text: "/keyboard_two"},
			{Text: "/keyboard_three"},
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
