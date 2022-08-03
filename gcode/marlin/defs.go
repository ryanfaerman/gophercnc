package marlin

import (
	g "github.com/ryanfaerman/gophercnc/gcode"
)

var (
	A = g.Code('A')
	B = g.Code('B')
	C = g.Code('C')
	D = g.Code('D')
	E = g.Code('E')
	F = g.Code('F')
	G = g.Code('G')
	H = g.Code('H')
	I = g.Code('I')
	J = g.Code('J')
	K = g.Code('K')
	L = g.Code('L')
	M = g.Code('M')
	N = g.Code('M')
	O = g.Code('O')
	P = g.Code('P')
	Q = g.Code('Q')
	R = g.Code('R')
	S = g.Code('S')
	T = g.Code('T')
	U = g.Code('U')
	V = g.Code('V')
	W = g.Code('W')
	X = g.Code('X')
	Y = g.Code('Y')
	Z = g.Code('Z')

	String  = g.Argument('_')
	Comment = g.Argument(';')

	G1 = g.Command(G, 0.0, X, Y, Z)
)
