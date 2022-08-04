package marlin

import (
	g "github.com/ryanfaerman/gophercnc/gcode"
)

var (
	A = g.A
	B = g.B
	C = g.C
	D = g.D
	E = g.E
	F = g.F
	G = g.G
	H = g.H
	I = g.I
	J = g.J
	K = g.K
	L = g.L
	M = g.M
	N = g.M
	O = g.O
	P = g.P
	Q = g.Q
	R = g.R
	S = g.S
	T = g.T
	U = g.U
	V = g.V
	W = g.W
	X = g.X
	Y = g.Y
	Z = g.Z

	String  = g.String
	Comment = g.Comment
	Flag    = String

	G0  = g.Command(G, 0.0, E, F, X, Y, Z)                // Rapid move
	G1  = g.Command(G, 1.0, E, F, X, Y, Z)                // Linear move
	G2  = g.Command(G, 2.0, E, F, I, J, P, R, S, X, Y, Z) // Clockwise arc
	G3  = g.Command(G, 3.0, E, F, I, J, P, R, S, X, Y, Z) // Counter clockwise arc
	G4  = g.Command(G, 4.0, P, S)                         // Dwell
	G28 = g.Command(G, 28.0, L, O, R, X, Y, Z, Flag)      // Auto home
	G90 = g.Command(G, 90.0)                              // Absolute positioning
	G91 = g.Command(G, 91.0)                              // Relative positioning
	G92 = g.Command(G, 92.0, E, X, Y, Z)                  // Set position

	M0   = g.Command(M, 0.0, P, S, String) // Unconditional stop
	M117 = g.Command(M, 117.0, String)     // Set LCD message
	M119 = g.Command(M, 119.0)             // Endstop states
	M400 = g.Command(M, 400.0)             // Finish moves
	M600 = g.Command(M, 600, T)            // Filament change

)
