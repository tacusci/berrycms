package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tacusci/berrycms/db"
	"github.com/tacusci/berrycms/web"
	"github.com/tacusci/berrycms/web/config"
	"github.com/tacusci/logging/v2"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	opts   config.Options
	dbConn *sql.DB
	router web.MutableRouter
}

func New(opts config.Options) *Server {
	return &Server{
		router: web.MutableRouter{},
		opts:   opts,
	}
}

func (s *Server) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.New("unable to start berrycms: bootup cancelled")
	default:
		displayInitStartupMsg()

		db, err := connectDB()
		if err != nil {
			return err
		}
		s.dbConn = db

		httpSvr, err := newHttpServer(ctx, s.opts)
		if err != nil {
			return err
		}

		err = s.initRouter(ctx, httpSvr)
		if err != nil {
			return err
		}

		displayInitStartupHTTPServerMsg(s.opts.Addr)
		if err := s.listenAndServe(httpSvr); err != nil {
			return fmt.Errorf("unable to launch berrycms: %w", err)
		}
		return nil
	}
}

func (s *Server) listenAndServe(svr *http.Server) (err error) {
	if svr.TLSConfig == nil {
		err = svr.ListenAndServe()
		return
	}
	err = svr.ListenAndServeTLS("", "")
	return
}

func (s *Server) initRouter(ctx context.Context, svr *http.Server) error {
	select {
	case <-ctx.Done():
		return errors.New("ctx cancelled did not init router")
	default:
		s.router = web.MutableRouter{
			Server:              svr,
			ActivityLogLoc:      s.opts.ActivityLogLoc,
			AdminOff:            s.opts.AdminPagesDisabled,
			AdminHidden:         len(s.opts.AdminHiddenPassword) > 0,
			AdminHiddenPassword: s.opts.AdminHiddenPassword,
			NoRobots:            s.opts.NoRobots,
			NoSitemap:           s.opts.NoSitemap,
			CpuProfile:          s.opts.CpuProfile,
		}
		s.router.Reload()
		return nil
	}
}

func connectDB() (*sql.DB, error) {
	return db.Connect(db.SQLITE, "", "berrycms")
}

func newHttpServer(ctx context.Context, opts config.Options) (*http.Server, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx cancelled returned nil server")
	default:
		httpSvr := http.Server{
			Addr:         fmt.Sprintf("%s:%d", opts.Addr, opts.Port),
			WriteTimeout: time.Second * 60,
			ReadTimeout:  time.Second * 60,
			IdleTimeout:  time.Second * 120,
		}

		if len(opts.AutoCertDomain) > 0 {
			c := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(opts.AutoCertDomain),
				Cache:      autocert.DirCache(opts.AutoCertDomain),
			}
			httpSvr.Addr = ":https"
			httpSvr.TLSConfig = c.TLSConfig()
		}

		return &httpSvr, nil
	}
}

func displayInitStartupMsg() {
	logging.Info("Berry CMS %s", db.VERSION)
}

func displayInitStartupHTTPServerMsg(addr string) {
	logging.Info("Starting http server @ %s 🌏 ...", addr)
}

func (s *Server) Shutdown() <-chan struct{} {
	done := make(chan struct{})
	defer close(done)

	return done
}
