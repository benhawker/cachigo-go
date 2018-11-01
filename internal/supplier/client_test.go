package supplier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benhawker/cachigo/internal/supplier"
)

type HttpClientMock struct{}

func TestSupplierClient_MakeRequest(t *testing.T) {
	supplierServer := SupplierResponseSuccessStub()
	defer supplierServer.Close()

	r, _ := supplier.NewClient()
	_, err := r.MakeRequest(supplierServer.URL)

	if err != nil {
		t.Error("unexpected http error", err)
	}
}

func SupplierResponseSuccessStub() *httptest.Server {
	var resp string
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp = "{\"abc\": 123.45, \"def\": 123.4534}"
		w.Write([]byte(resp))
	}))
}
