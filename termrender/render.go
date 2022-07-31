package termrender

import (
	"encoding/json"
	"encoding/xml"
	"io"

	"github.com/ryanfaerman/gophercnc/termrender/csv"
	"github.com/ryanfaerman/gophercnc/termrender/jsonlines"
	"github.com/ryanfaerman/gophercnc/termrender/table"
	"gopkg.in/yaml.v2"
)

// Engine is the workhorse of render abstraction. It handles configuration and
// controls which underlying rendering marshaler to use and where to send the
// output.
type Engine struct {
	RenderFormat RenderFormat
	Writer       io.Writer
}

// New instantiates a default rendering engine
func New() *Engine {
	return &Engine{
		RenderFormat: DefaultRenderFormat,
		Writer:       DefaultWriter,
	}
}

// SetRenderFormat is a helper to set the render format
func (e *Engine) SetRenderFormat(f RenderFormat) {
	e.RenderFormat = f
}

// WithFormat returns a copy of the Engine, overriding the engine's render
// format. Useful for one-time rendering where a different render format is
// required.
func (e *Engine) WithFormat(f RenderFormat) *Engine {
	other := *e
	other.RenderFormat = f
	return &other
}

// WithWriter returns a copy of the Engine, overriding the engine's writer.
// Useful for one-time rendering to a specified writer.
func (e *Engine) WithWriter(w io.Writer) *Engine {
	other := *e
	other.Writer = w
	return &other
}

// Marshal the provided data into one of the RenderFormats,
// returning a byte slice representation.
func (e *Engine) Marshal(data interface{}) ([]byte, error) {
	if e.RenderFormat == FormatDefault {
		e.RenderFormat = DefaultRenderFormat
	}

	switch e.RenderFormat {
	case FormatJSON:
		return json.MarshalIndent(data, "", "  ")
	case FormatCSV:
		return csv.Marshal(data)
	case FormatTable:
		return table.Marshal(data)
	case FormatYML:
		return yaml.Marshal(data)
	case FormatJSONLines:
		return jsonlines.Marshal(data)
	case FormatXML:
		return xml.Marshal(data)
	default:
		return nil, ErrInvalidRenderFormat
	}
}

// RenderString marshals the provided data into one of the RenderFormats,
// returning a string representation.
func (e *Engine) MarshalString(data interface{}) (string, error) {
	b, err := e.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Render marshals the provided data and writes it to the configured Writer.
func (e *Engine) Render(data interface{}) error {
	if e.Writer == nil {
		e.Writer = DefaultWriter
	}

	b, err := e.Marshal(data)
	if err != nil {
		return err
	}

	_, err = e.Writer.Write(b)

	return err
}
