package main

import (
	"log"
	"net/http"

	"bitbucket.org/sketchground/ajournal/journal"
	"bitbucket.org/sketchground/ajournal/postgres"
	"bitbucket.org/sketchground/ajournal/user"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/urfave/negroni"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=jzs dbname=journal sslmode=disable")
	if err != nil {
		log.Println("Could not connect to database!")
		log.Fatalln(err)
	}

	baserouter := mux.NewRouter()
	apirouter := baserouter.PathPrefix("/api").Subrouter().StrictSlash(true)

	jr := postgres.NewJournalRepo(db)
	js := journal.NewService(jr)
	journal.SetupHandler(apirouter, js)

	ur := postgres.NewUserRepo(db)
	us := user.NewService(ur)
	user.SetupHandler(apirouter, us)

	// Setup api router
	baserouter.PathPrefix("/api").Handler(negroni.New(negroni.Wrap(apirouter)))
	// Setup static file handler
	baserouter.PathPrefix("/").Handler(http.FileServer(http.Dir("www")))

	base := negroni.Classic()
	// Setup middleware that injects currently logged in user into the stack.
	base.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		nr := user.ApplyUserToRequest(r, us)
		next(w, nr)
	}))

	base.UseHandler(baserouter)
	// TODO Hook up base with a middleware that sets current user in context.

	err = http.ListenAndServe(":8080", base)
	if err != nil {
		panic(err)
	}
}
