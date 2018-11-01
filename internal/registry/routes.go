package registry

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/benhawker/cachigo/internal/caching"
	"github.com/benhawker/cachigo/internal/projection"
	"github.com/benhawker/cachigo/internal/supplier"
)

type Registry struct {
	Mux             *mux.Router
	Logger          *log.Logger
	Suppliers       map[string]string
	SuppliersClient supplier.Cli
	Cache           caching.Cache
}

func (r Registry) Register() {
	qh := &projection.QueryHandler{
		Supplier:        r.Suppliers,
		SuppliersClient: r.SuppliersClient,
		Cache:           r.Cache,
		Logger:          r.Logger,
	}

	r.Mux.HandleFunc("/api/hotels", qh.ServeHTTP)
}
