// Code generated by "stringer -linecomment -type=Material"; DO NOT EDIT.

package tool

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MaterialUnknown-0]
	_ = x[MaterialHSS-1]
	_ = x[MaterialCarbide-2]
	_ = x[MaterialTICoated-3]
	_ = x[MaterialCeramics-4]
}

const _Material_name = "UnspecifiedHSSCarbideTI CoatedCeramics"

var _Material_index = [...]uint8{0, 11, 14, 21, 30, 38}

func (i Material) String() string {
	if i < 0 || i >= Material(len(_Material_index)-1) {
		return "Material(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Material_name[_Material_index[i]:_Material_index[i+1]]
}
