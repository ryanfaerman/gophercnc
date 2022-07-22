package main

import (
	"fmt"

	"github.com/ryanfaerman/gophercnc/version"
	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Display the version info",
	Args:  cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(version.Version)
	},
}
