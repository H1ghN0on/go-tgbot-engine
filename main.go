package main

import (
	"log"
	"os"
	"strconv"

	"github.com/H1ghN0on/go-tgbot-engine/bot"
	"github.com/H1ghN0on/go-tgbot-engine/globalstate"
	"github.com/H1ghN0on/go-tgbot-engine/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/joho/godotenv"
)

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

	var gs globalstate.GlobalState

	bot := bot.NewBot(botAPI, &gs)
	bot.ListenMessages()
}
