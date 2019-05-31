package main

import (
	"os"

	"github.com/getsentry/raven-go"

	"github.com/dropseed/deps/internal/output"
)

func init() {
	// Don't use the regular SENTRY_DSN env var
	sentry_dsn := os.Getenv("DEPS_SENTRY_DSN")
	if sentry_dsn != "" {
		println("Sentry error reporting enabled")
	}
	raven.SetDSN(sentry_dsn)
}

func printErrAndExitFailure(err error) {
	output.Error(err.Error())
	os.Exit(1)
}

func main() {
	raven.CapturePanic(func() {
		if err := rootCmd.Execute(); err != nil {
			raven.CaptureErrorAndWait(err, nil)
			printErrAndExitFailure(err)
		}
	}, nil)
}
