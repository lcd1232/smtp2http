package main

import (
	"runtime"
	"time"

	"github.com/alash3al/go-smtpsrv"
	"github.com/sirupsen/logrus"
)

var logger logrus.FieldLogger

func initLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	if *flagDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logger = logrus.StandardLogger()
}

func main() {
	initLogger()
	srv := &smtpsrv.Server{
		Addr:        *flagListenAddr,
		MaxBodySize: *flagMaxMessageSize,
		Handler:     handler,
		Name:        *flagServerName,
	}
	logger.
		WithField("addr", *flagListenAddr).
		WithField("max-body-size", *flagMaxMessageSize).
		WithField("server-name", *flagServerName).
		WithField("webhook-url", *flagWebhook).
		WithField("strict-validation", *flagStrictValidation).
		WithField("go-version", runtime.Version()).
		Info("app started")
	defer func() {
		logger.WithField("go-version", runtime.Version()).Info("app shutdown")
	}()
	if err := srv.ListenAndServe(); err != nil {
		logger.WithError(err).Error("got error")
	}
}
