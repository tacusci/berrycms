package web

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//AdminPagesHandler contains response functions for pages admin page
type AdminPagesHandler struct {
	Router *MutableRouter
	route  string
}

//Get takes the web request and writes response to session
func (ph *AdminPagesHandler) Get(w http.ResponseWriter, r *http.Request) {
	pageroutes := make([]string, 0)

	pt := db.PagesTable{}
	row, err := pt.Select(db.Conn, "route", "")

	if err != nil {
		logging.ErrorAndExit(err.Error())
	}

	for row.Next() {
		p := db.Page{}
		row.Scan(&p.Route)
		pageroutes = append(pageroutes, p.Route)
	}

	pctx := plush.NewContext()
	pctx.Set("names", pageroutes)

	content, err := ioutil.ReadFile("res" + string(os.PathSeparator) + "admin.pages.html")
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	renderedContent, err := plush.Render(string(content), pctx)
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	w.Write([]byte(renderedContent))
}

func (ph *AdminPagesHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (ph *AdminPagesHandler) Route() string { return ph.route }

func (ph *AdminPagesHandler) HandlesGet() bool  { return true }
func (ph *AdminPagesHandler) HandlesPost() bool { return false }
