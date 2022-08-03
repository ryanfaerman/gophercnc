package gcode

import (
	"os"
	"testing"
)

func TestCommandString(t *testing.T) {

	focus := os.Getenv("FOCUS")

	examples := map[string]struct {
		fn     commandFunc
		input  []code
		output string
	}{
		"niladic": {
			fn:     Command(G, 0.0),
			input:  []code{},
			output: "G0",
		},
		"niladic-comment": {
			fn:     Command(G, 0.0),
			input:  []code{Comment("hello there")},
			output: "G0 ; hello there",
		},
		"niladic-extras": {
			fn:     Command(G, 0.0),
			input:  []code{A(1.3)},
			output: "G0",
		},
		"ternary": {
			fn:     Command(G, 0.0, X, Y, Z),
			input:  []code{X(1.1), Y(2.2), Z(3.3)},
			output: "G0 X1.1 Y2.2 Z3.3",
		},
		"ternary-out-of-order": {
			fn:     Command(G, 0.0, X, Y, Z),
			input:  []code{Y(2.2), Z(3.3), X(1.1)},
			output: "G0 X1.1 Y2.2 Z3.3",
		},
		"ternary-extras": {
			fn:     Command(G, 0.0, X, Y, Z),
			input:  []code{Y(2.2), Z(3.3), X(1.1), A(5.5)},
			output: "G0 X1.1 Y2.2 Z3.3",
		},
		"monadic-string-nil": {
			fn:     Command(M, 400.0, String),
			input:  []code{},
			output: "M400",
		},
		"monadic-string": {
			fn:     Command(M, 400.0, String),
			input:  []code{String("hello there")},
			output: "M400 hello there",
		},
		"monadic-string-comment": {
			fn:     Command(M, 400.0, String),
			input:  []code{String("hello there"), Comment("general kenobi")},
			output: "M400 hello there ; general kenobi",
		},
		"monadic-comment-string": {
			fn:     Command(M, 400.0, String),
			input:  []code{Comment("general kenobi"), String("hello there")},
			output: "M400 hello there ; general kenobi",
		},
		"limited-within": {
			fn:     LimitCommand(Command(G, 4.0, X, Y), X(0.0), X(50.0)),
			input:  []code{X(5.55)},
			output: "G4 X5.55",
		},
		"limited-above": {
			fn:     LimitCommand(Command(G, 4.0, X, Y), X(0.0), X(50.0)),
			input:  []code{X(55.4)},
			output: "G4 X50",
		},
		"limited-below": {
			fn:     LimitCommand(Command(G, 4.0, X, Y), X(44.4), X(60.0)),
			input:  []code{X(10.0)},
			output: "G4 X44.4",
		},
	}

	for name, example := range examples {
		example := example

		t.Run(name, func(t *testing.T) {
			if focus != "" && name != focus {
				t.Skipf("example '%s' is out of focus", name)
			}

			actual := example.fn(example.input...)

			if actual.String() != example.output {
				t.Errorf("output is incorrect\nactual:\n%s\n\nexpected:\n%s", actual.String(), example.output)
			}

		})

	}

}

func TestLimitedCode(t *testing.T) {
	focus := os.Getenv("FOCUS")

	examples := map[string]struct {
		fn     codeFunc
		input  float64
		output string
	}{

		"code-above": {
			fn:     LimitCode(X, 10.0, 50.5),
			input:  6000.0,
			output: "X50.5",
		},
		"code-within": {
			fn:     LimitCode(X, 10.0, 50.5),
			input:  25.25,
			output: "X25.25",
		},
		"code-below": {
			fn:     LimitCode(X, 10.0, 50.5),
			input:  5.5,
			output: "X10",
		},
	}

	for name, example := range examples {
		example := example

		t.Run(name, func(t *testing.T) {
			if focus != "" && name != focus {
				t.Skipf("example '%s' is out of focus", name)
			}

			actual := example.fn(example.input)

			if actual.String() != example.output {
				t.Errorf("output is incorrect\nactual:\n%s\n\nexpected:\n%s", actual.String(), example.output)
			}

		})

	}

}
