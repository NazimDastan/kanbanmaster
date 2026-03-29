package handlers

import (
	"encoding/json"
	"net/http"

	"kanbanmaster/cmd/api/middleware"
	"kanbanmaster/cmd/services"
)

type CommentHandler struct {
	commentService *services.CommentService
	notifService   *services.NotificationService
	authService    *services.AuthService
	taskService    *services.TaskService
	authz          *services.AuthzService
}

func NewCommentHandler(cs *services.CommentService, ns *services.NotificationService, as *services.AuthService, ts *services.TaskService, authz *services.AuthzService) *CommentHandler {
	return &CommentHandler{commentService: cs, notifService: ns, authService: as, taskService: ts, authz: authz}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("taskId")

	ok, _ := h.authz.UserCanAccessTask(userID, taskID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

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

	// Send notifications to task assignee and creator
	if h.notifService != nil {
		commenterName := "Someone"
		if user, err := h.authService.GetUser(userID); err == nil {
			commenterName = user.Name
		}
		task, _ := h.taskService.Get(taskID)
		if task != nil {
			taskTitle := task.Title
			// Notify assignee (if different from commenter)
			if task.AssigneeID != nil && *task.AssigneeID != userID {
				h.notifService.NotifyComment(*task.AssigneeID, commenterName, taskTitle, taskID)
			}
			// Notify creator (if different from commenter and assignee)
			if task.CreatorID != userID && (task.AssigneeID == nil || task.CreatorID != *task.AssigneeID) {
				h.notifService.NotifyComment(task.CreatorID, commenterName, taskTitle, taskID)
			}
		}
	}

	writeJSON(w, comment, http.StatusCreated)
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("taskId")

	ok, _ := h.authz.UserCanAccessTask(userID, taskID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

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

	ok, _ := h.authz.UserCanAccessComment(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.commentService.Delete(id, userID); err != nil {
		writeError(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
