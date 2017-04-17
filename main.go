package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"bitbucket.org/sketchground/ajournal/journal"
	"bitbucket.org/sketchground/ajournal/postgres"
	"bitbucket.org/sketchground/ajournal/profile"
	"bitbucket.org/sketchground/ajournal/services"
	"bitbucket.org/sketchground/ajournal/user"
	"bitbucket.org/sketchground/ajournal/utils/logger"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/urfave/negroni"
)

const (
	// BuildVersionDevel for devel setups
	BuildVersionDevel = "DEVEL"
	// BuildVersionStaging for staging setups
	BuildVersionStaging = "STAGING"
	// BuildVersionProd for prod setups
	BuildVersionProd = "PROD"
)

var (
	// BuildVersion is overwritten from build script when deploying
	BuildVersion = "next"
	// BuildType which kind of build we're at
	BuildType = BuildVersionDevel
	// BuildTime is overwritten from build script when deploying
	BuildTime = "2015-08-19"
)

func main() {

	stripeKey := os.Getenv("AJ_STRIPE_SK")
	if stripeKey == "" {
		log.Fatalf("Environment variable AJ_STRIPE_SK not set!\nRemember to set your stripe private key")
		return
	}

	dbuser := os.Getenv("AJ_DB_USER")
	if dbuser == "" {
		dbuser = "jzs"
	}
	dbname := os.Getenv("AJ_DB_NAME")
	if dbname == "" {
		dbname = "journal"
	}

	port := os.Getenv("AJ_PORT")
	if port == "" {
		port = ":8080"
	}

	// Set up sentry.io to log crashes!
	dsn := os.Getenv("AJ_SENTRY_DSN")
	if dsn == "" && BuildType != BuildVersionDevel {
		log.Fatalf("Environment variable AJ_SENTRY_DSN not set!\nRemember to set you sentry dsn")
		return
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%v dbname=%v sslmode=disable", dbuser, dbname))
	if err != nil {
		log.Fatalf("Could not connect to database! %v", err)
		return
	}

	baserouter := mux.NewRouter()
	apirouter := baserouter.PathPrefix("/api").Subrouter().StrictSlash(true)

	alogger := logger.New(BuildType == BuildVersionDevel, dsn)

	jr := postgres.NewJournalRepo(db, alogger)
	js := journal.NewService(jr)
	journal.SetupHandler(apirouter, js, alogger)

	ur := postgres.NewUserRepo(db)
	us := user.NewService(ur)
	user.SetupHandler(apirouter, us, alogger)

	pr := postgres.NewProfileRepo(db, alogger)
	sr := services.NewStripeSubscriptionRepo(stripeKey, db)
	ps := profile.NewService(pr, sr)
	profile.SetupHandler(apirouter, ps, alogger)

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
	baserouter.PathPrefix("/").Handler(http.FileServer(http.Dir("www")))

	base := negroni.New(negroni.NewRecovery(), alogger)

	// Setup middleware that injects currently logged in user into the stack.
	base.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		nr := user.ApplyUserToRequest(r, us)
		next(w, nr)
	}))

	base.UseHandler(baserouter)

	alogger.Printf(context.Background(), "Listening on: %v", port)
	server := &http.Server{Addr: port, Handler: base}

	// subscribe to SIGINT signals
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			alogger.Error(context.Background(), err)
		}
	}()

	<-sigchan // wait for SIGINT
	log.Println("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)

	log.Println("Server gracefully stopped")
}
