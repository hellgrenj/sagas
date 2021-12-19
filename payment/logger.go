package main

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
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
func (l *Log) Info(msg string) {
	l.info.Println(msg)
}
func (l *Log) Error(msg string) {
	l.error.Println(msg)
}
