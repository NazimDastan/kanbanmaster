package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"kanbanmaster/cmd/services"
)

type TeamHandler struct {
	teamService *services.TeamService
	authz       *services.AuthzService
}

func NewTeamHandler(teamService *services.TeamService, authz *services.AuthzService) *TeamHandler {
	return &TeamHandler{teamService: teamService, authz: authz}
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var input services.CreateTeamInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.Name == "" || input.OrganizationID == "" {
		writeError(w, "Name and organizationId are required", http.StatusBadRequest)
		return
	}

	// Check user has access to the organization
	ok, _ := h.authz.UserCanAccessOrg(userID, input.OrganizationID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	team, err := h.teamService.Create(userID, input)
	if err != nil {
		writeError(w, "Failed to create team", http.StatusInternalServerError)
		return
	}
	writeJSON(w, team, http.StatusCreated)
}

func (h *TeamHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	teams, err := h.teamService.List(userID)
	if err != nil {
		writeError(w, "Failed to list teams", http.StatusInternalServerError)
		return
	}
	if teams == nil {
		writeJSON(w, []interface{}{}, http.StatusOK)
		return
	}
	writeJSON(w, teams, http.StatusOK)
}

func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTeam(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	team, err := h.teamService.GetWithMembers(id)
	if err != nil {
		if errors.Is(err, services.ErrTeamNotFound) {
			writeError(w, "Team not found", http.StatusNotFound)
			return
		}
		writeError(w, "Failed to get team", http.StatusInternalServerError)
		return
	}
	writeJSON(w, team, http.StatusOK)
}

func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTeam(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeError(w, "Name is required", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.Update(id, userID, body.Name)
	if err != nil {
		if errors.Is(err, services.ErrNotTeamLeader) {
			writeError(w, "Insufficient permissions", http.StatusForbidden)
			return
		}
		writeError(w, "Failed to update team", http.StatusInternalServerError)
		return
	}
	writeJSON(w, team, http.StatusOK)
}

func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTeam(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	err := h.teamService.Delete(id, userID)
	if err != nil {
		if errors.Is(err, services.ErrNotTeamLeader) {
			writeError(w, "Insufficient permissions", http.StatusForbidden)
			return
		}
		writeError(w, "Failed to delete team", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) Invite(w http.ResponseWriter, r *http.Request) {
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

	member, err := h.teamService.InviteMember(teamID, userID, body.Email, body.Role)
	if err != nil {
		if errors.Is(err, services.ErrNotTeamLeader) {
			writeError(w, "Insufficient permissions", http.StatusForbidden)
			return
		}
		if errors.Is(err, services.ErrUserNotFound) {
			writeError(w, "User not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrAlreadyMember) {
			writeError(w, "User is already a member", http.StatusConflict)
			return
		}
		writeError(w, "Failed to invite member", http.StatusInternalServerError)
		return
	}
	writeJSON(w, member, http.StatusCreated)
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	teamID := r.PathValue("id")
	memberUserID := r.PathValue("userId")

	ok, _ := h.authz.UserCanAccessTeam(userID, teamID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	err := h.teamService.RemoveMember(teamID, userID, memberUserID)
	if err != nil {
		if errors.Is(err, services.ErrNotTeamLeader) {
			writeError(w, "Insufficient permissions", http.StatusForbidden)
			return
		}
		if errors.Is(err, services.ErrMemberNotFound) {
			writeError(w, "Member not found", http.StatusNotFound)
			return
		}
		writeError(w, "Failed to remove member", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	teamID := r.PathValue("id")
	memberUserID := r.PathValue("userId")

	ok, _ := h.authz.UserCanAccessTeam(userID, teamID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Role == "" {
		writeError(w, "Role is required", http.StatusBadRequest)
		return
	}

	err := h.teamService.UpdateMemberRole(teamID, userID, memberUserID, body.Role)
	if err != nil {
		if errors.Is(err, services.ErrNotTeamLeader) {
			writeError(w, "Insufficient permissions", http.StatusForbidden)
			return
		}
		if errors.Is(err, services.ErrMemberNotFound) {
			writeError(w, "Member not found", http.StatusNotFound)
			return
		}
		writeError(w, "Failed to update role", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
