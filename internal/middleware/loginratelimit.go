package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"

	req "github.com/doublehops/dh-go-framework/internal/request"
)

// loginRateLimitWindow and loginRateLimitMax bound how many login attempts a single client IP
// may make in a rolling window before further attempts are rejected with 429 — POST /auth/login
// previously had no brute-force protection at all.
const (
	loginRateLimitWindow = time.Minute
	loginRateLimitMax    = 10
)

type loginAttempts struct {
	count      int
	windowFrom time.Time
}

// LoginRateLimiter tracks login attempts per client IP within a rolling time window. It's
// in-process state, so it resets on deploy and isn't shared across multiple instances — good
// enough as a first line of defense, but should eventually be paired with (or replaced by) an
// edge-level control such as a WAF rate-based rule.
type LoginRateLimiter struct {
	mu       sync.Mutex
	attempts map[string]*loginAttempts
	// trustForwardedFor keys clients by the first X-Forwarded-For entry instead of RemoteAddr.
	// Only enable it behind a proxy that overwrites the header; otherwise a client can send a
	// different fake value on every request and never hit the limit.
	trustForwardedFor bool
}

// NewLoginRateLimiter creates a LoginRateLimiter and starts its background cleanup goroutine,
// which stops the attempts map from growing unbounded over the process's lifetime.
func NewLoginRateLimiter(trustForwardedFor bool) *LoginRateLimiter {
	l := &LoginRateLimiter{attempts: make(map[string]*loginAttempts), trustForwardedFor: trustForwardedFor}
	go l.cleanupLoop()

	return l
}

func (l *LoginRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(loginRateLimitWindow)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for ip, a := range l.attempts {
			if now.Sub(a.windowFrom) > loginRateLimitWindow {
				delete(l.attempts, ip)
			}
		}
		l.mu.Unlock()
	}
}

// Allow reports whether ip may make another attempt right now, incrementing its counter if so.
func (l *LoginRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	a, ok := l.attempts[ip]
	if !ok || now.Sub(a.windowFrom) > loginRateLimitWindow {
		l.attempts[ip] = &loginAttempts{count: 1, windowFrom: now}

		return true
	}

	if a.count >= loginRateLimitMax {
		return false
	}

	a.count++

	return true
}

// Middleware rejects requests from a client IP that has exceeded the login rate limit with
// 429 Too Many Requests, before they ever reach the handler (and its bcrypt comparison).
func (l *LoginRateLimiter) Middleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !l.Allow(clientIP(r, l.trustForwardedFor)) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(req.GeneralErrResp(req.ErrTooManyRequests.Error(), http.StatusTooManyRequests))

			return
		}

		next(w, r, ps)
	}
}

// clientIP extracts the caller's IP from the raw connection's RemoteAddr, or — when
// trustForwardedFor is set (e.g. behind a load balancer that overwrites the header) — from the
// first X-Forwarded-For entry when present.
func clientIP(r *http.Request, trustForwardedFor bool) string {
	if fwd := r.Header.Get("X-Forwarded-For"); trustForwardedFor && fwd != "" {
		if idx := strings.Index(fwd, ","); idx != -1 {
			return strings.TrimSpace(fwd[:idx])
		}

		return strings.TrimSpace(fwd)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
