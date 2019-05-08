package projection

import (
	"encoding/json"
	"net/http"
	"sync"
	"fmt"

	"github.com/benhawker/cachigo/internal/caching"
	"github.com/benhawker/cachigo/internal/supplier"
	log "github.com/sirupsen/logrus"
)

// QueryHandler definition.
type QueryHandler struct {
	Supplier        map[string]string
	Cache           caching.Cache
	Logger          *log.Logger
	SuppliersClient supplier.Cli
}

// Response defines the response body for the hotels endpoint
type Response struct {
	Data []HotelOffer `json:"data"`
}

// HotelOffer defines each element of the response body.
type HotelOffer struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Supplier string  `json:"supplier"`
}

func (h *QueryHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestParams, err := NewRequestParams(req)
	if err != nil {
		h.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	validatedSuppliers := h.validateRequestedSuppliers(requestParams)

	allResponses := []HotelOffer{}

	responsesCh := make(chan supplier.ResponseWithError)
	wg := &sync.WaitGroup{} // Declare a wait group

	fmt.Println("Before loop")

	// STEP 1: Make concurrent requests
	for name, url := range validatedSuppliers {
		resp := make(chan supplier.ResponseWithError)    // channel to collect response
		go h.SuppliersClient.MakeRequestXX(name, url, wg, resp) // make actual request
	
		go func() {
			select {
			case r := <-resp: // when the request finishes
				responsesCh <- r // sucess, distribute response to the return channel
			}
			wg.Done() // Now goroutine is complete.
		}()
	}

	// STEP 2: Wait until all responses are recieved/completed.
	go func() {
		wg.Wait()
		close(responsesCh)
	}()

	// STEP 3: Gather all responses. Process will block until all responses are received.
	var responses []*supplier.ResponseWithError
	for resp := range responsesCh {
		responses = append(responses, &resp)
	}

	fmt.Println("After....")

	for _, response := range responses {
		if response.Error != nil {
			fmt.Println(err)
			h.Logger.Errorf("calling %s raised error: %s", response.Name, err.Error())
		}

		transformedResponse := transformSupplierResponse(response.Name, response)
		allResponses = append(allResponses, transformedResponse...)
	}

	// for name, url := range validatedSuppliers {
	// 	cacheKey := h.Cache.BuildKey(requestParams, name)
	// 	supplierResponse, cacheHit := h.Cache.Get(cacheKey)
	// 	if !cacheHit {
	// 		supplierResponse, err = h.SuppliersClient.MakeRequest(url)
	// 		if err != nil {
	// 			h.Logger.Errorf("calling %s raised error: %s", name, err.Error())
	// 		}
	// 		h.Cache.Set(cacheKey, supplierResponse)
	// 	}

	// 	transformedResponse := transformSupplierResponse(name, supplierResponse.(supplier.Response))
	// 	allResponses = append(allResponses, transformedResponse...)
	// }

	bestPrices := findBestPriceByHotel(allResponses)
	response := Response{
		Data: transformBestPrices(bestPrices),
	}

	b, err := json.Marshal(response)
	if err != nil {
		h.Logger.Error(err)
		http.Error(w, "unable to decode data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (h *QueryHandler) validateRequestedSuppliers(requestParams RequestParams) map[string]string {
	if len(requestParams.Suppliers()) == 0 {
		return h.Supplier
	}

	suppliers := map[string]string{}
	for _, supplierName := range requestParams.Suppliers() {
		if supplierURL, ok := h.Supplier[supplierName]; ok {
			suppliers[supplierName] = supplierURL
		}
	}
	return suppliers
}

func findBestPriceByHotel(allResponses []HotelOffer) map[string]HotelOffer {
	byHotel := map[string]HotelOffer{}
	for _, resp := range allResponses {
		if val, ok := byHotel[resp.ID]; ok {
			if resp.Price < val.Price {
				byHotel[resp.ID] = resp
			}
		} else {
			byHotel[resp.ID] = resp
		}
	}
	return byHotel
}

func transformBestPrices(bestPrices map[string]HotelOffer) []HotelOffer {
	resp := []HotelOffer{}
	for _, v := range bestPrices {
		resp = append(resp, v)
	}
	return resp
}

func transformSupplierResponse(name string, response *supplier.ResponseWithError) []HotelOffer {
	offers := []HotelOffer{}

	for hotel, price := range response.Body {
		ho := HotelOffer{
			ID:       hotel,
			Price:    price,
			Supplier: name,
		}

		offers = append(offers, ho)
	}
	return offers
}
