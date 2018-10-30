package projection_test

import (
	"testing"

	"net/http"

	"github.com/benhawker/cachigo/internal/projection"
	"github.com/cheekybits/is"
)

func TestParams_ParseValidParamsWithoutSuppliers(t *testing.T) {
	is := is.New(t)
	req, _ := http.NewRequest("GET", "/api/hotels?checkin=123&checkout=456&destination=istanbul&guests=2", nil)

	params, err := projection.NewRequestParams(req)

	is.NoErr(err)
	is.Equal(params.Checkin(), "123")
	is.Equal(params.Checkout(), "456")
	is.Equal(params.Destination(), "istanbul")
	is.Equal(params.NumberOfGuests(), 2)
}

func TestParams_ParseValidParamsWithSuppliers(t *testing.T) {
	is := is.New(t)
	req, _ := http.NewRequest("GET", "/api/hotels?checkin=123&checkout=456&destination=istanbul&guests=2&suppliers=supplier1,supplier2", nil)

	params, err := projection.NewRequestParams(req)

	is.NoErr(err)
	is.Equal(params.Checkin(), "123")
	is.Equal(params.Checkout(), "456")
	is.Equal(params.Destination(), "istanbul")
	is.Equal(params.NumberOfGuests(), 2)
	is.Equal(params.Suppliers(), []string{"supplier1", "supplier2"})
}

func TestParams_ParseParams(t *testing.T) {
	is := is.New(t)
	req, _ := http.NewRequest("GET", "/api/hotels", nil)
	_, err := projection.NewRequestParams(req)
	is.Err(err)
}

func TestParams_ParseGuestsAsNonIntegerParams(t *testing.T) {
	is := is.New(t)
	req, _ := http.NewRequest("GET", "/api/hotels?checkin=123&checkout=456&destination=istanbul&guests=two", nil)
	_, err := projection.NewRequestParams(req)

	is.Err(err)
}
