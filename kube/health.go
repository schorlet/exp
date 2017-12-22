package main

import (
	"net/http"
	"sync/atomic"
)

// liveness is a liveness probe.
//
// The purpose of a liveness probe is to understand that
// the application is running. If the liveness probe fails,
// the service will be restarted.
func liveness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readiness is a readiness probe.
//
// The purpose of a readiness probe is to understand if
// the application is ready to serve traffic. If the readiness probe fails,
// the container will be removed from service load balancers.
func readiness(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
