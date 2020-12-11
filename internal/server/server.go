package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Log          *logrus.Logger
	Router       *chi.Mux
	AssetHandler StaticHandler
	DB           *badger.DB
	Http         *http.Server
}

func New(log *logrus.Logger, port int, sh StaticHandler, db *badger.DB) (*Server, error) {

	if db == nil {
		return nil, errors.New("No valid database provided")
	}

	router := chi.NewRouter()

	s := &Server{
		Log:          log,
		Router:       router,
		AssetHandler: sh,
		DB:           db,
		Http: &http.Server{
			Addr:           ":" + strconv.Itoa(port),
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   60 * time.Second,
			IdleTimeout:    300 * time.Second,
			MaxHeaderBytes: 1 << 16,
			// TLSConfig: tlsConfig,
			Handler: router,
		},
	}

	s.SetupRoutes()

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	errs := make(chan error)

	tps := []string{"tcp4", "tcp6"}
	for _, l := range tps {
		s.Log.WithField("transport", l).WithField("addr", s.Http.Addr).Debug("opening socket")
		n, err := net.Listen(l, s.Http.Addr)
		if err != nil {
			return err
		}
		go func() {
			errs <- s.Http.Serve(n)
		}()
	}
	s.Log.Info("listening on -> " + s.Http.Addr)

	var err error
	select {
	case err = <-errs:
	case <-ctx.Done():
	}

	return err
}

func (s *Server) Shutdown() error {
	s.Log.Info("gracefully shutting down server")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) // nolint not bothering to cancel since we're near process exit
	return s.Http.Shutdown(ctx)
}
