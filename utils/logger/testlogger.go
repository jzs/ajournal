package logger

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// NewTestLogger returns a logger for test purposes
func NewTestLogger() Logger {
	return &testlogger{msgs: []string{}}
}

type testlogger struct {
	msgs      []string
	hasErrors bool
}

func (l *testlogger) Error(ctx context.Context, err error) {
	l.msgs = append(l.msgs, err.Error())
	l.hasErrors = true
}
func (l *testlogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.msgs = append(l.msgs, fmt.Sprintf(format, args...))
	l.hasErrors = true
}
func (l *testlogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.msgs = append(l.msgs, fmt.Sprintf(format, args...))
	l.hasErrors = true
}
func (l *testlogger) Print(ctx context.Context, err error) {
	if err == nil {
		l.msgs = append(l.msgs, "Trying to log a nil error")
		return
	}
	l.msgs = append(l.msgs, err.Error())
}
func (l *testlogger) Printf(ctx context.Context, format string, args ...interface{}) {
	l.msgs = append(l.msgs, fmt.Sprintf(format, args...))
}
func (l *testlogger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// TODO: Add logger!!!!!!!!!
	next(w, r)
}

func (l *testlogger) Flush() {
	if !l.hasErrors {
		return
	}
	for _, m := range l.msgs {
		log.Println(m)
	}
}
