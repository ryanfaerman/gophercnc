package main

import (
	"fmt"
	"math"

	"github.com/ryanfaerman/gophercnc/gcode"
	"github.com/spf13/cobra"
)

var (
	cmdSuction = &cobra.Command{
		Use:     "suction",
		Short:   "Generate facing toolpaths",
		Aliases: []string{"suck"},
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

			// surfacingTool := gcode.Tool{
			// 	Number:             12,
			// 	Description:        "Surfacing Tool",
			// 	Diameter:           24.5,
			// 	StepoverPercentage: 35.0,
			// 	DepthOfCut:         0.75,
			// 	Speed:              35000,
			// }

			ballTool := gcode.Tool{
				Number:             2,
				Description:        "Ball Endmill",
				Diameter:           6.35,
				StepoverPercentage: 40.0,
				DepthOfCut:         2,
				Speed:              2000,
			}

			activeTool := ballTool

			doc.ChangeTool(ballTool)

			doc.Pause("Start Spindle")

			dimensionX := 560.0
			dimensionY := 540.0
			targetDepth := 2.0

			// passWidth := surfacingTool.Diameter * (surfacingTool.StepoverPercentage / 100.0)

			// passWidth := ballTool.Diameter/2 + 20

			passWidth := roundFloat(dimensionX/(ballTool.Diameter/2+20), 2)
			passHeight := roundFloat(dimensionY/(ballTool.Diameter/2+20), 2)

			// passes := math.Ceil(dimensionX / passWidth)
			doc.Comment(fmt.Sprintf("Pass width: %f mm", passWidth))

			currentDepth := 0.0
			// currentDepth = math.Min(targetDepth, surfacingTool.DepthOfCut)
			currentDepth = targetDepth

			doc.LCDMessage("Cutting Y Along X")
			doc.LinearMove(gcode.X(0), gcode.Y(0), gcode.Z(-currentDepth))

			for k := passWidth; k <= dimensionX; k += math.Min(passWidth, dimensionX-k) {
				doc.LinearMove(gcode.Z(-currentDepth))
				doc.LinearMove(gcode.Y(dimensionY), gcode.F(activeTool.Speed))
				doc.LinearMove(gcode.Z(1.0))
				doc.LinearMove(gcode.X(k), gcode.F(activeTool.Speed))
				doc.LinearMove(gcode.Z(-currentDepth))
				doc.LinearMove(gcode.Y(0))
				if k == dimensionX || (dimensionX-k < passWidth) {
					break
				}
				k += math.Min(passWidth, dimensionX-k)

				doc.LinearMove(gcode.Z(1.0))
				doc.LinearMove(gcode.X(k))
			}

			doc.LCDMessage("Returning to Origin")
			doc.LinearMove(gcode.Z(2))
			doc.LinearMove(gcode.X(0), gcode.Y(0), gcode.F(2*activeTool.Speed))
			doc.LinearMove(gcode.Z(-currentDepth))

			doc.LCDMessage("Facing X Along Y")
			for k := passHeight; k <= dimensionY; k += math.Min(passHeight, dimensionY-k) {
				doc.LinearMove(gcode.Z(-currentDepth))
				doc.LinearMove(gcode.X(dimensionX), gcode.F(activeTool.Speed))
				doc.LinearMove(gcode.Z(1.0))
				doc.LinearMove(gcode.Y(k), gcode.F(activeTool.Speed))
				doc.LinearMove(gcode.Z(-currentDepth))
				doc.LinearMove(gcode.X(0))
				if k == dimensionY || (dimensionY-k < passHeight) {
					break
				}
				k += math.Min(passHeight, dimensionY-k)
				doc.LinearMove(gcode.Z(1.0))
				doc.LinearMove(gcode.Y(k))
			}
			doc.LinearMove(gcode.Z(2))
			doc.LinearMove(gcode.X(0), gcode.Y(0), gcode.F(2*activeTool.Speed))
			doc.LCDMessage("Stop Spindle")

			fmt.Println(doc)
			return nil
		},
	}
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
