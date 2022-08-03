package gcode

type commandFunc func(...code) code

func (cf commandFunc) Address() byte { return cf().address }

func (cf commandFunc) Value() float64 {
	return cf().value
}

// Command is a slightly higher order GCode, that is, one that can accept
// arguments/parameters. This function will generate a function for creating
// functions (trippy, I know), that allows you to define what commands are
// available to your gcode program.
//
// Why would you even go down this route? This lets the compiler help enforce
// the available commands and it also lets us add constraints rather than it
// just being a simple string.
//
// The command returned by this function will only accept values that are
// provided and will consume them in the order defined in the argument list.
// That is to say, if you want to accept Z Y into your M9001 command, you can
// with `Command(M, 9001.0, Z, Y)`.
//
// The address is the GCode underwhich this command lives(M, G, etc.).
func Command(address Addressable, v float64, accepts ...Addressable) commandFunc {
	accepts = append(accepts, code{address: ';'})
	return func(params ...code) code {
		validParams := []code{}

		for _, allowed := range accepts {
			for _, p := range params {
				if p.address != allowed.Address() {
					continue
				}

				validParams = append(validParams, p)
			}
		}

		return code{
			address: address.Address(),
			value:   v,
			codes:   validParams,
			accepts: accepts,
		}

	}

}

var PanicOnLimitFailure bool

// LimitCommand wraps the given commandFunc and applies a series of limits to
// its inputs. The limits are expected to be in tupes (min, max) and share the
// same address (X, Y, etc.). If the address is not the same or the count of
// the limits in not even either:
// - the error is ignored and all limits are ignored
// - a panic is triggered
//
// This is controlled by the PanicOnLimitFailure variable.
func LimitCommand(c commandFunc, limits ...Parameter) commandFunc {
	if len(limits)%2 == 1 {
		if PanicOnLimitFailure {
			panic("must have even number of limits")
		}
		// TODO: log a warning that limits aren't applied
		return c
	}

	limitations := map[byte][]float64{}

	for len(limits) >= 2 {
		var (
			min Parameter
			max Parameter
		)
		min, limits = limits[0], limits[1:]
		max, limits = limits[0], limits[1:]

		if min.Address() != max.Address() {
			if PanicOnLimitFailure {
				panic("min/max are not on the same address")
			}
			// TODO: Log a warning that our limits aren't applied
			return c
		}

		limitations[min.Address()] = []float64{min.Value(), max.Value()}

	}

	return func(params ...code) code {
		for i, param := range params {
			if limit, ok := limitations[param.address]; ok {
				if param.value < limit[0] {
					params[i].value = limit[0]
				}
				if param.value > limit[1] {
					params[i].value = limit[1]
				}
			}
		}

		return c(params...)
	}
}

// LimitCode adds a limiter  to the possible values of a code. Any code wrapped
// with this function can only have values between the min and max.
func LimitCode(c codeFunc, min float64, max float64) codeFunc {
	return func(value float64) code {
		if value < min {
			value = min
		}
		if value > max {
			value = max
		}

		return c(value)
	}
}

// G0 = Limit(G0, X(50), X(-50), F(34))
