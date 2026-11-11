//go:build todo

package middleware

import "net/http"

func RequestID(next http.Handler) http.Handler {
	return next
}

func Logging(next http.Handler) http.Handler {
	return next
}

func Recovery(next http.Handler) http.Handler {
	return next
}

func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	return h
}
