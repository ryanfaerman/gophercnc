package gcode

import (
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

func (a Parameter) String() string {
	return a.Export(3)
}

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
