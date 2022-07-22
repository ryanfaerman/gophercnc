package gcode

import (
	"fmt"
	"strings"
)

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
