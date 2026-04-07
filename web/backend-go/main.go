package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	sampleOnce  sync.Once
	sampleCache *Replay
	sampleErr   error
)

func loadSample() (*Replay, error) {
	sampleOnce.Do(func() {
		dataDir := envOr("DATA_DIR", "data")
		raw, err := os.ReadFile(filepath.Join(dataDir, "sample.out"))
		if err != nil {
			sampleErr = err
			return
		}
		sampleCache, sampleErr = parseReplay(string(raw))
	})
	return sampleCache, sampleErr
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"detail": msg})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":           "ok",
		"engine_available": gameAvailable(),
	})
}

func handleSample(w http.ResponseWriter, r *http.Request) {
	replay, err := loadSample()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, replay)
}

func handleAIs(w http.ResponseWriter, r *http.Request) {
	ais := listAIs()
	if ais == nil {
		ais = []string{}
	}
	writeJSON(w, http.StatusOK, ais)
}

type generateRequest struct {
	Seed    *int     `json:"seed"`
	Players []string `json:"players"`
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if !gameAvailable() {
		writeError(w, http.StatusServiceUnavailable, "C++ engine not available in this deployment")
		return
	}

	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	seed := rand.Intn(999999) + 1
	if req.Seed != nil {
		seed = *req.Seed
	}

	raw, err := runGame(seed, req.Players)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	replay, err := parseReplay(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, replay)
}

type parseRequest struct {
	ReplayText string `json:"replay_text"`
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	var req parseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	replay, err := parseReplay(req.ReplayText)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse replay: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, replay)
}

type spaHandler struct {
	staticDir string
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticDir, filepath.Clean(r.URL.Path))
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		http.ServeFile(w, r, path)
		return
	}
	http.ServeFile(w, r, filepath.Join(h.staticDir, "index.html"))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := envOr("PORT", "8087")
	staticDir := envOr("STATIC_DIR", "frontend/dist")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/status", handleStatus)
	mux.HandleFunc("GET /api/sample", handleSample)
	mux.HandleFunc("GET /api/ais", handleAIs)
	mux.HandleFunc("POST /api/generate", handleGenerate)
	mux.HandleFunc("POST /api/parse", handleParse)

	if info, err := os.Stat(staticDir); err == nil && info.IsDir() {
		mux.Handle("/", &spaHandler{staticDir: staticDir})
	}

	handler := corsMiddleware(mux)

	log.Printf("listening on 0.0.0.0:%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
