package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sketchground/ajournal/blob"
	"github.com/sketchground/ajournal/journal"
	"github.com/sketchground/ajournal/oauth"
	"github.com/sketchground/ajournal/postgres"
	"github.com/sketchground/ajournal/profile"
	"github.com/sketchground/ajournal/services"
	"github.com/sketchground/ajournal/user"
	"github.com/sketchground/ajournal/utils"
	"github.com/sketchground/ajournal/utils/logger"

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
	BuildTime = "xxxx-xx-xx"
)

func main() {
	ctx := context.Background()
	log := logger.New(BuildType == BuildVersionDevel)

	tFolder := os.Getenv("AJ_TRANSLATE_FOLDER")
	if tFolder == "" {
		log.Fatalf(ctx, "Environment variable AJ_TRANSLATE_FOLDER not set!\nRemember to set the path to your translate folder")
		return
	}

	translator, err := utils.NewTranslator(tFolder, log)
	if err != nil {
		log.Fatalf(ctx, "Could not load translator. Reason : %v", err)
		return
	}

	stripeKey := os.Getenv("AJ_STRIPE_SK")
	if stripeKey == "" {
		log.Fatalf(ctx, "Environment variable AJ_STRIPE_SK not set!\nRemember to set your stripe private key")
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
	dbpass := os.Getenv("AJ_DB_PASS")

	port := os.Getenv("AJ_PORT")
	if port == "" {
		port = ":8080"
	}

	wwwdir := os.Getenv("AJ_WWW_DIR")
	if wwwdir == "" {
		wwwdir = "/var/www/ajournal"
	}
	log.Printf(ctx, "AJ_WWW_DIR is set to %v", wwwdir)

	passwordstr := ""
	if dbpass != "" {
		passwordstr = fmt.Sprintf("password=%v", dbpass)
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%v dbname=%v %v sslmode=disable", dbuser, dbname, passwordstr))
	if err != nil {
		log.Fatalf(ctx, "Could not connect to database! %v", err)
		return
	}

	baserouter := mux.NewRouter()
	apirouter := baserouter.PathPrefix("/api").Subrouter().StrictSlash(true)

	apirouter.Path("/version").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := map[string]string{
			"Version": BuildVersion,
			"Time":    BuildTime,
			"Type":    BuildType,
		}
		utils.JSONResp(r.Context(), log, r, w, version, nil)
	})

	endpoint := os.Getenv("AJ_S3_ENDPOINT")
	if endpoint == "" {
		log.Fatalf(ctx, "Environment variable AJ_S3_ENDPOINT not set!\nRemember to set key")
	}
	accessKey := os.Getenv("AJ_S3_ACCESSKEY")
	if endpoint == "" {
		log.Fatalf(ctx, "Environment variable AJ_S3_ACCESSKEY not set!\nRemember to set key")
	}
	secretKey := os.Getenv("AJ_S3_SECRETKEY")
	if endpoint == "" {
		log.Fatalf(ctx, "Environment variable AJ_S3_SECRETKEY not set!\nRemember to set key")
	}

	var br blob.Repository
	mockS3 := os.Getenv("AJ_S3_MOCK")
	if mockS3 == "true" {
		// Mock S3
		br = services.NewS3MockRepo()
	} else {
		br = services.NewS3Repo(endpoint, accessKey, secretKey, "ajournal")
	}
	bs := blob.NewService(br)

	jr := postgres.NewJournalRepo(db, log)
	js := journal.NewService(jr)
	journal.SetupHandler(apirouter, js, bs, log)

	ur := postgres.NewUserRepo(db)
	us := user.NewService(translator, ur)
	user.SetupHandler(apirouter, us, log)

	pr := postgres.NewProfileRepo(db, log)
	sr := services.NewStripeSubscriptionRepo(stripeKey, db)
	ps := profile.NewService(pr, sr)
	profile.SetupHandler(apirouter, ps, log)

	creds := oauth.Credentials{
		ClientID:     os.Getenv("AJ_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("AJ_OAUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AJ_OAUTH_REDIRECT_URL"),
		Provider:     oauth.ProviderGoogle,
	}
	or := postgres.NewOauthRepo(db)
	oas := oauth.NewService(or, ur, pr)
	oauth.SetupHandler(apirouter, oas, log, creds)

	// Setup api router
	baserouter.PathPrefix("/api").Handler(negroni.New(negroni.Wrap(apirouter)))

	// Setup helper routes that redirects to public journal page
	baserouter.HandleFunc("/journal/{journalid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["journalid"]

		http.Redirect(w, r, fmt.Sprintf("/#/view/%v", id), http.StatusFound)
	})
	// Setup helper routes that redirects to public user page with his/her journals
	baserouter.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		http.Redirect(w, r, fmt.Sprintf("/app#/users/%v", id), http.StatusFound)
	})

	baserouter.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, wwwdir+"/app.html")
	})

	// Setup static file handler
	baserouter.PathPrefix("/").Handler(http.FileServer(http.Dir(wwwdir)))

	base := negroni.New(negroni.NewRecovery(), log, translator)

	// Setup middleware that injects currently logged in user into the stack.
	base.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		nr := user.ApplyUserToRequest(r, us)
		next(w, nr)
	}))

	base.UseHandler(baserouter)

	log.Printf(context.Background(), "Starting server: %v.%v \tAt:%v", BuildType, BuildVersion, BuildTime)
	log.Printf(context.Background(), "Listening on: %v", port)
	server := &http.Server{Addr: port, Handler: base}

	// subscribe to SIGINT signals
	sigchan := make(chan os.Signal, 5)
	signal.Notify(sigchan, os.Interrupt)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Error(context.Background(), err)
		}
	}()

	<-sigchan // wait for SIGINT
	log.Printf(ctx, "Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = server.Shutdown(tctx)
	if err != nil {
		log.Fatalf(ctx, "Failed shutting down server %v", err)
	}

	log.Printf(ctx, "Server gracefully stopped")
}
