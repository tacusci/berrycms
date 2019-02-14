package main

import (
	"bytes"
	"errors"
)

var cache *bytes.Buffer

func Add(val *[]byte) error {
	if cache == nil {
		return errors.New("Sitemap cache immutable... User has likely disabled sitemap.xml")
	}
}

func Del(val *[]byte) error {}

func Generate(adminPagesDisabled bool) error {}

func CacheExists() bool {}

func CacheBytes() []byte {}

func Reset() {
	//we don't want to allocate memory each reset
	if cache == nil {
		cache = &bytes.Buffer{}
	}
	cache.Reset()
}
