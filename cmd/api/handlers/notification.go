package handlers

import (
	"net/http"

	"kanbanmaster/cmd/services"
)

type NotificationHandler struct {
	notifService *services.NotificationService
}

func NewNotificationHandler(ns *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: ns}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	notifications, err := h.notifService.ListByUser(userID)
	if err != nil {
		writeError(w, "Failed to list notifications", http.StatusInternalServerError)
		return
	}
	writeJSON(w, notifications, http.StatusOK)
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	if err := h.notifService.MarkRead(id, userID); err != nil {
		writeError(w, "Failed to mark as read", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	if err := h.notifService.MarkAllRead(userID); err != nil {
		writeError(w, "Failed to mark all as read", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
