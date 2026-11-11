package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			b := make([]byte, 16)
			rand.Read(b)
			id = hex.EncodeToString(b)
		}
		r.Header.Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logged := &recorder{w: w, code: http.StatusOK}
		next.ServeHTTP(logged, r)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, logged.code, time.Since(start))
	})
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				log.Printf("panic recovered")
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type recorder struct {
	w    http.ResponseWriter
	code int
}

func (rr *recorder) Header() http.Header         { return rr.w.Header() }
func (rr *recorder) WriteHeader(code int)        { rr.code = code; rr.w.WriteHeader(code) }
func (rr *recorder) Write(p []byte) (int, error) { return rr.w.Write(p) }
