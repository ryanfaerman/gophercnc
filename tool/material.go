package tool

import (
	"encoding/json"
	"strings"
)

//go:generate stringer -linecomment -type=Material
type Material int

const (
	MaterialUnknown  Material = iota // Unspecified
	MaterialHSS                      // HSS
	MaterialCarbide                  // Carbide
	MaterialTICoated                 // TI Coated
	MaterialCeramics                 // Ceramics
)

func (t *Material) FromString(value string) Material {
	switch strings.ToLower(value) {
	case "hss":
		return MaterialHSS
	case "carbide":
		return MaterialCarbide
	case "ti coated":
		return MaterialTICoated
	case "ceramics":
		return MaterialCeramics
	default:
		return MaterialHSS
	}
}

func (t Material) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Material) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*t = t.FromString(s)

	return nil
}
