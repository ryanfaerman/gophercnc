package log

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type key int

const (
	ctxLogID key = iota
	ctxLogger
)

var (
	data    = map[uint64][]interface{}{}
	dataMtx = sync.RWMutex{}
	count   uint64
)

func (l Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := statusRecorder{w, 200}

		id := atomic.AddUint64(&count, 1)
		defer func() {
			dataMtx.Lock()
			delete(data, id)
			dataMtx.Unlock()
		}()
		ctx := context.WithValue(r.Context(), ctxLogID, id)
		ctx = context.WithValue(ctx, ctxLogger, l)

		start := time.Now()
		next.ServeHTTP(&rw, r.WithContext(ctx))
		end := time.Now()

		entry := l.WithFields(
			"http_method", r.Method,
			"http_addr", r.RemoteAddr,
			"http_path", r.URL.Path,
			"http_status", rw.status,
			"http_latency", end.Sub(start),
		)

		dataMtx.RLock()
		if fields, ok := data[id]; ok {
			entry = entry.WithFields(fields...)
		}
		dataMtx.RUnlock()

		lfunc := entry.Info
		if rw.status >= 500 {
			lfunc = entry.Error
		}

		lfunc("canonical-log-line")
	})
}
