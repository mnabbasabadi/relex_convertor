// Package main is the entrypoint for the service.
package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mnabbasabadi/relex_convertor/service/config"
	httpAPI "github.com/mnabbasabadi/relex_convertor/service/internal/api/http"
	"github.com/mnabbasabadi/relex_convertor/service/internal/logic"
	"github.com/rs/zerolog"
)

func main() {

	// Load config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	s, logger := StartServer(cfg.Address, cfg.ServiceName, cfg.LogLevel)

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info().Msgf("starting server on %s", s.Addr)
		serverErrors <- s.ListenAndServe()
	}()

	shutdown := listenForShutdown()
	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("error starting server")
			return
		}
		logger.Fatal().Err(err).Msg("error starting server")
	case sig := <-shutdown:
		logger.Info().Msgf("caught signal %v: terminating", sig)

	}

}

// ListenForShutdown creates a channel and subscribes to specific signals to trigger a shutdown of the service.
func listenForShutdown() chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	return shutdown
}

// StartServer starts the http server
func StartServer(addr, serviceName, logLevel string) (*http.Server, zerolog.Logger) {

	logger, err := config.SetupLogger(serviceName, "main", logLevel, false)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to setup logger")
	}
	defer func() {
		if err != nil {
			logger.Err(err).Msg("run error")
		}
	}()

	defer config.RecoverAndLogPanic(logger)

	logicLogger := logger.With().Str(config.LogStrKeyModule, "logic").Logger()
	l := logic.New(logicLogger)

	serverLogger := logger.With().Str(config.LogStrKeyModule, "server").Logger()
	httpHandler := httpAPI.NewServer(l, serverLogger)

	httpLogger := logger.With().Str(config.LogStrKeyModule, "http").Logger()
	s := setupHTTPServer(httpLogger, httpHandler, addr)
	return s, logger
}

// ContextKey is the type used for context keys.
// TODO: should be moved to a common package
type ContextKey string

var (
	contextKey ContextKey = "logger"
)

func setupHTTPServer(logger zerolog.Logger, mux http.Handler, httpServerAddr string) *http.Server {
	httpServer := http.Server{
		Addr:              httpServerAddr,
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           mux,
		BaseContext: func(listener net.Listener) context.Context {
			return context.WithValue(context.Background(), contextKey, logger)
		},
	}
	return &httpServer
}
