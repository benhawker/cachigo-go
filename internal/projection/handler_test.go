package projection_test

import (
	"testing"

	"net/http"
	"net/http/httptest"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/cheekybits/is"

	"github.com/benhawker/cachigo/internal/caching"
	"github.com/benhawker/cachigo/internal/projection"
	"github.com/benhawker/cachigo/internal/supplier"
)

func NewQueryHandler(cache caching.Cache, sClient supplier.Cli) *projection.QueryHandler {
	return &projection.QueryHandler{
		Supplier:        mockSuppliers(),
		SuppliersClient: sClient,
		Cache:           caching.NewCache(),
		Logger:          log.StandardLogger(),
	}
}

func mockSuppliers() map[string]string {
	ms := map[string]string{
		"supplier1": "https://api.myjson.com/bins/2tlb8",
		"supplier2": "https://api.myjson.com/bins/42lok",
	}
	return ms
}

func TestQueryHandler_ServeHTTP(t *testing.T) {
	is := is.New(t)
	queryStr := "/api/hotels?checkin=12122018&checkout=16122019&destination=paris&guests=2&suppliers=supplier2,supplier1,supplier3"
	req, err := http.NewRequest("GET", queryStr, nil)
	if err != nil {
		t.Fatal(err)
	}

	qh := NewQueryHandler(
		caching.NewCache(),
		supplier.MockClient{
			MakeRequestFunc: func(url string) (supplier.Response, error) {
				return nil, nil
			},
		},
	)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(qh.ServeHTTP)
	handler.ServeHTTP(rr, req)

	is.Equal(rr.Code, http.StatusOK)
	is.NotNil(rr.Body.String())
}

func TestQueryHandler_ServeHTTPWithoutRequiredQueryParameters(t *testing.T) {
	is := is.New(t)
	req, err := http.NewRequest("GET", "/api/hotels?checkin=12122018", nil)
	if err != nil {
		t.Fatal(err)
	}
	qh := NewQueryHandler(
		caching.NewCache(),
		supplier.MockClient{
			MakeRequestFunc: func(url string) (supplier.Response, error) {
				return nil, nil
			},
		},
	)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(qh.ServeHTTP)
	handler.ServeHTTP(rr, req)

	is.Equal(rr.Code, http.StatusUnprocessableEntity)
	is.Equal(strings.Contains(rr.Body.String(), "you must pass the `destination` parameter"), true)
}
