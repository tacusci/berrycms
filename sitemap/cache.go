package sitemap

import (
	"bytes"
	"fmt"

	"github.com/tacusci/berrycms/db"
)

var cache *bytes.Buffer

func Generate() error {
	Reset()
	_, err := cache.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
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

		_, err = cache.WriteString(fmt.Sprintf("\t<url>\n\t<loc>%s</loc>\n</url>\n", pageRouteToAdd))
		if err != nil {
			return err
		}
	}

	cache.WriteString("</urlset>")

	return nil
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
