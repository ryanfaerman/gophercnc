package gcode

type Mode int

const (
	ModeUnknown Mode = iota
	ModeAbsolute
	ModeRelative
	ModeInches
	ModeMillimeters
)

func SetMode(m Mode) string {
	switch m {
	case ModeAbsolute:
		return "G90"
	case ModeRelative:
		return "G91"
	case ModeMillimeters:
		return "G21"
	case ModeInches:
		return "G20"
	default:
		panic("unknown mode")
	}
}

func (m Mode) String() string {
	return SetMode(m)
}
