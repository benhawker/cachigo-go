package caching

import (
	"bytes"
	"strconv"
)

type params interface {
	Checkin() string
	Checkout() string
	Destination() string
	NumberOfGuests() int
}

// BuildKey returns in the format: "01012017020012017istanbul2supplier1"
func (c *Cache) BuildKey(params params, supplierName string) string {
	guestsStr := strconv.Itoa(params.NumberOfGuests())

	var buf bytes.Buffer
	buf.WriteString(params.Checkin())
	buf.WriteString(params.Checkout())
	buf.WriteString(params.Destination())
	buf.WriteString(guestsStr)
	buf.WriteString(supplierName)

	return string(buf.String())
}
