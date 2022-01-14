package main

import (
	"context"
	"os/exec"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/onprem/muppet/pkg/api"
)

func fetchAndRun(ctx context.Context, client *api.ClientWithResponses, hostname string, logger log.Logger) error {
	resp, err := client.ListCommandQueueWithResponse(ctx, hostname)
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
			hostname,
			api.MarkCommandDoneJSONRequestBody{Uuid: v.Uuid, ExitStatus: float32(cmd.ProcessState.ExitCode())},
		)

		if err != nil {
			level.Error(logger).Log("msg", "marking command as done", "err", err)
		}
	}

	return nil
}
