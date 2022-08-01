package log

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Logger is a simple logger that has three levels: Debug, Info, and Error.
// Other than the name of the levels (and subsequent conditional output) the
// levels are identical in function. A logger without configuration writes at
// DebugLevel to os.Stderr.
//
// This may not be the most efficient or featureful logger out there -- it
// doesn't do colors, doesn't support multiple output formats, doesn't have
// much in the way of fancy hooks and such -- it just does the bare minimum
// that I need/want it to do. It *does* provide:
//
//  - structured logging with fields
//  - multiple log levels
//  - a simple interface
//  - no external dependencies
//
// If fancier/more-powerful features are needed/wanted... use a different
// package.
type Logger struct {
	Out   io.Writer
	Level Level

	fields []interface{}
	r      *http.Request
}

func (l Logger) WithLevel(lvl Level) Logger {
	l.Level = lvl
	return l
}

// WithFields provides some sugar to build the fields up over time until
// calling one of the Level-based functions (Debug, et. al.) and all uses of
// the Logger returned share the same fields.
func (l Logger) WithFields(args ...interface{}) Logger {
	l.fields = append(l.fields, args...)
	if l.r != nil {
		id, ok := l.r.Context().Value(ctxLogID).(uint64)
		if !ok {
			return l
		}

		// args = append(l.fields, args...)

		dataMtx.Lock()
		defer dataMtx.Unlock()
		if _, ok := data[id]; !ok {
			data[id] = args
		} else {
			data[id] = append(data[id], args...)
		}
	}
	return l
}

func (l Logger) WithError(err error) Logger {
	return l.WithFields("err", err.Error())
}

// Print emits a log at the current level
func (l Logger) Print(args ...interface{}) { l.log(LevelDebug, args) }

// Printf emits a formattable log at the current level
func (l Logger) Printf(format string, args ...any) {
	l.Print(fmt.Sprintf(format, args...))
}

// Println emits a log at the current level. This is a shim to comply with the
// Log interface from the stdlib
func (l Logger) Println(args ...any) { l.Print(args...) }

// Debug emits logs when the set level is Debug or lower
func (l Logger) Debug(args ...interface{}) { l.log(LevelDebug, args) }

// Info emits logs when the set level is Info or lower
func (l Logger) Info(args ...interface{}) { l.log(LevelInfo, args) }

// Warn emits logs when the set level is Warn or lower
func (l Logger) Warn(args ...interface{}) { l.log(LevelWarn, args) }

// Error emits logs when the set level is Error or lower
func (l Logger) Error(args ...interface{}) { l.log(LevelError, args) }

// Fatal emits logs when the set level is Fatal or lower
func (l Logger) Fatal(args ...interface{}) { l.log(LevelFatal, args) }

func (l Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

func (l Logger) log(lvl Level, args []interface{}) {
	if l.Out == nil {
		l.Out = os.Stderr
	}
	if l.Level <= lvl {
		lvl.Log(l.Out, append(args, l.fields...)...)
	}
	if l.Level == LevelFatal {
		os.Exit(1)
	}
}
