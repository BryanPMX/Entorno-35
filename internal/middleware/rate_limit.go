package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type fixedWindowEntry struct {
	windowStart time.Time
	count       int
}

// FixedWindowRateLimiter is a simple in-memory fixed-window limiter keyed by client identity.
// It is process-local (single-instance) and intended for basic abuse protection.
type FixedWindowRateLimiter struct {
	mu          sync.Mutex
	entries     map[string]fixedWindowEntry
	limit       int
	window      time.Duration
	lastCleanup time.Time
}

// NewFixedWindowRateLimiter creates a limiter with the given limit and window.
func NewFixedWindowRateLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	if limit <= 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &FixedWindowRateLimiter{
		entries: make(map[string]fixedWindowEntry),
		limit:   limit,
		window:  window,
	}
}

func (l *FixedWindowRateLimiter) Allow(key string, now time.Time) (bool, time.Duration) {
	if l == nil {
		return true, 0
	}

	now = now.UTC()
	if strings.TrimSpace(key) == "" {
		key = "anonymous"
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.cleanupExpiredLocked(now)

	entry, exists := l.entries[key]
	if !exists || now.Sub(entry.windowStart) >= l.window {
		l.entries[key] = fixedWindowEntry{windowStart: now, count: 1}
		return true, 0
	}

	if entry.count >= l.limit {
		retryAfter := l.window - now.Sub(entry.windowStart)
		if retryAfter < 0 {
			retryAfter = 0
		}
		return false, retryAfter
	}

	entry.count++
	l.entries[key] = entry
	return true, 0
}

func (l *FixedWindowRateLimiter) cleanupExpiredLocked(now time.Time) {
	// Opportunistic cleanup at most once per window to avoid unbounded map growth.
	if !l.lastCleanup.IsZero() && now.Sub(l.lastCleanup) < l.window {
		return
	}
	for key, entry := range l.entries {
		if now.Sub(entry.windowStart) >= l.window {
			delete(l.entries, key)
		}
	}
	l.lastCleanup = now
}

// CheckoutSessionRateLimitMiddleware applies IP-based rate limiting to public checkout creation.
func CheckoutSessionRateLimitMiddleware(limiter *FixedWindowRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter == nil {
			c.Next()
			return
		}

		key := "checkout:" + strings.TrimSpace(c.ClientIP())
		allowed, retryAfter := limiter.Allow(key, time.Now())
		if allowed {
			c.Next()
			return
		}

		retryAfterSeconds := int(retryAfter.Seconds())
		if retryAfterSeconds < 1 {
			retryAfterSeconds = 1
		}
		c.Header("Retry-After", strconv.Itoa(retryAfterSeconds))
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "too many checkout attempts. please try again later",
		})
		c.Abort()
	}
}
