package log

import (
	"os"
	"testing"
	"time"
)

type slowWriter struct {
	Delay time.Duration
}

func (s slowWriter) Write(p []byte) (int, error) {
	time.Sleep(s.Delay)
	return os.Stderr.Write(p)
}

func TestLogger(t *testing.T) {
	debug := LevelDebug

	sw := os.Stderr
	debug.Log(sw, "something happened", "hello", "world", "moom")
	debug.Log(sw, "something happened", "hello", "world", "moom")
	debug.Log(sw, "something happened", "moom", "boom")

	log := Logger{
		Out:   os.Stderr,
		Level: LevelDebug,
	}
	log = log.WithFields("foo", 7)

	log.Debug("banananana", "hello", 5)

	Info("hello there", "speaker", "obi-wan")
}
