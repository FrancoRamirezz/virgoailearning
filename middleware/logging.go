package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type LogEntry struct {
	Timestamp   string `json:"timestamp"`
	RequestID   string `json:"request_id"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	StatusCode  int    `json:"status_code"`
	Duration    string `json:"duration"`
	UserAgent   string `json:"user_agent"`
	IP          string `json:"ip"`
	UserID      string `json:"user_id,omitempty"`
	Error       string `json:"error,omitempty"`
	RequestSize int64  `json:"request_size"`
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture status code and response size
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Get request ID from context
		requestID := GetRequestIDFromContext(r.Context())

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		
		// Create structured log entry
		logEntry := LogEntry{
			Timestamp:   start.UTC().Format(time.RFC3339),
			RequestID:   requestID,
			Method:      r.Method,
			Path:        r.URL.Path,
			StatusCode:  wrapped.statusCode,
			Duration:    duration.String(),
			UserAgent:   r.UserAgent(),
			IP:          getClientIP(r),
			RequestSize: r.ContentLength,
		}

		// Add user ID if available from context
		if userID := getUserIDFromContext(r.Context()); userID != "" {
			logEntry.UserID = userID
		}

		// Add error if status code indicates failure
		if wrapped.statusCode >= 400 && wrapped.error != "" {
			logEntry.Error = wrapped.error
		}

		// Log as JSON for structured logging
		if logJSON, err := json.Marshal(logEntry); err == nil {
			log.Println(string(logJSON))
		} else {
			// Fallback to simple logging
			log.Printf("[%s] %s %s %d %v", requestID, r.Method, r.URL.Path, wrapped.statusCode, duration)
		}
	})
}

// LogError logs an error with context
func LogError(r *http.Request, err error, message string) {
	requestID := GetRequestIDFromContext(r.Context())
	
	errorLog := map[string]interface{}{
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"request_id": requestID,
		"level":      "ERROR",
		"message":    message,
		"error":      err.Error(),
		"method":     r.Method,
		"path":       r.URL.Path,
		"ip":         getClientIP(r),
	}

	if userID := getUserIDFromContext(r.Context()); userID != "" {
		errorLog["user_id"] = userID
	}

	if logJSON, err := json.Marshal(errorLog); err == nil {
		log.Println(string(logJSON))
	} else {
		log.Printf("[%s] ERROR: %s - %v", requestID, message, err)
	}
}

// LogInfo logs an informational message with context
func LogInfo(r *http.Request, message string, data map[string]interface{}) {
	requestID := GetRequestIDFromContext(r.Context())
	
	infoLog := map[string]interface{}{
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"request_id": requestID,
		"level":      "INFO",
		"message":    message,
		"method":     r.Method,
		"path":       r.URL.Path,
	}

	// Add additional data
	for k, v := range data {
		infoLog[k] = v
	}

	if userID := getUserIDFromContext(r.Context()); userID != "" {
		infoLog["user_id"] = userID
	}

	if logJSON, err := json.Marshal(infoLog); err == nil {
		log.Println(string(logJSON))
	} else {
		log.Printf("[%s] INFO: %s", requestID, message)
	}
}

// getUserIDFromContext retrieves user ID from request context
func getUserIDFromContext(ctx context.Context) string {
	// This assumes you have user information in context from auth middleware
	if user := ctx.Value("user"); user != nil {
		// Type assertion would depend on your User struct
		// For now, return empty string
		return ""
	}
	return ""
}

// SetLogOutput sets the log output destination
func SetLogOutput(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	error      string
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	// Capture error messages from response body for logging
	if rw.statusCode >= 400 && rw.error == "" {
		rw.error = string(data)
		if len(rw.error) > 200 {
			rw.error = rw.error[:200] + "..."
		}
	}
	return rw.ResponseWriter.Write(data)
}

func (rw *responseWriter) Flush() {
	// Check if the underlying ResponseWriter is a Flusher
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
