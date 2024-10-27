package logger

import (
	"fmt"
	"time"
)

type LogLevel int

// Logging levels
const (
	Info LogLevel = iota
	Warning
	Critical
)

type Logger struct {
	level          LogLevel
	category       string
	StateMachine   *Logger
	CommandHandler *Logger
}

// Creating a new main logger with subcategories
func NewLogger(level LogLevel) *Logger {
	rootLogger := &Logger{
		level:    level,
		category: "Root",
	}

	// init subcategories
	rootLogger.StateMachine = &Logger{
		level:    level,
		category: "StateMachine",
	}
	rootLogger.CommandHandler = &Logger{
		level:    level,
		category: "CommandHandler",
	}

	return rootLogger
}

// basic method for displaying log messages
func (l *Logger) logMessage(level LogLevel, message string) {
	if level < l.level {
		// Skip the message if its level is lower than the current logger level
		return
	}

	// log time
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Prefix logging level
	var levelStr, color string
	switch level {
	case Info:
		levelStr = "[INFO]"
		color = "\033[34m" // Blue color
	case Warning:
		levelStr = "[WARNING]"
		color = "\033[33m" // Yellow color
	case Critical:
		levelStr = "[CRITICAL]"
		color = "\033[31m" // Red color
	}

	resetColor := "\033[0m"

	// Forming and display a message
	logLine := fmt.Sprintf("%s %s%s [%s] %s%s", timestamp, color, levelStr, l.category, message, resetColor)
	fmt.Println(logLine)

}

func (l *Logger) Info(message string)     { l.logMessage(Info, message) }
func (l *Logger) Warning(message string)  { l.logMessage(Warning, message) }
func (l *Logger) Critical(message string) { l.logMessage(Critical, message) }
