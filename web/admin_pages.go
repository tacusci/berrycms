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
func (aph *AdminPagesHandler) Get(w http.ResponseWriter, r *http.Request) {
	pageroutes := make([]string, 0)

	pt := db.PagesTable{}
	rows, err := pt.Select(db.Conn, "route", "")
	defer rows.Close()

	if err != nil {
		logging.ErrorAndExit(err.Error())
	}

	for rows.Next() {
		p := db.Page{}
		rows.Scan(&p.Route)
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

func (aph *AdminPagesHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (aph *AdminPagesHandler) Route() string { return aph.route }

func (aph *AdminPagesHandler) HandlesGet() bool  { return true }
func (aph *AdminPagesHandler) HandlesPost() bool { return false }
