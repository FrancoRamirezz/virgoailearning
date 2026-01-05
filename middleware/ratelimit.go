package middleware

import (
	"backend/utils"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// RateLimiter stores rate limiting data
type RateLimiter struct {
	requests map[string]*ClientData
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// ClientData stores request count and window start time for each client
type ClientData struct {
	count     int
	window    time.Time
	lastSeen  time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*ClientData),
		limit:    requestsPerMinute,
		window:   time.Minute,
	}

	// Start cleanup goroutine to remove old entries
	go rl.cleanup()
	
	return rl
}

// RateLimitMiddleware limits requests per IP address
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(requestsPerMinute)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			clientIP := getClientIP(r)
			
			if !limiter.Allow(clientIP) {
				// Rate limit exceeded
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "60")
				
				utils.WriteErrorResponse(w, http.StatusTooManyRequests, 
					"Rate limit exceeded. Please try again later.")
				return
			}

			// Set rate limit headers
			remaining := limiter.GetRemaining(clientIP)
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			
			next.ServeHTTP(w, r)
		})
	}
}

// Allow checks if a request is allowed for the given client
func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	
	client, exists := rl.requests[clientIP]
	if !exists {
		rl.requests[clientIP] = &ClientData{
			count:    1,
			window:   now,
			lastSeen: now,
		}
		return true
	}

	client.lastSeen = now

	// Check if we're in a new window
	if now.Sub(client.window) >= rl.window {
		client.count = 1
		client.window = now
		return true
	}

	// Check if limit exceeded
	if client.count >= rl.limit {
		return false
	}

	client.count++
	return true
}

// GetRemaining returns remaining requests for the client
func (rl *RateLimiter) GetRemaining(clientIP string) int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	client, exists := rl.requests[clientIP]
	if !exists {
		return rl.limit - 1
	}

	// If in new window, return full limit minus 1
	if time.Since(client.window) >= rl.window {
		return rl.limit - 1
	}

	remaining := rl.limit - client.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// cleanup removes old entries to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute * 5) // Cleanup every 5 minutes
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for ip, client := range rl.requests {
			// Remove clients not seen in the last hour
			if now.Sub(client.lastSeen) > time.Hour {
				delete(rl.requests, ip)
			}
		}
		rl.mutex.Unlock()
	}
}

// getClientIP extracts the real client IP address
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Check X-Real-IP header (for reverse proxies like nginx)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to remote address
	return r.RemoteAddr
}