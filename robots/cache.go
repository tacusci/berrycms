package robots

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/coocood/freecache"
	"github.com/tacusci/berrycms/db"
)

var RobotsCache *freecache.Cache

func Generate() error {

	sb := bytes.Buffer{}

	sb.WriteString("User-agent: *\n")
	//block indexing admin pages
	//NOTE: if the admin pages URI has been hidden, we're deliberately omitting this from robots.txt
	sb.WriteString("Disallow: /admin\n")

	pt := db.PagesTable{}
	rows, err := pt.Select(db.Conn, "route", "roleprotected = '1'")

	if err != nil {
		return err
	}

	var pageRouteToDisallow string
	var writtenMoreRoutes bool

	for rows.Next() {
		err := rows.Scan(&pageRouteToDisallow)
		if err != nil {
			return err
		}

		c, err := sb.WriteString(fmt.Sprintf("Disallow: %s\n", pageRouteToDisallow))
		if err != nil {
			return err
		}

		if c > 0 {
			writtenMoreRoutes = true
		}
	}

	if writtenMoreRoutes {
		err := ioutil.WriteFile("./static/robots.txt", sb.Bytes(), 0644)
		if err != nil {
			return err
		}

	}

	return nil
}

func Cache(sb *bytes.Buffer) error {
	RobotsCache = freecache.NewCache(len(sb.Bytes()))
	key := []byte("robots")
	val := sb.Bytes()
	err := RobotsCache.Set(key, val, 0)
	if err != nil {
		return err
	}
	return nil
}
