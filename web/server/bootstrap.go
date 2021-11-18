package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tacusci/berrycms/web/config"
	"github.com/tacusci/logging"
)

func Bootup(opts config.Options) <-chan struct{} {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	ctx, cancelBootup := context.WithCancel(context.Background())
	proc := process{}

	startupErr := make(chan error)
	go proc.run(ctx, opts, startupErr)

	killsig := <-interrupt
	fmt.Print("\r")
	logging.Error(fmt.Sprintf("Received signal: %s", killsig))

	cancelBootup()
	return proc.stop()
}

type process struct {
	svr *Server
}

func (p *process) run(ctx context.Context, opts config.Options, err chan<- error) {
	p.svr = New(opts)
	select {
	case <-ctx.Done():
		err <- errors.New("startup cancelled")
		return
	default:
		err <- p.svr.Start(ctx)
	}
}

func (p *process) stop() <-chan struct{} {
	return p.svr.Shutdown()
}
