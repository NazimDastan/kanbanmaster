package handlers

import (
	"encoding/json"
	"net/http"

	"kanbanmaster/cmd/services"
)

type ColumnHandler struct {
	columnService *services.ColumnService
}

func NewColumnHandler(columnService *services.ColumnService) *ColumnHandler {
	return &ColumnHandler{columnService: columnService}
}

func (h *ColumnHandler) Create(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("boardId")

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeError(w, "Name is required", http.StatusBadRequest)
		return
	}

	col, err := h.columnService.Create(boardID, body.Name)
	if err != nil {
		writeError(w, "Failed to create column", http.StatusInternalServerError)
		return
	}
	writeJSON(w, col, http.StatusCreated)
}

func (h *ColumnHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeError(w, "Name is required", http.StatusBadRequest)
		return
	}

	col, err := h.columnService.Update(id, body.Name)
	if err != nil {
		writeError(w, "Failed to update column", http.StatusInternalServerError)
		return
	}
	writeJSON(w, col, http.StatusOK)
}

func (h *ColumnHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.columnService.Delete(id); err != nil {
		writeError(w, "Failed to delete column", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ColumnHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BoardID string                  `json:"boardId"`
		Items   []services.ReorderInput `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.BoardID == "" {
		writeError(w, "BoardId and items are required", http.StatusBadRequest)
		return
	}

	if err := h.columnService.Reorder(body.BoardID, body.Items); err != nil {
		writeError(w, "Failed to reorder columns", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
