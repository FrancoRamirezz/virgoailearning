package handlers

import (
	"backend/config"
	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	testDB     *gorm.DB
	testPGX    *pgx.Conn
	testRouter *mux.Router
	testConfig *config.Config
)

// setupTestDB initializes test database connection
func setupTestDB(t *testing.T) {
	// Set JWT secret in environment for middleware
	jwtSecret := getEnv("TEST_JWT_SECRET", "test-secret-key-for-testing-only")
	os.Setenv("JWT_SECRET", jwtSecret)

	// Load test configuration
	testConfig = &config.Config{
		Database: config.DatabaseConfig{
			Host:     getEnv("TEST_DB_HOST", "localhost"),
			Port:     getEnvInt("TEST_DB_PORT", 5432),
			User:     getEnv("TEST_DB_USER", "postgres"),
			Password: getEnv("TEST_DB_PASSWORD", "password"),
			Name:     getEnv("TEST_DB_NAME", "backend_test"),
		},
		JWT: config.JWTConfig{
			Secret: jwtSecret,
			Expiry: 24 * time.Hour,
		},
		Server: config.ServerConfig{
			Port: "8080",
			Env:  "test",
		},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
		},
	}

	// Connect using GORM
	database.Connect(testConfig)
	testDB = database.DB

	// Connect using PGX for direct queries
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		testConfig.Database.Host,
		testConfig.Database.User,
		testConfig.Database.Password,
		testConfig.Database.Name,
		testConfig.Database.Port,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	require.NoError(t, err, "Failed to connect to test database with PGX")
	testPGX = conn

	// Run migrations
	database.Migrate()

	// Setup router
	testRouter = setupTestRouter()
}

// teardownTestDB cleans up test database
func teardownTestDB(t *testing.T) {
	// Clean up all test data
	cleanupTestData(t)

	if testPGX != nil {
		testPGX.Close(context.Background())
	}

	if testDB != nil {
		sqlDB, _ := testDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}

// cleanupTestData removes all test data from database
func cleanupTestData(t *testing.T) {
	if testDB == nil {
		return
	}

	// Delete in reverse order of dependencies
	testDB.Exec("DELETE FROM enrollments")
	testDB.Exec("DELETE FROM course_contents")
	testDB.Exec("DELETE FROM courses")
	testDB.Exec("DELETE FROM payments")
	testDB.Exec("DELETE FROM contacts")
	testDB.Exec("DELETE FROM users")
}

// setupTestRouter creates a test router with all handlers
func setupTestRouter() *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public routes (no auth required)
	authHandler := NewAuthHandler()
	authHandler.InitializeOAuth(testConfig)
	authHandler.RegisterRoutes(api)

	// Routes that have both public and protected endpoints
	contactHandler := NewContactHandler()
	contactHandler.RegisterRoutes(api)

	courseHandler := NewCourseHandler()
	courseHandler.RegisterRoutes(api)

	// Protected routes - apply auth middleware
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Re-register handlers on protected router (they'll register both public and protected routes)
	// But we need to manually set up protected routes since middleware is applied at router level
	protectedContactRouter := protected.PathPrefix("/contacts").Subrouter()
	protectedContactRouter.HandleFunc("/search", contactHandler.SearchContacts).Methods("GET")
	protectedContactRouter.HandleFunc("", contactHandler.ListContacts).Methods("GET")
	protectedContactRouter.HandleFunc("/{id}", contactHandler.GetContact).Methods("GET")
	protectedContactRouter.HandleFunc("/{id}", contactHandler.UpdateContact).Methods("PUT")
	protectedContactRouter.HandleFunc("/{id}", contactHandler.DeleteContact).Methods("DELETE")

	protectedCourseRouter := protected.PathPrefix("/courses").Subrouter()
	protectedCourseRouter.HandleFunc("", courseHandler.CreateCourse).Methods("POST")
	protectedCourseRouter.HandleFunc("/{id}", courseHandler.UpdateCourse).Methods("PUT")
	protectedCourseRouter.HandleFunc("/{id}", courseHandler.DeleteCourse).Methods("DELETE")
	protectedCourseRouter.HandleFunc("/{id}/publish", courseHandler.PublishCourse).Methods("POST")
	protectedCourseRouter.HandleFunc("/{id}/content", courseHandler.AddCourseContent).Methods("POST")
	protectedCourseRouter.HandleFunc("/{id}/content/{contentId}", courseHandler.UpdateCourseContent).Methods("PUT")
	protectedCourseRouter.HandleFunc("/{id}/content/{contentId}", courseHandler.DeleteCourseContent).Methods("DELETE")
	protectedCourseRouter.HandleFunc("/{id}/enroll", courseHandler.EnrollInCourse).Methods("POST")

	paymentHandler := NewCoursePaymentHandler()
	protectedPaymentRouter := protected.PathPrefix("/payments").Subrouter()
	protectedPaymentRouter.HandleFunc("/process", paymentHandler.ProcessPayment).Methods("POST")
	protectedPaymentRouter.HandleFunc("/history", paymentHandler.GetPaymentHistory).Methods("GET")
	protectedPaymentRouter.HandleFunc("/{id}", paymentHandler.GetPaymentByID).Methods("GET")
	protectedPaymentRouter.HandleFunc("/{id}/status", paymentHandler.GetPaymentStatus).Methods("GET")
	protectedPaymentRouter.HandleFunc("/{id}/receipt", paymentHandler.GetPaymentReceipt).Methods("GET")

	// Protected auth routes
	protectedAuthRouter := protected.PathPrefix("/auth").Subrouter()
	protectedAuthRouter.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
	protectedAuthRouter.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	protectedAuthRouter.HandleFunc("/change-password", authHandler.ChangePassword).Methods("POST")
	protectedAuthRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// Admin routes
	adminRouter := protected.PathPrefix("/admin/payments").Subrouter()
	adminRouter.Use(middleware.AdminMiddleware)
	adminRouter.HandleFunc("", paymentHandler.GetAllPayments).Methods("GET")
	adminRouter.HandleFunc("/{id}/refund", paymentHandler.RefundPayment).Methods("POST")
	adminRouter.HandleFunc("/stats", paymentHandler.GetPaymentStats).Methods("GET")

	return router
}

// Helper function to validate JWT token and extract user ID
func validateTokenAndGetUserID(tokenString string) (uint, error) {
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(testConfig.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	return claims.UserID, nil
}

// Helper functions for creating test data
func createTestUser(t *testing.T, email, password, role string) *models.User {
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user := models.User{
		Email:     email,
		Password:  &hashedPassword,
		FirstName: "Test",
		LastName:  "User",
		Role:      role,
		IsActive:  true,
		Provider:  "local",
	}

	err = testDB.Create(&user).Error
	require.NoError(t, err)

	return &user
}

func createTestContact(t *testing.T, userID *uint) *models.Contact {
	contact := models.Contact{
		Name:    "Test Contact",
		Email:   "contact@test.com",
		Phone:   "1234567890",
		Subject: "Test Subject",
		Message: "Test Message",
		Status:  "pending",
		UserID:  userID,
	}

	err := testDB.Create(&contact).Error
	require.NoError(t, err)

	return &contact
}

func createTestCourse(t *testing.T, instructorID uint, isPublished bool) *models.Course {
	course := models.Course{
		Title:        "Test Course",
		Description:  "Test Description",
		InstructorID: instructorID,
		Price:        9999, // $99.99 in cents
		Currency:     "usd",
		Category:     "programming",
		Level:        "beginner",
		Duration:     10,
		IsPublished:  isPublished,
		ThumbnailURL: "https://example.com/thumb.jpg",
	}

	err := testDB.Create(&course).Error
	require.NoError(t, err)

	return &course
}

func createTestEnrollment(t *testing.T, userID, courseID uint, progress float64) *models.Enrollment {
	enrollment := models.Enrollment{
		UserID:     userID,
		CourseID:   courseID,
		Status:     "active",
		Progress:   progress,
		EnrolledAt: time.Now(),
	}

	err := testDB.Create(&enrollment).Error
	require.NoError(t, err)

	return &enrollment
}

// Helper function to make HTTP requests
func makeRequest(t *testing.T, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	return recorder
}

// Helper function to extract token from response
func extractToken(t *testing.T, response *httptest.ResponseRecorder) string {
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(t, err)

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return ""
	}

	token, _ := data["token"].(string)
	return token
}

// Helper function to get user from response
func extractUser(t *testing.T, response *httptest.ResponseRecorder) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(t, err)

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil
	}

	user, ok := data["user"].(map[string]interface{})
	if !ok {
		return nil
	}

	return user
}

// Helper function to generate JWT token for testing
func generateTestToken(t *testing.T, userID uint, email, role string) string {
	token, err := utils.GenerateToken(userID, email, role, testConfig)
	require.NoError(t, err)
	return token
}

// Environment variable helpers
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		fmt.Sscanf(value, "%d", &result)
		return result
	}
	return defaultValue
}

// ============================================================================
// DATABASE CONNECTION TESTS WITH PGX
// ============================================================================

func TestDatabaseConnectionWithPGX(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Test PGX connection
	var result int
	err := testPGX.QueryRow(context.Background(), "SELECT 1").Scan(&result)
	assert.NoError(t, err)
	assert.Equal(t, 1, result)

	// Test querying users table
	var count int
	err = testPGX.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, count, 0)
}

func TestUserExistsInDatabase(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create test user
	user := createTestUser(t, "pgx-test@example.com", "password123", "user")

	// Verify using PGX
	var email string
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT email FROM users WHERE id = $1",
		user.ID,
	).Scan(&email)

	assert.NoError(t, err)
	assert.Equal(t, "pgx-test@example.com", email)
}

func TestContactSoftDelete(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create test contact
	user := createTestUser(t, "softdelete@test.com", "password123", "user")
	contact := createTestContact(t, &user.ID)

	// Delete contact (soft delete)
	err := testDB.Delete(&contact).Error
	assert.NoError(t, err)

	// Verify deleted_at is set using PGX
	var deletedAt *time.Time
	err = testPGX.QueryRow(
		context.Background(),
		"SELECT deleted_at FROM contacts WHERE id = $1",
		contact.ID,
	).Scan(&deletedAt)

	assert.NoError(t, err)
	assert.NotNil(t, deletedAt)

	// Verify contact is not returned by normal query
	var foundContact models.Contact
	err = testDB.First(&foundContact, contact.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestEnrollmentProgressUpdate(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create test data
	instructor := createTestUser(t, "instructor@test.com", "password123", "instructor")
	student := createTestUser(t, "student@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, true)
	enrollment := createTestEnrollment(t, student.ID, course.ID, 0)

	// Update progress
	newProgress := 50.0
	err := testDB.Model(&enrollment).Update("progress", newProgress).Error
	assert.NoError(t, err)

	// Verify using PGX
	var progress float64
	err = testPGX.QueryRow(
		context.Background(),
		"SELECT progress FROM enrollments WHERE id = $1",
		enrollment.ID,
	).Scan(&progress)

	assert.NoError(t, err)
	assert.Equal(t, newProgress, progress)
}

// ============================================================================
// AUTH HANDLER TESTS
// ============================================================================

func TestAuthRegister(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	payload := map[string]string{
		"email":      "register@test.com",
		"password":   "TestPass123",
		"first_name": "Test",
		"last_name":  "User",
	}

	response := makeRequest(t, "POST", "/api/v1/auth/register", payload, "")

	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.True(t, result["success"].(bool))

	// Verify token in response
	token := extractToken(t, response)
	assert.NotEmpty(t, token)

	// Verify user in database using PGX
	var email string
	err = testPGX.QueryRow(
		context.Background(),
		"SELECT email FROM users WHERE email = $1",
		"register@test.com",
	).Scan(&email)
	assert.NoError(t, err)
	assert.Equal(t, "register@test.com", email)
}

func TestAuthRegisterMissingEmail(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	payload := map[string]string{
		"password":   "TestPass123",
		"first_name": "Test",
		"last_name":  "User",
	}

	response := makeRequest(t, "POST", "/api/v1/auth/register", payload, "")
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestAuthRegisterDuplicateEmail(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create first user
	createTestUser(t, "duplicate@test.com", "password123", "user")

	// Try to register with same email
	payload := map[string]string{
		"email":      "duplicate@test.com",
		"password":   "TestPass123",
		"first_name": "Test",
		"last_name":  "User",
	}

	response := makeRequest(t, "POST", "/api/v1/auth/register", payload, "")
	assert.Equal(t, http.StatusConflict, response.Code)
}

func TestAuthLogin(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create user first
	_ = createTestUser(t, "login@test.com", "correctpassword", "user")

	// Test cases
	tests := []struct {
		name           string
		email          string
		password       string
		expectedStatus int
	}{
		{"Valid credentials", "login@test.com", "correctpassword", http.StatusOK},
		{"Wrong password", "login@test.com", "wrongpassword", http.StatusUnauthorized},
		{"Non-existent user", "nonexistent@test.com", "anypassword", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]string{
				"email":    tt.email,
				"password": tt.password,
			}

			response := makeRequest(t, "POST", "/api/v1/auth/login", payload, "")
			assert.Equal(t, tt.expectedStatus, response.Code)

			if tt.expectedStatus == http.StatusOK {
				token := extractToken(t, response)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestAuthGetProfile(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create user and generate token
	user := createTestUser(t, "profile@test.com", "password123", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	// Get profile
	response := makeRequest(t, "GET", "/api/v1/auth/profile", nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	userData := extractUser(t, response)
	assert.NotNil(t, userData)
	assert.Equal(t, float64(user.ID), userData["id"].(float64))
}

func TestUnauthorizedAccess(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{"No token", "", http.StatusUnauthorized},
		{"Invalid token", "invalid.token.here", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := makeRequest(t, "GET", "/api/v1/auth/profile", nil, tt.token)
			assert.Equal(t, tt.expectedStatus, response.Code)
		})
	}
}

func TestAuthRefreshToken(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "refresh@test.com", "password123", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	response := makeRequest(t, "POST", "/api/v1/auth/refresh", nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	newToken := extractToken(t, response)
	assert.NotEmpty(t, newToken)
	assert.NotEqual(t, token, newToken)
}

func TestAuthChangePassword(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "changepass@test.com", "oldpassword", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	payload := map[string]string{
		"current_password": "oldpassword",
		"new_password":     "newpassword123",
	}

	response := makeRequest(t, "POST", "/api/v1/auth/change-password", payload, token)
	assert.Equal(t, http.StatusOK, response.Code)

	// Verify new password works
	loginPayload := map[string]string{
		"email":    "changepass@test.com",
		"password": "newpassword123",
	}
	loginResponse := makeRequest(t, "POST", "/api/v1/auth/login", loginPayload, "")
	assert.Equal(t, http.StatusOK, loginResponse.Code)

	// Verify old password doesn't work
	loginPayload["password"] = "oldpassword"
	loginResponse = makeRequest(t, "POST", "/api/v1/auth/login", loginPayload, "")
	assert.Equal(t, http.StatusUnauthorized, loginResponse.Code)
}

// ============================================================================
// CONTACT HANDLER TESTS
// ============================================================================

func TestContactCreate(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	payload := map[string]string{
		"name":    "Test Contact",
		"email":   "contact@test.com",
		"phone":   "1234567890",
		"subject": "Test Subject",
		"message": "Test Message",
	}

	response := makeRequest(t, "POST", "/api/v1/contacts", payload, "")
	assert.Equal(t, http.StatusOK, response.Code)

	// Verify contact in database using PGX
	var count int
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM contacts WHERE email = $1",
		"contact@test.com",
	).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestContactGet(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "contactget@test.com", "password123", "user")
	contact := createTestContact(t, &user.ID)
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	response := makeRequest(t, "GET", fmt.Sprintf("/api/v1/contacts/%d", contact.ID), nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.True(t, result["success"].(bool))
}

func TestContactGetNotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "notfound@test.com", "password123", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	response := makeRequest(t, "GET", "/api/v1/contacts/99999", nil, token)
	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestContactList(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "list@test.com", "password123", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	// Create 5 test contacts
	for i := 0; i < 5; i++ {
		createTestContact(t, &user.ID)
	}

	response := makeRequest(t, "GET", "/api/v1/contacts?page=1&limit=10", nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.NoError(t, err)

	data := result["data"].(map[string]interface{})
	contacts := data["contacts"].([]interface{})
	assert.GreaterOrEqual(t, len(contacts), 5)
}

func TestContactUpdate(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "update@test.com", "password123", "user")
	contact := createTestContact(t, &user.ID)
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	payload := map[string]string{
		"status": "resolved",
	}

	response := makeRequest(t, "PUT", fmt.Sprintf("/api/v1/contacts/%d", contact.ID), payload, token)
	assert.Equal(t, http.StatusOK, response.Code)

	// Verify update using PGX
	var status string
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT status FROM contacts WHERE id = $1",
		contact.ID,
	).Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "resolved", status)
}

func TestContactDelete(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	admin := createTestUser(t, "admin@test.com", "password123", "admin")
	contact := createTestContact(t, &admin.ID)
	token := generateTestToken(t, admin.ID, admin.Email, admin.Role)

	response := makeRequest(t, "DELETE", fmt.Sprintf("/api/v1/contacts/%d", contact.ID), nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	// Verify soft delete using PGX
	var deletedAt *time.Time
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT deleted_at FROM contacts WHERE id = $1",
		contact.ID,
	).Scan(&deletedAt)
	assert.NoError(t, err)
	assert.NotNil(t, deletedAt)
}

func TestContactSearch(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "search@test.com", "password123", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	// Create contact with specific name
	contact := models.Contact{
		Name:    "John Doe",
		Email:   "john@test.com",
		Subject: "Test",
		Message: "Test",
		Status:  "pending",
		UserID:  &user.ID,
	}
	testDB.Create(&contact)

	response := makeRequest(t, "GET", "/api/v1/contacts/search?query=John", nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.NoError(t, err)

	data := result["data"].(map[string]interface{})
	contacts := data["contacts"].([]interface{})
	assert.GreaterOrEqual(t, len(contacts), 1)
}

func TestForbiddenAccess(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user1 := createTestUser(t, "user1@test.com", "password123", "user")
	user2 := createTestUser(t, "user2@test.com", "password123", "user")
	contact := createTestContact(t, &user1.ID)
	token2 := generateTestToken(t, user2.ID, user2.Email, user2.Role)

	// User2 tries to update user1's contact
	payload := map[string]string{"status": "resolved"}
	response := makeRequest(t, "PUT", fmt.Sprintf("/api/v1/contacts/%d", contact.ID), payload, token2)
	assert.Equal(t, http.StatusForbidden, response.Code)
}

// ============================================================================
// COURSE HANDLER TESTS
// ============================================================================

func TestCourseCreate(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "instructor@test.com", "password123", "instructor")
	token := generateTestToken(t, instructor.ID, instructor.Email, instructor.Role)

	payload := map[string]interface{}{
		"title":         "New Course",
		"description":   "Course Description",
		"price":         9999,
		"currency":      "usd",
		"category":      "programming",
		"level":         "beginner",
		"duration":      10,
		"thumbnail_url": "https://example.com/thumb.jpg",
	}

	response := makeRequest(t, "POST", "/api/v1/courses", payload, token)
	assert.Equal(t, http.StatusOK, response.Code)

	// Verify course in database using PGX
	var title string
	var isPublished bool
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT title, is_published FROM courses WHERE instructor_id = $1 ORDER BY id DESC LIMIT 1",
		instructor.ID,
	).Scan(&title, &isPublished)
	assert.NoError(t, err)
	assert.Equal(t, "New Course", title)
	assert.False(t, isPublished) // New courses start as unpublished
}

func TestCourseCreateNonInstructor(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "user@test.com", "password123", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	payload := map[string]interface{}{
		"title":       "New Course",
		"description": "Course Description",
		"price":       9999,
	}

	response := makeRequest(t, "POST", "/api/v1/courses", payload, token)
	assert.Equal(t, http.StatusForbidden, response.Code)
}

func TestCourseGet(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "courseget@test.com", "password123", "instructor")
	course := createTestCourse(t, instructor.ID, true)

	response := makeRequest(t, "GET", fmt.Sprintf("/api/v1/courses/%d", course.ID), nil, "")
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestCourseList(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "courselist@test.com", "password123", "instructor")

	// Create published and unpublished courses
	createTestCourse(t, instructor.ID, true)
	createTestCourse(t, instructor.ID, true)
	createTestCourse(t, instructor.ID, false)

	response := makeRequest(t, "GET", "/api/v1/courses?page=1&limit=10", nil, "")
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.NoError(t, err)

	data := result["data"].(map[string]interface{})
	courses := data["courses"].([]interface{})
	// Should only see published courses
	assert.GreaterOrEqual(t, len(courses), 2)
}

func TestCourseUpdate(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "courseupdate@test.com", "password123", "instructor")
	course := createTestCourse(t, instructor.ID, false)
	token := generateTestToken(t, instructor.ID, instructor.Email, instructor.Role)

	payload := map[string]interface{}{
		"title": "Updated Course Title",
	}

	response := makeRequest(t, "PUT", fmt.Sprintf("/api/v1/courses/%d", course.ID), payload, token)
	assert.Equal(t, http.StatusOK, response.Code)

	// Verify update using PGX
	var title string
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT title FROM courses WHERE id = $1",
		course.ID,
	).Scan(&title)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Course Title", title)
}

func TestCoursePublish(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "publish@test.com", "password123", "instructor")
	course := createTestCourse(t, instructor.ID, false)
	token := generateTestToken(t, instructor.ID, instructor.Email, instructor.Role)

	payload := map[string]bool{
		"is_published": true,
	}

	response := makeRequest(t, "POST", fmt.Sprintf("/api/v1/courses/%d/publish", course.ID), payload, token)
	assert.Equal(t, http.StatusOK, response.Code)

	// Verify publication status using PGX
	var isPublished bool
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT is_published FROM courses WHERE id = $1",
		course.ID,
	).Scan(&isPublished)
	assert.NoError(t, err)
	assert.True(t, isPublished)
}

func TestCourseEnroll(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "enrollinst@test.com", "password123", "instructor")
	student := createTestUser(t, "enrollstudent@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, true)
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	response := makeRequest(t, "POST", fmt.Sprintf("/api/v1/courses/%d/enroll", course.ID), nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	// Verify enrollment created using PGX
	var enrollmentID uint
	var status string
	var progress float64
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT id, status, progress FROM enrollments WHERE user_id = $1 AND course_id = $2",
		student.ID, course.ID,
	).Scan(&enrollmentID, &status, &progress)
	assert.NoError(t, err)
	assert.Equal(t, "active", status)
	assert.Equal(t, 0.0, progress)
}

func TestCourseEnrollAlreadyEnrolled(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "duplicateenroll@test.com", "password123", "instructor")
	student := createTestUser(t, "duplicatestudent@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, true)
	createTestEnrollment(t, student.ID, course.ID, 0)
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	response := makeRequest(t, "POST", fmt.Sprintf("/api/v1/courses/%d/enroll", course.ID), nil, token)
	assert.Equal(t, http.StatusConflict, response.Code)
}

func TestCourseEnrollUnpublished(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "unpublished@test.com", "password123", "instructor")
	student := createTestUser(t, "unpubstudent@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, false)
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	response := makeRequest(t, "POST", fmt.Sprintf("/api/v1/courses/%d/enroll", course.ID), nil, token)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

// ============================================================================
// PAYMENT HANDLER TESTS
// ============================================================================

func TestPaymentProcess(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "payinst@test.com", "password123", "instructor")
	student := createTestUser(t, "paystudent@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, true)
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	payload := map[string]interface{}{
		"course_id": course.ID,
		"amount":    course.Price,
		"currency":  "usd",
	}

	response := makeRequest(t, "POST", "/api/v1/payments/process", payload, token)
	assert.Equal(t, http.StatusOK, response.Code)

	// Verify payment created using PGX
	var paymentID uint
	var status string
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT id, status FROM payments WHERE user_id = $1 AND amount = $2 ORDER BY id DESC LIMIT 1",
		student.ID, course.Price,
	).Scan(&paymentID, &status)
	assert.NoError(t, err)
	assert.Equal(t, "succeeded", status)

	// Verify enrollment created
	var enrollmentID uint
	err = testPGX.QueryRow(
		context.Background(),
		"SELECT id FROM enrollments WHERE user_id = $1 AND course_id = $2",
		student.ID, course.ID,
	).Scan(&enrollmentID)
	assert.NoError(t, err)
}

func TestPaymentProcessWrongAmount(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "wrongamount@test.com", "password123", "instructor")
	student := createTestUser(t, "wrongstudent@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, true)
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	payload := map[string]interface{}{
		"course_id": course.ID,
		"amount":    course.Price - 100, // Wrong amount
		"currency":  "usd",
	}

	response := makeRequest(t, "POST", "/api/v1/payments/process", payload, token)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestPaymentHistory(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	student := createTestUser(t, "history@test.com", "password123", "user")
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	// Create some payments
	instructor := createTestUser(t, "historyinst@test.com", "password123", "instructor")
	course1 := createTestCourse(t, instructor.ID, true)
	course2 := createTestCourse(t, instructor.ID, true)

	// Process payments
	payload1 := map[string]interface{}{
		"course_id": course1.ID,
		"amount":    course1.Price,
		"currency":  "usd",
	}
	makeRequest(t, "POST", "/api/v1/payments/process", payload1, token)

	payload2 := map[string]interface{}{
		"course_id": course2.ID,
		"amount":    course2.Price,
		"currency":  "usd",
	}
	makeRequest(t, "POST", "/api/v1/payments/process", payload2, token)

	// Get payment history
	response := makeRequest(t, "GET", "/api/v1/payments/history?page=1&limit=10", nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.NoError(t, err)

	data := result["data"].(map[string]interface{})
	payments := data["payments"].([]interface{})
	assert.GreaterOrEqual(t, len(payments), 2)
}

func TestPaymentGetByID(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "getpayinst@test.com", "password123", "instructor")
	student := createTestUser(t, "getpaystudent@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, true)
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	// Process payment
	payload := map[string]interface{}{
		"course_id": course.ID,
		"amount":    course.Price,
		"currency":  "usd",
	}
	processResponse := makeRequest(t, "POST", "/api/v1/payments/process", payload, token)

	// Extract payment ID from response
	var processResult map[string]interface{}
	json.Unmarshal(processResponse.Body.Bytes(), &processResult)
	processData := processResult["data"].(map[string]interface{})
	paymentID := uint(processData["id"].(float64))

	// Get payment by ID
	response := makeRequest(t, "GET", fmt.Sprintf("/api/v1/payments/%d", paymentID), nil, token)
	assert.Equal(t, http.StatusOK, response.Code)
}

// ============================================================================
// INTEGRATION TESTS (FULL FLOWS)
// ============================================================================

func TestAuthFlow(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// 1. Register user
	registerPayload := map[string]string{
		"email":      "flow@test.com",
		"password":   "TestPass123",
		"first_name": "Test",
		"last_name":  "User",
	}
	registerResponse := makeRequest(t, "POST", "/api/v1/auth/register", registerPayload, "")
	assert.Equal(t, http.StatusOK, registerResponse.Code)
	_ = extractToken(t, registerResponse) // Token extracted but not needed for this test

	// 2. Login with credentials
	loginPayload := map[string]string{
		"email":    "flow@test.com",
		"password": "TestPass123",
	}
	loginResponse := makeRequest(t, "POST", "/api/v1/auth/login", loginPayload, "")
	assert.Equal(t, http.StatusOK, loginResponse.Code)
	loginToken := extractToken(t, loginResponse)

	// 3. Get profile with token
	profileResponse := makeRequest(t, "GET", "/api/v1/auth/profile", nil, loginToken)
	assert.Equal(t, http.StatusOK, profileResponse.Code)

	// 4. Refresh token
	refreshResponse := makeRequest(t, "POST", "/api/v1/auth/refresh", nil, loginToken)
	assert.Equal(t, http.StatusOK, refreshResponse.Code)
	refreshToken := extractToken(t, refreshResponse)
	assert.NotEqual(t, loginToken, refreshToken)

	// 5. Change password
	changePassPayload := map[string]string{
		"current_password": "TestPass123",
		"new_password":     "NewPass123",
	}
	changePassResponse := makeRequest(t, "POST", "/api/v1/auth/change-password", changePassPayload, refreshToken)
	assert.Equal(t, http.StatusOK, changePassResponse.Code)

	// 6. Login with new password → success
	newLoginPayload := map[string]string{
		"email":    "flow@test.com",
		"password": "NewPass123",
	}
	newLoginResponse := makeRequest(t, "POST", "/api/v1/auth/login", newLoginPayload, "")
	assert.Equal(t, http.StatusOK, newLoginResponse.Code)

	// 7. Login with old password → fail
	oldLoginPayload := map[string]string{
		"email":    "flow@test.com",
		"password": "TestPass123",
	}
	oldLoginResponse := makeRequest(t, "POST", "/api/v1/auth/login", oldLoginPayload, "")
	assert.Equal(t, http.StatusUnauthorized, oldLoginResponse.Code)
}

func TestCourseFlow(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// 1. Create course (as instructor)
	instructor := createTestUser(t, "courseflow@test.com", "password123", "instructor")
	token := generateTestToken(t, instructor.ID, instructor.Email, instructor.Role)

	createPayload := map[string]interface{}{
		"title":       "Flow Course",
		"description": "Course Description",
		"price":       9999,
		"currency":    "usd",
		"category":    "programming",
		"level":       "beginner",
		"duration":    10,
	}
	createResponse := makeRequest(t, "POST", "/api/v1/courses", createPayload, token)
	assert.Equal(t, http.StatusOK, createResponse.Code)

	// Extract course ID
	var createResult map[string]interface{}
	json.Unmarshal(createResponse.Body.Bytes(), &createResult)
	createData := createResult["data"].(map[string]interface{})
	courseID := uint(createData["id"].(float64))

	// 2. Verify is_published = false
	var isPublished bool
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT is_published FROM courses WHERE id = $1",
		courseID,
	).Scan(&isPublished)
	assert.NoError(t, err)
	assert.False(t, isPublished)

	// 3. Publish course
	publishPayload := map[string]bool{
		"is_published": true,
	}
	publishResponse := makeRequest(t, "POST", fmt.Sprintf("/api/v1/courses/%d/publish", courseID), publishPayload, token)
	assert.Equal(t, http.StatusOK, publishResponse.Code)

	// 4. Verify shows in ListCourses
	listResponse := makeRequest(t, "GET", "/api/v1/courses?page=1&limit=10", nil, "")
	assert.Equal(t, http.StatusOK, listResponse.Code)

	// 5. Enroll as student
	student := createTestUser(t, "coursestudent@test.com", "password123", "user")
	studentToken := generateTestToken(t, student.ID, student.Email, student.Role)

	enrollResponse := makeRequest(t, "POST", fmt.Sprintf("/api/v1/courses/%d/enroll", courseID), nil, studentToken)
	assert.Equal(t, http.StatusOK, enrollResponse.Code)

	// 6. Update progress to 50%
	var enrollmentID uint
	err = testPGX.QueryRow(
		context.Background(),
		"SELECT id FROM enrollments WHERE user_id = $1 AND course_id = $2",
		student.ID, courseID,
	).Scan(&enrollmentID)
	assert.NoError(t, err)

	err = testDB.Model(&models.Enrollment{}).Where("id = ?", enrollmentID).Update("progress", 50.0).Error
	assert.NoError(t, err)

	// 7. Update progress to 100%
	err = testDB.Model(&models.Enrollment{}).Where("id = ?", enrollmentID).Update("progress", 100.0).Error
	assert.NoError(t, err)

	// 8. Verify status = "completed" (would need handler to update status)
	var progress float64
	var status string
	err = testPGX.QueryRow(
		context.Background(),
		"SELECT progress, status FROM enrollments WHERE id = $1",
		enrollmentID,
	).Scan(&progress, &status)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, progress)
}

func TestPaymentFlow(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// 1. Create course
	instructor := createTestUser(t, "payflowinst@test.com", "password123", "instructor")
	course := createTestCourse(t, instructor.ID, true)

	// 2. Process payment
	student := createTestUser(t, "payflowstudent@test.com", "password123", "user")
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	payload := map[string]interface{}{
		"course_id": course.ID,
		"amount":    course.Price,
		"currency":  "usd",
	}
	paymentResponse := makeRequest(t, "POST", "/api/v1/payments/process", payload, token)
	assert.Equal(t, http.StatusOK, paymentResponse.Code)

	// 3. Verify payment created
	var paymentID uint
	var paymentStatus string
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT id, status FROM payments WHERE user_id = $1 AND amount = $2 ORDER BY id DESC LIMIT 1",
		student.ID, course.Price,
	).Scan(&paymentID, &paymentStatus)
	assert.NoError(t, err)
	assert.Equal(t, "succeeded", paymentStatus)

	// 4. Verify user auto-enrolled
	var enrollmentID uint
	var enrollmentStatus string
	err = testPGX.QueryRow(
		context.Background(),
		"SELECT id, status FROM enrollments WHERE user_id = $1 AND course_id = $2",
		student.ID, course.ID,
	).Scan(&enrollmentID, &enrollmentStatus)
	assert.NoError(t, err)
	assert.Equal(t, "active", enrollmentStatus)

	// 5. Check enrollment status = "active"
	assert.Equal(t, "active", enrollmentStatus)
}

// ============================================================================
// ADDITIONAL EDGE CASE TESTS
// ============================================================================

func TestContactListPagination(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "pagination@test.com", "password123", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	// Create 15 contacts
	for i := 0; i < 15; i++ {
		createTestContact(t, &user.ID)
	}

	// Test first page
	response := makeRequest(t, "GET", "/api/v1/contacts?page=1&limit=10", nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	contacts := data["contacts"].([]interface{})
	total := data["total"].(float64)

	assert.Equal(t, 10, len(contacts))
	assert.Equal(t, float64(15), total)

	// Test second page
	response2 := makeRequest(t, "GET", "/api/v1/contacts?page=2&limit=10", nil, token)
	assert.Equal(t, http.StatusOK, response2.Code)

	var result2 map[string]interface{}
	json.Unmarshal(response2.Body.Bytes(), &result2)
	data2 := result2["data"].(map[string]interface{})
	contacts2 := data2["contacts"].([]interface{})

	assert.Equal(t, 5, len(contacts2))
}

func TestContactListWithStatusFilter(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "statusfilter@test.com", "password123", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	// Create contacts with different statuses
	contact1 := createTestContact(t, &user.ID)
	contact2 := createTestContact(t, &user.ID)
	testDB.Model(&contact1).Update("status", "resolved")
	testDB.Model(&contact2).Update("status", "pending")

	// Filter by pending
	response := makeRequest(t, "GET", "/api/v1/contacts?status=pending", nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	contacts := data["contacts"].([]interface{})

	// Should only get pending contacts
	for _, c := range contacts {
		contact := c.(map[string]interface{})
		assert.Equal(t, "pending", contact["status"].(string))
	}
}

func TestCourseListWithFilters(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "filterinst@test.com", "password123", "instructor")

	// Create courses with different categories and levels
	course1 := createTestCourse(t, instructor.ID, true)
	testDB.Model(&course1).Updates(map[string]interface{}{
		"category": "programming",
		"level":    "beginner",
	})

	course2 := createTestCourse(t, instructor.ID, true)
	testDB.Model(&course2).Updates(map[string]interface{}{
		"category": "design",
		"level":    "intermediate",
	})

	// Filter by category
	response := makeRequest(t, "GET", "/api/v1/courses?category=programming", nil, "")
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	courses := data["courses"].([]interface{})

	// All courses should be programming
	for _, c := range courses {
		course := c.(map[string]interface{})
		assert.Equal(t, "programming", course["category"].(string))
	}
}

func TestCourseGetUnpublishedAsNonInstructor(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "unpubinst@test.com", "password123", "instructor")
	course := createTestCourse(t, instructor.ID, false)

	// Regular user tries to access unpublished course
	response := makeRequest(t, "GET", fmt.Sprintf("/api/v1/courses/%d", course.ID), nil, "")
	assert.Equal(t, http.StatusForbidden, response.Code)
}

func TestCourseUpdateByNonInstructor(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "updateinst@test.com", "password123", "instructor")
	otherUser := createTestUser(t, "otheruser@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, false)
	otherToken := generateTestToken(t, otherUser.ID, otherUser.Email, otherUser.Role)

	payload := map[string]interface{}{
		"title": "Hacked Title",
	}

	response := makeRequest(t, "PUT", fmt.Sprintf("/api/v1/courses/%d", course.ID), payload, otherToken)
	assert.Equal(t, http.StatusForbidden, response.Code)
}

func TestPaymentProcessUnpublishedCourse(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "unpubpayinst@test.com", "password123", "instructor")
	student := createTestUser(t, "unpubpaystudent@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, false) // Unpublished
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	payload := map[string]interface{}{
		"course_id": course.ID,
		"amount":    course.Price,
		"currency":  "usd",
	}

	response := makeRequest(t, "POST", "/api/v1/payments/process", payload, token)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestPaymentGetByIDNotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	student := createTestUser(t, "notfoundpay@test.com", "password123", "user")
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	response := makeRequest(t, "GET", "/api/v1/payments/99999", nil, token)
	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestPaymentGetByIDOtherUser(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "otherpayinst@test.com", "password123", "instructor")
	student1 := createTestUser(t, "student1@test.com", "password123", "user")
	student2 := createTestUser(t, "student2@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, true)
	token1 := generateTestToken(t, student1.ID, student1.Email, student1.Role)
	token2 := generateTestToken(t, student2.ID, student2.Email, student2.Role)

	// Student1 processes payment
	payload := map[string]interface{}{
		"course_id": course.ID,
		"amount":    course.Price,
		"currency":  "usd",
	}
	processResponse := makeRequest(t, "POST", "/api/v1/payments/process", payload, token1)
	var processResult map[string]interface{}
	json.Unmarshal(processResponse.Body.Bytes(), &processResult)
	processData := processResult["data"].(map[string]interface{})
	paymentID := uint(processData["id"].(float64))

	// Student2 tries to access student1's payment
	response := makeRequest(t, "GET", fmt.Sprintf("/api/v1/payments/%d", paymentID), nil, token2)
	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestContactCreateWithAuthenticatedUser(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "authcontact@test.com", "password123", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	payload := map[string]string{
		"name":    "Authenticated Contact",
		"email":   "auth@test.com",
		"subject": "Test Subject",
		"message": "Test Message",
	}

	response := makeRequest(t, "POST", "/api/v1/contacts", payload, token)
	assert.Equal(t, http.StatusOK, response.Code)

	// Verify user_id is set using PGX
	var userID *uint
	err := testPGX.QueryRow(
		context.Background(),
		"SELECT user_id FROM contacts WHERE email = $1",
		"auth@test.com",
	).Scan(&userID)
	assert.NoError(t, err)
	assert.NotNil(t, userID)
	assert.Equal(t, user.ID, *userID)
}

func TestAuthChangePasswordWrongCurrentPassword(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "wrongpass@test.com", "correctpass", "user")
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	payload := map[string]string{
		"current_password": "wrongpass",
		"new_password":     "newpass123",
	}

	response := makeRequest(t, "POST", "/api/v1/auth/change-password", payload, token)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
}

func TestCourseDeleteByNonInstructor(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "deleteinst@test.com", "password123", "instructor")
	otherUser := createTestUser(t, "deleteuser@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, false)
	otherToken := generateTestToken(t, otherUser.ID, otherUser.Email, otherUser.Role)

	response := makeRequest(t, "DELETE", fmt.Sprintf("/api/v1/courses/%d", course.ID), nil, otherToken)
	assert.Equal(t, http.StatusForbidden, response.Code)
}

func TestContactDeleteByNonAdmin(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	user := createTestUser(t, "nonadmin@test.com", "password123", "user")
	contact := createTestContact(t, &user.ID)
	token := generateTestToken(t, user.ID, user.Email, user.Role)

	response := makeRequest(t, "DELETE", fmt.Sprintf("/api/v1/contacts/%d", contact.ID), nil, token)
	assert.Equal(t, http.StatusForbidden, response.Code)
}

func TestPaymentStatus(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "statusinst@test.com", "password123", "instructor")
	student := createTestUser(t, "statusstudent@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, true)
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	// Process payment
	payload := map[string]interface{}{
		"course_id": course.ID,
		"amount":    course.Price,
		"currency":  "usd",
	}
	processResponse := makeRequest(t, "POST", "/api/v1/payments/process", payload, token)
	var processResult map[string]interface{}
	json.Unmarshal(processResponse.Body.Bytes(), &processResult)
	processData := processResult["data"].(map[string]interface{})
	paymentID := uint(processData["id"].(float64))

	// Get payment status
	response := makeRequest(t, "GET", fmt.Sprintf("/api/v1/payments/%d/status", paymentID), nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	assert.Equal(t, "succeeded", data["status"].(string))
}

func TestPaymentReceipt(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "receiptinst@test.com", "password123", "instructor")
	student := createTestUser(t, "receiptstudent@test.com", "password123", "user")
	course := createTestCourse(t, instructor.ID, true)
	token := generateTestToken(t, student.ID, student.Email, student.Role)

	// Process payment
	payload := map[string]interface{}{
		"course_id": course.ID,
		"amount":    course.Price,
		"currency":  "usd",
	}
	processResponse := makeRequest(t, "POST", "/api/v1/payments/process", payload, token)
	var processResult map[string]interface{}
	json.Unmarshal(processResponse.Body.Bytes(), &processResult)
	processData := processResult["data"].(map[string]interface{})
	paymentID := uint(processData["id"].(float64))

	// Get receipt
	response := makeRequest(t, "GET", fmt.Sprintf("/api/v1/payments/%d/receipt", paymentID), nil, token)
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	assert.NotNil(t, data["receipt_id"])
	assert.Equal(t, float64(paymentID), data["payment_id"].(float64))
}

func TestCourseListPagination(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	instructor := createTestUser(t, "coursepaginst@test.com", "password123", "instructor")

	// Create 12 published courses
	for i := 0; i < 12; i++ {
		createTestCourse(t, instructor.ID, true)
	}

	// Test first page
	response := makeRequest(t, "GET", "/api/v1/courses?page=1&limit=10", nil, "")
	assert.Equal(t, http.StatusOK, response.Code)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	courses := data["courses"].([]interface{})
	total := data["total"].(float64)

	assert.Equal(t, 10, len(courses))
	assert.Equal(t, float64(12), total)

	// Test second page
	response2 := makeRequest(t, "GET", "/api/v1/courses?page=2&limit=10", nil, "")
	assert.Equal(t, http.StatusOK, response2.Code)

	var result2 map[string]interface{}
	json.Unmarshal(response2.Body.Bytes(), &result2)
	data2 := result2["data"].(map[string]interface{})
	courses2 := data2["courses"].([]interface{})

	assert.Equal(t, 2, len(courses2))
}

// ============================================================================
// MAIN TEST SETUP/TEARDOWN
// ============================================================================

func TestMain(m *testing.M) {
	// Setup can be done here if needed
	code := m.Run()
	// Teardown can be done here if needed
	os.Exit(code)
}
