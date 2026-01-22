package middleware

import (
	"backend/utils"
	"context"
	"net/http"
	"time"
)

// TimeoutMiddleware wraps handlers with a timeout context
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create context with timeout
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Create new request with timeout context
			r = r.WithContext(ctx)

			// Channel to signal completion
			done := make(chan struct{})
			
			go func() {
				defer func() {
					if err := recover(); err != nil {
						// Handle panics in timeout scenario
						select {
						case <-done:
							// Already completed, ignore
						default:
							close(done)
						}
					}
				}()
				
				next.ServeHTTP(w, r)
				close(done)
			}()

			// Wait for completion or timeout
			select {
			case <-done:
				// Request completed normally
				return
			case <-ctx.Done():
				// Request timed out
				if ctx.Err() == context.DeadlineExceeded {
					utils.WriteErrorResponse(w, http.StatusRequestTimeout, 
						"Request timeout exceeded")
				} else {
					utils.WriteErrorResponse(w, http.StatusInternalServerError, 
						"Request cancelled")
				}
				return
			}
		})
	}
}