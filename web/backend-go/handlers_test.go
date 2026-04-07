package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	w := httptest.NewRecorder()

	handleStatus(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", body["status"])
	}
	if _, ok := body["engine_available"]; !ok {
		t.Error("missing engine_available field")
	}
}

func TestHandleParseValid(t *testing.T) {
	raw := buildTestReplay(2, 3, 3, 1, 1)
	payload := `{"replay_text":"` + strings.ReplaceAll(raw, `"`, `\"`) + `"}`

	req := httptest.NewRequest(http.MethodPost, "/api/parse", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleParse(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body=%s", w.Code, w.Body.String())
	}

	var replay Replay
	if err := json.Unmarshal(w.Body.Bytes(), &replay); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if replay.Seed != 12345 {
		t.Errorf("expected seed 12345, got %d", replay.Seed)
	}
}

func TestHandleParseInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/parse",
		strings.NewReader(`{"replay_text":"not valid"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleParse(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleParseBadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/parse",
		strings.NewReader(`{bad json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleParse(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleGenerateNoEngine(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/generate",
		strings.NewReader(`{"seed": 42}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleGenerate(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 (no engine), got %d", w.Code)
	}
}

func TestCORSMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := corsMiddleware(inner)

	req := httptest.NewRequest(http.MethodOptions, "/api/status", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("missing CORS Allow-Origin header")
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected 200 for GET, got %d", w2.Code)
	}
	if w2.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("missing CORS header on GET")
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusCreated, map[string]string{"key": "val"})

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["key"] != "val" {
		t.Errorf("expected val, got %s", body["key"])
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusNotFound, "gone")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["detail"] != "gone" {
		t.Errorf("expected 'gone', got %s", body["detail"])
	}
}
