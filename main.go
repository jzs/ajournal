package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	ajournal "github.com/jzs/ajournal/app"
	"github.com/jzs/ajournal/utils/logger"
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

	var s ajournal.Configuration
	err := envconfig.Process("aj", &s)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}
	s.BuildVersion = BuildVersion
	s.BuildType = BuildType
	s.BuildTime = BuildTime

	log := logger.New(s.BuildType == BuildVersionDevel)
	base := ajournal.Setup(ctx, s, log)

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
