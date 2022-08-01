package main

import (
	_ "github.com/glebarez/go-sqlite"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/ryanfaerman/gophercnc/config"
	"github.com/ryanfaerman/gophercnc/tool"
	"github.com/spf13/cobra"
)

var (
	toolsLibraryPath = ""

	cmdTools = &cobra.Command{
		Use:     "tools",
		Short:   "view tool library",
		Aliases: []string{"tool"},
		RunE:    cmdToolsLibrary.RunE,
	}

	cmdToolsLibrary = &cobra.Command{
		Use:     "library",
		Short:   "list available tools",
		Aliases: []string{"list"},
		RunE: func(_ *cobra.Command, args []string) error {
			var (
				path string
				err  error
			)
			if toolsLibraryPath == "" {
				active, err := config.ActiveLibrary()
				if err != nil {
					return err
				}
				path = active.Path
			} else {

				path, err = homedir.Expand(toolsLibraryPath)
				if err != nil {
					return err
				}
			}

			lib, err := tool.LoadLibrary(path)

			type summary struct {
				Number      int
				Description string
			}
			rows := []summary{}

			for _, t := range lib.Tools {
				rows = append(rows, summary{
					Number:      t.Number(),
					Description: t.Description,
				})
			}

			if err := renderer.Render(rows); err != nil {
				return err
			}
			return err
		},
	}
)

func init() {
	cmdToolsLibrary.PersistentFlags().StringVar(
		&toolsLibraryPath,
		"path",
		"",
		"path to tools library",
	)

	cmdTools.AddCommand(
		cmdToolsLibrary,
		// cmdLibrary,
		cmdToolsLibraryImport,
		cmdToolsLibraryList,
		cmdToolsLibraryActivate,
	)

	cmdToolsLibrary.AddCommand(
		cmdToolsLibraryActivate,
		cmdToolsLibraryImport,
		cmdToolsLibraryList,
	)
}
