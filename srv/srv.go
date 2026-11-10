package main

import (
	"github.com/volts-dev/volts"
	"github.com/volts-dev/volts/logger"
	"github.com/volts-dev/volts/router"
	"github.com/volts-dev/volts/server"
)

func index(ctx *router.THttpContext) {
	ctx.RenderTemplate("/watermark/index.html", nil)
}

func xmain() {
	r := router.New()
	r.Url("GET", "/", index)
	srv := server.New(
		server.Router(r),
	)

	srv.Config().Init(
		server.Debug(),
	)

	if err := srv.Start(); err != nil {
		logger.Err(err)
	}

	// serve as a http server
	app := volts.New(volts.Server(srv))
	app.Run()
}
