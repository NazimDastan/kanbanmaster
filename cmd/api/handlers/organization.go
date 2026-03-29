package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"kanbanmaster/cmd/services"
)

type OrgHandler struct {
	orgService *services.OrgService
	authz      *services.AuthzService
}

func NewOrgHandler(orgService *services.OrgService, authz *services.AuthzService) *OrgHandler {
	return &OrgHandler{orgService: orgService, authz: authz}
}

func (h *OrgHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var input services.CreateOrgInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.Name == "" {
		writeError(w, "Name is required", http.StatusBadRequest)
		return
	}

	org, err := h.orgService.Create(userID, input)
	if err != nil {
		writeError(w, "Failed to create organization", http.StatusInternalServerError)
		return
	}
	writeJSON(w, org, http.StatusCreated)
}

func (h *OrgHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	orgs, err := h.orgService.List(userID)
	if err != nil {
		writeError(w, "Failed to list organizations", http.StatusInternalServerError)
		return
	}
	if orgs == nil {
		writeJSON(w, []interface{}{}, http.StatusOK)
		return
	}
	writeJSON(w, orgs, http.StatusOK)
}

func (h *OrgHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessOrg(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	org, err := h.orgService.Get(id)
	if err != nil {
		if errors.Is(err, services.ErrOrgNotFound) {
			writeError(w, "Organization not found", http.StatusNotFound)
			return
		}
		writeError(w, "Failed to get organization", http.StatusInternalServerError)
		return
	}
	writeJSON(w, org, http.StatusOK)
}

func (h *OrgHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessOrg(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	var input services.CreateOrgInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	org, err := h.orgService.Update(id, userID, input)
	if err != nil {
		if errors.Is(err, services.ErrNotOrgOwner) {
			writeError(w, "Not authorized", http.StatusForbidden)
			return
		}
		writeError(w, "Failed to update organization", http.StatusInternalServerError)
		return
	}
	writeJSON(w, org, http.StatusOK)
}

func (h *OrgHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessOrg(userID, id)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	err := h.orgService.Delete(id, userID)
	if err != nil {
		if errors.Is(err, services.ErrNotOrgOwner) {
			writeError(w, "Not authorized", http.StatusForbidden)
			return
		}
		writeError(w, "Failed to delete organization", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
