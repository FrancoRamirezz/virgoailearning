package handlers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PostHandler struct{}

func NewPostHandler() *PostHandler {
	return &PostHandler{}
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	if err := database.DB.Preload("Author").Find(&posts).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	utils.WriteSuccessResponse(w, "Posts retrieved successfully", posts)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.Post
	if err := database.DB.Preload("Author").First(&post, uint(postID)).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Post not found")
		return
	}

	utils.WriteSuccessResponse(w, "Post retrieved successfully", post)
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	var req models.CreatePostRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	post := models.Post{
		Title:       req.Title,
		Content:     req.Content,
		AuthorID:    user.ID,
		IsPublished: req.IsPublished,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// Fetch the created post with author information
	database.DB.Preload("Author").First(&post, post.ID)

	utils.WriteSuccessResponse(w, "Post created successfully", post)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	user := r.Context().Value("user").(models.User)

	var post models.Post
	if err := database.DB.First(&post, uint(postID)).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Post not found")
		return
	}

	// Check if user is the author or admin
	if post.AuthorID != user.ID && user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "You can only update your own posts")
		return
	}

	var req models.UpdatePostRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Update only provided fields
	updateData := make(map[string]interface{})
	if req.Title != nil {
		updateData["title"] = *req.Title
	}
	if req.Content != nil {
		updateData["content"] = *req.Content
	}
	if req.IsPublished != nil {
		updateData["is_published"] = *req.IsPublished
	}

	if err := database.DB.Model(&post).Updates(updateData).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update post")
		return
	}

	// Fetch updated post with author information
	database.DB.Preload("Author").First(&post, uint(postID))

	utils.WriteSuccessResponse(w, "Post updated successfully", post)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	user := r.Context().Value("user").(models.User)

	var post models.Post
	if err := database.DB.First(&post, uint(postID)).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Post not found")
		return
	}

	// Check if user is the author or admin
	if post.AuthorID != user.ID && user.Role != "admin" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "You can only delete your own posts")
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	utils.WriteSuccessResponse(w, "Post deleted successfully", nil)
}

func (h *PostHandler) RegisterRoutes(router *mux.Router) {
	postRouter := router.PathPrefix("/posts").Subrouter()
	postRouter.HandleFunc("", h.GetPosts).Methods("GET")
	postRouter.HandleFunc("/{id}", h.GetPost).Methods("GET")
	postRouter.HandleFunc("", h.CreatePost).Methods("POST")
	postRouter.HandleFunc("/{id}", h.UpdatePost).Methods("PUT")
	postRouter.HandleFunc("/{id}", h.DeletePost).Methods("DELETE")
}


