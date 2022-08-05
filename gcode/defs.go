package gcode

var (
	A = Code('A')
	B = Code('B')
	C = Code('C')
	D = Code('D')
	E = Code('E')
	F = Code('F')
	G = Code('G')
	H = Code('H')
	I = Code('I')
	J = Code('J')
	K = Code('K')
	L = Code('L')
	M = Code('M')
	N = Code('M')
	O = Code('O')
	P = Code('P')
	Q = Code('Q')
	R = Code('R')
	S = Code('S')
	T = Code('T')
	U = Code('U')
	V = Code('V')
	W = Code('W')
	X = Code('X')
	Y = Code('Y')
	Z = Code('Z')

	String  = Argument('_')
	Comment = Argument(';')
)

var (
	Initializer = code{
		address: '!',
	}

	Finalizer = code{
		address: '$',
	}
)
