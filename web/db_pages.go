package web

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gobuffalo/plush"
	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/logging"
)

//SavedPageHandler handler to contain pointer to core router and the URI string
type SavedPageHandler struct {
	Router *MutableRouter
	route  string
}

//Get handles get requests to URI
func (sph *SavedPageHandler) Get(w http.ResponseWriter, r *http.Request) {
	pt := db.PagesTable{}
	//JUST FOR LIVE/HOT ROUTE REMAPPING TESTING
	if r.RequestURI == "/addnew" {
		for i := 0; i < 51; i++ {
			pt.Insert(db.Conn, db.Page{
				CreatedDateTime: time.Now().Unix(),
				Title:           fmt.Sprintf("Carbon %d", i),
				Route:           fmt.Sprintf("/carbonite-%d", i),
				Content:         fmt.Sprintf("<h2>Carbonite %d</h2>", i),
				Roleprotected:   true,
			})
		}
		sph.Router.Reload()
	}
	rows, err := pt.Select(db.Conn, "content", fmt.Sprintf("route = '%s'", r.RequestURI))
	defer rows.Close()
	if err != nil {
		logging.Error(err.Error())
		w.Write([]byte("<h1>500 Server Error</h1>"))
		return
	}
	p := &db.Page{}
	for rows.Next() {
		rows.Scan(&p.Content)
	}

	ctx := plush.NewContext()
	ctx.Set("pagecontent", template.HTML(p.Content))
	Render(w, p, ctx)
}

//Post handles post requests to URI
func (sph *SavedPageHandler) Post(w http.ResponseWriter, r *http.Request) {}

//Route get URI route for handler
func (sph *SavedPageHandler) Route() string { return sph.route }

//HandlesGet retrieve whether this handler handles get requests
func (sph *SavedPageHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (sph *SavedPageHandler) HandlesPost() bool { return false }
