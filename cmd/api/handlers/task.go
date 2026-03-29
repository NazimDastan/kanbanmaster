package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"kanbanmaster/cmd/api/middleware"
	"kanbanmaster/cmd/services"
)

type TaskHandler struct {
	taskService  *services.TaskService
	notifService *services.NotificationService
	authService  *services.AuthService
	authz        *services.AuthzService
}

func NewTaskHandler(taskService *services.TaskService, notifService *services.NotificationService, authService *services.AuthService, authz *services.AuthzService) *TaskHandler {
	return &TaskHandler{taskService: taskService, notifService: notifService, authService: authService, authz: authz}
}

func (h *TaskHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	filter := r.URL.Query().Get("filter") // "assigned", "completed", "overdue", "in_progress", or "" for all

	tasks, err := h.taskService.ListByUser(userID, filter)
	if err != nil {
		writeError(w, "Failed to list tasks", http.StatusInternalServerError)
		return
	}
	writeJSON(w, tasks, http.StatusOK)
}

func (h *TaskHandler) Search(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, []struct{}{}, http.StatusOK)
		return
	}

	tasks, err := h.taskService.Search(userID, q)
	if err != nil {
		writeError(w, "Failed to search tasks", http.StatusInternalServerError)
		return
	}
	writeJSON(w, tasks, http.StatusOK)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var input services.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	input.Title = middleware.SanitizeString(input.Title)
	if input.Description != nil {
		sanitized := middleware.SanitizeString(*input.Description)
		input.Description = &sanitized
	}

	if input.ColumnID == "" || input.Title == "" {
		writeError(w, "ColumnId and title are required", http.StatusBadRequest)
		return
	}

	// Check user has access to the column
	ok, _ := h.authz.UserCanAccessColumn(userID, input.ColumnID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	task, err := h.taskService.Create(userID, input)
	if err != nil {
		writeError(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	writeJSON(w, task, http.StatusCreated)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTask(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	task, err := h.taskService.Get(id)
	if err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			writeError(w, "Task not found", http.StatusNotFound)
			return
		}
		writeError(w, "Failed to get task", http.StatusInternalServerError)
		return
	}
	writeJSON(w, task, http.StatusOK)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTask(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	var input services.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.Update(id, input)
	if err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			writeError(w, "Task not found", http.StatusNotFound)
			return
		}
		writeError(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	writeJSON(w, task, http.StatusOK)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTask(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.taskService.Delete(id); err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			writeError(w, "Task not found", http.StatusNotFound)
			return
		}
		writeError(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) Move(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTask(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	var input services.MoveTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.ColumnID == "" {
		writeError(w, "ColumnId is required", http.StatusBadRequest)
		return
	}

	if err := h.taskService.Move(id, input); err != nil {
		writeError(w, "Failed to move task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) Assign(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTask(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	var body struct {
		AssigneeID string `json:"assigneeId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.AssigneeID == "" {
		writeError(w, "AssigneeId is required", http.StatusBadRequest)
		return
	}

	// Get task title for notification
	task, _ := h.taskService.Get(id)
	taskTitle := "Task"
	if task != nil {
		taskTitle = task.Title
	}

	if err := h.taskService.Assign(id, body.AssigneeID); err != nil {
		writeError(w, "Failed to assign task", http.StatusInternalServerError)
		return
	}

	// Send notification to assignee with assigner info
	if h.notifService != nil && body.AssigneeID != userID {
		assignerName := "Someone"
		if assigner, err := h.authService.GetUser(userID); err == nil {
			assignerName = assigner.Name
		}
		h.notifService.NotifyTaskAssigned(body.AssigneeID, taskTitle+" (by "+assignerName+")", id)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) AddAssignee(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTask(userID, taskID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	var body struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.UserID == "" {
		writeError(w, "userId is required", http.StatusBadRequest)
		return
	}

	if err := h.taskService.AddAssignee(taskID, body.UserID); err != nil {
		writeError(w, "Failed to add assignee", http.StatusInternalServerError)
		return
	}

	// Send notification to new assignee
	if h.notifService != nil && body.UserID != userID {
		task, _ := h.taskService.Get(taskID)
		taskTitle := "Task"
		if task != nil {
			taskTitle = task.Title
		}
		assignerName := "Someone"
		if assigner, err := h.authService.GetUser(userID); err == nil {
			assignerName = assigner.Name
		}
		h.notifService.NotifyTaskAssigned(body.UserID, taskTitle+" (by "+assignerName+")", taskID)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) RemoveAssignee(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTask(userID, taskID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	targetUserID := r.PathValue("userId")
	if targetUserID == "" {
		writeError(w, "userId is required", http.StatusBadRequest)
		return
	}

	if err := h.taskService.RemoveAssignee(taskID, targetUserID); err != nil {
		writeError(w, "Failed to remove assignee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) CreateSubtask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("taskId")

	ok, _ := h.authz.UserCanAccessTask(userID, taskID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	var body struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
		writeError(w, "Title is required", http.StatusBadRequest)
		return
	}

	sub, err := h.taskService.CreateSubtask(taskID, body.Title)
	if err != nil {
		writeError(w, "Failed to create subtask", http.StatusInternalServerError)
		return
	}
	writeJSON(w, sub, http.StatusCreated)
}

func (h *TaskHandler) ToggleSubtask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessSubtask(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.taskService.ToggleSubtask(id); err != nil {
		writeError(w, "Failed to toggle subtask", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) DeleteSubtask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessSubtask(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.taskService.DeleteSubtask(id); err != nil {
		writeError(w, "Failed to delete subtask", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
