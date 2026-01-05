package routes

import (
	"backend/config"
	"backend/handlers"
	"backend/middleware"
	"backend/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func SetupRoutes(cfg *config.Config) *mux.Router {
	// Create the main router that will handle all HTTP requests
	router := mux.NewRouter()

	// Add production-grade middleware in order of execution
	// 1. Request ID generation (must be first to track all requests)
	router.Use(middleware.RequestIDMiddleware)
	
	// 2. Panic recovery (protects against crashes)
	router.Use(middleware.RecoveryMiddleware)
	
	// 3. Security headers (adds security protections)
	router.Use(middleware.SecurityHeadersMiddleware)
	
	// 4. Request size limits (prevents large payloads)
	router.Use(middleware.RequestSizeLimitMiddleware(10 * 1024 * 1024)) // 10MB limit
	
	// 5. Rate limiting (prevents abuse) - 100 requests per minute
	router.Use(middleware.RateLimitMiddleware(100))
	
	// 6. Request timeout (prevents hanging requests) - 30 second timeout
	router.Use(middleware.TimeoutMiddleware(30 * time.Second))
	
	// 7. CORS (handles cross-origin requests)
	router.Use(middleware.SetupCORS(cfg.CORS.AllowedOrigins))
	
	// 8. Structured logging (must be after RequestID)
	router.Use(middleware.LoggingMiddleware)

	// Simple endpoint to check if the server is up
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteSuccessResponse(w, "Server is running", map[string]string{
			"status": "healthy",
			"env":    cfg.Server.Env,
		})
	}).Methods("GET")

	// Group all API endpoints under /api/v1 (version 1 of the API)
	api := router.PathPrefix("/api/v1").Subrouter()

	// ============================================================================
	// PUBLIC ROUTES (No authentication required)
	// ============================================================================

	// Authentication routes (register, login, OAuth)
	authHandler := handlers.NewAuthHandler()
	authHandler.InitializeOAuth(cfg)
	authHandler.RegisterRoutes(api)

	// Course routes - public endpoints (list, view courses)
	courseHandler := handlers.NewCourseHandler()
	courseHandler.RegisterRoutes(api)

	// Contact routes - public contact form
	contactHandler := handlers.NewContactHandler()
	contactHandler.RegisterRoutes(api)

	// ============================================================================
	// PROTECTED ROUTES (Require JWT authentication)
	// ============================================================================
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// User management (profile, settings, etc.)
	userHandler := handlers.NewUserHandler()
	userHandler.RegisterRoutes(protected)

	// Post/Blog management
	postHandler := handlers.NewPostHandler()
	postHandler.RegisterRoutes(protected)

	// Course payment system (process payments, history, etc.)
	coursePaymentHandler := handlers.NewCoursePaymentHandler()
	coursePaymentHandler.RegisterRoutes(protected)

	// Stripe payment intents (original payment system)
	paymentHandler := handlers.NewPaymentHandler(cfg)
	paymentHandler.RegisterRoutes(protected)

	// Stripe webhook endpoint (public): Stripe calls this to notify us
	// about payment events. It must be public so Stripe can reach it.
	api.HandleFunc("/payments/webhook", paymentHandler.StripeWebhook).Methods("POST")

	// ============================================================================
	// ADMIN ROUTES (Require admin role)
	// ============================================================================
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminMiddleware)

	// Admin can access all contact management
	// Admin can access payment statistics and refunds (already in coursePaymentHandler)
	// Additional admin-specific endpoints can be added here

	return router
}
