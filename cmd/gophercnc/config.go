package main

import (
	"fmt"

	"github.com/ryanfaerman/gophercnc/config"
	"github.com/spf13/cobra"
)

var (
	cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "work with system configuration",
		RunE: func(_ *cobra.Command, args []string) error {

			all, err := config.GetAll()
			if err != nil {
				return err
			}

			if err := renderer.Render(all); err != nil {
				return err
			}

			return nil
		},
	}

	cmdConfigGet = &cobra.Command{
		Use:   "get",
		Short: "read a specific config uri",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {

			data, err := config.Get(args[0])
			if err != nil {
				return err
			}

			fmt.Println(data)

			return nil
		},
	}

	cmdConfigSet = &cobra.Command{
		Use:     "set",
		Short:   "set a specific config uri",
		Aliases: []string{"add"},
		Args:    cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return config.Set(args[0], args[1])
		},
	}

	cmdConfigUnset = &cobra.Command{
		Use:     "remove",
		Short:   "set a specific config uri",
		Aliases: []string{"unset", "rm"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return config.Unset(args[0])
		},
	}
)

func init() {
	cmdConfig.AddCommand(
		cmdConfigGet,
		cmdConfigSet,
		cmdConfigUnset,
	)
}
