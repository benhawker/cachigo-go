package caching

import (
	"time"

	"github.com/benhawker/cachigo/internal/supplier"
)

const (
	expire_in_minutes = 5
)

type Cache struct {
	Store map[string]CacheValue
}

type CacheValue struct {
	Data   supplier.Response
	Expiry int64
}

func NewCache() Cache {
	return Cache{
		Store: map[string]CacheValue{},
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	if val, ok := c.Store[key]; ok {
		if val.Expiry > time.Now().Unix() {
			return val.Data, true
		}
	}

	return supplier.Response{}, false
}

func (c *Cache) Set(key string, value interface{}) error {
	sv := CacheValue{
		Data:   value.(supplier.Response),
		Expiry: expiryTime(),
	}

	c.Store[key] = sv
	return nil
}

func expiryTime() int64 {
	return time.Now().Unix() + (expire_in_minutes * 60)
}
