package gcode

import (
	"strconv"
	"strings"
)

var Precision = 2

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

// code is the underlying _thing_ that holds all the actual important bits.
// This is a common tool to represent any of the various codes (Z,G, etc). It
// can track what it will accept as a parameter, and most importantly, it can
// turn itself into GCode strings for consumption by an actual machine of some
// sort.
//
// It does this in a generic way without concern to the 'flavor' of gcode in
// use. There are a few non-standard character classes (that address: a, z, g,
// etc.). These are used to denote that we have a non-gcode item. Namely, the
// ';' address indicates we have a comment. The '_' address indicates a string.
//
type code struct {
	address byte
	value   float64
	codes   []code
	accepts []Addressable
	comment string
}

type mapKey struct {
	addr byte
	val  float64
}

func (c code) Address() byte      { return c.address }
func (c code) Arity() int         { return len(c.codes) }
func (c code) Value() float64     { return c.value }
func (c code) Parameters() []code { return c.codes }
func (c code) Key() mapKey {
	return mapKey{
		addr: c.address,
		val:  c.value,
	}
}

func (c code) String() string {
	var b strings.Builder

	switch c.address {
	case String.Address():
		b.WriteString(c.comment)
		return b.String()
	case Comment.Address():
		b.WriteString("; ")
		b.WriteString(c.comment)
		return b.String()
	case Finalizer.Address(), Initializer.Address():
		return ""
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

// Code creates a code func for the given address. For example: `banana :
// =Code('X')` returns a code function, `banana()` that can be used to
// represent an X code in the final output.
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

// Argument is similar to Code, but rather than have a function that takes a
// float64 as its value, it accepts a string. This way you can have commands
// like M400 that optionally have a string value to write to an LCD.
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
