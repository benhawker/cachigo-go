package caching_test

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/benhawker/cachigo/internal/caching"
)

type MockParams struct {
	CheckinFn 		 func() string
  	CheckoutFn 		 func() string
  	DestinationFn 	 func() string
  	NumberOfGuestsFn func() int
}

func (m MockParams) Checkin() string {
	return m.CheckinFn()
}

func (m MockParams) Checkout() string {
	return m.CheckoutFn()
}

func (m MockParams) Destination() string {
	return m.DestinationFn()
}

func (m MockParams) NumberOfGuests() int {
	return m.NumberOfGuestsFn()
}


func TestKeyBuilder_ReturnsTheExpectedKey(t *testing.T) {
	is := is.New(t)
	cache := caching.NewCache()

	mockParams := MockParams{
		CheckinFn: func() string {
			return "01012018"
		},
		CheckoutFn: func() string {
			return "08012018"
		},
		DestinationFn: func() string {
			return "Lisbon"
		},
		NumberOfGuestsFn: func() int {
			return 2
		},
	}

	key := cache.BuildKey(mockParams, "supplier1")
	expectedKey := "0101201808012018Lisbon2supplier1"

	is.Equal(key, expectedKey)

}