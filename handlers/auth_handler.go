package handlers

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// AuthHandler handles all authentication and account security operations
// This handler manages user registration, login, token management, and password operations
type AuthHandler struct {
	config      *config.Config
	store       *sessions.CookieStore
	initialized bool
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		initialized: false,
	}
}

// InitializeOAuth sets up Goth with Google provider
// This should be called once with the config before using OAuth handlers
func (h *AuthHandler) InitializeOAuth(cfg *config.Config) {
	if h.initialized {
		return
	}

	h.config = cfg
	//!===========================================================================
	// Create session store for OAuth state management
	// In production, use a secure random key (32 or 64 bytes)
	sessionKey := []byte(cfg.JWT.Secret) // Reusing JWT secret for session key
	if len(sessionKey) < 32 {
		// Pad or use a default if JWT secret is too short
		sessionKey = []byte("mAQ6vL0acGN+aB5+dJtQCN7fSdN9UBbXZSLiA99XfeM=") //franco's: gKWugzXkW5Ia8o3pWSWJ7nl/uylFzsf0I3mkuWWRQdQ, not sure if the "=" at the end is needed
	}
	h.store = sessions.NewCookieStore(sessionKey)
	h.store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
		Secure:   cfg.Server.Env == "production",
		SameSite: http.SameSiteLaxMode,
	}

	// Set up gothic to use our session store
	gothic.Store = h.store

	// Register Google provider with Goth
	googleProvider := google.New(
		cfg.OAuth.GoogleClientID,
		cfg.OAuth.GoogleClientSecret,
		cfg.OAuth.CallbackURL,
		"email", "profile",
	)

	goth.UseProviders(googleProvider)
	h.initialized = true
}

// 1. Register - Creates a new user account
// POST /api/v1/auth/register
// Request: { email, password, first_name, last_name }
// Response: { token, user }
//
// Why: This is the entry point for new users. It validates the input, checks for duplicates,
// hashes the password securely, creates the user record, and immediately returns a JWT token
// so the user can start using the app without needing to log in separately.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Check if user already exists - prevents duplicate accounts
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.WriteErrorResponse(w, http.StatusConflict, "User already exists")
		return
	}

	// Hash password using bcrypt - never store plain text passwords
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user with default role "user" and active status
	user := models.User{
		Email:     req.Email,
		Password:  &hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		IsActive:  true,
		Provider:  "local",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token immediately so user doesn't need to log in after registration
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, utils.GetConfig())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return token and user info (excluding sensitive data like password)
	response := models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	utils.WriteSuccessResponse(w, "User registered successfully", response)

	//handle db/hashing error
}

// 2. Login - Authenticates an existing user and returns a JWT token
// POST /api/v1/auth/login
// Request: { email, password }
// Response: { token, user }
//
// Why: This validates user credentials and issues a JWT token. The token is used for
// subsequent API calls to prove the user's identity without sending password each time.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal if email exists or not (security best practice)
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Verify password using bcrypt comparison
	// Check if user has a password (OAuth users might not have one)
	if user.Password == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if !utils.CheckPasswordHash(req.Password, *user.Password) {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check if account is active - prevents login for deactivated accounts
	if !user.IsActive {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Account is deactivated")
		return
	}

	// Generate new JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, utils.GetConfig())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return token and user info
	response := models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	utils.WriteSuccessResponse(w, "Login successful", response)
}

// 3. GetProfile - Returns the current authenticated user's profile
// GET /api/v1/auth/profile
// Headers: Authorization: Bearer <token>
// Response: { user }
//
// Why: Allows users to view their own profile information. The user is extracted from
// the JWT token via the AuthMiddleware, so no user ID needs to be passed.
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// User is already validated and added to context by AuthMiddleware
	user, ok := r.Context().Value("user").(models.User)

	if !ok {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "User context missing or invalid")
		return
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Profile retrieved successfully", userResponse)
}

// 4. RefreshToken - Generates a new JWT token for the current user
// POST /api/v1/auth/refresh
// Headers: Authorization: Bearer <token>
// Response: { token, user }
//
// Why: JWT tokens expire for security. Instead of forcing users to log in again,
// they can refresh their token while it's still valid. This improves user experience.
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// User is already validated by AuthMiddleware
	user := r.Context().Value("user").(models.User)

	// Generate a new token with extended expiry
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, utils.GetConfig())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	utils.WriteSuccessResponse(w, "Token refreshed successfully", response)
}

// 5. ChangePassword - Allows authenticated users to change their password
// POST /api/v1/auth/change-password
// Headers: Authorization: Bearer <token>
// Request: { current_password, new_password }
// Response: { message }
//
// Why: Users need to update passwords periodically for security. This requires
// the current password to prevent unauthorized changes if someone gains access to the account.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Get authenticated user from context
	user := r.Context().Value("user").(models.User)

	// Verify current password
	// OAuth users don't have passwords, so they can't change password this way
	if user.Password == nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "OAuth users cannot change password using this endpoint")
		return
	}
	if !utils.CheckPasswordHash(req.CurrentPassword, *user.Password) {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update password in database
	if err := database.DB.Model(&user).Update("password", &hashedPassword).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	utils.WriteSuccessResponse(w, "Password changed successfully", nil)
}

// 6. ForgotPassword - Initiates password reset process
// POST /api/v1/auth/forgot-password
// Request: { email }
// Response: { message }
//
// Why: When users forget their password, they need a way to reset it. This endpoint
// would typically send an email with a reset link. For now, it just validates the email exists.
// In production, you'd generate a reset token and send it via email.
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Check if user exists
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal if email exists (security best practice)
		// Still return success to prevent email enumeration attacks
		utils.WriteSuccessResponse(w, "If the email exists, a password reset link has been sent", nil)
		return
	}

	// TODO: In production, generate a secure reset token, store it with expiry,
	// and send an email with the reset link
	// For now, just return success message
	utils.WriteSuccessResponse(w, "If the email exists, a password reset link has been sent", nil)
}

// 7. ResetPassword - Completes password reset using a reset token
// POST /api/v1/auth/reset-password
// Request: { token, new_password }
// Response: { message }
//
// Why: This completes the password reset flow. Users click the link from their email
// which contains a token, and they can set a new password. The token should be single-use
// and expire after a short time (e.g., 1 hour).
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// TODO: In production, validate the reset token:
	// 1. Check if token exists in database (e.g., PasswordResetToken table)
	// 2. Check if token hasn't expired
	// 3. Check if token hasn't been used
	// 4. Get user associated with token
	// 5. Update password
	// 6. Mark token as used

	// For now, return error indicating feature not fully implemented
	utils.WriteErrorResponse(w, http.StatusNotImplemented, "Password reset token validation not implemented")
}

// 8. Logout - Logs out the current user
// POST /api/v1/auth/logout
// Headers: Authorization: Bearer <token>
// Response: { message }
//
// Why: While JWTs are stateless and can't be "revoked" server-side without a token blacklist,
// this endpoint provides a clear logout action. In production, you might want to implement
// a token blacklist or use refresh tokens that can be invalidated.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: In production with token blacklist:
	// 1. Extract token from request
	// 2. Add token to blacklist (Redis or database) with expiry matching token expiry
	// 3. AuthMiddleware would check blacklist before validating token

	// For now, just return success - client should discard the token
	utils.WriteSuccessResponse(w, "Logged out successfully", nil)
}

// 9. BeginAuth - Initiates OAuth flow with Google
// GET /api/v1/auth/google
// Response: Redirects to Google OAuth consent screen
//
// Why: This starts the OAuth flow. Users are redirected to Google to authenticate,
// and Google will redirect back to the callback URL with an authorization code.
func (h *AuthHandler) BeginAuth(w http.ResponseWriter, r *http.Request) {
	if !h.initialized {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "OAuth not initialized")
		return
	}

	// Set provider in context for gothic
	r = r.WithContext(context.WithValue(r.Context(), "provider", "google"))

	// Begin the OAuth flow - gothic will redirect to Google
	gothic.BeginAuthHandler(w, r)
}

// 10. Callback - Handles OAuth callback from Google
// GET /api/v1/auth/google/callback
// Response: { token, user }
//
// Why: After user authenticates with Google, Google redirects here with an auth code.
// We exchange the code for user info, then create or log in the user and return a JWT token.
func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	if !h.initialized {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "OAuth not initialized")
		return
	}

	// Set provider in context for gothic
	r = r.WithContext(context.WithValue(r.Context(), "provider", "google"))

	// Complete the OAuth flow and get user info from Google
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("OAuth authentication failed: %v", err))
		return
	}

	// Check if user already exists by email
	var dbUser models.User
	err = database.DB.Where("email = ?", user.Email).First(&dbUser).Error

	if err != nil {
		// User doesn't exist, create new user from OAuth data
		// Parse name from Google user data
		firstName := user.FirstName
		lastName := user.LastName

		// If names are empty, try to parse from Name field
		if firstName == "" && lastName == "" && user.Name != "" {
			parts := strings.Fields(user.Name)
			if len(parts) > 0 {
				firstName = parts[0]
			}
			if len(parts) > 1 {
				lastName = strings.Join(parts[1:], " ")
			}
		}

		// Default values if still empty
		if firstName == "" {
			firstName = "User"
		}
		if lastName == "" {
			lastName = "Name"
		}

		// Create new user (no password for OAuth users)
		dbUser = models.User{
			Email:     user.Email,
			Password:  nil, // OAuth users don't have passwords
			FirstName: firstName,
			LastName:  lastName,
			Role:      "user",
			IsActive:  true,
			Provider:  "google",
		}

		if err := database.DB.Create(&dbUser).Error; err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
			return
		}
	} else {
		// User exists - check if they're using the same provider
		if dbUser.Provider != "google" && dbUser.Provider != "local" {
			// User exists with different provider - could merge accounts or return error
			utils.WriteErrorResponse(w, http.StatusConflict, "An account with this email already exists with a different login method")
			return
		}

		// Update provider if it was local (allowing both local and OAuth login)
		if dbUser.Provider == "local" {
			// Optionally update provider or allow both
			// For now, we'll allow both methods
		} else if dbUser.Provider == "" {
			// Legacy user without provider set
			database.DB.Model(&dbUser).Update("provider", "google")
		}

		// Check if account is active
		if !dbUser.IsActive {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Account is deactivated")
			return
		}
	}

	// Generate JWT token for the user
	token, err := utils.GenerateToken(dbUser.ID, dbUser.Email, dbUser.Role, utils.GetConfig())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return token and user info
	response := models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        dbUser.ID,
			Email:     dbUser.Email,
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			Role:      dbUser.Role,
			IsActive:  dbUser.IsActive,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
		},
	}

	utils.WriteSuccessResponse(w, "Google login successful", response)
}

// RegisterRoutes sets up all authentication routes
func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()

	// Public routes (no authentication required)
	authRouter.HandleFunc("/register", h.Register).Methods("POST")
	authRouter.HandleFunc("/login", h.Login).Methods("POST")
	authRouter.HandleFunc("/forgot-password", h.ForgotPassword).Methods("POST")
	authRouter.HandleFunc("/reset-password", h.ResetPassword).Methods("POST")

	// OAuth routes
	authRouter.HandleFunc("/google", h.BeginAuth).Methods("GET")
	authRouter.HandleFunc("/google/callback", h.Callback).Methods("GET")

	// Protected routes (require authentication via AuthMiddleware)
	protected := authRouter.PathPrefix("").Subrouter()
	// Note: AuthMiddleware should be applied in routes.go for protected routes
	protected.HandleFunc("/profile", h.GetProfile).Methods("GET")
	protected.HandleFunc("/refresh", h.RefreshToken).Methods("POST")
	protected.HandleFunc("/change-password", h.ChangePassword).Methods("POST")
	protected.HandleFunc("/logout", h.Logout).Methods("POST")
}
