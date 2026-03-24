package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"kanbanmaster/cmd/services"
)

type BoardHandler struct {
	boardService *services.BoardService
}

func NewBoardHandler(boardService *services.BoardService) *BoardHandler {
	return &BoardHandler{boardService: boardService}
}

func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name   string `json:"name"`
		TeamID string `json:"teamId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" || body.TeamID == "" {
		writeError(w, "Name and teamId are required", http.StatusBadRequest)
		return
	}

	board, err := h.boardService.Create(body.Name, body.TeamID)
	if err != nil {
		writeError(w, "Failed to create board", http.StatusInternalServerError)
		return
	}
	writeJSON(w, board, http.StatusCreated)
}

func (h *BoardHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	boards, err := h.boardService.List(userID)
	if err != nil {
		writeError(w, "Failed to list boards", http.StatusInternalServerError)
		return
	}
	if boards == nil {
		writeJSON(w, []interface{}{}, http.StatusOK)
		return
	}
	writeJSON(w, boards, http.StatusOK)
}

func (h *BoardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	board, err := h.boardService.GetWithColumns(id)
	if err != nil {
		if errors.Is(err, services.ErrBoardNotFound) {
			writeError(w, "Board not found", http.StatusNotFound)
			return
		}
		writeError(w, "Failed to get board", http.StatusInternalServerError)
		return
	}
	writeJSON(w, board, http.StatusOK)
}

func (h *BoardHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeError(w, "Name is required", http.StatusBadRequest)
		return
	}

	board, err := h.boardService.Update(id, body.Name)
	if err != nil {
		writeError(w, "Failed to update board", http.StatusInternalServerError)
		return
	}
	writeJSON(w, board, http.StatusOK)
}

func (h *BoardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.boardService.Delete(id); err != nil {
		writeError(w, "Failed to delete board", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
