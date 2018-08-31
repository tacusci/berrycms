package web

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

type Handler interface {
	Route() string
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	HandlesGet() bool
	HandlesPost() bool
}

func GetDefaultHandlers(router *MutableRouter) []Handler {
	return []Handler{
		&LoginHandler{
			route:  "/login",
			Router: router,
		},
		&LogoutHandler{
			route:  "/logout",
			Router: router,
		},
		&AdminHandler{
			route:  "/admin",
			Router: router,
		},
		&AdminUsersHandler{
			route:  "/admin/users",
			Router: router,
		},
		&AdminPagesHandler{
			route:  "/admin/pages",
			Router: router,
		},
		&AdminPagesNewHandler{
			route:  "/admin/pages/new",
			Router: router,
		},
		&AdminPagesEditHandler{
			route:  "/admin/pages/edit/{uuid}",
			Router: router,
		},
	}
}

func UnixToTimeString(unix int64) string {
	return time.Unix(unix, 0).Format("15:04:05 02-01-2006")
}

func RenderDefault(w http.ResponseWriter, template string, pctx *plush.Context) error {
	header, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "header.snip")

	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return err
	}

	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + template)
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return err
	}

	renderedContent, err := plush.Render(string(append(append(header, []byte("\n")...), content...))+"\n</html>", pctx)

	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return err
	}
	_, err = w.Write([]byte(renderedContent))
	return err
}

func Render(w http.ResponseWriter, p *db.Page, ctx *plush.Context) error {
	html, err := plush.Render("<html><%= pagecontent %></html>", ctx)
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return err
	}
	w.Write([]byte(html))
	return err
}
