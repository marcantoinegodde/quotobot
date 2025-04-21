package main

import (
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

func newLogger() *Logger {
	return &Logger{
		Info:  infoLog(),
		Error: errorLog(),
	}
}

func infoLog() *log.Logger {
	return log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func errorLog() *log.Logger {
	return log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
