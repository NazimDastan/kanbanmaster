package handlers

import (
	"net/http"

	"kanbanmaster/cmd/services"
)

type DashboardHandler struct {
	perfService *services.PerformanceService
}

func NewDashboardHandler(ps *services.PerformanceService) *DashboardHandler {
	return &DashboardHandler{perfService: ps}
}

func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	summary, err := h.perfService.GetDashboardSummary(userID)
	if err != nil {
		writeError(w, "Failed to get summary", http.StatusInternalServerError)
		return
	}
	writeJSON(w, summary, http.StatusOK)
}

func (h *DashboardHandler) TeamPerformance(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")

	performance, err := h.perfService.GetTeamPerformance(teamID)
	if err != nil {
		writeError(w, "Failed to get performance", http.StatusInternalServerError)
		return
	}
	writeJSON(w, performance, http.StatusOK)
}

func (h *DashboardHandler) OverdueTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	tasks, err := h.perfService.GetOverdueTasks(userID)
	if err != nil {
		writeError(w, "Failed to get overdue tasks", http.StatusInternalServerError)
		return
	}
	writeJSON(w, tasks, http.StatusOK)
}
