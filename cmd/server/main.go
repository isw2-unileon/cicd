// Package main is the entry point for the calculator HTTP API server.
// This is a deliberately simple application so students can focus on
// the CI/CD workflows rather than the application code itself.
package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/isw2-unileon/cicd/internal/calculator"
)

// CalculateRequest is the JSON body for POST /calculate.
type CalculateRequest struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

// CalculateResponse is the JSON response for POST /calculate.
type CalculateResponse struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

func main() {
	// slog.New creates a structured logger. JSONHandler emits one JSON object
	// per log line — easy to ingest in Datadog, CloudWatch, Loki, etc.
	// Use slog.NewTextHandler for human-readable output during development.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("server starting", "port", port)
	if err := http.ListenAndServe(":"+port, newHandler()); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

// newHandler builds and returns the application's HTTP handler.
// Extracted from main so tests can create the handler without starting a server.
func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/calculate", calculateHandler)
	return mux
}

// healthHandler returns a simple 200 OK to signal the app is running.
// This is used by load balancers and deployment pipelines to check readiness.
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// calculateHandler handles POST /calculate requests.
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("invalid request body", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := calculator.Calculate(req.Operation, req.A, req.B)

	w.Header().Set("Content-Type", "application/json")
	resp := CalculateResponse{Result: result}
	if err != nil {
		slog.Warn("calculation error", "operation", req.Operation, "a", req.A, "b", req.B, "error", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp.Error = err.Error()
	}
	json.NewEncoder(w).Encode(resp)
}
