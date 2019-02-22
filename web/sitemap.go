package web

import (
	"fmt"
	"net/http"

	"github.com/tacusci/berrycms/sitemap"
	"github.com/tacusci/logging"
)

type SitemapHandler struct {
	Router *MutableRouter
	route  string
}

func (sh *SitemapHandler) Get(w http.ResponseWriter, r *http.Request) {
	//if the sitemap.xml has been disabled, don't continue
	if sh.Router.NoSitemap {
		fourOhFour(w, r)
		return
	}

	//if the sitemap.xml page hasn't been visited before cache won't have been generated
	if !sitemap.CacheExists() {
		logging.Debug(fmt.Sprintf("Sitemap.xml cache doesn't exist yet, creating it with URL hostname: %s", r.Host))
		//creates a sitemap string and loads into in-memory cache
		err := sitemap.Generate(r.URL.Scheme, r.Host)
		if err != nil {
			logging.Error(err.Error())
		}
	}

	w.Write(sitemap.CacheBytes())
}

func (rh *SitemapHandler) Post(w http.ResponseWriter, r *http.Request) {}

//Route get URI route for handler
func (rh *SitemapHandler) Route() string { return rh.route }

//HandlesGet retrieve whether this handler handles get requests
func (rh *SitemapHandler) HandlesGet() bool { return true }

//HandlesPost retrieve whether this handler handles post requests
func (rh *SitemapHandler) HandlesPost() bool { return false }
