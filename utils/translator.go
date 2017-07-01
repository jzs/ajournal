package utils

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/jzs/translate-i18-go"
	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/sketchground/ajournal/utils/logger"
	"golang.org/x/text/language"
)

func init() {
	translatorCtx = translatorContext("translatorcontext")
}

type translatorContext string

var translatorCtx translatorContext

var b bundle.Bundle

// T interface for translator instance
type T interface {
	With(args interface{}) T
	Zero() T
	Other() T
	Plural(count, many uint64) T
	String() string
}

// Translator is our translator
type Translator struct {
	tr    *translate.Translator
	langs language.Matcher
}

// NewTranslator returns a new translator with languages loaded from the given folder
func NewTranslator(folderpath string, log logger.Logger) (*Translator, error) {
	files, err := ioutil.ReadDir(folderpath)
	if err != nil {
		return nil, err
	}

	tags := []language.Tag{}

	t := &Translator{}
	langs := []*translate.Language{}
	for _, file := range files {
		lang := strings.TrimSuffix(file.Name(), ".yaml")
		r, err := os.Open(path.Join(folderpath, file.Name()))
		if err != nil {
			return nil, err
		}
		tag := language.MustParse(lang)
		l, err := translate.LoadYaml(r, tag.String())
		if err != nil {
			return nil, err
		}
		langs = append(langs, l)
		tags = append(tags, tag)
	}
	t.tr = translate.New(langs...)
	t.tr.SetLog(func(s string, args ...interface{}) {
		log.Errorf(context.Background(), s, args)
	})
	t.langs = language.NewMatcher(tags)

	return t, nil
}

// T returns a func for translating a string to the language based on what is set in the context.
// Defaults to en-us
func (t *Translator) T(ctx context.Context) func(string) translate.T {
	val := ctx.Value(translatorCtx)
	if langs, ok := val.([]string); ok {
		return t.tr.Tfunc(langs...)
	}
	return t.tr.Tfunc("en-us")
}

// ServeHTTP Method for supporting injection into Negroni
func (t *Translator) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	acceptLang := r.Header.Get("Accept-Language")
	preferred, _, err := language.ParseAcceptLanguage(acceptLang)
	if err != nil {
		preferred = []language.Tag{language.English}
	}
	lang, _, _ := t.langs.Match(preferred...)

	langs := []string{lang.String()}
	ctx := context.WithValue(r.Context(), translatorCtx, langs)
	nr := r.WithContext(ctx)
	next(w, nr)
}

// NewTestTranslator returns a new translator with languages loaded from the given folder
func NewTestTranslator() *Translator {
	t := &Translator{
		langs: language.NewMatcher([]language.Tag{language.English}),
	}

	d := bytes.NewReader([]byte(""))

	l, _ := translate.LoadYaml(d, "en-us")
	t.tr = translate.New(l)
	return t
}
