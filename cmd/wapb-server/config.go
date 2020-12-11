package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type Config struct {
	Port    int
	DevMode bool
}

func setup() (Config, context.Context, context.CancelFunc, *logrus.Logger) {
	// CLI arg handling
	verbose := pflag.CountP("verbose", "v", "increased logging. Use multiple times for more info")
	j := pflag.BoolP("json", "j", false, "output logs in JSON format")
	port := pflag.IntP("port", "p", 7473, "Listening port")
	if *port < 1 {
		*port = 7473
	}
	dev := pflag.BoolP("dev", "d", false, "enable development mode. Listens to npm dev server for static assets")

	pflag.Parse()

	// log setup based on args
	log := logrus.New()
	setLogLevel(log, *verbose)
	setLogMode(log, *j)

	// signal handling & shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		log.Info("got exit signal. Cleaning up")
		cancel()
	}()

	// done
	return Config{
		Port:    *port,
		DevMode: *dev,
	}, ctx, cancel, log

}

func setLogLevel(log *logrus.Logger, level int) {
	lvls := []logrus.Level{
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
	// restrict to bounds
	if level > 3 {
		level = 3
	} else if level < 0 {
		level = 0
	}
	log.SetLevel(lvls[level])
}

func setLogMode(log *logrus.Logger, useJSON bool) {
	if useJSON {
		log.Formatter = UTCFormatter{&logrus.JSONFormatter{
			TimestampFormat: time.RFC1123,
		}}
	} else {
		log.Formatter = UTCFormatter{&logrus.TextFormatter{
			TimestampFormat:  time.RFC1123,
			QuoteEmptyFields: true,
		}}
	}
}

type UTCFormatter struct{ logrus.Formatter }

func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}
