package registry

import (
	"github.com/gorilla/mux"

	"github.com/benhawker/cachigo/internal/caching"
	"github.com/benhawker/cachigo/internal/errors"
	"github.com/benhawker/cachigo/internal/projection"
)

type Registry struct {
	Mux        *mux.Router
	ErrHandler errors.Handler
	Suppliers  map[string]string
	Cache      caching.Cache
}

func (r Registry) Register() {
	qh := &projection.QueryHandler{
		Supplier:   r.Suppliers,
		Cache:      r.Cache,
		ErrHandler: r.ErrHandler,
	}

	r.Mux.HandleFunc("/api/hotels", qh.ServeHTTP)
}
