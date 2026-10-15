package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginRateLimiter_AllowsUpToMaxThenBlocks(t *testing.T) {
	l := NewLoginRateLimiter(false)

	for i := 0; i < loginRateLimitMax; i++ {
		assert.True(t, l.Allow("203.0.113.1"), "attempt %d should be allowed", i+1)
	}

	assert.False(t, l.Allow("203.0.113.1"), "attempt beyond the max should be blocked")
}

func TestLoginRateLimiter_TracksIPsIndependently(t *testing.T) {
	l := NewLoginRateLimiter(false)

	for i := 0; i < loginRateLimitMax; i++ {
		assert.True(t, l.Allow("203.0.113.1"))
	}
	assert.False(t, l.Allow("203.0.113.1"), "first IP should now be blocked")

	assert.True(t, l.Allow("203.0.113.2"), "a different IP must have its own, unaffected budget")
}

func TestClientIP(t *testing.T) {
	tests := []struct {
		name         string
		remoteAddr   string
		forwardedFor string
		trustFwd     bool
		expectedIP   string
	}{
		{name: "no X-Forwarded-For falls back to RemoteAddr", remoteAddr: "198.51.100.5:54321", trustFwd: true, expectedIP: "198.51.100.5"},
		{name: "single X-Forwarded-For value is used", remoteAddr: "10.0.0.1:1234", forwardedFor: "203.0.113.9", trustFwd: true, expectedIP: "203.0.113.9"},
		{name: "last of multiple X-Forwarded-For values is used", remoteAddr: "10.0.0.1:1234", forwardedFor: "198.51.100.77, 203.0.113.9", trustFwd: true, expectedIP: "203.0.113.9"},
		{name: "untrusted X-Forwarded-For is ignored", remoteAddr: "10.0.0.1:1234", forwardedFor: "203.0.113.9", trustFwd: false, expectedIP: "10.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/v1/auth/login", http.NoBody)
			r.RemoteAddr = tt.remoteAddr
			if tt.forwardedFor != "" {
				r.Header.Set("X-Forwarded-For", tt.forwardedFor)
			}

			assert.Equal(t, tt.expectedIP, clientIP(r, tt.trustFwd))
		})
	}
}
