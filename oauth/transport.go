package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gorilla/mux"
	"github.com/sketchground/ajournal/profile"
	"github.com/sketchground/ajournal/user"
	"github.com/sketchground/ajournal/utils/logger"
)

// SetupHandler sets up routes for the journal service
func SetupHandler(router *mux.Router, os Service, l logger.Logger, creds Credentials) {
	if creds.Provider != ProviderGoogle {
		panic("Unsupported oauth provider")
	}

	googleCfg := &oauth2.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		RedirectURL:  creds.RedirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	}

	router.Path("/oauth/google").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := googleCfg.AuthCodeURL("random string", oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	router.Path("/oauth/google/callback").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		state := r.FormValue("state")
		if state != "random string" {
			RenderInternalError(ctx, w, l, errors.New("Invalid state"))
			return
		}

		code := r.FormValue("code")
		token, err := googleCfg.Exchange(oauth2.NoContext, code)
		if err != nil {
			RenderInternalError(ctx, w, l, err)
			return
		}

		client := googleCfg.Client(ctx, token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			RenderInternalError(ctx, w, l, err)
			return
		}
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)

		var uinfo GoogleUserInfo
		err = json.Unmarshal(data, &uinfo)
		if err != nil {
			RenderInternalError(ctx, w, l, err)
			return
		}

		// Check if we already have a user. If so, then log in, otherwise register.
		tok, err := os.Login(ctx, uinfo.Email, ProviderGoogle)
		if err != nil {
			err = os.Register(ctx, &user.User{Username: uinfo.Email}, &profile.Profile{Name: uinfo.Name, Email: uinfo.Email})
			if err != nil {
				RenderInternalError(ctx, w, l, err)
				return
			}
			tok, err = os.Login(ctx, uinfo.Email, ProviderGoogle)
			if err != nil {
				RenderInternalError(ctx, w, l, err)
				return
			}
		}

		http.SetCookie(w, user.CreateCookie(tok))
		// TODO: Download profile picture (uinfo.Picture)
		// TODO: Redirect to a sane url!
		http.Redirect(w, r, "/app", http.StatusFound)
	})
}

// RenderInternalError renders an oops page.
func RenderInternalError(ctx context.Context, w http.ResponseWriter, l logger.Logger, err error) {
	// TODO: Render error page. And make this func a utility function.
	l.Error(ctx, err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
