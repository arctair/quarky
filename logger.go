package main

import "log"

// Logger ...
type Logger interface {
	error(err error)
}

// LoggerConsole ...
type LoggerConsole struct{}

func (l *LoggerConsole) error(err error) {
	log.Print(err)
}
