package main

import (
	"expvar"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/e-inwork-com/go-profile-search-service/api"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if available
	err := godotenv.Load()
	if err != nil {
		log.Println("Enviroment file .env is not found!")
	}

	// Set Configuration
	var cfg api.Config

	// Read environment  from a command line and OS
	flag.IntVar(&cfg.Port, "port", 4002, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.StringVar(&cfg.SolrURL, "solr-url", os.Getenv("SOLRURL"), "Solr URL")
	flag.StringVar(&cfg.SolrProfile, "solr-profile", os.Getenv("SOLRPROFILE"), "Solr Profile Path")
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.Cors.TrustedOrigins = strings.Fields(val)
		return nil
	})
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	// Show version on the terminal
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", api.Version)
		fmt.Printf("Build time:\t%s\n", api.BuildTime)
		os.Exit(0)
	}

	// Set logger
	logger := api.New(os.Stdout, api.LevelInfo)

	// Publish variables
	expvar.NewString("version").Set(api.Version)
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	// Set the application
	app := &api.Application{
		Config: cfg,
		Logger: logger,
	}

	// Run the application
	err = app.Serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
