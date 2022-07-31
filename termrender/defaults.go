package termrender

import (
	"io"
	"os"
)

var (
	// These defaults provide a base level of functionality and can be overridden
	// at the engine level.
	DefaultWriter       = os.Stdout
	DefaultRenderFormat = FormatTable
	DefaultEngine       = New()
)

// SetRenderFormat changes the render format for the DefaultEngine
func SetRenderFormat(f RenderFormat) { DefaultEngine.SetRenderFormat(f) }

// WithFormat returns an engine with the specified format based on the DefaultEngine
func WithFormat(f RenderFormat) *Engine { return DefaultEngine.WithFormat(f) }

// WithWriter returns an engine with the specified writer based on the DefaultEngine
func WithWriter(w io.Writer) *Engine { return DefaultEngine.WithWriter(w) }

// Marshal returns the bytes for the given data using the DefaultEngine
func Marshal(data interface{}) ([]byte, error) { return DefaultEngine.Marshal(data) }

// MarshalString returns a string for the given data using the DefaultEngine
func MarshalString(data interface{}) (string, error) { return DefaultEngine.MarshalString(data) }

// Render writes marshalled data to the writer of the DefaultEngine
func Render(data interface{}) error { return DefaultEngine.Render(data) }
