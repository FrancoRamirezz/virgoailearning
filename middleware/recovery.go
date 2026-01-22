package middleware

import (
	"backend/utils"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

// RecoveryMiddleware recovers from panics and returns a 500 error instead of crashing the server
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Printf("PANIC: %v\nStack trace:\n%s", err, debug.Stack())

				// Get request ID from context if available
				requestID := GetRequestIDFromContext(r.Context())
				if requestID != "" {
					log.Printf("Request ID: %s", requestID)
				}

				// Return 500 error to client
				utils.WriteErrorResponse(w, http.StatusInternalServerError, 
					fmt.Sprintf("Internal server error (Request ID: %s)", requestID))
			}
		}()

		next.ServeHTTP(w, r)
	})
}