package structuredlogging

import (
	"fmt"
	"io"
	"strings"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	w      io.Writer
	level  Level
	fields map[string]any
}

func NewLogger(w io.Writer, level Level) *Logger {
	return &Logger{w: w, level: level, fields: map[string]any{}}
}

func (l *Logger) With(key string, value any) *Logger {
	fs := make(map[string]any, len(l.fields)+1)
	for k, v := range l.fields {
		fs[k] = v
	}
	fs[key] = value
	return &Logger{w: l.w, level: l.level, fields: fs}
}

func (l *Logger) Info(msg string, extra ...string) {
	if l.level > Info {
		return
	}
	l.write(Info, msg)
}

func (l *Logger) Error(msg string, err error) {
	if l.level > Error {
		return
	}
	if err != nil {
		l.fields = merge(l.fields, map[string]any{"error": err.Error()})
	}
	l.write(Error, msg)
}

func (l *Logger) write(level Level, msg string) {
	var parts []string
	parts = append(parts, fmt.Sprintf("[%s]", level))
	for k, v := range l.fields {
		parts = append(parts, fmt.Sprintf("%v=%v", k, v))
	}
	parts = append(parts, msg)
	fmt.Fprintln(l.w, strings.Join(parts, " "))
}

func merge(a, b map[string]any) map[string]any {
	m := make(map[string]any, len(a)+len(b))
	for k, v := range a {
		m[k] = v
	}
	for k, v := range b {
		m[k] = v
	}
	return m
}
