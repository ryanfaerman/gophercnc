/*
jsonlines is a wrapper around encoding.JSON that marshals data into the JSON
Lines format. See: http://jsonlines.org/
*/
package jsonlines

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// Marshal receives an arbitrary input (usually a slice of some values) and
// marshals each entry, one per line, into JSON seperated by a new line
// character.
func Marshal(input interface{}) ([]byte, error) {
	switch reflect.TypeOf(input).Kind() {
	case reflect.Slice, reflect.Array:
		var out bytes.Buffer
		val := reflect.ValueOf(input)
		for i := 0; i < val.Len(); i++ {
			d, err := json.Marshal(val.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			out.Write(d)
			out.WriteString("\n")
		}
		return out.Bytes(), nil
	default:
		return json.Marshal(input)
	}
}
