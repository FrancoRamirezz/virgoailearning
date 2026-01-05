package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func SetupCORS(allowedOrigins []string) func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
			"X-Requested-With",
			"Accept",
			"Origin",
			"Cache-Control",
			"X-File-Name",
		},
		ExposedHeaders:   []string{"X-Total-Count", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
		Debug:            false,
	})

	return c.Handler
}
