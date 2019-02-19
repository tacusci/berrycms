package web

import (
	"net/http"

	"github.com/tacusci/berrycms/sitemap"
)

type SitemapHandler struct {
	Router *MutableRouter
	route  string
}

func (sh *SitemapHandler) Get(w http.ResponseWriter, r *http.Request) {
	//if the sitemap.xml file cache hasn't been created then technically there is no sitemap page
	if !sitemap.CacheExists() {
		fourOhFour(w, r)
		return
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
