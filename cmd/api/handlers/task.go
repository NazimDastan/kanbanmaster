package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"kanbanmaster/cmd/api/middleware"
	"kanbanmaster/cmd/services"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
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

	task, err := h.taskService.Create(userID, input)
	if err != nil {
		writeError(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	writeJSON(w, task, http.StatusCreated)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

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
	id := r.PathValue("id")

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
	id := r.PathValue("id")

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
	id := r.PathValue("id")

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
	id := r.PathValue("id")

	var body struct {
		AssigneeID string `json:"assigneeId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.AssigneeID == "" {
		writeError(w, "AssigneeId is required", http.StatusBadRequest)
		return
	}

	if err := h.taskService.Assign(id, body.AssigneeID); err != nil {
		writeError(w, "Failed to assign task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) CreateSubtask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")

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
	id := r.PathValue("id")

	if err := h.taskService.ToggleSubtask(id); err != nil {
		writeError(w, "Failed to toggle subtask", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) DeleteSubtask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.taskService.DeleteSubtask(id); err != nil {
		writeError(w, "Failed to delete subtask", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
