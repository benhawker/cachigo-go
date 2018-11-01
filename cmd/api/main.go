package main

import (
	"net/http"
	"os"
	"os/signal"

	"context"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/benhawker/cachigo/internal/caching"
	"github.com/benhawker/cachigo/internal/initialize"
	"github.com/benhawker/cachigo/internal/registry"
	"github.com/benhawker/cachigo/internal/supplier"
	log "github.com/sirupsen/logrus"
)

const (
	suppliersFile = "suppliers.yml"
	swaggerFile   = "index.html"
	port          = ":1111"
)

var (
	wg              sync.WaitGroup
	suppliers       map[string]string
	suppliersClient supplier.Cli
	cache           caching.Cache
)

func init() {
	var err error

	if suppliers, err = initialize.ReadYAML(suppliersFile); err != nil {
		log.Fatalf("no suppliers file was found: %v", err)
	}

	suppliersClient, err = supplier.NewClient()
	if err != nil {
		log.Fatalf("error initializing the supplier client: %v", err)
	}
}

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	mux := mux.NewRouter()
	loggedRouter := handlers.LoggingHandler(os.Stdout, mux)

	reg := registry.Registry{
		Mux:             mux,
		Suppliers:       suppliers,
		SuppliersClient: suppliersClient,
		Cache:           caching.NewCache(),
		Logger:          log.StandardLogger(),
	}
	reg.Register()

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/swagger", swagger)

	server := &http.Server{
		Addr:         port,
		Handler:      loggedRouter,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		IdleTimeout:  time.Second,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Sever starting. Listening on port ", port)
		if err := server.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	<-stop
	log.Info("Shutting down...")
	server.Shutdown(context.Background())
	wg.Wait()
	log.Info("🏁 Server stopped")
}

func swagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, swaggerFile)
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}
