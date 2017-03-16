package logger

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/urfave/negroni"
)

type loggertype int

const (
	loggercontext = iota
)

func Error(ctx context.Context, err error) {
	// TODO: Log stack trace.
	uid := ctx.Value(loggercontext).(string)
	log.Printf("[ERROR] \t | %v | %v", uid, err.Error())
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	// TODO: Log stack trace.
	uid := ctx.Value(loggercontext).(string)
	str := fmt.Sprintf(format, args...)
	log.Printf("[ERROR] \t | %v | %v", uid, str)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	uid := ctx.Value(loggercontext).(string)
	str := fmt.Sprintf(format, args...)
	log.Printf("[INFO] \t | %v | %v", uid, str)
}

func NewLogger() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()

		// TODO Consider logging other things like the users ip.
		uid := uuid.NewV4()

		log.Printf("[ajournal] | %v | %v | %v %v \n", uid.String(), r.Host, r.Method, r.URL.Path)

		ctx := context.WithValue(r.Context(), loggercontext, uid.String())
		nr := r.WithContext(ctx)
		next(w, nr)

		res := w.(negroni.ResponseWriter)

		log.Printf("[ajournal] | %v | %v | %v \t | %v | %v %v \n", uid.String(), res.Status(), time.Since(start), r.Host, r.Method, r.URL.Path)

	})
}
