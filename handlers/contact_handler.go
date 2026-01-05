package handlers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// ContactHandler handles all contact management operations
// This handler manages contact submissions, inquiries, and support tickets
type ContactHandler struct{}

func NewContactHandler() *ContactHandler {
	return &ContactHandler{}
}

// 1. CreateContact - Creates a new contact entry (support inquiry, feedback, etc.)
// POST /api/v1/contacts
// Request: { name, email, phone, subject, message }
// Response: { contact }
//
// Why: Allows users (or anonymous visitors) to submit inquiries, support requests, or feedback.
// The user_id is optional - logged-in users can be linked, but anonymous submissions are also allowed.
func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var req models.CreateContactRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Create contact entry
	contact := models.Contact{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Subject: req.Subject,
		Message: req.Message,
		Status:  "pending", // New contacts start as pending
	}

	// If user is authenticated, link the contact to their account
	if user, ok := r.Context().Value("user").(models.User); ok {
		contact.UserID = &user.ID
	}

	if err := database.DB.Create(&contact).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create contact")
		return
	}

	response := models.ContactResponse{
		ID:        contact.ID,
		UserID:    contact.UserID,
		Name:      contact.Name,
		Email:     contact.Email,
		Phone:     contact.Phone,
		Subject:   contact.Subject,
		Message:   contact.Message,
		Status:    contact.Status,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Contact created successfully", response)
}

// 2. GetContact - Retrieves a specific contact by ID
// GET /api/v1/contacts/{id}
// Headers: Authorization: Bearer <token>
// Response: { contact }
//
// Why: Allows viewing details of a specific contact entry. Typically used by admins
// or support staff to view and respond to inquiries. Regular users can only view their own contacts.
func (h *ContactHandler) GetContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contactID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	var contact models.Contact
	if err := database.DB.First(&contact, contactID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Contact not found")
		return
	}

	// Check authorization: users can only view their own contacts, admins can view all
	user := r.Context().Value("user").(models.User)
	if user.Role != "admin" && (contact.UserID == nil || *contact.UserID != user.ID) {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	response := models.ContactResponse{
		ID:        contact.ID,
		UserID:    contact.UserID,
		Name:      contact.Name,
		Email:     contact.Email,
		Phone:     contact.Phone,
		Subject:   contact.Subject,
		Message:   contact.Message,
		Status:    contact.Status,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Contact retrieved successfully", response)
}

// 3. ListContacts - Retrieves a paginated list of contacts
// GET /api/v1/contacts?page=1&limit=10&status=pending
// Headers: Authorization: Bearer <token>
// Response: { contacts: [], total, page, limit }
//
// Why: Provides a way to browse and manage multiple contacts. Regular users see only
// their own contacts, while admins see all contacts. Supports filtering and pagination.
func (h *ContactHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	status := r.URL.Query().Get("status")

	user := r.Context().Value("user").(models.User)
	offset := (page - 1) * limit

	var contacts []models.Contact
	var total int64
	query := database.DB.Model(&models.Contact{})

	// Regular users only see their own contacts, admins see all
	if user.Role != "admin" {
		query = query.Where("user_id = ?", user.ID)
	}

	// Filter by status if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count for pagination
	if err := query.Count(&total).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to count contacts")
		return
	}

	// Get paginated results
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&contacts).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve contacts")
		return
	}

	// Convert to response format
	contactResponses := make([]models.ContactResponse, len(contacts))
	for i, contact := range contacts {
		contactResponses[i] = models.ContactResponse{
			ID:        contact.ID,
			UserID:    contact.UserID,
			Name:      contact.Name,
			Email:     contact.Email,
			Phone:     contact.Phone,
			Subject:   contact.Subject,
			Message:   contact.Message,
			Status:    contact.Status,
			CreatedAt: contact.CreatedAt,
			UpdatedAt: contact.UpdatedAt,
		}
	}

	response := map[string]interface{}{
		"contacts": contactResponses,
		"total":    total,
		"page":     page,
		"limit":    limit,
	}

	utils.WriteSuccessResponse(w, "Contacts retrieved successfully", response)
}

// 4. UpdateContact - Updates a contact entry (typically status or response)
// PUT /api/v1/contacts/{id}
// Headers: Authorization: Bearer <token>
// Request: { name?, email?, phone?, subject?, message?, status? }
// Response: { contact }
//
// Why: Allows updating contact information, typically used by admins to mark contacts
// as "responded" or "resolved" after handling them. Users can update their own pending contacts.
func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contactID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	var contact models.Contact
	if err := database.DB.First(&contact, contactID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Contact not found")
		return
	}

	// Check authorization
	user := r.Context().Value("user").(models.User)
	if user.Role != "admin" && (contact.UserID == nil || *contact.UserID != user.ID) {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	// Parse update request
	var req models.UpdateContactRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Update fields if provided
	if req.Name != nil {
		contact.Name = *req.Name
	}
	if req.Email != nil {
		contact.Email = *req.Email
	}
	if req.Phone != nil {
		contact.Phone = *req.Phone
	}
	if req.Subject != nil {
		contact.Subject = *req.Subject
	}
	if req.Message != nil {
		contact.Message = *req.Message
	}
	if req.Status != nil {
		// Validate status value
		validStatuses := []string{"pending", "responded", "resolved", "archived"}
		isValid := false
		for _, s := range validStatuses {
			if *req.Status == s {
				isValid = true
				break
			}
		}
		if isValid {
			contact.Status = *req.Status
		}
	}

	if err := database.DB.Save(&contact).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update contact")
		return
	}

	response := models.ContactResponse{
		ID:        contact.ID,
		UserID:    contact.UserID,
		Name:      contact.Name,
		Email:     contact.Email,
		Phone:     contact.Phone,
		Subject:   contact.Subject,
		Message:   contact.Message,
		Status:    contact.Status,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Contact updated successfully", response)
}

// 5. DeleteContact - Soft deletes a contact entry
// DELETE /api/v1/contacts/{id}
// Headers: Authorization: Bearer <token>
// Response: { message }
//
// Why: Allows removing contacts from active view. Uses soft delete (GORM DeletedAt)
// so data is preserved but hidden. Typically only admins can delete contacts.
func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contactID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	var contact models.Contact
	if err := database.DB.First(&contact, contactID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Contact not found")
		return
	}

	// Check authorization - typically only admins can delete
	user := r.Context().Value("user").(models.User)
	if user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Soft delete (GORM will set DeletedAt)
	if err := database.DB.Delete(&contact).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete contact")
		return
	}

	utils.WriteSuccessResponse(w, "Contact deleted successfully", nil)
}

// 6. SearchContacts - Searches contacts by name or email
// GET /api/v1/contacts/search?query=john&page=1&limit=10
// Headers: Authorization: Bearer <token>
// Response: { contacts: [], total, page, limit }
//
// Why: Provides search functionality to quickly find specific contacts. Useful for
// support staff looking for previous inquiries from a specific person or email.
func (h *ContactHandler) SearchContacts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Search query is required")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	status := r.URL.Query().Get("status")

	user := r.Context().Value("user").(models.User)
	offset := (page - 1) * limit

	var contacts []models.Contact
	var total int64
	dbQuery := database.DB.Model(&models.Contact{})

	// Regular users only see their own contacts
	if user.Role != "admin" {
		dbQuery = dbQuery.Where("user_id = ?", user.ID)
	}

	// Search in name and email fields
	searchPattern := "%" + query + "%"
	dbQuery = dbQuery.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)

	// Filter by status if provided
	if status != "" {
		dbQuery = dbQuery.Where("status = ?", status)
	}

	// Get total count
	if err := dbQuery.Count(&total).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to count contacts")
		return
	}

	// Get paginated results
	if err := dbQuery.Order("created_at DESC").Offset(offset).Limit(limit).Find(&contacts).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to search contacts")
		return
	}

	// Convert to response format
	contactResponses := make([]models.ContactResponse, len(contacts))
	for i, contact := range contacts {
		contactResponses[i] = models.ContactResponse{
			ID:        contact.ID,
			UserID:    contact.UserID,
			Name:      contact.Name,
			Email:     contact.Email,
			Phone:     contact.Phone,
			Subject:   contact.Subject,
			Message:   contact.Message,
			Status:    contact.Status,
			CreatedAt: contact.CreatedAt,
			UpdatedAt: contact.UpdatedAt,
		}
	}

	response := map[string]interface{}{
		"contacts": contactResponses,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"query":    query,
	}

	utils.WriteSuccessResponse(w, "Contacts found successfully", response)
}

// RegisterRoutes sets up all contact management routes
func (h *ContactHandler) RegisterRoutes(router *mux.Router) {
	contactRouter := router.PathPrefix("/contacts").Subrouter()

	// Public route (no auth required for creating contacts)
	contactRouter.HandleFunc("", h.CreateContact).Methods("POST")

	// Protected routes (require authentication)
	// Note: AuthMiddleware should be applied in routes.go
	contactRouter.HandleFunc("/search", h.SearchContacts).Methods("GET")
	contactRouter.HandleFunc("", h.ListContacts).Methods("GET")
	contactRouter.HandleFunc("/{id}", h.GetContact).Methods("GET")
	contactRouter.HandleFunc("/{id}", h.UpdateContact).Methods("PUT")
	contactRouter.HandleFunc("/{id}", h.DeleteContact).Methods("DELETE")
}
