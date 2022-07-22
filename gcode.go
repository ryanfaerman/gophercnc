package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Parameter struct {
	Name    string
	Value   float64
	IsNull  bool
	Maximum float64
	Minimum float64
}

func (a Parameter) Export(precision int) string {
	if a.IsNull {
		return a.Name
	}
	x := strconv.FormatFloat(a.Value, 'f', precision, 64)

	// Hacky way to remove silly zeroes
	if strings.IndexRune(x, '.') != -1 {
		for x[len(x)-1] == '0' {
			x = x[:len(x)-1]
		}
		if x[len(x)-1] == '.' {
			x = x[:len(x)-1]
		}
	}

	return a.Name + x
}

func NewParameter(n string, v ...float64) Parameter {
	p := Parameter{
		Name: n,
	}
	if len(v) == 0 {
		p.IsNull = true
	}
	if len(v) > 0 {
		p.Value = v[0]
	}
	return p

}

func (a Parameter) String() string {
	return a.Export(3)
}

func X(v ...float64) Parameter { return NewParameter("X", v...) }
func Y(v ...float64) Parameter { return NewParameter("Y", v...) }
func Z(v ...float64) Parameter { return NewParameter("Z", v...) }
func E(v ...float64) Parameter { return NewParameter("E", v...) }
func F(v ...float64) Parameter { return NewParameter("F", v...) }

type Mode int

const (
	ModeUnknown Mode = iota
	ModeAbsolute
	ModeRelative
	ModeInches
	ModeMillimeters
)

func SetMode(m Mode) string {
	switch m {
	case ModeAbsolute:
		return "G90"
	case ModeRelative:
		return "G91"
	case ModeMillimeters:
		return "G21"
	case ModeInches:
		return "G20"
	default:
		panic("unknown mode")
	}
}

func (m Mode) String() string {
	return SetMode(m)
}

// func Positioning(m Mode) string { return SetMode(m) }
// func Units(m Mode) string       { return SetMode(m) }

func processParams(params []Parameter, names ...string) string {
	var out strings.Builder

	for _, n := range names {
		for _, p := range params {
			if p.Name != n {
				continue
			}

			out.WriteString(p.String())
			out.WriteString(" ")
		}
	}

	return out.String()
}

func SetPosition(params ...Parameter) string {
	return "G92 " + processParams(params, "E", "X", "Y", "Z")
}

type Tool struct {
	Number             int
	Description        string
	Diameter           float64
	StepoverPercentage float64
	DepthOfCut         float64
	FeedRate           float64
	Speed              float64
}

type Document struct {
	Instructions []string
	Units        Mode
	Positioning  Mode
	Preamble     string

	Tools       []Tool
	currentTool Tool
}

func (doc *Document) Append(i string) {
	doc.Instructions = append(doc.Instructions, i)
}

func (doc *Document) Prepend(i string) {
	doc.Instructions = append([]string{i}, doc.Instructions...)
}

func (doc *Document) String() string {
	var out strings.Builder

	if doc.Units == ModeUnknown {
		doc.Units = ModeMillimeters
	}
	if doc.Positioning == ModeUnknown {
		doc.Positioning = ModeAbsolute
	}

	out.WriteString(stringToComment(doc.Preamble))
	out.WriteString("\n")
	out.WriteString(doc.Units.String())
	out.WriteString("\n")
	out.WriteString(doc.Positioning.String())
	out.WriteString("\n")
	out.WriteString(strings.Join(doc.Instructions, "\n"))

	return out.String()
}

func (doc *Document) SetPositioning(m Mode) { doc.Append(SetMode(m)) }
func (doc *Document) SetUnits(m Mode)       { doc.Append(SetMode(m)) }
func (doc *Document) Comment(m string)      { doc.Append(stringToComment(m)) }
func (doc *Document) SetPosition(params ...Parameter) {
	doc.Append("G92 " + processParams(params, "E", "X", "Y", "Z"))
}
func (doc *Document) Pause(message string) { doc.Append("M0 " + message) }
func (doc *Document) Home(params ...Parameter) {
	doc.Append("G28 " + processParams(params, "X", "Y", "Z"))
}

func (doc *Document) LinearMove(params ...Parameter) {
	doc.Append("G0 " + processParams(params, "X", "Y", "Z", "F"))
}

func (doc *Document) LCDMessage(msg string) { doc.Append("M117 " + msg) }

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

func (doc *Document) ChangeTool(t Tool) {
	if t == doc.currentTool {
		return
	}
	doc.Tools = append(doc.Tools, t)
	doc.currentTool = doc.Tools[len(doc.Tools)-1]
	if len(doc.Tools) > 1 {
		doc.Comment("COMMENCE TOOL CHANGE PROCEDURE")
		doc.Comment(fmt.Sprintf("Change to: %d", t.Number))
	}
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
