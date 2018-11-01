package caching

import (
	"time"

	"github.com/benhawker/cachigo/internal/supplier"
)

const (
	expireInMinutes = 5
)

// Cache defines the Store for the Cache
type Cache struct {
	Store map[string]CacheValue
}

// CacheValue defines the struct to hold the value stored in the Cache
type CacheValue struct {
	Data   supplier.Response
	Expiry int64
}

// NewCache Constructor function for a new Cache
func NewCache() Cache {
	return Cache{
		Store: map[string]CacheValue{},
	}
}

// Get retrieves a value from the Cache
func (c *Cache) Get(key string) (interface{}, bool) {
	if val, ok := c.Store[key]; ok {
		if val.Expiry > time.Now().Unix() {
			return val.Data, true
		}
	}

	return supplier.Response{}, false
}

// Set sets a key/value pair in the Cache
func (c *Cache) Set(key string, value interface{}) error {
	sv := CacheValue{
		Data:   value.(supplier.Response),
		Expiry: expiryTime(),
	}

	c.Store[key] = sv
	return nil
}

func expiryTime() int64 {
	return time.Now().Unix() + (expireInMinutes * 60)
}
