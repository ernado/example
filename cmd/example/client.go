package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ernado/example/internal/oas"
	"github.com/go-faster/errors"
	"github.com/go-faster/sdk/app"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func Client() *cobra.Command {
	return &cobra.Command{
		Use: "client",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.Run(func(ctx context.Context, lg *zap.Logger, t *app.Telemetry) error {
				// TODO: Instrument http client with OpenTelemetry.
				client, err := oas.NewClient("http://server:8080",
					oas.WithMeterProvider(t.MeterProvider()),
					oas.WithTracerProvider(t.TracerProvider()),
				)
				if err != nil {
					return errors.Wrap(err, "create client")
				}

				for i := range 10 {
					if _, err := client.CreateTask(ctx, &oas.CreateTaskRequest{
						Title: fmt.Sprintf("Task %d", i),
					}); err != nil {
						return errors.Wrap(err, "create task")
					}
				}

				g, ctx := errgroup.WithContext(ctx)
				g.Go(func() error {
					ticker := time.NewTicker(500 * time.Millisecond)
					defer ticker.Stop()
					for {
						select {
						case <-ctx.Done():
							return ctx.Err()
						case <-ticker.C:
							tasks, err := client.ListTasks(ctx)
							if err != nil {
								return errors.Wrap(err, "list tasks")
							}
							lg.Info("Tasks", zap.Int("count", len(tasks.Tasks)))
						}
					}
				})

				return g.Wait()
			})
			return nil
		},
	}
}
