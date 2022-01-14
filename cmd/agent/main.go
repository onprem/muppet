package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	flag "github.com/spf13/pflag"

	"github.com/onprem/muppet/pkg/api"
)

type config struct {
	serviceURL string
	hostname   string
	interval   uint
}

func parseFlags() (*config, error) {
	cfg := &config{}

	flag.StringVar(&cfg.serviceURL, "service-url", "http://localhost:8080", "The URL of muppet service to fetch commands from.")
	flag.StringVar(&cfg.hostname, "hostname", "", "The hostname to fetch commands for.")
	flag.UintVar(&cfg.interval, "interval", 60, "The interval at which to poll the muppet service for commands to execute, given in seconds.")

	flag.Parse()

	if cfg.hostname == "" {
		return nil, errors.New("hostname is required")
	}

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

	client, err := api.NewClientWithResponses(cfg.serviceURL)
	if err != nil {
		level.Error(logger).Log("msg", "creating API client", "err", err)

		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	var g run.Group

	g.Add(run.SignalHandler(ctx, os.Interrupt))

	{
		level.Info(logger).Log("msg", "starting command agent", "hostname", cfg.hostname, "serviceURL", cfg.serviceURL)

		if err := fetchAndRun(ctx, client, cfg.hostname, logger); err != nil {
			level.Error(logger).Log("msg", "fetch and run commands", "err", err)
		}

		g.Add(func() error {
			ticker := time.NewTicker(time.Duration(cfg.interval) * time.Second)
			for {
				select {
				case <-ticker.C:
					if err := fetchAndRun(ctx, client, cfg.hostname, logger); err != nil {
						level.Error(logger).Log("msg", "fetch and run commands", "err", err)
					}
				case <-ctx.Done():
					return nil
				}
			}
		}, func(e error) {
			cancel()
		})
	}

	g.Run()
}
