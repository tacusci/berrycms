package robots

import (
	"fmt"
	"github.com/tacusci/berrycms/db"
	"strings"
)

func GenerateFile() error {

	sb := strings.Builder{}

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
		fmt.Println(sb.String())
	}

	return nil
}
