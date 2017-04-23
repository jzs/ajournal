package logger

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	raven "github.com/getsentry/raven-go"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/negroni"
)

type loggertype string

var loggercontext loggertype

func init() {
	loggercontext = "loggertype"
}

// Logger interface describes functions available on a logger
type Logger interface {
	Error(ctx context.Context, err error)
	Errorf(ctx context.Context, format string, args ...interface{})
	Print(ctx context.Context, err error)
	Printf(ctx context.Context, format string, args ...interface{})
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type logger struct {
	isDevel bool
	log     *log.Logger
}

// New creates a new logger
// isDevel determines whether errors should be logged to sentry.io or not
func New(isDevel bool, dsn string) Logger {
	if !isDevel {
		raven.SetDSN(dsn) // Set DSN up for sentry.io (To log crashes!)
	}

	stdlogger := log.New(os.Stderr, "", 0)

	return &logger{
		isDevel: isDevel,
		log:     stdlogger,
	}
}

func (l *logger) Error(ctx context.Context, err error) {
	uid, ok := ctx.Value(loggercontext).(string)
	if !ok {
		uid = ""
	}
	l.log.Printf("[ERROR] | %v | %+v", uid, err)
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
	l.log.Printf("[ERROR] | %v | %+v", uid, str)
}

func (l *logger) Print(ctx context.Context, err error) {
	uid, ok := ctx.Value(loggercontext).(string)
	if !ok {
		uid = ""
	}
	l.log.Printf("[INFO] | %v | %+v", uid, err)
}

func (l *logger) Printf(ctx context.Context, format string, args ...interface{}) {
	uid, ok := ctx.Value(loggercontext).(string)
	if !ok {
		uid = ""
	}
	str := fmt.Sprintf(format, args...)
	l.log.Printf("[INFO] | %v | %v", uid, str)
}

// ServeHTTP Method for supporting injection into Negroni
func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	// TODO Consider logging other things like the users ip.
	uid := uuid.NewV4()

	l.log.Printf("[INFO] | %v | %v | %v %v \n", uid.String(), r.Host, r.Method, r.URL.Path)

	ctx := context.WithValue(r.Context(), loggercontext, uid.String())
	nr := r.WithContext(ctx)
	next(w, nr)

	res := w.(negroni.ResponseWriter)

	l.log.Printf("[INFO] | %v | %v | %v \t\n", uid.String(), res.Status(), time.Since(start))
}
