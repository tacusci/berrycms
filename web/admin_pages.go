package web

import (
	"net/http"

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
	pages := make([]db.Page, 0)

	pt := db.PagesTable{}
	rows, err := pt.Select(db.Conn, "createddatetime, uuid, title, route", "")
	defer rows.Close()

	if err != nil {
		logging.ErrorAndExit(err.Error())
	}

	for rows.Next() {
		p := db.Page{}
		rows.Scan(&p.CreatedDateTime, &p.UUID, &p.Title, &p.Route)
		pages = append(pages, p)
	}

	pctx := plush.NewContext()
	pctx.Set("unixtostring", UnixToTimeString)
	pctx.Set("title", "Pages")
	pctx.Set("pages", pages)

	RenderDefault(w, "admin.pages.html", pctx)
}

func (aph *AdminPagesHandler) Post(w http.ResponseWriter, r *http.Request) {}

func (aph *AdminPagesHandler) Route() string { return aph.route }

func (aph *AdminPagesHandler) HandlesGet() bool  { return true }
func (aph *AdminPagesHandler) HandlesPost() bool { return false }
