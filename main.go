package main

import (
	"log"
	"os"
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot"
	"github.com/H1ghN0on/go-tgbot-engine/globalstate"
	"github.com/H1ghN0on/go-tgbot-engine/handlers"
	cmd "github.com/H1ghN0on/go-tgbot-engine/handlers/commands"
	"github.com/H1ghN0on/go-tgbot-engine/logger"
	"github.com/H1ghN0on/go-tgbot-engine/statemachine"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/joho/godotenv"
)

func configurateStateMachine(sm *statemachine.StateMachine) {
	startState := statemachine.NewState(
		"start-state",

		cmd.ShowCommandsCommand,

		cmd.LevelOneCommand,
		cmd.LevelTwoCommand,
		cmd.LevelThreeCommand,
		cmd.ShowCommandsCommand,
		cmd.KeyboardStartCommand,
		cmd.LevelFourStartCommand,
		cmd.BigMessagesCommand,
		cmd.SetInfoStartCommand,
		cmd.CheckboxStartCommand,
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
		cmd.AnyCommand,
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

	var sm statemachine.StateMachine
	configurateStateMachine(&sm)

	var gs globalstate.GlobalState

	commandHandler := handlers.NewCommandHandler(&sm, &gs)

	client := bot.NewClient(botAPI, commandHandler)
	client.ListenMessages()
}
