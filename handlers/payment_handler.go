package handlers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// CoursePaymentHandler handles all course payment processing operations
// This handler manages course payments, payment history, receipts, refunds, and statistics
type CoursePaymentHandler struct{}

func NewCoursePaymentHandler() *CoursePaymentHandler {
	return &CoursePaymentHandler{}
}

// 1. ProcessPayment - Process course payment
// POST /api/v1/payments/process
// Headers: Authorization: Bearer <token>
// Request: { course_id, amount, currency, payment_method_id }
// Response: { payment }
//
// Why: Allows users to purchase courses. Creates a payment record and links it to the course.
// In a full implementation, this would integrate with Stripe or another payment processor.
func (h *CoursePaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	var req struct {
		CourseID        uint   `json:"course_id" validate:"required"`
		Amount          int64  `json:"amount" validate:"required,min=1"`
		Currency        string `json:"currency" validate:"required"`
		PaymentMethodID string `json:"payment_method_id"`
		Description     string `json:"description"`
	}

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Verify course exists
	var course models.Course
	if err := database.DB.First(&course, req.CourseID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course not found")
		return
	}

	// Verify course is published
	if !course.IsPublished {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Course is not available for purchase")
		return
	}

	// Verify amount matches course price
	if req.Amount != course.Price {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Payment amount does not match course price")
		return
	}

	// Check if user is already enrolled
	var existingEnrollment models.Enrollment
	if err := database.DB.Where("user_id = ? AND course_id = ?", user.ID, req.CourseID).First(&existingEnrollment).Error; err == nil {
		utils.WriteErrorResponse(w, http.StatusConflict, "Already enrolled in this course")
		return
	}

	// Create payment record
	description := req.Description
	if description == "" {
		description = fmt.Sprintf("Payment for course: %s", course.Title)
	}

	payment := models.Payment{
		UserID:          user.ID,
		StripePaymentID: fmt.Sprintf("pay_%d_%d", user.ID, time.Now().Unix()),
		Amount:          req.Amount,
		Currency:        req.Currency,
		Status:          "succeeded",
		Description:     description,
		PaymentMethodID: req.PaymentMethodID,
	}

	// Store course_id in metadata
	metadata := map[string]interface{}{
		"course_id":    req.CourseID,
		"course_title": course.Title,
	}
	metadataBytes, _ := json.Marshal(metadata)
	payment.Metadata = string(metadataBytes)

	if err := database.DB.Create(&payment).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to process payment")
		return
	}

	// Create enrollment after successful payment
	enrollment := models.Enrollment{
		UserID:     user.ID,
		CourseID:   req.CourseID,
		Status:     "active",
		Progress:   0,
		EnrolledAt: time.Now(),
	}

	if err := database.DB.Create(&enrollment).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Payment processed but enrollment failed")
		return
	}

	// Load user relationship for response
	database.DB.Preload("User").First(&payment, payment.ID)

	response := models.PaymentResponse{
		ID:              payment.ID,
		StripePaymentID: payment.StripePaymentID,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		Status:          payment.Status,
		Description:     payment.Description,
		CreatedAt:       payment.CreatedAt,
	}

	utils.WriteSuccessResponse(w, "Payment processed successfully", response)
}

// 2. GetPaymentHistory - Get user's payment history
// GET /api/v1/payments/history
// Headers: Authorization: Bearer <token>
// Response: { payments: [] }
//
// Why: Allows users to view their complete payment history. Essential for tracking purchases
// and accessing receipts for past transactions.
func (h *CoursePaymentHandler) GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	// Parse query parameters for pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var payments []models.Payment
	var total int64

	// Get total count
	if err := database.DB.Model(&models.Payment{}).Where("user_id = ?", user.ID).Count(&total).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch payment history")
		return
	}

	// Get paginated payments
	if err := database.DB.Where("user_id = ?", user.ID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&payments).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch payment history")
		return
	}

	// Convert to response format
	paymentResponses := make([]models.PaymentResponse, len(payments))
	for i, payment := range payments {
		paymentResponses[i] = models.PaymentResponse{
			ID:              payment.ID,
			StripePaymentID: payment.StripePaymentID,
			Amount:          payment.Amount,
			Currency:        payment.Currency,
			Status:          payment.Status,
			Description:     payment.Description,
			CreatedAt:       payment.CreatedAt,
		}
	}

	response := map[string]interface{}{
		"payments": paymentResponses,
		"total":    total,
		"page":     page,
		"limit":    limit,
	}

	utils.WriteSuccessResponse(w, "Payment history retrieved successfully", response)
}

// 3. GetPaymentByID - Get specific payment
// GET /api/v1/payments/{id}
// Headers: Authorization: Bearer <token>
// Response: { payment }
//
// Why: Allows users to retrieve detailed information about a specific payment. Useful for
// viewing payment details, checking status, or accessing receipt information.
func (h *CoursePaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	vars := mux.Vars(r)
	paymentID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	var payment models.Payment
	if err := database.DB.Where("id = ? AND user_id = ?", paymentID, user.ID).First(&payment).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Payment not found")
		return
	}

	response := models.PaymentResponse{
		ID:              payment.ID,
		StripePaymentID: payment.StripePaymentID,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		Status:          payment.Status,
		Description:     payment.Description,
		CreatedAt:       payment.CreatedAt,
	}

	utils.WriteSuccessResponse(w, "Payment retrieved successfully", response)
}

// 4. GetPaymentStatus - Check payment status
// GET /api/v1/payments/{id}/status
// Headers: Authorization: Bearer <token>
// Response: { status, payment_id }
//
// Why: Provides a quick way to check the current status of a payment without retrieving
// all payment details. Useful for polling payment status after initiating a transaction.
func (h *CoursePaymentHandler) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	vars := mux.Vars(r)
	paymentID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	var payment models.Payment
	if err := database.DB.Where("id = ? AND user_id = ?", paymentID, user.ID).
		Select("id", "status", "stripe_payment_id").
		First(&payment).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Payment not found")
		return
	}

	response := map[string]interface{}{
		"payment_id": payment.ID,
		"status":     payment.Status,
	}

	utils.WriteSuccessResponse(w, "Payment status retrieved successfully", response)
}

// 5. GetPaymentReceipt - Generate receipt
// GET /api/v1/payments/{id}/receipt
// Headers: Authorization: Bearer <token>
// Response: { receipt }
//
// Why: Generates a detailed receipt for a completed payment. Includes all relevant
// information needed for accounting, tax purposes, or customer records.
func (h *CoursePaymentHandler) GetPaymentReceipt(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	vars := mux.Vars(r)
	paymentID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	var payment models.Payment
	if err := database.DB.Preload("User").Where("id = ? AND user_id = ?", paymentID, user.ID).First(&payment).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Payment not found")
		return
	}

	// Parse metadata to get course information
	var metadata map[string]interface{}
	if payment.Metadata != "" {
		json.Unmarshal([]byte(payment.Metadata), &metadata)
	}

	receipt := map[string]interface{}{
		"receipt_id":     fmt.Sprintf("RCP-%d", payment.ID),
		"payment_id":     payment.ID,
		"transaction_id": payment.StripePaymentID,
		"amount":         payment.Amount,
		"currency":       payment.Currency,
		"status":         payment.Status,
		"description":    payment.Description,
		"customer": map[string]interface{}{
			"id":    payment.User.ID,
			"email": payment.User.Email,
			"name":  fmt.Sprintf("%s %s", payment.User.FirstName, payment.User.LastName),
		},
		"payment_date": payment.CreatedAt.Format(time.RFC3339),
		"metadata":     metadata,
	}

	utils.WriteSuccessResponse(w, "Receipt generated successfully", receipt)
}

// 6. GetAllPayments - Get all payments (admin)
// GET /api/v1/admin/payments
// Headers: Authorization: Bearer <token>
// Response: { payments: [], total, page, limit }
//
// Why: Allows administrators to view all payments across the platform. Essential for
// financial oversight, auditing, and customer support.
func (h *CoursePaymentHandler) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	// Check admin access
	if user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	status := r.URL.Query().Get("status")
	userIDStr := r.URL.Query().Get("user_id")

	var payments []models.Payment
	var total int64
	query := database.DB.Model(&models.Payment{}).Preload("User")

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if userIDStr != "" {
		userID, _ := strconv.ParseUint(userIDStr, 10, 32)
		query = query.Where("user_id = ?", userID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch payments")
		return
	}

	// Get paginated results
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&payments).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch payments")
		return
	}

	// Convert to response format
	paymentResponses := make([]models.PaymentResponse, len(payments))
	for i, payment := range payments {
		paymentResponses[i] = models.PaymentResponse{
			ID:              payment.ID,
			StripePaymentID: payment.StripePaymentID,
			Amount:          payment.Amount,
			Currency:        payment.Currency,
			Status:          payment.Status,
			Description:     payment.Description,
			CreatedAt:       payment.CreatedAt,
		}
	}

	response := map[string]interface{}{
		"payments": paymentResponses,
		"total":    total,
		"page":     page,
		"limit":    limit,
	}

	utils.WriteSuccessResponse(w, "Payments retrieved successfully", response)
}

// 7. RefundPayment - Process refund (admin)
// POST /api/v1/admin/payments/{id}/refund
// Headers: Authorization: Bearer <token>
// Request: { amount?, reason }
// Response: { refund }
//
// Why: Allows administrators to process refunds for payments. Supports full or partial
// refunds. In production, this would integrate with Stripe's refund API.
func (h *CoursePaymentHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	// Check admin access
	if user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	paymentID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	var req struct {
		Amount *int64 `json:"amount"` // Optional: nil means full refund
		Reason string `json:"reason"`
	}

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	var payment models.Payment
	if err := database.DB.First(&payment, paymentID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Payment not found")
		return
	}

	// Check if payment is eligible for refund
	if payment.Status != "succeeded" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Only succeeded payments can be refunded")
		return
	}

	// Determine refund amount
	refundAmount := payment.Amount
	if req.Amount != nil {
		if *req.Amount <= 0 || *req.Amount > payment.Amount {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid refund amount")
			return
		}
		refundAmount = *req.Amount
	}

	// Update payment status
	// In production, you would call Stripe's refund API here
	// For now, we'll mark it as refunded in the database
	payment.Status = "refunded"
	if err := database.DB.Save(&payment).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to process refund")
		return
	}

	refund := map[string]interface{}{
		"refund_id":   fmt.Sprintf("REF-%d-%d", payment.ID, time.Now().Unix()),
		"payment_id":  payment.ID,
		"amount":      refundAmount,
		"currency":    payment.Currency,
		"reason":      req.Reason,
		"refunded_at": time.Now().Format(time.RFC3339),
		"refunded_by": user.ID,
		"original_payment": map[string]interface{}{
			"id":     payment.ID,
			"amount": payment.Amount,
		},
	}

	utils.WriteSuccessResponse(w, "Refund processed successfully", refund)
}

// 8. GetPaymentStats - Payment statistics (admin)
// GET /api/v1/admin/payments/stats
// Headers: Authorization: Bearer <token>
// Query params: start_date?, end_date?
// Response: { stats }
//
// Why: Provides comprehensive payment statistics for administrators. Useful for
// financial reporting, business intelligence, and understanding payment trends.
func (h *CoursePaymentHandler) GetPaymentStats(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	// Check admin access
	if user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Parse date filters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
			return
		}
	} else {
		// Default to 30 days ago
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
			return
		}
		// Set to end of day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		endDate = time.Now()
	}

	query := database.DB.Model(&models.Payment{}).Where("created_at BETWEEN ? AND ?", startDate, endDate)

	// Total payments count
	var totalPayments int64
	query.Count(&totalPayments)

	// Total revenue (succeeded payments only)
	var totalRevenue struct {
		Total int64
	}
	database.DB.Model(&models.Payment{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "succeeded").
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&totalRevenue)

	// Payments by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	database.DB.Model(&models.Payment{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts)

	statusMap := make(map[string]int64)
	for _, sc := range statusCounts {
		statusMap[sc.Status] = sc.Count
	}

	// Refunded amount
	var refundedAmount struct {
		Total int64
	}
	database.DB.Model(&models.Payment{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "refunded").
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&refundedAmount)

	// Average payment amount
	var avgAmount struct {
		Avg float64
	}
	database.DB.Model(&models.Payment{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "succeeded").
		Select("COALESCE(AVG(amount), 0) as avg").
		Scan(&avgAmount)

	stats := map[string]interface{}{
		"period": map[string]interface{}{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
		"total_payments":     totalPayments,
		"total_revenue":      totalRevenue.Total,
		"refunded_amount":    refundedAmount.Total,
		"average_payment":    avgAmount.Avg,
		"payments_by_status": statusMap,
	}

	utils.WriteSuccessResponse(w, "Payment statistics retrieved successfully", stats)
}

// RegisterRoutes sets up all course payment management routes
func (h *CoursePaymentHandler) RegisterRoutes(router *mux.Router) {
	paymentRouter := router.PathPrefix("/payments").Subrouter()

	// User payment routes (require authentication)
	paymentRouter.HandleFunc("/process", h.ProcessPayment).Methods("POST")
	paymentRouter.HandleFunc("/history", h.GetPaymentHistory).Methods("GET")
	paymentRouter.HandleFunc("/{id}", h.GetPaymentByID).Methods("GET")
	paymentRouter.HandleFunc("/{id}/status", h.GetPaymentStatus).Methods("GET")
	paymentRouter.HandleFunc("/{id}/receipt", h.GetPaymentReceipt).Methods("GET")

	// Admin routes (require admin middleware - should be applied in routes.go)
	adminRouter := router.PathPrefix("/admin/payments").Subrouter()
	adminRouter.HandleFunc("", h.GetAllPayments).Methods("GET")
	adminRouter.HandleFunc("/{id}/refund", h.RefundPayment).Methods("POST")
	adminRouter.HandleFunc("/stats", h.GetPaymentStats).Methods("GET")
}
