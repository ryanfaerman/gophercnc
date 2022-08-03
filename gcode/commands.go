package gcode

var Precision = 2

// func G(v float64, params ...Code) Code {
// 	return Code{
// 		Address: 'G',
// 		Value:   v,
// 		Codes:   params,
// 		Arity:   len(params),
// 	}
// }

type ParameterList []Addressable

func (pl ParameterList) Names() []string {
	out := []string{}

	for _, p := range pl {
		out = append(out, string(p.Address()))
	}
	return out
}

func (pl ParameterList) Arity() int {
	return len(pl)
}

/*

Accepts indicates the parameters *and their order* that the command will use.

x = param('X')
y = param('Y')
Command('G', 0.0, x, y)

*/

type commandFunc func(...code) code

func (cf commandFunc) Address() byte { return cf().address }

func (cf commandFunc) Value() float64 {
	return cf().value
}

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

func Limit(c commandFunc, limits ...Parameter) commandFunc {
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

// G0 = Limit(G0, X(50), X(-50), F(34))
