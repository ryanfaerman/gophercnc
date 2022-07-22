package main

import (
	"os"

	"github.com/ryanfaerman/gophercnc/log"
	"github.com/ryanfaerman/gophercnc/version"
	"github.com/spf13/cobra"
)

var (
	logger              = log.WithFields("app", "gophercnc")
	globalLogLevel      = "info"
	globalLogFormat     = "logfmt"
	globalDisableColors = false

	root = &cobra.Command{
		Use:     "gophercnc",
		Version: version.Version.String(),
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			switch globalLogFormat {
			case "pretty":
				logger.Out = log.PrettyFormatter{
					Out:           os.Stderr,
					DisableColors: globalDisableColors,
				}
			case "json":
				logger.Out = log.JSONFormatter{Out: os.Stderr}
			}

			if l, err := log.ParseLevel(globalLogLevel); err != nil {
				logger.WithError(err).
					Warn("got error when parsing log level, defaulting to INFO")
			} else {
				logger = logger.WithLevel(l)
				if l != log.LevelInfo {
					logger.Info("changed log level", "level", l)
				}
			}

		},
	}
)

func init() {
	root.PersistentFlags().StringVar(&globalLogLevel, "log-level", "info", "minimum level of logs to print to STDERR")
	root.PersistentFlags().StringVar(&globalLogFormat, "log-format", "pretty", "show logs as: pretty, logfmt, json")
	root.PersistentFlags().BoolVar(&globalDisableColors, "no-color", false, "disable colorized output")
	root.AddCommand(
		cmdVersion,
		// cmdWeb,
	)
}

func main() {
	root.Execute()
}
