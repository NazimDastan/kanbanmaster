package handlers

import (
	"encoding/json"
	"net/http"

	"kanbanmaster/cmd/services"
)

type AttachmentHandler struct {
	attachmentService *services.AttachmentService
}

func NewAttachmentHandler(as *services.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{attachmentService: as}
}

func (h *AttachmentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("taskId")

	var body struct {
		Filename    string `json:"filename"`
		ContentType string `json:"contentType"`
		Size        int    `json:"size"`
		Data        string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Filename == "" || body.Data == "" {
		writeError(w, "Filename and data are required", http.StatusBadRequest)
		return
	}

	// Max 5MB
	if len(body.Data) > 5*1024*1024 {
		writeError(w, "File too large (max 5MB)", http.StatusBadRequest)
		return
	}

	a, err := h.attachmentService.Create(taskID, userID, body.Filename, body.ContentType, body.Size, body.Data)
	if err != nil {
		writeError(w, "Failed to upload", http.StatusInternalServerError)
		return
	}
	writeJSON(w, a, http.StatusCreated)
}

func (h *AttachmentHandler) List(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")
	attachments, err := h.attachmentService.ListByTask(taskID)
	if err != nil {
		writeError(w, "Failed to list", http.StatusInternalServerError)
		return
	}
	writeJSON(w, attachments, http.StatusOK)
}

func (h *AttachmentHandler) Download(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := h.attachmentService.Get(id)
	if err != nil {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}
	writeJSON(w, a, http.StatusOK)
}

func (h *AttachmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")
	if err := h.attachmentService.Delete(id, userID); err != nil {
		writeError(w, "Failed to delete", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
