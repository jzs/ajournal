package logger

import (
	"context"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
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
	Fatalf(ctx context.Context, format string, args ...interface{})
	Print(ctx context.Context, err error)
	Printf(ctx context.Context, format string, args ...interface{})
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type logger struct {
	isDevel bool
	log     *logrus.Logger
}

// New creates a new logger
// isDevel determines whether errors should be logged to sentry.io or not
func New(isDevel bool) Logger {
	stdlogger := logrus.New()
	if !isDevel {
		stdlogger.Formatter = &logrus.JSONFormatter{}
	}
	// TODO Install hook that sends everything to logstash maybe?

	return &logger{
		isDevel: isDevel,
		log:     stdlogger,
	}
}

func (l *logger) Error(ctx context.Context, err error) {
	l.getLogger(ctx).WithFields(logrus.Fields{"error": fmt.Sprintf("%+v", err)}).Error(err)
}

func (l *logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.getLogger(ctx).Errorf(format, args...)
}

func (l *logger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.getLogger(ctx).Fatalf(format, args...)
}

func (l *logger) Print(ctx context.Context, err error) {
	l.getLogger(ctx).Info(err)
}

func (l *logger) Printf(ctx context.Context, format string, args ...interface{}) {
	l.getLogger(ctx).Infof(format, args...)
}

// ServeHTTP Method for supporting injection into Negroni
func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	// TODO Consider logging other things like the users ip.
	uid := uuid.NewV4()

	entry := l.log.WithFields(logrus.Fields{
		"requestid": uid.String(),
		"host":      r.Host,
		"method":    r.Method,
		"path":      r.URL.Path,
	})
	entry.Info()

	ctx := context.WithValue(r.Context(), loggercontext, entry)
	nr := r.WithContext(ctx)
	next(w, nr)

	res := w.(negroni.ResponseWriter)

	entry.WithFields(logrus.Fields{
		"status":   res.Status(),
		"duration": time.Since(start),
	}).Info()
}

func (l *logger) getLogger(ctx context.Context) logrus.FieldLogger {
	log, ok := ctx.Value(loggercontext).(*logrus.Entry)
	if !ok {
		return l.log
	}
	return log
}
