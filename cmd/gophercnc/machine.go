package main

import (
	"fmt"
	"sort"

	"github.com/davecgh/go-spew/spew"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/ryanfaerman/gophercnc/config"
	"github.com/ryanfaerman/gophercnc/gcode"
	. "github.com/ryanfaerman/gophercnc/gcode/marlin"
	"github.com/ryanfaerman/gophercnc/machine"
	"github.com/ryanfaerman/gophercnc/tool"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	cmdMachine = &cobra.Command{
		Use:   "machine",
		Short: "machine stuff",
		RunE: func(_ *cobra.Command, _ []string) error {
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
			fmt.Println(len(lib.Tools))

			if err != nil {
				return err
			}
			// m := machine.Machine{
			// 	Procedures: map[string]string{
			// 		"M600": `
			// 			M117 hello there
			// 			G28 X
			// 			G92 X0.34
			// 		`,
			// 	},
			// }

			m, err := machine.LoadMachine("/Users/rfaerman/repos/github.com/ryanfaerman/gophercnc/mpcnc.yml")
			if err != nil {
				return err
			}
			spew.Dump(m)

			y, err := yaml.Marshal(m)
			if err != nil {
				return err
			}
			fmt.Println(string(y))

			spw := spew.ConfigState{ContinueOnMethod: true}

			spw.Dump(gcode.ParseString("M600"))

			for k, p := range m.Procedures {
				cp, err := gcode.ParseString(k)
				if err != nil {
					return err
				}
				fmt.Println("kkk")
				spw.Dump(cp[0])
				spw.Dump(cp[0].Key())

				frag, err := gcode.ParseString(p)
				if err != nil {
					return err
				}

				doc.AddHook(cp[0], func(_ gcode.CodePoint) gcode.Fragment {
					return frag
				})
			}

			fmt.Println("--- OUTPUT ---")
			fmt.Println(doc)

			return nil
		},
	}
	cmdMachineImport = &cobra.Command{
		Use:     "import PATH",
		Short:   "load a machine",
		Aliases: []string{"add"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var err error
			path := args[0]

			path, err = homedir.Expand(path)
			if err != nil {
				return err
			}

			return machine.Import(path)
		},
	}

	cmdMachines = &cobra.Command{
		Use:    "machines",
		Short:  cmdMachineList.Short,
		RunE:   cmdMachineList.RunE,
		Hidden: true,
	}

	cmdMachineList = &cobra.Command{
		Use:     "available",
		Short:   "show all tool libraries",
		Aliases: []string{"list"},
		RunE: func(_ *cobra.Command, _ []string) error {

			type row struct {
				Name   string
				Active bool
			}

			active, err := config.ActiveMachine()
			if err != nil {
				return nil
			}

			rows := []row{}

			libs, err := config.Machines()
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

	cmdMachineActivate = &cobra.Command{
		Use:     "activate",
		Short:   "mark a machine as active",
		Aliases: []string{"use"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return config.ActivateMachine(args[0])
		},
	}
)

func init() {
	root.AddCommand(cmdMachines)
	cmdMachine.AddCommand(
		cmdMachineList,
		cmdMachineImport,
		cmdMachineActivate,
	)
}
