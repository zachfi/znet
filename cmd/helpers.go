package cmd

import (
	"os"

	kitlog "github.com/go-kit/log"
	log "github.com/sirupsen/logrus"
)

func initLogger() {
	formatter := log.TextFormatter{
		DisableQuote:     true,
		DisableTimestamp: true,
	}

	log.SetFormatter(&formatter)
	if trace {
		log.SetLevel(log.TraceLevel)
	} else if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.WithFields(log.Fields{
		"version": Version,
	}).Info("starting")
}

func newLogger() kitlog.Logger {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	return logger
}
