package main

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/ryanfaerman/gophercnc/config"
	"github.com/ryanfaerman/gophercnc/gcode"
	. "github.com/ryanfaerman/gophercnc/gcode/marlin"
	"github.com/ryanfaerman/gophercnc/tool"
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
			active, err := config.ActiveLibrary()
			if err != nil {
				return err
			}
			lib, err := tool.LoadLibrary(active.Path)
			if err != nil {
				return err
			}

			fmt.Println(lib)

			doc.AddHook(M600, func(c gcode.CodePoint) gcode.Fragment {
				fmt.Println("borsht")

				// fmt.Println(c.String())

				// spw := spew.ConfigState{ContinueOnMethod: true}
				// spw.Dump(c)

				t := 0
				for _, p := range c.Parameters() {
					fmt.Println("changing!:", p.String(), "tool:", p.Value())
					t = int(p.Value())
				}

				tool, err := lib.FindByNumber(t)
				if err != nil {
					panic(err.Error())
				}

				return gcode.Fragment{
					M0(String(fmt.Sprintf("Change to: (%d) %s", tool.Number(), tool.Description))),
					G0(X(900.0)),
					G0(Y(90.8)),
					M0(String("End tool change")),
				}
			})

			setup := `
			; Hello, this is my gcode file, there are many like it, 
			; but this is mine
			G92
			G28
			M117 Let's do this thing
			G0 X12.3 Y45.6 Z19
			`

			doc.Initalizer(func(_ gcode.CodePoint) gcode.Fragment {
				frag, _ := gcode.ParseString(setup)
				return frag
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
	cmdCodeParse = &cobra.Command{
		Use:   "parse",
		Short: "code stuff",
		RunE: func(_ *cobra.Command, _ []string) error {
			// const src = `
			// ; COMMAND_SPINDLE_CLOCKWISE
			// M0 Turn ON 5000RPM
			// ; COMMAND_COOLANT_ON
			// M117  Bore3 ; display message
			// G0 Z15
			// G0 X-24.237 Y0.635 F2500
			// G0 Z1 F300
			// `
			const src = `
			;COMMAND_SPINDLE_CLOCKWISE
M400 Turn On 5000RPM ; Do the thing

M28 X


G0 X-24.237 Y0.635 F2500 ; Needful

`

			out, err := gcode.Parse(strings.NewReader(src))
			fmt.Println("--- OUTPUT ---")
			fmt.Println(out)

			return err

		},
	}
)

func init() {
	cmdCode.PersistentFlags().IntVarP(&codeToolNumber, "tool", "t", codeToolNumber, "tool number to use")
	cmdCode.PersistentFlags().StringVarP(&codeToolLibrary, "library", "L", codeToolLibrary, "name of active tool library")
}
