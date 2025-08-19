package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/nguyentk31/url-shortening/internal/config"
	"github.com/nguyentk31/url-shortening/internal/database"
	"github.com/nguyentk31/url-shortening/internal/routers"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

func main() {
	// Load config path from command line
	configPath := flag.String("c", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Instantiates the database
	postgres, err := database.NewPostgres(config.Database)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer postgres.Close()

	r := routers.NewRouter(postgres.DB)

	srv := newServer(r, config.Server)

	startServer(srv)
	defer closeServer(srv)
}

func newServer(router http.Handler, config config.Server) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
}

func startServer(srv *http.Server) {
	log.Printf("starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func closeServer(srv *http.Server) {
	log.Printf("closing server on %s", srv.Addr)
	if err := srv.Close(); err != nil {
		log.Fatal(err)
	}
}
