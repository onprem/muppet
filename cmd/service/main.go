package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"

	"github.com/onprem/muppet/pkg/api"
	"github.com/onprem/muppet/pkg/server"
	"github.com/onprem/muppet/pkg/store"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	ctx := context.Background()

	var g run.Group

	g.Add(run.SignalHandler(ctx, os.Interrupt))

	{
		srv := &http.Server{
			Addr:    ":8080",
			Handler: api.Handler(server.NewServer(store.NewInMemStore(), logger)),
		}

		g.Add(func() error {
			level.Info(logger).Log("msg", "starting http server")

			return srv.ListenAndServe()
		}, func(err error) {
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			// Gracefully shutdown the HTTP server.
			_ = srv.Shutdown(ctx)
			cancel()
		})
	}

	if err := g.Run(); err != nil {
		level.Error(logger).Log("err", err)
	}
}
