package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// service interface with start and stop methods
type Service interface {
	Start() error
	Stop() error
}

// lifecycle for service starting and stopping
type Lifecycle struct {
	svc         Service
	name        string
	metricsPort int
	msrv        *http.Server
}

type Option func(lc *Lifecycle)

func Name(name string) Option {
	return func(args *Lifecycle) {
		args.name = name
	}
}

func MetricsPort(p int) Option {
	return func(args *Lifecycle) {
		args.metricsPort = p
	}
}

func NewLifecycle(svc Service, opts ...Option) *Lifecycle {
	lc := &Lifecycle{
		svc:         svc,
		metricsPort: 9090,
	}

	for _, setter := range opts {
		setter(lc)
	}

	return lc
}

// service runner
func (lc *Lifecycle) Run() error {
	// run the service. It is assumed that method will not block
	err := lc.svc.Start()
	if err != nil {
		return errors.Wrap(err, "can't start service")
	}

	// run metrics server
	go lc.runMetrics()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// wait for interrupt
	sig := <-sigc

	log.Info().Msgf("Caught signal %s: shutting down.", sig)

	log.Info().Msg("Stopping metrics server...")
	err = lc.msrv.Shutdown(context.Background())

	if err != nil {
		log.Error().Msg("Error stopping metrics server")
	}

	// stopping the service
	err = lc.svc.Stop()

	if err != nil {
		return errors.Wrap(err, "can't stop service")
	}

	return nil
}

// prometheus http metrics server
func (lc *Lifecycle) runMetrics() {
	addr := ":" + strconv.Itoa(lc.metricsPort)
	lc.msrv = &http.Server{Addr: addr}
	http.Handle("/metrics", promhttp.Handler())
	log.Info().Msgf("Starting metrics server on %s", addr)

	err := lc.msrv.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start metrics http server")
		}
	}
}
