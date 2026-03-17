// Integration tests for the calculator HTTP API.
//
// These tests use net/http/httptest to exercise the full HTTP stack —
// routing, request parsing, business logic, and response encoding —
// without binding to a real network port.
//
// Two httptest helpers are demonstrated:
//
//   - httptest.NewRecorder() — lightweight, in-memory response writer.
//     Call the handler directly. Fast; no network involved.
//
//   - httptest.NewServer()   — starts a real TCP server on a random port.
//     Use a real http.Client to make requests. Closer to production.
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ── /health ───────────────────────────────────────────────────────────────────

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		wantStatusCode int
		wantStatus     string
	}{
		{
			name:           "GET returns 200 with status ok",
			method:         http.MethodGet,
			wantStatusCode: http.StatusOK,
			wantStatus:     "ok",
		},
		// The health endpoint accepts any method — it is a readiness probe and
		// load balancers may use HEAD. Documenting the current behaviour here.
		{
			name:           "POST also returns 200",
			method:         http.MethodPost,
			wantStatusCode: http.StatusOK,
			wantStatus:     "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			rec := httptest.NewRecorder()

			healthHandler(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatusCode)
			}

			var body map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("could not decode response body: %v", err)
			}

			if got := body["status"]; got != tt.wantStatus {
				t.Errorf("body[status] = %q, want %q", got, tt.wantStatus)
			}

			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}
		})
	}
}

// ── /calculate ────────────────────────────────────────────────────────────────

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           any // marshalled to JSON; use a raw string for malformed JSON
		wantStatusCode int
		wantResult     float64
		wantError      string // non-empty when an error response is expected
	}{
		// ── Happy path ──────────────────────────────────────────────────────
		{
			name:           "add two positive numbers",
			method:         http.MethodPost,
			body:           CalculateRequest{Operation: "add", A: 5, B: 3},
			wantStatusCode: http.StatusOK,
			wantResult:     8,
		},
		{
			name:           "subtract",
			method:         http.MethodPost,
			body:           CalculateRequest{Operation: "subtract", A: 10, B: 4},
			wantStatusCode: http.StatusOK,
			wantResult:     6,
		},
		{
			name:           "multiply",
			method:         http.MethodPost,
			body:           CalculateRequest{Operation: "multiply", A: 3, B: 7},
			wantStatusCode: http.StatusOK,
			wantResult:     21,
		},
		{
			name:           "divide",
			method:         http.MethodPost,
			body:           CalculateRequest{Operation: "divide", A: 9, B: 3},
			wantStatusCode: http.StatusOK,
			wantResult:     3,
		},
		{
			name:           "add with negative numbers",
			method:         http.MethodPost,
			body:           CalculateRequest{Operation: "add", A: -5, B: -3},
			wantStatusCode: http.StatusOK,
			wantResult:     -8,
		},
		{
			name:           "multiply by zero",
			method:         http.MethodPost,
			body:           CalculateRequest{Operation: "multiply", A: 999, B: 0},
			wantStatusCode: http.StatusOK,
			wantResult:     0,
		},
		// ── Error: domain errors (422 Unprocessable Entity) ─────────────────
		{
			name:           "divide by zero returns error",
			method:         http.MethodPost,
			body:           CalculateRequest{Operation: "divide", A: 10, B: 0},
			wantStatusCode: http.StatusUnprocessableEntity,
			wantError:      "division by zero",
		},
		{
			name:           "unknown operation returns error",
			method:         http.MethodPost,
			body:           CalculateRequest{Operation: "modulo", A: 10, B: 3},
			wantStatusCode: http.StatusUnprocessableEntity,
			wantError:      "unknown operation",
		},
		// ── Error: bad HTTP usage ────────────────────────────────────────────
		{
			name:           "GET method not allowed",
			method:         http.MethodGet,
			body:           nil,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "PUT method not allowed",
			method:         http.MethodPut,
			body:           nil,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "malformed JSON returns 400",
			method:         http.MethodPost,
			body:           `{not valid json`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "empty body returns 400",
			method:         http.MethodPost,
			body:           ``,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ── Build request body ──────────────────────────────────────────
			var reqBody *bytes.Buffer
			switch v := tt.body.(type) {
			case string:
				reqBody = bytes.NewBufferString(v)
			case nil:
				reqBody = bytes.NewBuffer(nil)
			default:
				data, err := json.Marshal(v)
				if err != nil {
					t.Fatalf("could not marshal request body: %v", err)
				}
				reqBody = bytes.NewBuffer(data)
			}

			req := httptest.NewRequest(tt.method, "/calculate", reqBody)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			calculateHandler(rec, req)

			// ── Assert status code ──────────────────────────────────────────
			if rec.Code != tt.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatusCode)
			}

			// ── Assert JSON body for success and domain-error cases ─────────
			if tt.wantStatusCode == http.StatusOK || tt.wantStatusCode == http.StatusUnprocessableEntity {
				var resp CalculateResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("could not decode response: %v", err)
				}

				if tt.wantError != "" {
					if resp.Error != tt.wantError {
						t.Errorf("error = %q, want %q", resp.Error, tt.wantError)
					}
				} else {
					if resp.Result != tt.wantResult {
						t.Errorf("result = %v, want %v", resp.Result, tt.wantResult)
					}
					if resp.Error != "" {
						t.Errorf("unexpected error in response: %q", resp.Error)
					}
				}
			}
		})
	}
}

// ── Full round-trip with httptest.NewServer ───────────────────────────────────

// TestServerRoundTrip uses httptest.NewServer to spin up a real TCP listener
// and exercises the full stack with a genuine http.Client.
// This catches issues that in-memory recorder tests cannot — e.g. middleware
// that depends on the actual net.Conn, TLS, or real HTTP/1.1 framing.
func TestServerRoundTrip(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()

	client := srv.Client() // pre-configured to trust the test server's TLS cert

	tests := []struct {
		name           string
		path           string
		method         string
		body           any
		wantStatusCode int
		wantResult     float64
		wantError      string
	}{
		{
			name:           "health endpoint",
			path:           "/health",
			method:         http.MethodGet,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "add via real HTTP",
			path:           "/calculate",
			method:         http.MethodPost,
			body:           CalculateRequest{Operation: "add", A: 100, B: 23},
			wantStatusCode: http.StatusOK,
			wantResult:     123,
		},
		{
			name:           "divide by zero via real HTTP",
			path:           "/calculate",
			method:         http.MethodPost,
			body:           CalculateRequest{Operation: "divide", A: 1, B: 0},
			wantStatusCode: http.StatusUnprocessableEntity,
			wantError:      "division by zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBuf *bytes.Buffer
			if tt.body != nil {
				data, err := json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("marshal: %v", err)
				}
				bodyBuf = bytes.NewBuffer(data)
			} else {
				bodyBuf = bytes.NewBuffer(nil)
			}

			req, err := http.NewRequest(tt.method, srv.URL+tt.path, bodyBuf)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatusCode)
			}

			if tt.wantStatusCode == http.StatusOK && tt.path == "/calculate" {
				var body CalculateResponse
				if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
					t.Fatalf("decode: %v", err)
				}
				if body.Result != tt.wantResult {
					t.Errorf("result = %v, want %v", body.Result, tt.wantResult)
				}
			}

			if tt.wantError != "" {
				var body CalculateResponse
				if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
					t.Fatalf("decode: %v", err)
				}
				if body.Error != tt.wantError {
					t.Errorf("error = %q, want %q", body.Error, tt.wantError)
				}
			}
		})
	}
}
