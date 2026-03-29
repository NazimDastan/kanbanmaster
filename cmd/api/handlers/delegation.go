package handlers

import (
	"encoding/json"
	"net/http"

	"kanbanmaster/cmd/services"
)

type DelegationHandler struct {
	delegationService *services.DelegationService
	authz             *services.AuthzService
}

func NewDelegationHandler(ds *services.DelegationService, authz *services.AuthzService) *DelegationHandler {
	return &DelegationHandler{delegationService: ds, authz: authz}
}

func (h *DelegationHandler) Delegate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTask(userID, taskID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	var input services.DelegateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.ToUserID == "" {
		writeError(w, "toUserId is required", http.StatusBadRequest)
		return
	}

	d, err := h.delegationService.Delegate(taskID, userID, input)
	if err != nil {
		writeError(w, "Failed to delegate task", http.StatusInternalServerError)
		return
	}
	writeJSON(w, d, http.StatusCreated)
}

func (h *DelegationHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("id")

	ok, _ := h.authz.UserCanAccessTask(userID, taskID)
	if !ok {
		writeError(w, "Not found", http.StatusNotFound)
		return
	}

	logs, err := h.delegationService.GetActivityLog(taskID)
	if err != nil {
		writeError(w, "Failed to get activity log", http.StatusInternalServerError)
		return
	}
	writeJSON(w, logs, http.StatusOK)
}
