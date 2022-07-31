package termrender

//go:generate stringer -type=RenderFormat

// RenderFormat determines the output format of the rendering engine
type RenderFormat int

const (
	// These are the available rendering formats
	FormatDefault RenderFormat = iota
	FormatJSON
	FormatJSONLines
	FormatCSV
	FormatTable
	FormatYML
	FormatXML
)

// ParseFormat attempts to parse a string into a RenderFormat, returning an
// error should there be no corresponding format for the string.
func ParseFormat(f string) (RenderFormat, error) {
	switch f {
	case "json":
		return FormatJSON, nil
	case "jsonlines", "jsonl", "json-lines":
		return FormatJSONLines, nil
	case "csv":
		return FormatCSV, nil
	case "table":
		return FormatTable, nil
	case "yml", "yaml":
		return FormatYML, nil
	case "xml":
		return FormatXML, nil
	}

	return FormatDefault, ErrInvalidRenderFormat
}
