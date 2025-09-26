package cmd

import (
	"fmt"
	"os"

	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/logging"
	"github.com/spf13/cobra"
)

var BinaryDownloadURL string
var VersionsAPI string
var RootCmd = &cobra.Command{}
var BinaryToInstall string
var windowsSuffix string = ".exe"
var logLevel string
var verbose bool

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level: debug, info, warn, error")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging (shorthand for --log-level=debug)")

	cobra.OnInitialize(func() {
		if verbose {
			logLevel = "debug"
		}
		logging.Setup(logLevel)
		logging.Debug("logger initialized", "level", logLevel)
	})
}
