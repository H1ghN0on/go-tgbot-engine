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

type LoggerSettings struct {
	Level LogLevel
}

type Logger struct {
	category string
}

var globalLoggerSettings LoggerSettings

func InitGlobalLoggerSettings(settings LoggerSettings) {

	switch settings.Level {
	case 0:
		globalLoggerSettings.Level = Info
	case 1:
		globalLoggerSettings.Level = Warning
	case 2:
		globalLoggerSettings.Level = Critical
	default:
		globalLoggerSettings.Level = Info
		fmt.Println("Invalid value for LogLevel. Choosed default LogLevel `INFO`")
	}
}

func (l *Logger) logMessage(level LogLevel, messages ...string) {
	if level < globalLoggerSettings.Level {
		return
	}

	if len(messages) == 0 {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var levelStr, color string
	switch level {
	case Info:
		levelStr = "INFO"
		color = "\033[34m" // Blue color
	case Warning:
		levelStr = "WARNING"
		color = "\033[33m" // Yellow color
	case Critical:
		levelStr = "CRITICAL"
		color = "\033[31m" // Red color
	default:
		fmt.Println("Invalid value for LogLevel. Choosed default LogLevel `INFO`")
		levelStr = "INFO"
		color = "\033[34m" // Blue color
	}

	resetColor := "\033[0m" // Reset color

	var fullMessage string
	for _, message := range messages {
		fullMessage += message + " "
	}
	fullMessage = fullMessage[:len(fullMessage)-1]

	logLine := fmt.Sprintf("%s %s[%s] [%s] %s%s", timestamp, color, levelStr, l.category, fullMessage, resetColor)
	fmt.Println(logLine)
}

func (l *Logger) Info(messages ...string)     { l.logMessage(Info, messages...) }
func (l *Logger) Warning(messages ...string)  { l.logMessage(Warning, messages...) }
func (l *Logger) Critical(messages ...string) { l.logMessage(Critical, messages...) }

var GlobalLogger = &Logger{}
