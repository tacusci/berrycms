package server

import (
	"context"

	"github.com/tacusci/berrycms/web"
	"github.com/tacusci/berrycms/web/config"
)

type Server struct {
	opts   config.Options
	router web.MutableRouter
}

func New(opts config.Options) *Server {
	return &Server{
		router: web.MutableRouter{},
		opts:   opts,
	}
}

func (s *Server) Start(ctx context.Context) {}

func (s *Server) Shutdown() <-chan struct{} {
	done := make(chan struct{})
	defer close(done)

	return done
}
