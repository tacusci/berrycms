package main

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/tacusci/berrycms/db"
)

var cache *bytes.Buffer

func Add(val *[]byte) error {
	if cache == nil {
		return errors.New("Sitemap cache immutable... User has likely disabled sitemap.xml")
	}
}

func Del(val *[]byte) error {
	return nil
}

func Generate(adminPagesDisabled bool) error {
	Reset()
	_, err := cache.WriteString("<?xml version\"1.0\" encoding\"UTF-8\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n\t")
	if err != nil {
		return err
	}

	pt := db.PagesTable{}
	rows, err := pt.Select(db.Conn, "route", "roleprotected = '0'")

	if err != nil {
		return err
	}

	var pageRouteToAdd string

	for rows.Next() {
		err := rows.Scan(&pageRouteToAdd)
		if err != nil {
			return err
		}

		_, err := cache.WriteString(fmt.Sprintf("<url>%s"))
	}
}

func CacheExists() bool {
	return cache != nil
}

func CacheBytes() []byte {
	return cache.Bytes()
}

func Reset() {
	//we don't want to allocate memory each reset
	if cache == nil {
		cache = &bytes.Buffer{}
	}
	cache.Reset()
}
