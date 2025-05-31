package handler

import (
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func RegisterHealthSystemRoutes(mux *runtime.ServeMux) {
	mux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))

		slog.Info("/health: ok")
	})

	mux.HandlePath("GET", "/ready", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))

		slog.Info("/ready: ok")
	})
}
