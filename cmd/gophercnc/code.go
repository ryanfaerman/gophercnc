package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/ryanfaerman/gophercnc/gcode"
	. "github.com/ryanfaerman/gophercnc/gcode/marlin"
	"github.com/spf13/cobra"
)

var (
	codeToolNumber  = 0
	codeToolLibrary = ""
	cmdCode         = &cobra.Command{
		Use:   "code",
		Short: "code stuff",
		RunE: func(_ *cobra.Command, _ []string) error {

			fmt.Println(Comment("hello"))

			machine := struct {
				MinX float64
				MaxX float64
				MinY float64
				MaxY float64
				MinZ float64
				MaxZ float64
			}{
				MinX: 0.0,
				MaxX: 540.0,
				MinY: 0.0,
				MaxY: 540.0,
				MinZ: -40.0,
				MaxZ: 30.0,
			}
			spew.Dump(machine)

			X = gcode.LimitCode(X, machine.MinX, machine.MaxX)
			doc := &gcode.Document{}
			doc.AddLines(
				G92(X(0.0), Y(0.0), Z(0.0)),
				Comment("** Tool Measure BEGIN**"),
				M0(String("Attach Z Probe")),
				G28(Flag("Z")),
				M600(T(3)),
			)

			doc.AddHook(M600, func(c gcode.CodePoint) gcode.Fragment {
				fmt.Println("bananas in pajamas")
				spew.Dump(c.Address())
				spew.Dump(c.Arity())
				for _, p := range c.Parameters() {
					fmt.Println("changing!:", p.String(), "tool:", p.Value())
				}
				return gcode.Fragment{
					M0(String("Tool Change")),
					G0(X(900.0)),
					G0(Y(90.8)),
					M0(String("End tool change")),
				}
			})

			doc.Initalizer(func(_ gcode.CodePoint) gcode.Fragment {
				return nil
				return gcode.Fragment{
					Comment("This is a story, of a lovely lady!"),
				}
			})

			doc.Finalizer(func(_ gcode.CodePoint) gcode.Fragment {
				return gcode.Fragment{
					Comment("this was fun, right?"),
				}
			})

			fmt.Println("DOCUMENT START")
			fmt.Println(doc)
			return nil
		},
	}
)

func init() {
	cmdCode.PersistentFlags().IntVarP(&codeToolNumber, "tool", "t", codeToolNumber, "tool number to use")
	cmdCode.PersistentFlags().StringVarP(&codeToolLibrary, "library", "L", codeToolLibrary, "name of active tool library")
}
