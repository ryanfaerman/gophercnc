package main

import (
	"fmt"
	"math"

	"github.com/ryanfaerman/gophercnc/gcode"
	"github.com/spf13/cobra"
)

var (
	cmdFacing = &cobra.Command{
		Use:     "facing",
		Short:   "Generate facing toolpaths",
		Aliases: []string{"face"},
		RunE: func(_ *cobra.Command, args []string) error {

			doc := &gcode.Document{}
			doc.SetPosition(gcode.X(0), gcode.Y(0), gcode.Z(0))
			doc.Preamble = `
	Sometimes things happen, you know?
	Sometimes they don't.
	It's just the way the cookie crumbles
	Cause you never know what you're gonna get.
	`

			probeThickness := 0.8
			doc.Comment("*** Tool Measure BEGIN ***")
			doc.Pause("Attach Z Probe")
			doc.Home(gcode.Z())
			doc.SetPosition(gcode.Z(probeThickness))
			doc.LinearMove(gcode.Z(30.0), gcode.F(300))
			doc.Pause("Detach Z Probe")
			doc.Comment("*** Tool Measure END***")

			surfacingTool := gcode.Tool{
				Number:             12,
				Description:        "Surfacing Tool",
				Diameter:           24.5,
				StepoverPercentage: 35.0,
				DepthOfCut:         0.75,
				Speed:              35000,
			}

			// ballTool := Tool{
			// 	Number:             2,
			// 	Description:        "Ball Endmill",
			// 	Diameter:           6.35,
			// 	StepoverPercentage: 40.0,
			// 	DepthOfCut:         2,
			// 	Speed:              1000,
			// }

			doc.ChangeTool(surfacingTool)

			doc.Pause("Start Spindle")

			dimensionX := 560.0
			dimensionY := 540.0
			targetDepth := 0.75

			passWidth := surfacingTool.Diameter * (surfacingTool.StepoverPercentage / 100.0)

			// passes := math.Ceil(dimensionX / passWidth)
			doc.Comment(fmt.Sprintf("Pass width: %f mm", passWidth))

			currentDepth := 0.0
			// currentDepth = math.Min(targetDepth, surfacingTool.DepthOfCut)
			currentDepth = targetDepth

			/*
				doc.LCDMessage("Facing Y Along X")
				doc.LinearMove(gcode.X(0), gcode.Y(0), gcode.Z(-currentDepth))

				for k := passWidth; k <= dimensionX; k += math.Min(passWidth, dimensionX-k) {
					doc.LinearMove(gcode.Y(dimensionY), gcode.F(surfacingTool.Speed))
					doc.LinearMove(gcode.X(k), gcode.F(surfacingTool.Speed))
					doc.LinearMove(gcode.Y(0))
					if k == dimensionX {
						break
					}
					k += math.Min(passWidth, dimensionX-k)
					doc.LinearMove(gcode.X(k))
				}

				doc.LCDMessage("Returning to Origin")
				doc.LinearMove(gcode.Z(2))
			*/
			doc.LinearMove(gcode.X(0), gcode.Y(0), gcode.F(2*surfacingTool.Speed))
			doc.LinearMove(gcode.Z(-currentDepth))

			doc.LCDMessage("Facing X Along Y")
			for k := passWidth; k <= dimensionY; k += math.Min(passWidth, dimensionY-k) {
				doc.LinearMove(gcode.X(dimensionX), gcode.F(surfacingTool.Speed))
				doc.LinearMove(gcode.Y(k), gcode.F(surfacingTool.Speed))
				doc.LinearMove(gcode.X(0))
				if k == dimensionY {
					break
				}
				k += math.Min(passWidth, dimensionY-k)
				doc.LinearMove(gcode.Y(k))
			}
			doc.LinearMove(gcode.Z(2))
			doc.LinearMove(gcode.X(0), gcode.Y(0), gcode.F(2*surfacingTool.Speed))
			doc.LCDMessage("Stop Spindle")

			fmt.Println(doc)
			return nil
		},
	}
)
