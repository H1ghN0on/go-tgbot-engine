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

func StateMachine() *Logger {
	logger := GlobalLogger
	logger.category = "StateMachine"
	return logger
}
