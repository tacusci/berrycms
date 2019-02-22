package sitemap

import (
	"bytes"
	"fmt"

	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/util"
)

var cache *bytes.Buffer
var additionalRoutes *[]string

func Add(val *string) error {
	if additionalRoutes == nil {
		additionalRoutes = &[]string{}
	}

	for _, v := range *additionalRoutes {
		if *val == v {
			return nil
		}
	}

	*additionalRoutes = append(*additionalRoutes, *val)

	return nil
}

func Del(val *string) error {
	if additionalRoutes == nil {
		return nil
	}

	*additionalRoutes = util.RemoveStringFromSlice(*additionalRoutes, *val)

	return nil
}

func Generate(httpScheme string, urlDomainPrefix string) error {
	Reset()

	if httpScheme == "" {
		httpScheme = "http"
	}

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

		var alreadyExists bool
		for _, v := range *additionalRoutes {
			if pageRouteToAdd == v {
				alreadyExists = true
				break
			}
		}

		if alreadyExists {
			continue
		}

		_, err = cache.WriteString(fmt.Sprintf("\t<url>\n\t\t<loc>%s://%s%s</loc>\n\t</url>\n", httpScheme, urlDomainPrefix, pageRouteToAdd))
		if err != nil {
			return err
		}
	}

	if additionalRoutes != nil {
		for _, v := range *additionalRoutes {
			_, err = cache.WriteString(fmt.Sprintf("\t<url>\n\t\t<loc>%s://%s%s</loc>\n\t</url>\n", httpScheme, urlDomainPrefix, v))
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
