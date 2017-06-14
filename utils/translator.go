package utils

import (
	"context"
	"log"
	"net/http"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
)

type translatorContext string

var translatorCtx translatorContext

var b bundle.Bundle

// T is a function type for translating a message
type T func(translationID string, args ...interface{}) string

// Translator is our translator
type Translator struct {
	b *bundle.Bundle
}

// NewTranslator returns a new translator
func NewTranslator() *Translator {
	return &Translator{
		b: bundle.New(),
	}
}

// AddTranslationFromFile loads a translation from file and adds it to the translator
func (t *Translator) AddTranslationFromFile(path string) error {
	return t.b.LoadTranslationFile(path)
}

// Middleware returns a func for injecting into negroni middleware stack making the translator
// available in the request context
func (t *Translator) Middleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		cookieLang, _ := r.Cookie("lang")
		acceptLang := r.Header.Get("Accept-Language")
		defaultLang := "en-us" // known valid language
		cl := acceptLang
		if cookieLang != nil {
			cl = cookieLang.Value
		}
		T, err := t.b.Tfunc(cl, acceptLang, defaultLang)
		if err != nil {
			log.Println(err)
		}

		ctx := context.WithValue(r.Context(), translatorCtx, T)
		nr := r.WithContext(ctx)
		next(w, nr)
	}
}

func init() {
	translatorCtx = translatorContext("translatorcontext")
}

// TranslatorFromCtx returns a translator from the given context
func TranslatorFromCtx(ctx context.Context) func(string, ...interface{}) string {
	val := ctx.Value(translatorCtx)
	if t, ok := val.(func(string, ...interface{}) string); ok {
		return t
	}
	return nil
}

// TestContextWithTranslator returns a context with a translator set
func TestContextWithTranslator(c context.Context) context.Context {
	ctx := context.WithValue(c, translatorCtx, func(s string, args ...interface{}) string {
		return "Mock translator"
	})
	return ctx
}
