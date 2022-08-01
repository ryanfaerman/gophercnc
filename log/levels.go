package log

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

//go:generate stringer -type=Level -linecomment

// Level is the log level we're to be writing at. Using an int allows the
// levels to be compared, so we can check if the log message should be written
// with easy comparison.
type Level int

const (
	LevelDebug Level = iota // debug
	LevelInfo               // info
	LevelWarn               // warn
	LevelError              // error
	LevelFatal              // fatal
)

// ParseLevel parses a string, matching it to a Level. If the input fails to
// match any known level, it defaults to LevelInfo.
func ParseLevel(in string) (Level, error) {
	switch strings.ToLower(in) {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	case "fatal":
		return LevelFatal, nil
	default:
		return LevelInfo, errors.New(fmt.Sprintf("cannot parse %s as a log level", in))
	}
}

// tuple is a key/value pair that provides context for any given log message
type tuple struct {
	key string
	val interface{}
}

// String creates the string form of the tuple
func (t tuple) String() string {
	return fmt.Sprintf(`%s="%v"`, t.key, t.val)
}

// tuples is a sortable collection of tuple that can also be printed as a
// string for a log line
type tuples []tuple

// These implement sort.Interface
func (t tuples) Len() int           { return len(t) }
func (t tuples) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t tuples) Less(i, j int) bool { return t[i].key < t[j].key }

// String builds an output string for the collection of key/value pairs.
func (t tuples) String() string {
	switch len(t) {
	case 0:
		return ""
	case 1:
		return t[0].String()
	}

	var b strings.Builder

	b.WriteString(t[0].String())
	for _, tuple := range t[1:] {
		b.WriteString(" ")
		b.WriteString(tuple.String())
	}

	return b.String()
}

// Log writes the log line to the provider io.Writer. The first in the args is
// assumed to be the message, all other pairs are assumed to be key/value pairs
// of additional detail. If there are an odd number of args (after the first
// becomes the message) it is dropped. When the args are empty, this function
// is a noop.
//
// The key/value pairs are all written in alphabetical order. The time, level,
// and message keys are added to the output.
//
// The time will always be in UTC, formated as RFC3339.
func (l Level) Log(w io.Writer, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	vals := tuples{
		{key: "time", val: FormatRFC3339(time.Now())},
		{key: "lvl", val: l.String()},
	}

	if len(args) >= 1 {
		var msg interface{}
		msg, args = args[0], args[1:]
		vals = append(vals, tuple{key: "msg", val: msg})
	}

	for len(args) >= 2 {
		var (
			key interface{}
			val interface{}
		)

		key, args = args[0], args[1:]
		val, args = args[0], args[1:]
		if keyStr, ok := key.(string); ok {
			vals = append(vals, tuple{key: keyStr, val: val})
		}
	}

	sort.Sort(vals)

	fmt.Fprint(w, vals.String())
	fmt.Fprint(w, "\n")
}
