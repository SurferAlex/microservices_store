package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type visitor struct {
	tokens     int
	lastRefill time.Time
}

type limiter struct {
	mu     sync.Mutex
	store  map[string]*visitor
	max    int
	window time.Duration
	refill int
}

func newLimiter(max int, window time.Duration) *limiter {
	return &limiter{
		store:  make(map[string]*visitor),
		max:    max,
		window: window,
		refill: max,
	}
}

func (l *limiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	v, ok := l.store[key]
	if !ok {
		l.store[key] = &visitor{tokens: l.max - 1, lastRefill: now}
		return true
	}

	// Рефил по окну
	if now.Sub(v.lastRefill) >= l.window {
		v.tokens = l.max
		v.lastRefill = now
	}

	if v.tokens <= 0 {
		return false
	}
	v.tokens--
	return true
}

func clientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func RateLimitIP(max int, window time.Duration) func(http.Handler) http.Handler {
	l := newLimiter(max, window)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			if !l.allow(ip) {
				w.Header().Set("Retry-After", strconv.Itoa(int(window.Seconds())))
				http.Error(w, "Слишком много запросов", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
