package caching_test

import (
	"testing"

	"github.com/benhawker/cachigo/internal/caching"
	"github.com/benhawker/cachigo/internal/supplier"
	"github.com/cheekybits/is"
)

func TestCache_SetsAndGetsKeyValuePairCorrectly(t *testing.T) {
	is := is.New(t)

	cache := caching.NewCache()
	sr := supplier.Response{"abc": 123.45, "def": 789.12}
	err := cache.Set("test", sr)

	is.NoErr(err)
	val, hit := cache.Get("test")

	is.Equal(hit, true)
	is.Equal(val, sr)
}

func TestCache_GettingANonExistentKeyReturnsFalse(t *testing.T) {
	is := is.New(t)

	cache := caching.NewCache()
	sr := supplier.Response{"abc": 123.45, "def": 789.12}
	err := cache.Set("test", sr)

	is.NoErr(err)
	val, hit := cache.Get("somethingElse")

	is.Equal(hit, false)
	is.Equal(val, supplier.Response{})
}
