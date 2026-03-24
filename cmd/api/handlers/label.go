package handlers

import (
	"encoding/json"
	"net/http"

	"kanbanmaster/cmd/services"
)

type LabelHandler struct {
	labelService *services.LabelService
}

func NewLabelHandler(ls *services.LabelService) *LabelHandler {
	return &LabelHandler{labelService: ls}
}

func (h *LabelHandler) Create(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("boardId")
	var body struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeError(w, "Name is required", http.StatusBadRequest)
		return
	}
	label, err := h.labelService.Create(boardID, body.Name, body.Color)
	if err != nil {
		writeError(w, "Failed to create label", http.StatusInternalServerError)
		return
	}
	writeJSON(w, label, http.StatusCreated)
}

func (h *LabelHandler) List(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("boardId")
	labels, err := h.labelService.ListByBoard(boardID)
	if err != nil {
		writeError(w, "Failed to list labels", http.StatusInternalServerError)
		return
	}
	writeJSON(w, labels, http.StatusOK)
}

func (h *LabelHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeError(w, "Name is required", http.StatusBadRequest)
		return
	}
	label, err := h.labelService.Update(id, body.Name, body.Color)
	if err != nil {
		writeError(w, "Failed to update label", http.StatusInternalServerError)
		return
	}
	writeJSON(w, label, http.StatusOK)
}

func (h *LabelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.labelService.Delete(id); err != nil {
		writeError(w, "Failed to delete label", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *LabelHandler) AddToTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")
	var body struct {
		LabelID string `json:"labelId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.LabelID == "" {
		writeError(w, "labelId is required", http.StatusBadRequest)
		return
	}
	if err := h.labelService.AddToTask(taskID, body.LabelID); err != nil {
		writeError(w, "Failed to add label", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *LabelHandler) RemoveFromTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")
	labelID := r.PathValue("labelId")
	if err := h.labelService.RemoveFromTask(taskID, labelID); err != nil {
		writeError(w, "Failed to remove label", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
