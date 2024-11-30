package logger

func CommandHandler() *Logger {
	logger := GlobalLogger
	logger.category = "CommandHandler"
	return logger
}

func Bot() *Logger {
	logger := GlobalLogger
	logger.category = "Bot"
	return logger
}

func Client() *Logger {
	logger := GlobalLogger
	logger.category = "Client"
	return logger
}

func StateMachine() *Logger {
	logger := GlobalLogger
	logger.category = "StateMachine"
	return logger
}

func Notificator() *Logger {
	logger := GlobalLogger
	logger.category = "Notificator"
	return logger
}