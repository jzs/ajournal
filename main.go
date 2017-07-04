package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sketchground/ajournal/app"
	"github.com/sketchground/ajournal/utils"
	"github.com/sketchground/ajournal/utils/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	params := app.Params{
		StripeKey:    stripeKey,
		WWWDir:       wwwdir,
		BuildVersion: BuildVersion,
		BuildTime:    BuildTime,
		BuildType:    BuildType,
	}
	base := app.SetupRouter(db, log, translator, params)

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
