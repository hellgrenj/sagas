package main

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
}

func NewLogger() *Logger {
	i := log.New(os.Stdout, "INFO: ", 3)
	e := log.New(os.Stderr, "ERROR: ", 3)
	return &Logger{info: i, error: e}
}
