package main

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

//go:generate go run assets_gen.go

func main() {
	cfg, ctx, cancel, log := setup()
	defer cancel()

	router := makeRouter(log, cfg)
	srv := makeServer(cfg, router)

	err := Start(ctx, log, srv)
	defer func() {
		if err := Shutdown(log, srv); err != nil {
			log.WithError(err).Error("unable to gracefully shutdown")
		}
	}()

	if err != nil && err != http.ErrServerClosed {
		log.WithError(err).Error("server error")
	}
}

func makeServer(cfg Config, router *chi.Mux) *http.Server {
	return &http.Server{
		Addr:           ":" + strconv.Itoa(cfg.Port),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    300 * time.Second,
		MaxHeaderBytes: 1 << 16,
		// TLSConfig: tlsConfig,
		Handler: router,
	}
}

func Start(ctx context.Context, log *logrus.Logger, s *http.Server) error {
	errs := make(chan error)

	tps := []string{"tcp4", "tcp6"}
	for _, l := range tps {
		log.WithField("transport", l).WithField("addr", s.Addr).Debug("opening socket")
		n, err := net.Listen(l, s.Addr)
		if err != nil {
			return err
		}
		go func() {
			errs <- s.Serve(n)
		}()
	}
	log.Info("listening on :" + s.Addr)

	var err error
	select {
	case err = <-errs:
	case <-ctx.Done():
	}

	return err
}

// shuts down the provided server, given 2 seconds
func Shutdown(log *logrus.Logger, s *http.Server) error {
	log.Info("gracefully shutting down server")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) // nolint not bothering to cancel since we're near process exit
	return s.Shutdown(ctx)
}
