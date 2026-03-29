package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"kanbanmaster/cmd/api/middleware"
	"kanbanmaster/cmd/services"
)

type InvitationHandler struct {
	invService *services.InvitationService
	authz      *services.AuthzService
}

func NewInvitationHandler(is *services.InvitationService, authz *services.AuthzService) *InvitationHandler {
	return &InvitationHandler{invService: is, authz: authz}
}

func (h *InvitationHandler) Send(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	teamID := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTeam(userID, teamID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	var body struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Email == "" {
		writeError(w, "Email is required", http.StatusBadRequest)
		return
	}
	body.Email = middleware.SanitizeString(body.Email)

	inv, err := h.invService.Send(teamID, userID, body.Email, body.Role)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			writeError(w, "User not found — they must register first", http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrAlreadyMember) {
			writeError(w, "Already a team member", http.StatusConflict)
			return
		}
		if errors.Is(err, services.ErrAlreadyInvited) {
			writeError(w, "Already invited (pending)", http.StatusConflict)
			return
		}
		writeError(w, "Failed to send invitation", http.StatusInternalServerError)
		return
	}
	writeJSON(w, inv, http.StatusCreated)
}

func (h *InvitationHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	invitations, err := h.invService.GetPending(userID)
	if err != nil {
		writeError(w, "Failed to get invitations", http.StatusInternalServerError)
		return
	}
	writeJSON(w, invitations, http.StatusOK)
}

func (h *InvitationHandler) Accept(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	if err := h.invService.Accept(id, userID); err != nil {
		if errors.Is(err, services.ErrInvitationNotFound) {
			writeError(w, "Invitation not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrNotInvitee) {
			writeError(w, "Not your invitation", http.StatusForbidden)
			return
		}
		writeError(w, "Failed to accept", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"message": "Invitation accepted"}, http.StatusOK)
}

func (h *InvitationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	if err := h.invService.Reject(id, userID); err != nil {
		writeError(w, "Failed to reject", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"message": "Invitation rejected"}, http.StatusOK)
}

func (h *InvitationHandler) GetTeamInvitations(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	teamID := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTeam(userID, teamID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	invitations, err := h.invService.GetTeamInvitations(teamID)
	if err != nil {
		writeError(w, "Failed to get invitations", http.StatusInternalServerError)
		return
	}
	writeJSON(w, invitations, http.StatusOK)
}
