package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	entdb "github.com/ernado/example/internal/db/ent"
	instrumentationdb "github.com/ernado/example/internal/db/instrumentation"
	"github.com/ernado/example/internal/handler"
	"github.com/ernado/example/internal/oas"
	"github.com/go-faster/errors"
	"github.com/go-faster/sdk/app"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func Server() *cobra.Command {
	return &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.Run(func(ctx context.Context, lg *zap.Logger, t *app.Telemetry) error {
				// TODO: Refactor into Application.
				entClient, err := entdb.Open(ctx, os.Getenv("DATABASE_URL"), t)
				if err != nil {
					return errors.Wrap(err, "connect to db")
				}

				// TODO: HACK
				if err := entClient.Schema.Create(ctx); err != nil {
					return errors.Wrap(err, "migrate schema")
				}

				db := entdb.New(entClient)
				instrumentedDB, err := instrumentationdb.New(
					db,
					t.TracerProvider(),
					t.MeterProvider(),
				)
				if err != nil {
					return errors.Wrap(err, "create db instrumentation layer")
				}
				h, err := handler.New(
					instrumentedDB,
					t.TracerProvider(),
					t.MeterProvider(),
				)
				if err != nil {
					return errors.Wrap(err, "create handler")
				}
				s, err := oas.NewServer(
					h,
					oas.WithMeterProvider(t.MeterProvider()),
					oas.WithTracerProvider(t.TracerProvider()),
				)
				if err != nil {
					return errors.Wrap(err, "create server")
				}

				// TODO: Instrument with OpenTelemetry.
				svc := &http.Server{
					Handler:           s,
					Addr:              ":8080",     // TODO: configure
					ReadHeaderTimeout: time.Second, // Prevent Slowloris attacks.
					BaseContext: func(_ net.Listener) context.Context {
						// NB: Using BaseContext is important to properly execute graceful shutdown.
						// BaseContext is canceled when graceful shutdown is completed.
						return t.BaseContext()
					},
				}

				g, gCtx := errgroup.WithContext(ctx)
				g.Go(func() error {
					lg.Info("Starting server", zap.String("addr", svc.Addr))
					if err := svc.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						return errors.Wrap(err, "listen and serve")
					}
					return nil
				})
				g.Go(func() error {
					// NB: Using ShutdownContext is important to properly execute graceful shutdown.
					shutdownContext := t.ShutdownContext()
					select {
					case <-gCtx.Done():
						// Non-graceful shutdown.
						lg.Warn("Context done before shutdown")
					case <-shutdownContext.Done():
						lg.Info("Shutting down server")
					}
					// NB: Explicitly using t.BaseContext() to ensure that server
					// is properly shut down before application exits.
					//
					// This context is canceled when shutdown is completed.
					return svc.Shutdown(t.BaseContext())
				})

				return g.Wait()
			})

			return nil
		},
	}
}
