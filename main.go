package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kelseyhightower/envconfig"
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

// Configuration contains all Environment settings that has to be present to run the service
type Configuration struct {
	TranslateFolder   string `split_words:"true" required:"true"`
	StripeSK          string `split_words:"true" `
	DBUser            string `envconfig:"db_user" default:"jzs"`
	DBName            string `envconfig:"db_name" default:"ajournal"`
	DBPass            string `envconfig:"db_pass"`
	Port              string `default:":8080"`
	WWWDir            string `envconfig:"www_dir" default:"/var/www/ajournal"`
	S3Endpoint        string `split_words:"true" required:"true"`
	S3AccessKey       string `envconfig:"s3_accesskey" required:"true"`
	S3SecretKey       string `envconfig:"s3_secretkey" required:"true"`
	S3Mock            bool   `split_words:"true"`
	OauthClientID     string `split_words:"true"`
	OauthClientSecret string `split_words:"true"`
	OauthRedirectURL  string `split_words:"true"`
}

func main() {
	ctx := context.Background()
	log := logger.New(BuildType == BuildVersionDevel)

	var s Configuration
	err := envconfig.Process("aj", &s)
	if err != nil {
		log.Fatalf(ctx, "%v", err.Error())
	}

	translator, err := utils.NewTranslator(s.TranslateFolder, log)
	if err != nil {
		log.Fatalf(ctx, "Could not load translator. Reason : %v", err)
		return
	}

	log.Printf(ctx, "AJ_WWW_DIR is set to %v", s.WWWDir)

	passwordstr := ""
	if s.DBPass != "" {
		passwordstr = fmt.Sprintf("password=%v", s.DBPass)
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%v dbname=%v %v sslmode=disable", s.DBUser, s.DBName, passwordstr))
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

	var br blob.Repository
	if s.S3Mock {
		br = services.NewS3MockRepo()
	} else {
		br = services.NewS3Repo(s.S3Endpoint, s.S3AccessKey, s.S3SecretKey, "ajournal")
	}
	bs := blob.NewService(br)

	jr := postgres.NewJournalRepo(db, log)
	js := journal.NewService(jr)
	journal.SetupHandler(apirouter, js, bs, log)

	ur := postgres.NewUserRepo(db)
	us := user.NewService(translator, ur)
	user.SetupHandler(apirouter, us, log)

	pr := postgres.NewProfileRepo(db, log)
	sr := services.NewStripeSubscriptionRepo(s.StripeSK, db)
	ps := profile.NewService(pr, sr)
	profile.SetupHandler(apirouter, ps, log)

	creds := oauth.Credentials{
		ClientID:     s.OauthClientID,
		ClientSecret: s.OauthClientSecret,
		RedirectURL:  s.OauthRedirectURL,
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
		http.ServeFile(w, r, s.WWWDir+"/app.html")
	})

	// Setup static file handler
	baserouter.PathPrefix("/").Handler(http.FileServer(http.Dir(s.WWWDir)))

	base := negroni.New(negroni.NewRecovery(), log, translator)

	// Setup middleware that injects currently logged in user into the stack.
	base.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		nr := user.ApplyUserToRequest(r, us)
		next(w, nr)
	}))

	base.UseHandler(baserouter)

	log.Printf(context.Background(), "Starting server: %v.%v \tAt:%v", BuildType, BuildVersion, BuildTime)
	log.Printf(context.Background(), "Listening on: %v", s.Port)
	server := &http.Server{Addr: s.Port, Handler: base}

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
