package utils

import (
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, prefix+" ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.logger.Printf("[INFO] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.logger.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.logger.Printf("[DEBUG] "+format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.logger.Printf("[WARN] "+format, v...)
}

