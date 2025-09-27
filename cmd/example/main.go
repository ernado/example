package main

import (
	"context"

	"github.com/go-faster/sdk/app"
	"go.uber.org/zap"
)

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger, t *app.Telemetry) error {
		<-t.ShutdownContext().Done()
		lg.Info("Shutting down")
		return nil
	})
}
