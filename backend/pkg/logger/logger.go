package logger

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
	debug *log.Logger
}

var defaultLogger *Logger

func init() {
	defaultLogger = New()
}

func New() *Logger {
	return &Logger{
		info:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		debug: log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func Info(format string, v ...interface{}) {
	defaultLogger.info.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	defaultLogger.error.Printf(format, v...)
}

func Debug(format string, v ...interface{}) {
	defaultLogger.debug.Printf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	defaultLogger.error.Fatalf(format, v...)
}
