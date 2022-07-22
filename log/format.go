package log

import "time"

func FormatRFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}
