package handlers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CourseHandler handles all course management operations
// This handler manages courses, course content (lessons), and course-related operations
type CourseHandler struct{}

func NewCourseHandler() *CourseHandler {
	return &CourseHandler{}
}

// 1. CreateCourse - Creates a new course
// POST /api/v1/courses
// Headers: Authorization: Bearer <token>
// Request: { title, description, price, currency, category, level, duration, thumbnail_url }
// Response: { course }
//
// Why: Allows instructors (or admins) to create new educational courses. The instructor
// is automatically set to the authenticated user. Courses start as unpublished.
func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCourseRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Get authenticated user (will be the instructor)
	user := r.Context().Value("user").(models.User)

	// Only instructors and admins can create courses
	if user.Role != "instructor" && user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Instructor or admin access required")
		return
	}

	// Create course
	course := models.Course{
		Title:        req.Title,
		Description:  req.Description,
		InstructorID: user.ID,
		Price:        req.Price,
		Currency:     req.Currency,
		Category:     req.Category,
		Level:        req.Level,
		Duration:     req.Duration,
		ThumbnailURL: req.ThumbnailURL,
		IsPublished:  false, // New courses start as unpublished
	}

	if err := database.DB.Create(&course).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create course")
		return
	}

	// Load instructor relationship for response
	database.DB.Preload("Instructor").First(&course, course.ID)

	response := models.CourseResponse{
		ID:           course.ID,
		Title:        course.Title,
		Description:  course.Description,
		InstructorID: course.InstructorID,
		Instructor: models.UserResponse{
			ID:        course.Instructor.ID,
			Email:     course.Instructor.Email,
			FirstName: course.Instructor.FirstName,
			LastName:  course.Instructor.LastName,
			Role:      course.Instructor.Role,
			IsActive:  course.Instructor.IsActive,
			CreatedAt: course.Instructor.CreatedAt,
			UpdatedAt: course.Instructor.UpdatedAt,
		},
		Price:        course.Price,
		Currency:     course.Currency,
		Category:     course.Category,
		Level:        course.Level,
		Duration:     course.Duration,
		IsPublished:  course.IsPublished,
		ThumbnailURL: course.ThumbnailURL,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Course created successfully", response)
}

// 2. GetCourse - Retrieves a specific course by ID
// GET /api/v1/courses/{id}
// Response: { course }
//
// Why: Allows viewing course details. Published courses are visible to everyone,
// unpublished courses are only visible to the instructor or admins.
func (h *CourseHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	var course models.Course
	if err := database.DB.Preload("Instructor").First(&course, courseID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course not found")
		return
	}

	// Check if course is published or user has permission
	user, ok := r.Context().Value("user").(models.User)
	if !course.IsPublished {
		// Unpublished courses only visible to instructor or admin
		if !ok || (user.ID != course.InstructorID && user.Role != "admin") {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Course not available")
			return
		}
	}

	response := models.CourseResponse{
		ID:           course.ID,
		Title:        course.Title,
		Description:  course.Description,
		InstructorID: course.InstructorID,
		Instructor: models.UserResponse{
			ID:        course.Instructor.ID,
			Email:     course.Instructor.Email,
			FirstName: course.Instructor.FirstName,
			LastName:  course.Instructor.LastName,
			Role:      course.Instructor.Role,
			IsActive:  course.Instructor.IsActive,
			CreatedAt: course.Instructor.CreatedAt,
			UpdatedAt: course.Instructor.UpdatedAt,
		},
		Price:        course.Price,
		Currency:     course.Currency,
		Category:     course.Category,
		Level:        course.Level,
		Duration:     course.Duration,
		IsPublished:  course.IsPublished,
		ThumbnailURL: course.ThumbnailURL,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Course retrieved successfully", response)
}

// 3. ListCourses - Retrieves a paginated list of courses with filters
// GET /api/v1/courses?category=programming&level=beginner&page=1&limit=10
// Response: { courses: [], total, page, limit }
//
// Why: Provides browsing functionality for courses. Supports filtering by category,
// level, instructor, and publication status. Regular users only see published courses.
func (h *CourseHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	category := r.URL.Query().Get("category")
	level := r.URL.Query().Get("level")
	instructorIDStr := r.URL.Query().Get("instructor_id")
	isPublishedStr := r.URL.Query().Get("is_published")

	offset := (page - 1) * limit
	var courses []models.Course
	var total int64
	query := database.DB.Model(&models.Course{}).Preload("Instructor")

	// Regular users only see published courses
	user, ok := r.Context().Value("user").(models.User)
	if !ok || (user.Role != "admin" && user.Role != "instructor") {
		query = query.Where("is_published = ?", true)
	} else {
		// Admins and instructors can filter by publication status
		if isPublishedStr != "" {
			isPublished, _ := strconv.ParseBool(isPublishedStr)
			query = query.Where("is_published = ?", isPublished)
		}
	}

	// Apply filters
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if instructorIDStr != "" {
		instructorID, _ := strconv.ParseUint(instructorIDStr, 10, 32)
		query = query.Where("instructor_id = ?", instructorID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to count courses")
		return
	}

	// Get paginated results
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve courses")
		return
	}

	// Convert to response format
	courseResponses := make([]models.CourseResponse, len(courses))
	for i, course := range courses {
		courseResponses[i] = models.CourseResponse{
			ID:           course.ID,
			Title:        course.Title,
			Description:  course.Description,
			InstructorID: course.InstructorID,
			Instructor: models.UserResponse{
				ID:        course.Instructor.ID,
				Email:     course.Instructor.Email,
				FirstName: course.Instructor.FirstName,
				LastName:  course.Instructor.LastName,
				Role:      course.Instructor.Role,
				IsActive:  course.Instructor.IsActive,
				CreatedAt: course.Instructor.CreatedAt,
				UpdatedAt: course.Instructor.UpdatedAt,
			},
			Price:        course.Price,
			Currency:     course.Currency,
			Category:     course.Category,
			Level:        course.Level,
			Duration:     course.Duration,
			IsPublished:  course.IsPublished,
			ThumbnailURL: course.ThumbnailURL,
			CreatedAt:    course.CreatedAt,
			UpdatedAt:    course.UpdatedAt,
		}
	}

	response := map[string]interface{}{
		"courses": courseResponses,
		"total":   total,
		"page":    page,
		"limit":   limit,
	}

	utils.WriteSuccessResponse(w, "Courses retrieved successfully", response)
}

// 4. UpdateCourse - Updates course details
// PUT /api/v1/courses/{id}
// Headers: Authorization: Bearer <token>
// Request: { title?, description?, price?, currency?, category?, level?, duration?, thumbnail_url?, is_published? }
// Response: { course }
//
// Why: Allows instructors to update their course information. Only the course instructor
// or admins can update courses. This is essential for keeping course content current.
func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course not found")
		return
	}

	// Check authorization - only instructor or admin can update
	user := r.Context().Value("user").(models.User)
	if user.ID != course.InstructorID && user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	// Parse update request
	var req models.UpdateCourseRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Update fields if provided
	if req.Title != nil {
		course.Title = *req.Title
	}
	if req.Description != nil {
		course.Description = *req.Description
	}
	if req.Price != nil {
		course.Price = *req.Price
	}
	if req.Currency != nil {
		course.Currency = *req.Currency
	}
	if req.Category != nil {
		course.Category = *req.Category
	}
	if req.Level != nil {
		course.Level = *req.Level
	}
	if req.Duration != nil {
		course.Duration = *req.Duration
	}
	if req.ThumbnailURL != nil {
		course.ThumbnailURL = *req.ThumbnailURL
	}
	if req.IsPublished != nil {
		course.IsPublished = *req.IsPublished
	}

	if err := database.DB.Save(&course).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update course")
		return
	}

	// Reload with instructor relationship
	database.DB.Preload("Instructor").First(&course, course.ID)

	response := models.CourseResponse{
		ID:           course.ID,
		Title:        course.Title,
		Description:  course.Description,
		InstructorID: course.InstructorID,
		Instructor: models.UserResponse{
			ID:        course.Instructor.ID,
			Email:     course.Instructor.Email,
			FirstName: course.Instructor.FirstName,
			LastName:  course.Instructor.LastName,
			Role:      course.Instructor.Role,
			IsActive:  course.Instructor.IsActive,
			CreatedAt: course.Instructor.CreatedAt,
			UpdatedAt: course.Instructor.UpdatedAt,
		},
		Price:        course.Price,
		Currency:     course.Currency,
		Category:     course.Category,
		Level:        course.Level,
		Duration:     course.Duration,
		IsPublished:  course.IsPublished,
		ThumbnailURL: course.ThumbnailURL,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Course updated successfully", response)
}

// 5. DeleteCourse - Soft deletes a course
// DELETE /api/v1/courses/{id}
// Headers: Authorization: Bearer <token>
// Response: { message }
//
// Why: Allows removing courses. Uses soft delete to preserve data. Only instructors
// of the course or admins can delete. This prevents accidental data loss.
func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course not found")
		return
	}

	// Check authorization
	user := r.Context().Value("user").(models.User)
	if user.ID != course.InstructorID && user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&course).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete course")
		return
	}

	utils.WriteSuccessResponse(w, "Course deleted successfully", nil)
}

// 6. PublishCourse - Publishes or unpublishes a course
// POST /api/v1/courses/{id}/publish
// Headers: Authorization: Bearer <token>
// Request: { is_published: true/false }
// Response: { course }
//
// Why: Provides a simple way to toggle course visibility. Instructors can publish
// courses when ready and unpublish to make changes. Only published courses are
// visible to regular users.
type PublishCourseRequest struct {
	IsPublished bool `json:"is_published"`
}

func (h *CourseHandler) PublishCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course not found")
		return
	}

	// Check authorization
	user := r.Context().Value("user").(models.User)
	if user.ID != course.InstructorID && user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	var req PublishCourseRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	course.IsPublished = req.IsPublished
	if err := database.DB.Save(&course).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update course")
		return
	}

	// Reload with instructor
	database.DB.Preload("Instructor").First(&course, course.ID)

	response := models.CourseResponse{
		ID:           course.ID,
		Title:        course.Title,
		Description:  course.Description,
		InstructorID: course.InstructorID,
		Instructor: models.UserResponse{
			ID:        course.Instructor.ID,
			Email:     course.Instructor.Email,
			FirstName: course.Instructor.FirstName,
			LastName:  course.Instructor.LastName,
			Role:      course.Instructor.Role,
			IsActive:  course.Instructor.IsActive,
			CreatedAt: course.Instructor.CreatedAt,
			UpdatedAt: course.Instructor.UpdatedAt,
		},
		Price:        course.Price,
		Currency:     course.Currency,
		Category:     course.Category,
		Level:        course.Level,
		Duration:     course.Duration,
		IsPublished:  course.IsPublished,
		ThumbnailURL: course.ThumbnailURL,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Course publication status updated", response)
}

// 7. GetCourseContent - Retrieves all content (lessons) for a course
// GET /api/v1/courses/{id}/content
// Response: { contents: [] }
//
// Why: Returns all lessons/modules within a course. Enrolled users can see all content,
// others can only see published content. This is essential for course navigation.
func (h *CourseHandler) GetCourseContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	// Check if course exists
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course not found")
		return
	}

	// Check if user can access content
	user, ok := r.Context().Value("user").(models.User)
	canAccessAll := ok && (user.ID == course.InstructorID || user.Role == "admin")

	// Check if user is enrolled
	var enrollment models.Enrollment
	isEnrolled := ok && database.DB.Where("user_id = ? AND course_id = ?", user.ID, courseID).First(&enrollment).Error == nil

	var contents []models.CourseContent
	query := database.DB.Where("course_id = ?", courseID)

	// If not instructor/admin and not enrolled, only show published content
	if !canAccessAll && !isEnrolled {
		query = query.Where("is_published = ?", true)
	}

	if err := query.Order("\"order\" ASC").Find(&contents).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve course content")
		return
	}

	// Convert to response format
	contentResponses := make([]models.CourseContentResponse, len(contents))
	for i, content := range contents {
		contentResponses[i] = models.CourseContentResponse{
			ID:          content.ID,
			CourseID:    content.CourseID,
			Title:       content.Title,
			Description: content.Description,
			Content:     content.Content,
			ContentType: content.ContentType,
			Order:       content.Order,
			Duration:    content.Duration,
			IsPublished: content.IsPublished,
			VideoURL:    content.VideoURL,
			CreatedAt:   content.CreatedAt,
			UpdatedAt:   content.UpdatedAt,
		}
	}

	response := map[string]interface{}{
		"contents":  contentResponses,
		"course_id": courseID,
	}

	utils.WriteSuccessResponse(w, "Course content retrieved successfully", response)
}

// 8. AddCourseContent - Adds new content (lesson) to a course
// POST /api/v1/courses/{id}/content
// Headers: Authorization: Bearer <token>
// Request: { title, description, content, content_type, order, duration, video_url }
// Response: { content }
//
// Why: Allows instructors to add lessons/modules to their courses. Content is ordered
// and can be different types (text, video, quiz, assignment).
func (h *CourseHandler) AddCourseContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	// Check if course exists and user has permission
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course not found")
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.ID != course.InstructorID && user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.CreateCourseContentRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	content := models.CourseContent{
		CourseID:    uint(courseID),
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		ContentType: req.ContentType,
		Order:       req.Order,
		Duration:    req.Duration,
		VideoURL:    req.VideoURL,
		IsPublished: false, // New content starts as unpublished
	}

	if err := database.DB.Create(&content).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create course content")
		return
	}

	response := models.CourseContentResponse{
		ID:          content.ID,
		CourseID:    content.CourseID,
		Title:       content.Title,
		Description: content.Description,
		Content:     content.Content,
		ContentType: content.ContentType,
		Order:       content.Order,
		Duration:    content.Duration,
		IsPublished: content.IsPublished,
		VideoURL:    content.VideoURL,
		CreatedAt:   content.CreatedAt,
		UpdatedAt:   content.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Course content added successfully", response)
}

// 9. UpdateCourseContent - Updates existing course content
// PUT /api/v1/courses/{courseId}/content/{contentId}
// Headers: Authorization: Bearer <token>
// Request: { title?, description?, content?, content_type?, order?, duration?, video_url?, is_published? }
// Response: { content }
//
// Why: Allows instructors to edit lessons. Essential for fixing errors, updating
// information, or improving course materials over time.
func (h *CourseHandler) UpdateCourseContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contentID, err := strconv.ParseUint(vars["contentId"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid content ID")
		return
	}

	var content models.CourseContent
	if err := database.DB.First(&content, contentID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course content not found")
		return
	}

	// Check authorization via course
	var course models.Course
	if err := database.DB.First(&course, content.CourseID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course not found")
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.ID != course.InstructorID && user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.UpdateCourseContentRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Update fields if provided
	if req.Title != nil {
		content.Title = *req.Title
	}
	if req.Description != nil {
		content.Description = *req.Description
	}
	if req.Content != nil {
		content.Content = *req.Content
	}
	if req.ContentType != nil {
		content.ContentType = *req.ContentType
	}
	if req.Order != nil {
		content.Order = *req.Order
	}
	if req.Duration != nil {
		content.Duration = *req.Duration
	}
	if req.VideoURL != nil {
		content.VideoURL = *req.VideoURL
	}
	if req.IsPublished != nil {
		content.IsPublished = *req.IsPublished
	}

	if err := database.DB.Save(&content).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update course content")
		return
	}

	response := models.CourseContentResponse{
		ID:          content.ID,
		CourseID:    content.CourseID,
		Title:       content.Title,
		Description: content.Description,
		Content:     content.Content,
		ContentType: content.ContentType,
		Order:       content.Order,
		Duration:    content.Duration,
		IsPublished: content.IsPublished,
		VideoURL:    content.VideoURL,
		CreatedAt:   content.CreatedAt,
		UpdatedAt:   content.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Course content updated successfully", response)
}

// 10. DeleteCourseContent - Removes content from a course
// DELETE /api/v1/courses/{courseId}/content/{contentId}
// Headers: Authorization: Bearer <token>
// Response: { message }
//
// Why: Allows instructors to remove lessons that are no longer needed or were
// added by mistake. Uses soft delete to preserve data.
func (h *CourseHandler) DeleteCourseContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contentID, err := strconv.ParseUint(vars["contentId"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid content ID")
		return
	}

	var content models.CourseContent
	if err := database.DB.First(&content, contentID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course content not found")
		return
	}

	// Check authorization via course
	var course models.Course
	if err := database.DB.First(&course, content.CourseID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course not found")
		return
	}

	user := r.Context().Value("user").(models.User)
	if user.ID != course.InstructorID && user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&content).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete course content")
		return
	}

	utils.WriteSuccessResponse(w, "Course content deleted successfully", nil)
}

// 11. EnrollInCourse - Enrolls the authenticated user in a course
// POST /api/v1/courses/{id}/enroll
// Headers: Authorization: Bearer <token>
// Response: { enrollment }
//
// Why: Allows users to enroll in courses. This creates an enrollment record that
// tracks progress. In a full implementation, this would also handle payment if
// the course is paid. This is a quick enrollment action from the course page.
func (h *CourseHandler) EnrollInCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid course ID")
		return
	}

	// Check if course exists and is published
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Course not found")
		return
	}

	if !course.IsPublished {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Course is not published")
		return
	}

	user := r.Context().Value("user").(models.User)

	// Check if already enrolled
	var existingEnrollment models.Enrollment
	if err := database.DB.Where("user_id = ? AND course_id = ?", user.ID, courseID).First(&existingEnrollment).Error; err == nil {
		utils.WriteErrorResponse(w, http.StatusConflict, "Already enrolled in this course")
		return
	}

	// TODO: In production, check if course is paid and handle payment
	// For now, create enrollment directly

	enrollment := models.Enrollment{
		UserID:     user.ID,
		CourseID:   uint(courseID),
		Status:     "active",
		Progress:   0,
		EnrolledAt: database.DB.NowFunc(),
	}

	if err := database.DB.Create(&enrollment).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to enroll in course")
		return
	}

	// Load relationships for response
	database.DB.Preload("User").Preload("Course").Preload("Course.Instructor").First(&enrollment, enrollment.ID)

	response := models.EnrollmentResponse{
		ID:     enrollment.ID,
		UserID: enrollment.UserID,
		User: models.UserResponse{
			ID:        enrollment.User.ID,
			Email:     enrollment.User.Email,
			FirstName: enrollment.User.FirstName,
			LastName:  enrollment.User.LastName,
			Role:      enrollment.User.Role,
			IsActive:  enrollment.User.IsActive,
			CreatedAt: enrollment.User.CreatedAt,
			UpdatedAt: enrollment.User.UpdatedAt,
		},
		CourseID: enrollment.CourseID,
		Course: models.CourseResponse{
			ID:           enrollment.Course.ID,
			Title:        enrollment.Course.Title,
			Description:  enrollment.Course.Description,
			InstructorID: enrollment.Course.InstructorID,
			Instructor: models.UserResponse{
				ID:        enrollment.Course.Instructor.ID,
				Email:     enrollment.Course.Instructor.Email,
				FirstName: enrollment.Course.Instructor.FirstName,
				LastName:  enrollment.Course.Instructor.LastName,
				Role:      enrollment.Course.Instructor.Role,
				IsActive:  enrollment.Course.Instructor.IsActive,
				CreatedAt: enrollment.Course.Instructor.CreatedAt,
				UpdatedAt: enrollment.Course.Instructor.UpdatedAt,
			},
			Price:        enrollment.Course.Price,
			Currency:     enrollment.Course.Currency,
			Category:     enrollment.Course.Category,
			Level:        enrollment.Course.Level,
			Duration:     enrollment.Course.Duration,
			IsPublished:  enrollment.Course.IsPublished,
			ThumbnailURL: enrollment.Course.ThumbnailURL,
			CreatedAt:    enrollment.Course.CreatedAt,
			UpdatedAt:    enrollment.Course.UpdatedAt,
		},
		Status:     enrollment.Status,
		Progress:   enrollment.Progress,
		EnrolledAt: enrollment.EnrolledAt,
		CreatedAt:  enrollment.CreatedAt,
		UpdatedAt:  enrollment.UpdatedAt,
	}

	utils.WriteSuccessResponse(w, "Enrolled in course successfully", response)
}

// RegisterRoutes sets up all course management routes
func (h *CourseHandler) RegisterRoutes(router *mux.Router) {
	courseRouter := router.PathPrefix("/courses").Subrouter()

	// Public routes (no auth required)
	courseRouter.HandleFunc("", h.ListCourses).Methods("GET")
	courseRouter.HandleFunc("/{id}", h.GetCourse).Methods("GET")
	courseRouter.HandleFunc("/{id}/content", h.GetCourseContent).Methods("GET")

	// Protected routes (require authentication)
	// Note: AuthMiddleware should be applied in routes.go
	courseRouter.HandleFunc("", h.CreateCourse).Methods("POST")
	courseRouter.HandleFunc("/{id}", h.UpdateCourse).Methods("PUT")
	courseRouter.HandleFunc("/{id}", h.DeleteCourse).Methods("DELETE")
	courseRouter.HandleFunc("/{id}/publish", h.PublishCourse).Methods("POST")
	courseRouter.HandleFunc("/{id}/content", h.AddCourseContent).Methods("POST")
	courseRouter.HandleFunc("/{id}/content/{contentId}", h.UpdateCourseContent).Methods("PUT")
	courseRouter.HandleFunc("/{id}/content/{contentId}", h.DeleteCourseContent).Methods("DELETE")
	courseRouter.HandleFunc("/{id}/enroll", h.EnrollInCourse).Methods("POST")
}
