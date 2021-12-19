package main

import (
	"log"
	"os"
)

type Logger interface {
	Info() *log.Logger
	Error() *log.Logger
}
type Log struct {
	info  *log.Logger
	error *log.Logger
}

func NewLogger() *Log {
	i := log.New(os.Stdout, "INFO: ", 3)
	e := log.New(os.Stderr, "ERROR: ", 3)
	return &Log{info: i, error: e}
}
