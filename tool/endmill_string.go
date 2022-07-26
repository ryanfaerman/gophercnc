// Code generated by "stringer -linecomment -type=Endmill"; DO NOT EDIT.

package tool

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EndmillUnknown-0]
	_ = x[EndmillBall-1]
	_ = x[EndmillBullNose-2]
	_ = x[EndmillEngraveChamfer-3]
	_ = x[EndmillDovetail-4]
	_ = x[EndmillFace-5]
	_ = x[EndmillFlat-6]
	_ = x[EndmillForm-7]
	_ = x[EndmillLollipop-8]
	_ = x[EndmillRadius-9]
	_ = x[EndmillSlot-10]
	_ = x[EndmillTapered-11]
	_ = x[EndmillThread-12]
}

const _Endmill_name = "EndmillUnknownBall end millBull nose end millEngrave/Chamfer millDovetail millFace millFlat end millForm millLollipop millRadius millSlot millTapered millThread mill"

var _Endmill_index = [...]uint8{0, 14, 27, 45, 65, 78, 87, 100, 109, 122, 133, 142, 154, 165}

func (i Endmill) String() string {
	if i < 0 || i >= Endmill(len(_Endmill_index)-1) {
		return "Endmill(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Endmill_name[_Endmill_index[i]:_Endmill_index[i+1]]
}
