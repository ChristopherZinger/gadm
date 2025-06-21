package logger

import (
	"log"
	"os"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] "+msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	log.Printf("[WARNING] "+msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	log.Printf("[FATAL] "+msg, args...)
	os.Exit(1)
}

var logger = NewLogger()

func Debug(msg string, args ...interface{}) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

func Warning(msg string, args ...interface{}) {
	logger.Warning(msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatal(msg, args...)
}
