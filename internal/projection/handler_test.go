package projection_test

import (
	"testing"

	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cheekybits/is"

	"github.com/benhawker/cachigo/internal/caching"
	"github.com/benhawker/cachigo/internal/errors"
	"github.com/benhawker/cachigo/internal/projection"
)

func NewQueryHandler(cache caching.Cache, errHandler errors.Handler) *projection.QueryHandler {
	return &projection.QueryHandler{
		Supplier:   mockSuppliers(),
		Cache:      caching.NewCache(),
		ErrHandler: errHandler,
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
	req, err := http.NewRequest("GET", "/api/hotels?checkin=12122018&checkout=16122019&destination=paris&guests=2&suppliers=supplier2,supplier1,supplier3", nil)
	if err != nil {
		t.Fatal(err)
	}

	qh := NewQueryHandler(
		caching.NewCache(),
		mockError{},
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
		mockError{},
	)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(qh.ServeHTTP)
	handler.ServeHTTP(rr, req)

	is.Equal(rr.Code, http.StatusUnprocessableEntity)
	is.Equal(strings.Contains(rr.Body.String(), "You must pass the `destination` parameter."), true)
}

type mockError struct{}

func (m mockError) Handle(err error) {}
