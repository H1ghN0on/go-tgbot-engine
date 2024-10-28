package logger

import (
	"fmt"
	"time"
)

type LogLevel int

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
	Bot            *Logger
}

var GlobalLogger = &Logger{}

func NewLogger(levelInt int) *Logger {
	var level LogLevel

	switch levelInt {
	case 0:
		level = Info
	case 1:
		level = Warning
	case 2:
		level = Critical
	default:
		level = Info
		fmt.Println("Invalid value for LogLevel. Choosed default LogLevel `INFO`")
	}

	rootLogger := &Logger{
		level:    level,
		category: "Root",
	}

	rootLogger.StateMachine = &Logger{
		level:    level,
		category: "StateMachine",
	}
	rootLogger.CommandHandler = &Logger{
		level:    level,
		category: "CommandHandler",
	}
    rootLogger.Bot = &Logger{
		level:    level,
		category: "Bot",
	}

	return rootLogger
}

func (l *Logger) logMessage(level LogLevel, message string) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

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
	default:
		fmt.Println("Invalid value for LogLevel. Choosed default LogLevel `INFO`")
		levelStr = "[INFO]"
		color = "\033[34m" // Blue color
	}

	resetColor := "\033[0m" // Reset color

	logLine := fmt.Sprintf("%s %s%s [%s] %s%s", timestamp, color, levelStr, l.category, message, resetColor)
	fmt.Println(logLine)

}

func (l *Logger) Info(message string)     { l.logMessage(Info, message) }
func (l *Logger) Warning(message string)  { l.logMessage(Warning, message) }
func (l *Logger) Critical(message string) { l.logMessage(Critical, message) }
