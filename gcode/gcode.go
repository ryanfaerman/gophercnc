package gcode

import (
	"strconv"
	"strings"
)

type Addressable interface {
	Address() byte
}

type Valuable interface {
	Value() float64
}

type Parameter interface {
	Addressable
	Valuable
}

type code struct {
	address byte
	value   float64
	codes   []code
	accepts []Addressable
	comment string
}

func (c code) Address() byte  { return c.address }
func (c code) Arity() int     { return len(c.codes) }
func (c code) Value() float64 { return c.value }

func (c code) String() string {
	var b strings.Builder

	switch c.address {
	case '_':
		b.WriteString(c.comment)
		return b.String()
	case ';':
		b.WriteString("; ")
		b.WriteString(c.comment)
		return b.String()
	default:
		b.WriteByte(c.address)

	}

	v := strconv.FormatFloat(c.value, 'f', Precision, 64)

	// Hacky way to remove silly zeroes
	if strings.IndexRune(v, '.') != -1 {
		for v[len(v)-1] == '0' {
			v = v[:len(v)-1]
		}
		if v[len(v)-1] == '.' {
			v = v[:len(v)-1]
		}
	}
	b.WriteString(v)

	if c.Arity() > 0 {
		for _, param := range c.codes {
			b.WriteString(" ")
			b.WriteString(param.String())
		}
	}

	return b.String()
}

func Code(address byte) codeFunc {
	return func(v float64) code {
		return code{
			address: address,
			value:   v,
		}
	}
}

type stringFunc func(string) code

func (sf stringFunc) Address() byte {
	return sf("").address
}

func Argument(address byte) stringFunc {
	return func(s string) code {
		return code{
			address: address,
			comment: s,
		}
	}
}

type codeFunc func(float64) code

func (cf codeFunc) Address() byte {
	return cf(0.0).address
}
