package gcode

import (
	"fmt"
	"math"
	"strings"
)

func X(v ...float64) Parameter { return NewParameter("X", v...) }
func Y(v ...float64) Parameter { return NewParameter("Y", v...) }
func Z(v ...float64) Parameter { return NewParameter("Z", v...) }
func E(v ...float64) Parameter { return NewParameter("E", v...) }
func F(v ...float64) Parameter { return NewParameter("F", v...) }

// func Positioning(m Mode) string { return SetMode(m) }
// func Units(m Mode) string       { return SetMode(m) }

func SetPosition(params ...Parameter) string {
	return "G92 " + processParams(params, "E", "X", "Y", "Z")
}

func stringToComment(m string) string {
	var out strings.Builder
	lines := strings.Split(m, "\n")
	for _, l := range lines {
		l := strings.TrimSpace(l)
		if l == "" {
			continue
		}
		out.WriteString("; " + strings.TrimSpace(l))
		if len(lines) > 1 {
			out.WriteString("\n")
		}
	}

	return out.String()

}

func main() {

	// fmt.Println(Positioning(ModeAbsolute))
	// fmt.Println(Z(0.8))

	// fmt.Println(SetPosition(X(0), Y(0), Z(0)))

	doc := &Document{}
	doc.SetPosition(X(0), Y(0), Z(0))
	doc.Preamble = `
	Sometimes things happen, you know?
	Sometimes they don't.
	It's just the way the cookie crumbles
	Cause you never know what you're gonna get.
	`

	probeThickness := 0.8
	doc.Comment("*** Tool Measure BEGIN ***")
	doc.Pause("Attach Z Probe")
	doc.Home(Z())
	doc.SetPosition(Z(probeThickness))
	doc.LinearMove(Z(40.0), F(300))
	doc.Pause("Detach Z Probe")
	doc.Comment("*** Tool Measure END***")

	surfacingTool := Tool{
		Number:             12,
		Description:        "Surfacing Tool",
		Diameter:           24.5,
		StepoverPercentage: 65.0,
		DepthOfCut:         0.5,
		Speed:              25000,
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
	targetDepth := 0.5

	passWidth := doc.currentTool.Diameter * (doc.currentTool.StepoverPercentage / 100.0)

	// passes := math.Ceil(dimensionX / passWidth)
	doc.Comment(fmt.Sprintf("Pass width: %f mm", passWidth))

	currentDepth := 0.0
	currentDepth = math.Min(targetDepth, doc.currentTool.DepthOfCut)

	doc.LCDMessage("Facing Y Along X")
	doc.LinearMove(X(0), Y(0), Z(-currentDepth))

	for k := passWidth; k <= dimensionX; k += math.Min(passWidth, dimensionX-k) {
		doc.LinearMove(Y(dimensionY), F(doc.currentTool.Speed))
		doc.LinearMove(X(k), F(doc.currentTool.Speed))
		doc.LinearMove(Y(0))
		if k == dimensionX {
			break
		}
		k += math.Min(passWidth, dimensionX-k)
		doc.LinearMove(X(k))
	}

	doc.LCDMessage("Returning to Origin")
	doc.LinearMove(Z(2))
	doc.LinearMove(X(0), Y(0), F(2*doc.currentTool.Speed))
	doc.LinearMove(Z(-currentDepth))

	doc.LCDMessage("Facing X Along Y")
	for k := passWidth; k <= dimensionY; k += math.Min(passWidth, dimensionY-k) {
		doc.LinearMove(X(dimensionX), F(doc.currentTool.Speed))
		doc.LinearMove(Y(k), F(doc.currentTool.Speed))
		doc.LinearMove(X(0))
		if k == dimensionY {
			break
		}
		k += math.Min(passWidth, dimensionY-k)
		doc.LinearMove(Y(k))
	}
	doc.LinearMove(Z(2))
	doc.LinearMove(X(0), Y(0), F(2*doc.currentTool.Speed))
	doc.LCDMessage("Stop Spindle")

	fmt.Println(doc)

}
