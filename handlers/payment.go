package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"backend/config"
	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v78"
	stripeCustomer "github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/stripe/stripe-go/v78/webhook"
)

type PaymentHandler struct {
	config *config.Config
}

func NewPaymentHandler(cfg *config.Config) *PaymentHandler {
	// Set Stripe API key
	stripe.Key = cfg.Stripe.SecretKey
	return &PaymentHandler{config: cfg}
}

// RegisterRoutes registers all payment routes
func (h *PaymentHandler) RegisterRoutes(router *mux.Router) {
	// Payment intent routes
	router.HandleFunc("/payments/intent", h.CreatePaymentIntent).Methods("POST")
	router.HandleFunc("/payments/intent/{id}", h.GetPaymentIntent).Methods("GET")
	router.HandleFunc("/payments/confirm", h.ConfirmPayment).Methods("POST")

	// Payment history
	router.HandleFunc("/payments/history", h.GetPaymentHistory).Methods("GET")
}

// CreatePaymentIntent creates a new payment intent
func (h *PaymentHandler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePaymentIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from JWT token
	userID, err := utils.GetUserIDFromToken(r)
	if err != nil {
		utils.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get or create Stripe customer
	stripeCustomerID, err := h.getOrCreateStripeCustomer(userID)
	if err != nil {
		utils.SendErrorResponse(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	// Create payment intent with Stripe
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(req.Currency),
		Customer: stripe.String(stripeCustomerID),
		Metadata: map[string]string{
			"user_id": strconv.Itoa(int(userID)),
		},
	}

	if req.Description != "" {
		params.Description = stripe.String(req.Description)
	}

	if req.Metadata != "" {
		var metadata map[string]string
		if err := json.Unmarshal([]byte(req.Metadata), &metadata); err == nil {
			for k, v := range metadata {
				params.Metadata[k] = v
			}
		}
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		utils.SendErrorResponse(w, "Failed to create payment intent: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save payment intent to database
	paymentIntent := models.PaymentIntent{
		UserID:         userID,
		StripeIntentID: pi.ID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         string(pi.Status),
		Description:    req.Description,
		ClientSecret:   pi.ClientSecret,
	}

	if err := database.DB.Create(&paymentIntent).Error; err != nil {
		utils.SendErrorResponse(w, "Failed to save payment intent", http.StatusInternalServerError)
		return
	}

	response := models.CreatePaymentIntentResponse{
		ClientSecret: pi.ClientSecret,
		IntentID:     pi.ID,
	}

	utils.SendJSONResponse(w, response, http.StatusCreated)
}

// ConfirmPayment confirms a payment intent
func (h *PaymentHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	var req models.ConfirmPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from JWT token
	userID, err := utils.GetUserIDFromToken(r)
	if err != nil {
		utils.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Find payment intent in database
	var paymentIntent models.PaymentIntent
	if err := database.DB.Where("stripe_intent_id = ? AND user_id = ?", req.PaymentIntentID, userID).First(&paymentIntent).Error; err != nil {
		utils.SendErrorResponse(w, "Payment intent not found", http.StatusNotFound)
		return
	}

	// Confirm payment intent with Stripe
	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(req.PaymentMethodID),
	}

	pi, err := paymentintent.Confirm(req.PaymentIntentID, params)
	if err != nil {
		utils.SendErrorResponse(w, "Failed to confirm payment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update payment intent in database
	paymentIntent.Status = string(pi.Status)
	paymentIntent.PaymentMethodID = req.PaymentMethodID
	database.DB.Save(&paymentIntent)

	// If payment succeeded, create payment record
	if pi.Status == stripe.PaymentIntentStatusSucceeded {
		payment := models.Payment{
			UserID:          userID,
			StripePaymentID: pi.ID,
			Amount:          pi.Amount,
			Currency:        string(pi.Currency),
			Status:          "succeeded",
			Description:     pi.Description,
			PaymentMethodID: req.PaymentMethodID,
			ClientSecret:    pi.ClientSecret,
		}

		if pi.Metadata != nil {
			metadataBytes, _ := json.Marshal(pi.Metadata)
			payment.Metadata = string(metadataBytes)
		}

		database.DB.Create(&payment)
	}

	response := models.PaymentResponse{
		ID:              paymentIntent.ID,
		StripePaymentID: pi.ID,
		Amount:          pi.Amount,
		Currency:        string(pi.Currency),
		Status:          string(pi.Status),
		Description:     pi.Description,
		CreatedAt:       paymentIntent.CreatedAt,
	}

	utils.SendJSONResponse(w, response, http.StatusOK)
}

// GetPaymentHistory returns payment history for the authenticated user
func (h *PaymentHandler) GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT token
	userID, err := utils.GetUserIDFromToken(r)
	if err != nil {
		utils.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payments []models.Payment
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&payments).Error; err != nil {
		utils.SendErrorResponse(w, "Failed to fetch payment history", http.StatusInternalServerError)
		return
	}

	var response []models.PaymentResponse
	for _, payment := range payments {
		response = append(response, models.PaymentResponse{
			ID:              payment.ID,
			StripePaymentID: payment.StripePaymentID,
			Amount:          payment.Amount,
			Currency:        payment.Currency,
			Status:          payment.Status,
			Description:     payment.Description,
			CreatedAt:       payment.CreatedAt,
		})
	}

	utils.SendJSONResponse(w, response, http.StatusOK)
}

// GetPaymentIntent returns a specific payment intent
func (h *PaymentHandler) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	intentID := vars["id"]

	// Get user ID from JWT token
	userID, err := utils.GetUserIDFromToken(r)
	if err != nil {
		utils.SendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var paymentIntent models.PaymentIntent
	if err := database.DB.Where("stripe_intent_id = ? AND user_id = ?", intentID, userID).First(&paymentIntent).Error; err != nil {
		utils.SendErrorResponse(w, "Payment intent not found", http.StatusNotFound)
		return
	}

	utils.SendJSONResponse(w, paymentIntent, http.StatusOK)
}

// StripeWebhook handles Stripe webhook events
func (h *PaymentHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		utils.SendErrorResponse(w, "Failed to read webhook payload", http.StatusBadRequest)
		return
	}

	// Verify webhook signature
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), h.config.Stripe.WebhookSecret)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid webhook signature", http.StatusBadRequest)
		return
	}

	// Check if we've already processed this event
	var existingEvent models.WebhookEvent
	if err := database.DB.Where("event_id = ?", event.ID).First(&existingEvent).Error; err == nil {
		// Event already processed
		w.WriteHeader(http.StatusOK)
		return
	}

	// Save webhook event

	webhookEvent := models.WebhookEvent{
		EventID:   event.ID,
		EventType: string(event.Type), // <- Error here. Change-fix: stripe.EventType to string
		Data:      string(event.Data.Raw),
	}

	if err := database.DB.Create(&webhookEvent).Error; err != nil {
		utils.SendErrorResponse(w, "Failed to save webhook event", http.StatusInternalServerError)
		return
	}

	// Process the event
	switch event.Type {
	case "payment_intent.succeeded":
		h.handlePaymentIntentSucceeded(event)
	case "payment_intent.payment_failed":
		h.handlePaymentIntentFailed(event)
	case "payment_intent.canceled":
		h.handlePaymentIntentCanceled(event)
	}

	// Mark event as processed
	webhookEvent.Processed = true
	database.DB.Save(&webhookEvent)

	w.WriteHeader(http.StatusOK)
}

// Helper function to get or create Stripe customer
func (h *PaymentHandler) getOrCreateStripeCustomer(userID uint) (string, error) {

	// Check if customer already exists in database
	var customer models.Customer
	if err := database.DB.Where("user_id = ?", userID).First(&customer).Error; err == nil {
		return customer.StripeCustomerID, nil
	}

	// Get user details
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}

	// Create Stripe customer
	params := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(user.FirstName + " " + user.LastName),
		Metadata: map[string]string{
			"user_id": strconv.Itoa(int(userID)),
		},
	}

	stripeCustomerObj, err := stripeCustomer.New(params) // <- Error here:
	if err != nil {
		return "", err
	}

	// Save customer to database
	customer = models.Customer{
		UserID:           userID,
		StripeCustomerID: stripeCustomerObj.ID,
		Email:            user.Email,
	}

	if err := database.DB.Create(&customer).Error; err != nil {
		return "", err
	}

	return stripeCustomerObj.ID, nil
}

// Webhook event handlers
func (h *PaymentHandler) handlePaymentIntentSucceeded(event stripe.Event) {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		return
	}

	// Update payment intent status
	database.DB.Model(&models.PaymentIntent{}).Where("stripe_intent_id = ?", pi.ID).Update("status", pi.Status)

	// Create payment record if it doesn't exist
	var existingPayment models.Payment
	if err := database.DB.Where("stripe_payment_id = ?", pi.ID).First(&existingPayment).Error; err != nil {
		// Extract user ID from metadata
		userIDStr, exists := pi.Metadata["user_id"]
		if !exists {
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return
		}

		payment := models.Payment{
			UserID:          uint(userID),
			StripePaymentID: pi.ID,
			Amount:          pi.Amount,
			Currency:        string(pi.Currency),
			Status:          "succeeded",
			Description:     pi.Description,
			PaymentMethodID: pi.PaymentMethod.ID,
			ClientSecret:    pi.ClientSecret,
		}

		if pi.Metadata != nil {
			metadataBytes, _ := json.Marshal(pi.Metadata)
			payment.Metadata = string(metadataBytes)
		}

		database.DB.Create(&payment)
	}
}

func (h *PaymentHandler) handlePaymentIntentFailed(event stripe.Event) {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		return
	}

	// Update payment intent status
	database.DB.Model(&models.PaymentIntent{}).Where("stripe_intent_id = ?", pi.ID).Update("status", pi.Status)
}

func (h *PaymentHandler) handlePaymentIntentCanceled(event stripe.Event) {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		return
	}

	// Update payment intent status
	database.DB.Model(&models.PaymentIntent{}).Where("stripe_intent_id = ?", pi.ID).Update("status", pi.Status)
}
