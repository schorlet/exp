package http

import (
	"expvar"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func statsHandler(name string, next http.Handler) http.Handler {
	stats := expvar.NewMap(name)
	stats.Set("requests", new(expvar.Int))
	stats.Set("errors", new(expvar.Int))
	duration := new(timing)
	stats.Set("duration", duration)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			begin := time.Now()

			defer func() {
				duration.Observe(time.Since(begin))
				stats.Add("requests", 1)
				if err := recover(); err != nil {
					stats.Add("errors", 1)
				}
			}()

			next.ServeHTTP(w, r)
		},
	)
}

// https://pocketgophers.com/10-to-instrument/

type timing struct {
	mu    sync.RWMutex
	count int64
	sum   time.Duration
}

func (t *timing) Observe(d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.count++
	t.sum += d
}

func (t timing) String() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var avg time.Duration
	if t.count > 0 {
		avg = time.Duration(t.sum.Nanoseconds() / t.count)
	}
	return fmt.Sprintf("%f", avg.Seconds())
}
