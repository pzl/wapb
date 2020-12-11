package main

import (
	"net/http"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/pzl/wapb/internal/server"
)

//go:generate go run assets_gen.go

func main() {
	cfg, ctx, cancel, log := setup()
	defer cancel()

	dbopts := badger.DefaultOptions(cfg.DBPath).WithLogger(log)
	if cfg.DBPath == ":MEMORY:" {
		dbopts.InMemory = true
		dbopts.Dir = ""
		dbopts.ValueDir = ""
	}

	db, err := badger.Open(dbopts)
	if err != nil {
		log.WithError(err).Error("unable to open database")
		panic(err)
	}
	defer db.Close()

	srv, err := server.New(log, cfg.Port, cfg.Handler, db)
	if err != nil {
		log.WithError(err).Error("error creating server")
		panic(err)
	}

	err = srv.Start(ctx)
	defer func() {
		if err := srv.Shutdown(); err != nil {
			log.WithError(err).Error("unable to gracefully shutdown")
		}
	}()

	if err != nil && err != http.ErrServerClosed {
		log.WithError(err).Error("server error")
	}
}
