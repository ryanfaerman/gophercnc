package tool

import (
	"encoding/json"
	"strings"
)

//go:generate stringer -linecomment -type=Endmill
type Endmill int

const (
	EndmillUnknown        Endmill = iota
	EndmillBall                   // Ball end mill
	EndmillBullNose               // Bull nose end mill
	EndmillEngraveChamfer         // Engrave/Chamfer mill
	EndmillDovetail               // Dovetail mill
	EndmillFace                   // Face mill
	EndmillFlat                   // Flat end mill
	EndmillForm                   // Form mill
	EndmillLollipop               // Lollipop mill
	EndmillRadius                 // Radius mill
	EndmillSlot                   // Slot mill
	EndmillTapered                // Tapered mill
	EndmillThread                 // Thread mill
)

func (t *Endmill) FromString(value string) Endmill {
	switch strings.ToLower(value) {
	case "ball end mill":
		return EndmillBall
	case "bull nose end mill":
		return EndmillBullNose
	case "engrave/chamfer mill":
		return EndmillEngraveChamfer
	case "dovetail":
		return EndmillDovetail
	case "face mill":
		return EndmillFace
	case "flat end mill":
		return EndmillFlat
	case "form mill":
		return EndmillForm
	case "lollipop mill":
		return EndmillLollipop
	case "radius mill":
		return EndmillRadius
	case "slot mill":
		return EndmillSlot
	case "tapered mill":
		return EndmillTapered
	case "thread":
		return EndmillTapered
	default:
		return EndmillUnknown
	}
}

func (t Endmill) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Endmill) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*t = t.FromString(s)

	return nil
}
