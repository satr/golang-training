//go:build todo

package structuredlogging

import "io"

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type Logger struct {
	w      io.Writer
	level  Level
	fields map[string]any
}

func NewLogger(w io.Writer, level Level) *Logger {
	// TODO: Initialize logger with target writer, level, and empty fields map.
	return &Logger{w: w, level: level}
}

func (l *Logger) With(key string, value any) *Logger {
	// TODO: Return a new Logger that inherits parent fields plus the new field.
	return l
}

func (l *Logger) Info(msg string, extra ...string) {
	// TODO: Write if level permits. Format: [INFO] key=val msg
}

func (l *Logger) Error(msg string, err error) {
	// TODO: Write if level permits. Include error value in output.
}
