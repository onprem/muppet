package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
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
}

func parseFlags() (*config, error) {
	cfg := &config{}

	flag.StringVar(&cfg.serviceURL, "service-url", "http://localhost:8080", "The URL of muppet service to fetch commands from.")
	flag.StringVar(&cfg.hostname, "hostname", "", "The hostname to fetch commands for.")

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
	}

	ctx, cancel := context.WithCancel(context.Background())

	var g run.Group

	g.Add(run.SignalHandler(ctx, os.Interrupt))

	{
		fn := func(ctx context.Context) error {
			resp, err := client.ListCommandQueueWithResponse(ctx, cfg.hostname)
			if err != nil {
				return err
			}

			for _, v := range *resp.JSON200 {
				logger = log.With(logger, "uuid", v.Uuid, "command", v.ShellCommand)

				cmd := exec.CommandContext(ctx, "sh", "-c", v.ShellCommand)
				_ = cmd.Run()

				level.Info(logger).Log("msg", "ran command", "exitcode", cmd.ProcessState.ExitCode())

				_, err := client.MarkCommandDone(
					ctx,
					cfg.hostname,
					api.MarkCommandDoneJSONRequestBody{Uuid: v.Uuid, ExitStatus: float32(cmd.ProcessState.ExitCode())},
				)

				if err != nil {
					level.Error(logger).Log("msg", "marking command as done", "err", err)
				}
			}

			return nil
		}

		g.Add(func() error {
			ticker := time.NewTicker(time.Minute)
			for {
				select {
				case <-ticker.C:
					if err := fn(ctx); err != nil {
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
