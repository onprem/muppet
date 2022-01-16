package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	flag "github.com/spf13/pflag"

	"github.com/onprem/muppet/pkg/api"
	"github.com/onprem/muppet/pkg/server"
	"github.com/onprem/muppet/pkg/store"
)

type config struct {
	address string
}

func parseFlags() (*config, error) {
	cfg := &config{}

	flag.StringVar(&cfg.address, "address", "0.0.0.0:8080", "The address to start the HTTP server on.")
	flag.Parse()

	return cfg, nil
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	cfg, err := parseFlags()
	if err != nil {
		level.Error(logger).Log("msg", "parsing flags", "err", err)

		return
	}

	ctx := context.Background()

	var g run.Group

	g.Add(run.SignalHandler(ctx, os.Interrupt))

	{
		srv := &http.Server{
			Addr:    cfg.address,
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
