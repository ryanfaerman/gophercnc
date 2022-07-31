package table

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/fatih/structs"
)

// And Encoder writes tablular data to the output stream
type Encoder struct {
	table *simpletable.Table
}

// Encode creates the tabular encoding of the input
//
// With structs, the following struct tags are supported using the `table` key:
//
//  - omitempty : ignores the column, if the value is the zero value
//  - noprefix  : prevents prefixing columns for nested structs with their field name
//  - right     : align right
//  - left      : align left
//  - center    : align center
//
// Fields with the tag `-` are not rendered e.g. `table:"-"`.
func (e *Encoder) Encode(input interface{}) error {
	if e.table == nil {
		e.table = simpletable.New()
	}

	switch reflect.TypeOf(input).Kind() {
	case reflect.Slice, reflect.Array:
		val := reflect.ValueOf(input)
		for i := 0; i < val.Len(); i++ {
			if err := e.Encode(val.Index(i).Interface()); err != nil {
				return err
			}
		}
	case reflect.String:
		row := []*simpletable.Cell{&simpletable.Cell{
			Text: fmt.Sprint(input),
		}}
		e.table.Body.Cells = append(e.table.Body.Cells, row)
	case reflect.Int, reflect.Float32, reflect.Float64:
		row := []*simpletable.Cell{&simpletable.Cell{
			Align: simpletable.AlignRight,
			Text:  fmt.Sprint(input),
		}}
		e.table.Body.Cells = append(e.table.Body.Cells, row)
	case reflect.Struct:
		headerCells, rowCells := e.encodeStruct(input, nil)
		e.table.Header = &simpletable.Header{
			Cells: headerCells,
		}
		e.table.Body.Cells = append(e.table.Body.Cells, rowCells)
	case reflect.Map:
		iter := reflect.ValueOf(input).MapRange()
		headerCells := []*simpletable.Cell{}
		rowCells := []*simpletable.Cell{}
		data := map[string]string{}
		keys := []string{}

		for iter.Next() {
			data[fmt.Sprint(iter.Key())] = join(fmt.Sprintf("%v", iter.Value()), ", ")
			keys = append(keys, fmt.Sprint(iter.Key()))
		}

		sort.Strings(keys)

		for _, key := range keys {

			headerCells = append(headerCells, &simpletable.Cell{
				Align: simpletable.AlignCenter,
				Text:  key,
			})

			rowCells = append(rowCells, &simpletable.Cell{
				Text: data[key],
			})
		}

		e.table.Header = &simpletable.Header{
			Cells: headerCells,
		}
		e.table.Body.Cells = append(e.table.Body.Cells, rowCells)
	case reflect.Ptr:
		e.Encode(reflect.Indirect(reflect.ValueOf(input)))
	default:
		fmt.Println(reflect.TypeOf(input).Kind())

	}

	return nil
}

func (e *Encoder) encodeStruct(input interface{}, parent *structs.Field) ([]*simpletable.Cell, []*simpletable.Cell) {
	headerCells := []*simpletable.Cell{}
	rowCells := []*simpletable.Cell{}
	s := structs.New(input)
	for _, field := range s.Fields() {
		if !field.IsExported() {
			continue
		}
		name := field.Name()

		tagName, tagOpts := parseTag(field.Tag("table"))
		if tagName != "" {
			if tagName == "-" {
				continue
			}
			name = tagName
		}

		if parent != nil {
			parentName := parent.Name()
			parentTagName, parentTagOpts := parseTag(parent.Tag("table"))
			if parentTagName != "" {
				parentName = parentTagName
			}
			if !tagOpts.Has("omitempty") {
				if field.IsZero() {
					rowCells = append(rowCells, &simpletable.Cell{})
					continue
				}
			}
			if !parentTagOpts.Has("noprefix") {
				name = fmt.Sprintf("%s.%s", parentName, name)
			}
		}

		if !structs.IsStruct(field.Value()) {
			headerCells = append(headerCells, &simpletable.Cell{
				Align: simpletable.AlignCenter,
				Text:  name,
			})
		}
		if tagOpts.Has("omitempty") {
			if field.IsZero() {
				rowCells = append(rowCells, &simpletable.Cell{})
				continue
			}
		}

		alignment := simpletable.AlignLeft
		if tagOpts.Has("left") {
			alignment = simpletable.AlignLeft
		}
		if tagOpts.Has("right") {
			alignment = simpletable.AlignRight
		}
		if tagOpts.Has("center") {
			alignment = simpletable.AlignCenter
		}

		if structs.IsStruct(field.Value()) {
			h, r := e.encodeStruct(field.Value(), field)
			headerCells = append(headerCells, h...)
			rowCells = append(rowCells, r...)
		} else {
			rowCells = append(rowCells, &simpletable.Cell{
				Align: alignment,
				Text:  join(field.Value(), ", "),
			})

		}
	}

	return headerCells, rowCells
}

// String renders the internal representation of the table to a string
func (e *Encoder) String() string {
	return e.table.String() + "\n"
}

// Bytes renders the internal representation of the table to a byte slice
func (e *Encoder) Bytes() []byte {
	return []byte(e.String())
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
		if s, ok := input.(fmt.Stringer); ok {
			return s.String()
		}
		return fmt.Sprint(input)
	}
}
