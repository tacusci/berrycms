package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tacusci/berrycms/web/config"
	"github.com/tacusci/logging"
)

func Bootup(opts config.Options) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	ctx, cancelBootup := context.WithCancel(context.Background())
	proc := process{}
	go proc.run(ctx, opts)

	killsig := <-interrupt
	fmt.Print("\r")
	logging.Error(fmt.Sprintf("Received signal: %s", killsig))

	cancelBootup()
	<-proc.stop()
}

type process struct {
	svr *Server
}

func (p *process) run(ctx context.Context, opts config.Options) {
	p.svr = New(opts)
	p.svr.Start(ctx)
}

func (p *process) stop() <-chan struct{} {
	return p.svr.Shutdown()
}
