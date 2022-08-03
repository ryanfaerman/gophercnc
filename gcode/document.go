package gcode

import "strings"

type Document []code

func (d Document) String() string {
	var out strings.Builder
	for _, c := range d {
		out.WriteString(c.String())
		out.WriteString("\n")

	}
	return out.String()
}
