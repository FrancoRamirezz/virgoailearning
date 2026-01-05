package middleware

import (
	"net/http"
)

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		headers := map[string]string{
			// Prevent MIME type sniffing
			"X-Content-Type-Options": "nosniff",
			
			// Prevent clickjacking
			"X-Frame-Options": "DENY",
			
			// Enable XSS protection
			"X-XSS-Protection": "1; mode=block",
			
			// Prevent information disclosure
			"X-Powered-By": "", // Remove server info
			"Server":       "", // Remove server info
			
			// Referrer policy
			"Referrer-Policy": "strict-origin-when-cross-origin",
			
			// Content Security Policy
			"Content-Security-Policy": "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https:; frame-src 'none'; object-src 'none';",
			
			// Permissions Policy (Feature Policy replacement)
			"Permissions-Policy": "geolocation=(), microphone=(), camera=(), fullscreen=(self), payment=()",
		}

		// Apply HSTS only for HTTPS requests
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			headers["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains; preload"
		}

		// Set all security headers
		for key, value := range headers {
			w.Header().Set(key, value)
		}

		next.ServeHTTP(w, r)
	})
}

// RequestSizeLimitMiddleware limits the size of request bodies
func RequestSizeLimitMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Limit request body size
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}