package main

import (
	"fmt"

	"github.com/ryanfaerman/gophercnc/gcode"
	. "github.com/ryanfaerman/gophercnc/gcode/marlin"
	"github.com/spf13/cobra"
)

var (
	cmdCode = &cobra.Command{
		Use:   "code",
		Short: "code stuff",
		RunE: func(_ *cobra.Command, _ []string) error {

			fmt.Println(Comment("hello"))

			doc := gcode.Document{
				Comment("hello there"),
				G1(X(45.5), F(1555.55)),
				G1(Z(15.4)),
			}

			fmt.Println(doc)
			return nil
		},
	}
)
