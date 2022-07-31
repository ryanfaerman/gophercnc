package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/structs"
)

// And Encoder writes tablular data to the output stream
type Encoder struct {
	headers []string
	rows    [][]string
}

// Encode creates the csv encoding of the input
//
// With structs, the following struct tags are supported using the `render` key:
//
//  - omitempty : ignores the column, if the value is the zero value
//  - noprefix  : prevents prefixing columns for nested structs with their field name
//
// Fields with the tag `-` are not rendered e.g. `render:"-"`.
func (e *Encoder) Encode(input interface{}) error {
	switch reflect.TypeOf(input).Kind() {
	case reflect.Slice, reflect.Array:
		return e.encodeSlice(input)
	case reflect.String, reflect.Int, reflect.Float32, reflect.Float64:
		e.rows = append(e.rows, []string{fmt.Sprint(input)})
	case reflect.Map:
		return e.encodeMap(input)
	case reflect.Struct:
		headers, row := e.encodeStruct(input, nil)
		e.headers = headers
		e.rows = append(e.rows, row)
	default:
		fmt.Println("yoyo")
		fmt.Println(reflect.TypeOf(input).Kind())

	}

	return nil
}

func (e *Encoder) encodeSlice(input interface{}) error {
	kind := reflect.TypeOf(input).Kind()
	if kind != reflect.Slice {
		panic("encodeSlice can only encode an input of type Slice")
	}

	// TODO: iterate over every map in the slice, getting a unique list of
	// available keys, then using that to drive the headers and map encoding

	row := []string{}
	val := reflect.ValueOf(input)

	mapKeys := []string{}
	if val.Len() > 0 && val.Index(0).Kind() == reflect.Map {
		keys := map[string]bool{}
		for i := 0; i < val.Len(); i++ {
			if val.Index(i).Kind() != reflect.Map {
				return errors.New("mixed-slice with map")
			}

			iter := reflect.ValueOf(val.Index(i).Interface()).MapRange()

			for iter.Next() {
				entry := fmt.Sprint(iter.Key())
				if _, ok := keys[entry]; !ok {
					keys[entry] = true
					mapKeys = append(mapKeys, entry)
				}
			}
		}
		sort.Strings(mapKeys)
	}

	for i := 0; i < val.Len(); i++ {
		switch val.Index(i).Type().Kind() {
		case reflect.String, reflect.Int, reflect.Float32, reflect.Float64, reflect.Interface:
			row = append(row, fmt.Sprint(val.Index(i).Interface()))
		case reflect.Slice, reflect.Array:
			e.encodeSlice(val.Index(i).Interface())
		case reflect.Map:
			e.encodeMap(val.Index(i).Interface(), mapKeys...)

		case reflect.Struct:
			h, r := e.encodeStruct(val.Index(i).Interface(), nil)
			e.headers = h
			e.rows = append(e.rows, r)

		}
	}
	if len(row) != 0 {
		e.rows = append(e.rows, row)
	}
	return nil
}

func (e *Encoder) encodeMap(input interface{}, keys ...string) error {
	iter := reflect.ValueOf(input).MapRange()
	data := map[string]string{}

	findKeys := len(keys) == 0

	for iter.Next() {
		data[fmt.Sprint(iter.Key())] = join(iter.Value(), ", ")
		if findKeys {
			keys = append(keys, fmt.Sprint(iter.Key()))
		}
	}

	row := []string{}

	sort.Strings(keys)

	for _, key := range keys {
		d, ok := data[key]
		if !ok {
			d = ""
		}
		row = append(row, d)
	}

	e.headers = keys
	e.rows = append(e.rows, row)

	return nil
}

func (e *Encoder) encodeStruct(input interface{}, parent *structs.Field) ([]string, []string) {
	headers := []string{}
	rows := []string{}

	s := structs.New(input)
	for _, field := range s.Fields() {
		if !field.IsExported() {
			continue
		}
		name := field.Name()

		tagName, tagOpts := parseTag(field.Tag("render"))
		if tagName != "" {
			if tagName == "-" {
				continue
			}
			name = tagName
		}

		if parent != nil {
			parentName := parent.Name()
			parentTagName, parentTagOpts := parseTag(parent.Tag("render"))
			if parentTagName != "" {
				parentName = parentTagName
			}
			if !tagOpts.Has("omitempty") {
				if field.IsZero() {
					rows = append(rows, "")
					continue
				}
			}
			if !parentTagOpts.Has("noprefix") {
				name = fmt.Sprintf("%s.%s", parentName, name)
			}
		}

		if tagOpts.Has("omitempty") {
			if field.IsZero() {
				// rows = append(rows, "")
				continue
			}
		}

		switch reflect.TypeOf(field.Value()).Kind() {
		case reflect.Struct:
			h, r := e.encodeStruct(field.Value(), field)
			headers = append(headers, h...)
			rows = append(rows, r...)
		case reflect.Map:
			iter := reflect.ValueOf(field.Value()).MapRange()
			data := map[string]string{}
			keys := []string{}

			for iter.Next() {
				data[fmt.Sprint(iter.Key())] = join(iter.Value(), ", ")
				keys = append(keys, fmt.Sprint(iter.Key()))
			}

			row := []string{}

			sort.Strings(keys)
			h := []string{}

			for _, key := range keys {

				kh := key

				if !tagOpts.Has("noprefix") {
					kh = fmt.Sprintf("%s.%s", field.Name(), key)
				}
				h = append(h, kh)
				row = append(row, data[key])
			}

			headers = append(headers, h...)
			rows = append(rows, row...)
		default:
			headers = append(headers, name)
			rows = append(rows, join(field.Value(), ", "))
		}

	}

	return headers, rows
}

// String renders the internal representation of the table to a string
func (e *Encoder) String() string {
	return string(e.Bytes())
}

// Bytes renders the internal representation of the table to a byte slice
func (e *Encoder) Bytes() []byte {
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	if len(e.headers) > 0 {
		w.Write(e.headers)
	}
	w.WriteAll(e.rows)
	w.Flush()
	return b.Bytes()
}

func join(input interface{}, sep string) string {
	switch v := input.(type) {
	case []string:
		return strings.Join(v, sep)
	case []int:
		ints := []string{}
		for _, i := range v {
			ints = append(ints, strconv.Itoa(i))
		}
		return strings.Join(ints, sep)
	case []float32:
		floats := []string{}
		for _, f := range v {
			floats = append(floats, fmt.Sprint(f))
		}
		return strings.Join(floats, sep)
	case []float64:
		floats := []string{}
		for _, f := range v {
			floats = append(floats, fmt.Sprint(f))
		}
		return strings.Join(floats, sep)
	default:
		return fmt.Sprint(input)
	}
}
