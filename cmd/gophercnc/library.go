package main

import (
	"fmt"
	"sort"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/ryanfaerman/gophercnc/config"
	"github.com/ryanfaerman/gophercnc/tool"
	"github.com/spf13/cobra"
)

var (
	cmdLibrary = &cobra.Command{
		Use:   "library",
		Short: "view tool library",
		// Aliases: []string{"list", "library"},
		RunE: func(_ *cobra.Command, _ []string) error {
			active, err := config.ActiveLibrary()
			if err != nil {
				return nil
			}
			fmt.Println(active.Name)

			// lib, err := tool.LoadLibrary(active.Path)

			// type summary struct {
			// 	Number      int
			// 	Description string
			// }
			// rows := []summary{}

			// for _, t := range lib.Tools {
			// 	rows = append(rows, summary{
			// 		Number:      t.Number(),
			// 		Description: t.Description,
			// 	})
			// }

			// if err := renderer.Render(rows); err != nil {
			// 	return err
			// }

			return nil
		},
	}
	cmdToolsLibraryImport = &cobra.Command{
		Use:     "import NAME PATH",
		Short:   "load a tool library",
		Aliases: []string{"add"},
		Args:    cobra.MinimumNArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			var err error
			name, path := args[0], args[1]

			path, err = homedir.Expand(path)
			if err != nil {
				return err
			}

			return tool.ImportLibrary(name, path)
		},
	}

	cmdToolsLibraryList = &cobra.Command{
		Use:     "libraries",
		Short:   "show all tool libraries",
		Aliases: []string{"list"},
		RunE: func(_ *cobra.Command, _ []string) error {

			type row struct {
				Name   string
				Active bool
			}

			active, err := config.ActiveLibrary()
			if err != nil {
				return nil
			}

			rows := []row{}

			libs, err := config.Libraries()
			if err != nil {
				return err
			}

			for _, lib := range libs {
				r := row{
					Name:   lib.Name,
					Active: lib.Name == active.Name,
				}
				rows = append(rows, r)
			}

			sort.Slice(rows, func(i, j int) bool {
				return rows[i].Name < rows[j].Name
			})

			if err := renderer.Render(rows); err != nil {
				return err
			}
			return nil
		},
	}

	cmdToolsLibraryActivate = &cobra.Command{
		Use:     "activate",
		Short:   "mark a libary as active",
		Aliases: []string{"use"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return config.ActivateLibrary(args[0])
		},
	}
)

func init() {

	cmdLibrary.AddCommand(
		cmdToolsLibraryImport,
		cmdToolsLibraryActivate,
	)

}
