package projection

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// RequestParams ...
type RequestParams struct {
	DateParams       DateParams
	DestinationParam string
	GuestsParam      int
	SuppliersParam   []string
}

// DateParams ...
type DateParams struct {
	Checkin  string
	Checkout string
}

// NewRequestParams ...
func NewRequestParams(req *http.Request) (RequestParams, error) {
	p := RequestParams{}

	dateParams, err := parseDateParams(req)
	if err != nil {
		return RequestParams{}, err
	}
	p.DateParams = dateParams

	destination, err := parseDestinationParam(req)
	if err != nil {
		return RequestParams{}, err
	}
	p.DestinationParam = destination

	numberOfGuests, err := parseGuestsParam(req)
	if err != nil {
		return RequestParams{}, err
	}
	p.GuestsParam = numberOfGuests

	suppliers, err := parseSuppliersParam(req)
	if err != nil {
		return RequestParams{}, err
	}
	p.SuppliersParam = suppliers

	return p, nil
}

// Checkin ...
func (p RequestParams) Checkin() string {
	return p.DateParams.Checkin
}

// Checkout ...
func (p RequestParams) Checkout() string {
	return p.DateParams.Checkout
}

// Destination ...
func (p RequestParams) Destination() string {
	return p.DestinationParam
}

// NumberOfGuests ...
func (p RequestParams) NumberOfGuests() int {
	return p.GuestsParam
}

// Suppliers ...
func (p RequestParams) Suppliers() []string {
	return p.SuppliersParam
}

func parseDateParams(req *http.Request) (DateParams, error) {
	checkin := req.FormValue("checkin")
	checkout := req.FormValue("checkout")

	if checkin == "" && checkout == "" {
		return DateParams{}, errors.New("you must pass `checkin` & `checkout` parameters")
	}

	return DateParams{
		Checkin:  checkin,
		Checkout: checkout,
	}, nil
}

func parseSuppliersParam(req *http.Request) ([]string, error) {
	suppliers := req.FormValue("suppliers")

	if suppliers == "" {
		return []string{}, nil
	}

	return strings.Split(suppliers, ","), nil
}

func parseDestinationParam(req *http.Request) (string, error) {
	destination := req.FormValue("destination")

	if destination == "" {
		return "", errors.New("you must pass the `destination` parameter")
	}

	return destination, nil
}

func parseGuestsParam(req *http.Request) (int, error) {
	numberOfGuests := req.FormValue("guests")

	if numberOfGuests == "" {
		return 0, errors.New("you must pass the `guests` parameter")
	}

	i, err := strconv.Atoi(numberOfGuests)
	if err != nil {
		return 0, errors.New("you must pass the `guests` parameter as an integer")
	}

	return i, nil
}
