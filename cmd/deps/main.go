package main

import (
	"os"

	raven "github.com/getsentry/raven-go"

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
	// panickedErr, _ := raven.CapturePanicAndWait(func() {
	// 	if err := rootCmd.Execute(); err != nil {
	// 		raven.CaptureErrorAndWait(err, nil)
	// 		printErrAndExitFailure(err)
	// 	} else {
	// 		os.Exit(0)
	// 	}
	// }, nil)

	// if panickedErr != nil {
	// 	output.Error(fmt.Sprintf("Panic: %v", panickedErr))
	// 	os.Exit(1)
	// }

	if err := rootCmd.Execute(); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		printErrAndExitFailure(err)
	}
}
