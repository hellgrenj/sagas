package main

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}
type logger struct {
	info   *log.Logger
	error  *log.Logger
	caller string
}

func NewLogger(caller string) *logger {
	i := log.New(os.Stdout, "INFO: ", 3)
	e := log.New(os.Stderr, "ERROR: ", 3)
	return &logger{info: i, error: e, caller: caller}
}
func (l *logger) Info(msg string) {
	msg = fmt.Sprintf("%s %s", l.caller, msg)
	l.info.Println(msg)
}
func (l *logger) Error(msg string) {
	msg = fmt.Sprintf("%s %s", l.caller, msg)
	l.error.Println(msg)
}
