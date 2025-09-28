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
	"github.com/ernado/example/internal/service"
	"github.com/ernado/example/internal/serviceinstrument"
	"github.com/go-faster/errors"
	"github.com/go-faster/sdk/app"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
				instrumentedDB, err := instrumentationdb.NewDBInstrumentation(
					db,
					t.TracerProvider(),
					t.MeterProvider(),
				)
				if err != nil {
					return errors.Wrap(err, "create db instrumentation layer")
				}
				instrumentedService, err := serviceinstrument.NewServiceInstrumentation(
					service.New(instrumentedDB),
					t.TracerProvider(),
					t.MeterProvider(),
				)
				h, err := handler.New(
					instrumentedService,
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

				svc := &http.Server{
					Handler: otelhttp.NewHandler(s, "",
						otelhttp.WithPropagators(t.TextMapPropagator()),
						otelhttp.WithMeterProvider(t.MeterProvider()),
						otelhttp.WithTracerProvider(t.TracerProvider()),
						otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
							op, ok := s.FindPath(r.Method, r.URL)
							if ok {
								return "http" + "." + op.OperationID()
							}
							return operation
						}),
					),
					Addr:              ":8080",
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
