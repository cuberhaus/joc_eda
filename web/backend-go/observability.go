// Phase 14 (Option A) — Sentry SDK + JSON-line stdout for the joc-eda
// backend. The DSN is empty by default so the SDK is a no-op until you
// set it via PersonalPortfolio/.env.shared.
//
// `service` tag distinguishes this backend from other demos in the
// shared Sentry project. The custom log writer emits one JSON object
// per line so the PersonalPortfolio log-relay can forward backend logs
// verbatim into the in-page debug overlay.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func envOrEmpty(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return ""
}

// defaultSampleRate picks a safe default when the operator hasn't pinned an
// explicit rate via SENTRY_TRACES_SAMPLE_RATE: 0.1 in production (free-tier
// safe), 1.0 in dev / staging (high signal). Mirrors the same convention the
// Python helper uses so all backends behave consistently.
func defaultSampleRate(env string) float64 {
	if env == "production" {
		return 0.1
	}
	return 1.0
}

// initSentry sets up Sentry. Returns a function the caller should defer
// for graceful flush, plus a flag indicating whether init succeeded.
func initSentry(service string) (func(), bool) {
	dsn := envOrEmpty("SENTRY_DSN")
	if dsn == "" {
		return func() {}, false
	}
	env := envOrEmpty("SENTRY_ENVIRONMENT")
	if env == "" {
		env = "local-dev"
	}
	tracesSample := defaultSampleRate(env)
	if v := envOrEmpty("SENTRY_TRACES_SAMPLE_RATE"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			tracesSample = parsed
		}
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      env,
		Release:          envOrEmpty("SENTRY_RELEASE"),
		TracesSampleRate: tracesSample,
		SendDefaultPII:   false,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "sentry.Init: %v\n", err)
		return func() {}, false
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("service", service)
	})
	return func() { sentry.Flush(2 * time.Second) }, true
}

// jsonStdoutWriter rewrites lines passed to log.* into JSON objects
// `{level, ns, msg, ts}` so the relay can pipe them straight through.
type jsonStdoutWriter struct {
	ns string
}

func (w jsonStdoutWriter) Write(p []byte) (int, error) {
	msg := strings.TrimRight(string(p), "\r\n")
	if msg == "" {
		return len(p), nil
	}
	level := "info"
	if strings.Contains(msg, "ERROR") || strings.Contains(msg, "error:") {
		level = "error"
	}
	line, err := json.Marshal(map[string]any{
		"level": level,
		"ns":    w.ns,
		"msg":   msg,
		"ts":    float64(time.Now().UnixNano()) / 1e9,
	})
	if err != nil {
		return os.Stdout.Write(p)
	}
	line = append(line, '\n')
	if _, err := os.Stdout.Write(line); err != nil {
		return 0, err
	}
	return len(p), nil
}

// withSentryHTTP wraps a handler with sentry-go's HTTP middleware when
// Sentry was initialised, plus a thin per-request layer that stamps the
// X-Session-Id header (set by the parent page's network tap, see
// `PersonalPortfolio/src/lib/debug-network.ts`) onto the request's Sentry
// scope so backend events join the same session as frontend / iframe
// events.
//
// The session-id middleware runs INSIDE the sentry middleware so the hub
// is already attached to the request context when we look it up — the
// other order would silently drop the tag because GetHubFromContext
// returns nil before sentryhttp has populated the context.
func withSentryHTTP(handler http.Handler, enabled bool) http.Handler {
	if !enabled {
		return handler
	}
	sessionMiddleware := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sid := r.Header.Get("X-Session-Id"); sid != "" {
			if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
				hub.Scope().SetTag("session_id", sid)
			}
		}
		handler.ServeHTTP(w, r)
	})
	return sentryhttp.New(sentryhttp.Options{Repanic: true}).Handle(sessionMiddleware)
}
