package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/oklog/run"
	"github.com/onprem/muppet/pkg/api"
)

func main() {
	client, err := api.NewClientWithResponses("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	hostname := "host001"

	ctx, cancel := context.WithCancel(context.Background())

	var g run.Group

	g.Add(run.SignalHandler(ctx, os.Interrupt))

	{
		fn := func(ctx context.Context) error {
			resp, err := client.ListCommandQueueWithResponse(ctx, hostname)
			if err != nil {
				return err
			}

			for _, v := range *resp.JSON200 {
				cmd := exec.CommandContext(ctx, "sh", "-c", v.ShellCommand)
				_ = cmd.Run()
				log.Printf("ran commad: %s; exit status: %d\n", v.ShellCommand, cmd.ProcessState.ExitCode())
				_, err := client.MarkCommandDone(
					ctx,
					hostname,
					api.MarkCommandDoneJSONRequestBody{Uuid: v.Uuid, ExitStatus: float32(cmd.ProcessState.ExitCode())},
				)

				if err != nil {
					log.Printf("err: %v", err)
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
						log.Printf("err: %v", err)
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
