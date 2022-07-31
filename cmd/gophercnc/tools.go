package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/glebarez/go-sqlite"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/ryanfaerman/gophercnc/config/tutorial"
	"github.com/ryanfaerman/gophercnc/tool"
	"github.com/spf13/cobra"
)

var (
	toolsLibraryPath = ""
	cmdTools         = &cobra.Command{
		Use:     "tools",
		Short:   "view tool library",
		Aliases: []string{"tool"},
	}

	cmdToolsView = &cobra.Command{
		Use:     "view",
		Short:   "view tool library",
		Aliases: []string{"list", "library"},
		RunE: func(_ *cobra.Command, args []string) error {
			if toolsLibraryPath == "" {
				return errors.New("missing tools library path")
			}

			//"~/Downloads/mpcnc.tools"
			path, err := homedir.Expand(toolsLibraryPath)
			if err != nil {
				return err
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
	cmdToolsImport = &cobra.Command{
		Use:     "import",
		Short:   "load a tool library",
		Aliases: []string{"load"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Println(args)
			fmt.Println(os.UserConfigDir())

			ctx := context.Background()
			db, err := sql.Open("sqlite", "./test.db")
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()

			// // get SQLite version
			// var version string
			// err = db.QueryRow("select sqlite_version()").Scan(&version)
			// if err != nil {
			// 	return err
			// }

			// fmt.Println(version)

			queries := tutorial.New(db)
			// list all authors
			authors, err := queries.ListAuthors(ctx)
			if err != nil {
				return err
			}
			log.Println(authors)

			// create an author
			insertedAuthor, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
				Name: "Brian Kernighan",
				Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
			})
			if err != nil {
				return err
			}
			log.Println(insertedAuthor)

			// get the author we just inserted
			fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor)
			if err != nil {
				return err
			}

			spew.Dump(fetchedAuthor)

			return nil
		},
	}
)

func init() {
	cmdToolsView.PersistentFlags().StringVar(
		&toolsLibraryPath,
		"path",
		"",
		"path to tools library",
	)

	cmdTools.AddCommand(
		cmdToolsView,
		cmdToolsImport,
	)
}
