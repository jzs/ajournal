package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sketchground/ajournal/journal"
	"github.com/sketchground/ajournal/postgres"
	"github.com/sketchground/ajournal/profile"
	"github.com/sketchground/ajournal/services"
	"github.com/sketchground/ajournal/user"
	"github.com/sketchground/ajournal/utils"
	"github.com/sketchground/ajournal/utils/logger"
	"github.com/urfave/negroni"
)

// Params are parameters for SetupRouter
type Params struct {
	BuildVersion string
	BuildTime    string
	BuildType    string

	StripeKey string
	WWWDir    string
}

// SetupRouter sets up a router with all repositories etc. initialized and returns a handler.
func SetupRouter(db *sqlx.DB, log logger.Logger, t *utils.Translator, params Params) http.Handler {
	baserouter := mux.NewRouter()
	apirouter := baserouter.PathPrefix("/api").Subrouter().StrictSlash(true)

	apirouter.Path("/version").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := map[string]string{
			"Version": params.BuildVersion,
			"Time":    params.BuildTime,
			"Type":    params.BuildType,
		}
		utils.JSONResp(r.Context(), log, r, w, version, nil)
	})

	jr := postgres.NewJournalRepo(db, log)
	js := journal.NewService(jr)
	journal.SetupHandler(apirouter, js, log)

	ur := postgres.NewUserRepo(db)
	us := user.NewService(t, ur)
	user.SetupHandler(apirouter, us, log)

	pr := postgres.NewProfileRepo(db, log)
	sr := services.NewStripeSubscriptionRepo(params.StripeKey, db)
	ps := profile.NewService(pr, sr)
	profile.SetupHandler(apirouter, ps, log)

	// Setup api router
	baserouter.PathPrefix("/api").Handler(negroni.New(negroni.Wrap(apirouter)))

	// Setup helper routes that redirects to public journal page
	baserouter.HandleFunc("/journal/{journalid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["journalid"]

		http.Redirect(w, r, fmt.Sprintf("/#/view/%v", id), http.StatusFound)
	})
	// Setup helper routes that redirects to public user page with his/her journals
	baserouter.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		http.Redirect(w, r, fmt.Sprintf("/#/viewuser/%v", id), http.StatusFound)
	})

	// Setup static file handler
	baserouter.PathPrefix("/").Handler(http.FileServer(http.Dir(params.WWWDir)))

	base := negroni.New(negroni.NewRecovery(), log, t)

	// Setup middleware that injects currently logged in user into the stack.
	base.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		nr := user.ApplyUserToRequest(r, us)
		next(w, nr)
	}))

	base.UseHandler(baserouter)
	return base
}
