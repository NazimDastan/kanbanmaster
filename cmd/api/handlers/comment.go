package handlers

import (
	"encoding/json"
	"net/http"

	"kanbanmaster/cmd/api/middleware"
	"kanbanmaster/cmd/services"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(cs *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: cs}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("taskId")

	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, "Content is required", http.StatusBadRequest)
		return
	}
	body.Content = middleware.SanitizeString(body.Content)
	if body.Content == "" {
		writeError(w, "Content is required", http.StatusBadRequest)
		return
	}

	comment, err := h.commentService.Create(taskID, userID, body.Content)
	if err != nil {
		writeError(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}
	writeJSON(w, comment, http.StatusCreated)
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")

	comments, err := h.commentService.ListByTask(taskID)
	if err != nil {
		writeError(w, "Failed to list comments", http.StatusInternalServerError)
		return
	}
	writeJSON(w, comments, http.StatusOK)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	if err := h.commentService.Delete(id, userID); err != nil {
		writeError(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
