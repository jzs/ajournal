package logger

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	raven "github.com/getsentry/raven-go"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/negroni"
)

type loggertype int

const (
	loggercontext = iota
)

type Logger interface {
	Error(ctx context.Context, err error)
	Errorf(ctx context.Context, format string, args ...interface{})
	Printf(ctx context.Context, format string, args ...interface{})
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type logger struct {
	isDevel bool
}

// New creates a new logger
// isDevel determines whether errors should be logged to sentry.io or not
func New(isDevel bool, dsn string) Logger {
	if !isDevel {
		raven.SetDSN(dsn) // Set DSN up for sentry.io (To log crashes!)
	}

	return &logger{
		isDevel: isDevel,
	}
}

func (l *logger) Error(ctx context.Context, err error) {
	uid, ok := ctx.Value(loggercontext).(string)
	if !ok {
		uid = ""
	}
	log.Printf("[ERROR] | %v | %v", uid, err.Error())
	debug.PrintStack()
	if !l.isDevel {
		raven.CaptureError(err, nil)
	}
}

func (l *logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	uid, ok := ctx.Value(loggercontext).(string)
	if !ok {
		uid = ""
	}
	str := fmt.Sprintf(format, args...)
	log.Printf("[ajournal] | [ERROR] | %v | %v", uid, str)
	debug.PrintStack()
}

func (l *logger) Printf(ctx context.Context, format string, args ...interface{}) {
	uid, ok := ctx.Value(loggercontext).(string)
	if !ok {
		uid = ""
	}
	str := fmt.Sprintf(format, args...)
	log.Printf("[ajournal] | [INFO] | %v | %v", uid, str)
}

// ServeHTTP Method for supporting injection into Negroni
func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	// TODO Consider logging other things like the users ip.
	uid := uuid.NewV4()

	log.Printf("[ajournal] | [INFO] | %v | %v | %v %v \n", uid.String(), r.Host, r.Method, r.URL.Path)

	ctx := context.WithValue(r.Context(), loggercontext, uid.String())
	nr := r.WithContext(ctx)
	next(w, nr)

	res := w.(negroni.ResponseWriter)

	log.Printf("[ajournal] | [INFO] | %v | %v | %v \t | %v | %v %v \n", uid.String(), res.Status(), time.Since(start), r.Host, r.Method, r.URL.Path)
}
