package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"kanbanmaster/cmd/services"
)

type ReportHandler struct {
	reportService *services.ReportService
}

func NewReportHandler(rs *services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: rs}
}

func (h *ReportHandler) RequestReport(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var input services.CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.TargetUserID == "" || input.TeamID == "" {
		writeError(w, "targetUserId and teamId are required", http.StatusBadRequest)
		return
	}

	rr, err := h.reportService.RequestReport(userID, input)
	if err != nil {
		writeError(w, "Failed to create report request", http.StatusInternalServerError)
		return
	}
	writeJSON(w, rr, http.StatusCreated)
}

func (h *ReportHandler) GetIncoming(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	reports, err := h.reportService.GetIncoming(userID)
	if err != nil {
		writeError(w, "Failed to get incoming requests", http.StatusInternalServerError)
		return
	}
	writeJSON(w, reports, http.StatusOK)
}

func (h *ReportHandler) GetSent(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	reports, err := h.reportService.GetSent(userID)
	if err != nil {
		writeError(w, "Failed to get sent requests", http.StatusInternalServerError)
		return
	}
	writeJSON(w, reports, http.StatusOK)
}

func (h *ReportHandler) Respond(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	reportID := r.PathValue("id")

	var body struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Response == "" {
		writeError(w, "Response is required", http.StatusBadRequest)
		return
	}

	rr, err := h.reportService.Respond(reportID, userID, body.Response)
	if err != nil {
		if errors.Is(err, services.ErrNotAuthorized) {
			writeError(w, "Not authorized", http.StatusForbidden)
			return
		}
		writeError(w, "Failed to respond", http.StatusInternalServerError)
		return
	}
	writeJSON(w, rr, http.StatusOK)
}

func (h *ReportHandler) Review(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	reportID := r.PathValue("id")

	rr, err := h.reportService.Review(reportID, userID)
	if err != nil {
		if errors.Is(err, services.ErrNotAuthorized) {
			writeError(w, "Not authorized", http.StatusForbidden)
			return
		}
		writeError(w, "Failed to review", http.StatusInternalServerError)
		return
	}
	writeJSON(w, rr, http.StatusOK)
}
