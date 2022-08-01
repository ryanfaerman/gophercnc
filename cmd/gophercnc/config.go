package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/ryanfaerman/gophercnc/config"
	"github.com/spf13/cobra"
)

var (
	cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "work with system configuration",
		RunE: func(_ *cobra.Command, args []string) error {

			fmt.Println("workin on da konfig")

			configDir, err := os.UserConfigDir()
			if err != nil {
				return err
			}
			path := filepath.Join(configDir, ApplicationName)
			if err := os.MkdirAll(path, 0750); err != nil {
				return err
			}

			fmt.Println(config.LibraryCachePath())
			spew.Dump(config.Libraries())

			return nil
		},
	}
)
