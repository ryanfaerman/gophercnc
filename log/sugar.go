package log

import (
	"net/http"
	"os"
)

// Default is a shared default logger that logs to os.Stderr at LevelInfo. It
// is used by the package level functions for their logger.
var Default = Logger{
	Out:   os.Stderr,
	Level: LevelInfo,
}

func Debug(args ...interface{}) { Default.Debug(args...) }
func Info(args ...interface{})  { Default.Info(args...) }
func Warn(args ...interface{})  { Default.Warn(args...) }
func Error(args ...interface{}) { Default.Error(args...) }

func WithFields(args ...interface{}) Logger { return Default.WithFields(args...) }
func WithLevel(lvl Level) Logger            { return Default.WithLevel(lvl) }
func WithError(err error) Logger            { return Default.WithError(err) }

func Middleware(next http.Handler) http.Handler { return Default.Middleware(next) }

func ForRequest(r *http.Request, args ...interface{}) Logger {
	l, ok := r.Context().Value(ctxLogger).(Logger)
	if !ok {
		return Default.WithFields(args...)
	}
	l.r = r

	return l.WithFields(args...)
}
