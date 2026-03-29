package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"kanbanmaster/cmd/api/handlers"
	"kanbanmaster/cmd/api/middleware"
	"kanbanmaster/cmd/config"
	"kanbanmaster/cmd/models"
	"kanbanmaster/cmd/services"
	ws "kanbanmaster/cmd/websocket"
)

func SetupRouter(cfg *config.Config, db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	// Services
	authService := services.NewAuthService(db, cfg)
	authzService := services.NewAuthzService(db)
	orgService := services.NewOrgService(db)
	teamService := services.NewTeamService(db)
	boardService := services.NewBoardService(db)
	columnService := services.NewColumnService(db)
	taskService := services.NewTaskService(db)
	delegationService := services.NewDelegationService(db)
	reportService := services.NewReportService(db)
	labelService := services.NewLabelService(db)
	commentService := services.NewCommentService(db)
	notifService := services.NewNotificationService(db)
	perfService := services.NewPerformanceService(db)
	attachmentService := services.NewAttachmentService(db)
	notifService.SetOnNotify(func(userID string, n models.Notification) {
		hub.SendToUser(userID, ws.Message{Type: "notification", Payload: n})
	})
	invitationService := services.NewInvitationService(db, notifService)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	orgHandler := handlers.NewOrgHandler(orgService, authzService)
	teamHandler := handlers.NewTeamHandler(teamService, authzService)
	boardHandler := handlers.NewBoardHandler(boardService, authzService)
	columnHandler := handlers.NewColumnHandler(columnService, authzService)
	taskHandler := handlers.NewTaskHandler(taskService, notifService, authService, authzService)
	delegationHandler := handlers.NewDelegationHandler(delegationService, authzService)
	reportHandler := handlers.NewReportHandler(reportService)
	labelHandler := handlers.NewLabelHandler(labelService, authzService)
	commentHandler := handlers.NewCommentHandler(commentService, notifService, authService, taskService, authzService)
	notifHandler := handlers.NewNotificationHandler(notifService)
	dashHandler := handlers.NewDashboardHandler(perfService, authzService)
	attachHandler := handlers.NewAttachmentHandler(attachmentService, authzService)
	invHandler := handlers.NewInvitationHandler(invitationService, authzService)

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "kanbanmaster-api",
		})
	})

	// Auth routes (public)
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/refresh", authHandler.Refresh)

	// Protected routes
	protected := http.NewServeMux()

	// Auth
	protected.HandleFunc("GET /api/auth/me", authHandler.Me)
	protected.HandleFunc("PUT /api/auth/profile", authHandler.UpdateProfile)
	protected.HandleFunc("PUT /api/auth/avatar", authHandler.UpdateAvatar)
	protected.HandleFunc("PUT /api/auth/password", authHandler.ChangePassword)

	// Organizations
	protected.HandleFunc("POST /api/organizations", orgHandler.Create)
	protected.HandleFunc("GET /api/organizations", orgHandler.List)
	protected.HandleFunc("GET /api/organizations/{id}", orgHandler.Get)
	protected.HandleFunc("PUT /api/organizations/{id}", orgHandler.Update)
	protected.HandleFunc("DELETE /api/organizations/{id}", orgHandler.Delete)

	// Teams
	protected.HandleFunc("POST /api/teams", teamHandler.Create)
	protected.HandleFunc("GET /api/teams", teamHandler.List)
	protected.HandleFunc("GET /api/teams/{id}", teamHandler.Get)
	protected.HandleFunc("PUT /api/teams/{id}", teamHandler.Update)
	protected.HandleFunc("DELETE /api/teams/{id}", teamHandler.Delete)
	protected.HandleFunc("POST /api/teams/{id}/invite", invHandler.Send)
	protected.HandleFunc("GET /api/teams/{id}/invitations", invHandler.GetTeamInvitations)

	// Invitations (user)
	protected.HandleFunc("GET /api/invitations", invHandler.GetPending)
	protected.HandleFunc("POST /api/invitations/{id}/accept", invHandler.Accept)
	protected.HandleFunc("POST /api/invitations/{id}/reject", invHandler.Reject)
	protected.HandleFunc("DELETE /api/teams/{id}/members/{userId}", teamHandler.RemoveMember)
	protected.HandleFunc("PATCH /api/teams/{id}/members/{userId}/role", teamHandler.UpdateMemberRole)

	// Boards
	protected.HandleFunc("POST /api/boards", boardHandler.Create)
	protected.HandleFunc("GET /api/boards", boardHandler.List)
	protected.HandleFunc("GET /api/boards/{id}", boardHandler.Get)
	protected.HandleFunc("PUT /api/boards/{id}", boardHandler.Update)
	protected.HandleFunc("DELETE /api/boards/{id}", boardHandler.Delete)

	// Columns
	protected.HandleFunc("POST /api/boards/{boardId}/columns", columnHandler.Create)
	protected.HandleFunc("PUT /api/columns/{id}", columnHandler.Update)
	protected.HandleFunc("DELETE /api/columns/{id}", columnHandler.Delete)
	protected.HandleFunc("PATCH /api/columns/reorder", columnHandler.Reorder)

	// Tasks
	protected.HandleFunc("GET /api/tasks", taskHandler.ListByUser)
	protected.HandleFunc("POST /api/tasks", taskHandler.Create)
	protected.HandleFunc("GET /api/tasks/search", taskHandler.Search)
	protected.HandleFunc("GET /api/tasks/{id}", taskHandler.Get)
	protected.HandleFunc("PUT /api/tasks/{id}", taskHandler.Update)
	protected.HandleFunc("DELETE /api/tasks/{id}", taskHandler.Delete)
	protected.HandleFunc("PATCH /api/tasks/{id}/move", taskHandler.Move)
	protected.HandleFunc("PATCH /api/tasks/{id}/assign", taskHandler.Assign)
	protected.HandleFunc("POST /api/tasks/{id}/assignees", taskHandler.AddAssignee)
	protected.HandleFunc("DELETE /api/tasks/{id}/assignees/{userId}", taskHandler.RemoveAssignee)

	// Subtasks
	protected.HandleFunc("POST /api/tasks/{taskId}/subtasks", taskHandler.CreateSubtask)
	protected.HandleFunc("PATCH /api/subtasks/{id}/toggle", taskHandler.ToggleSubtask)
	protected.HandleFunc("DELETE /api/subtasks/{id}", taskHandler.DeleteSubtask)

	// Delegation & Activity
	protected.HandleFunc("POST /api/tasks/{id}/delegate", delegationHandler.Delegate)
	protected.HandleFunc("GET /api/tasks/{id}/activity", delegationHandler.GetActivity)

	// Labels
	protected.HandleFunc("POST /api/boards/{boardId}/labels", labelHandler.Create)
	protected.HandleFunc("GET /api/boards/{boardId}/labels", labelHandler.List)
	protected.HandleFunc("PUT /api/labels/{id}", labelHandler.Update)
	protected.HandleFunc("DELETE /api/labels/{id}", labelHandler.Delete)
	protected.HandleFunc("POST /api/tasks/{taskId}/labels", labelHandler.AddToTask)
	protected.HandleFunc("DELETE /api/tasks/{taskId}/labels/{labelId}", labelHandler.RemoveFromTask)

	// Comments
	protected.HandleFunc("POST /api/tasks/{taskId}/comments", commentHandler.Create)
	protected.HandleFunc("GET /api/tasks/{taskId}/comments", commentHandler.List)
	protected.HandleFunc("DELETE /api/comments/{id}", commentHandler.Delete)

	// Attachments
	protected.HandleFunc("POST /api/tasks/{taskId}/attachments", attachHandler.Upload)
	protected.HandleFunc("GET /api/tasks/{taskId}/attachments", attachHandler.List)
	protected.HandleFunc("GET /api/attachments/{id}", attachHandler.Download)
	protected.HandleFunc("DELETE /api/attachments/{id}", attachHandler.Delete)

	// Reports
	protected.HandleFunc("POST /api/reports/request", reportHandler.RequestReport)
	protected.HandleFunc("GET /api/reports/requests", reportHandler.GetIncoming)
	protected.HandleFunc("GET /api/reports/requests/sent", reportHandler.GetSent)
	protected.HandleFunc("POST /api/reports/requests/{id}/respond", reportHandler.Respond)
	protected.HandleFunc("PATCH /api/reports/requests/{id}/review", reportHandler.Review)

	// Dashboard & Performance
	protected.HandleFunc("GET /api/dashboard/summary", dashHandler.Summary)
	protected.HandleFunc("GET /api/dashboard/team/{teamId}/performance", dashHandler.TeamPerformance)
	protected.HandleFunc("GET /api/dashboard/overdue", dashHandler.OverdueTasks)

	// Notifications
	protected.HandleFunc("GET /api/notifications", notifHandler.List)
	protected.HandleFunc("PATCH /api/notifications/{id}/read", notifHandler.MarkRead)
	protected.HandleFunc("PATCH /api/notifications/read-all", notifHandler.MarkAllRead)

	// WebSocket
	protected.HandleFunc("GET /ws/notifications", hub.HandleWS)

	// Mount protected routes with auth middleware
	authMiddleware := middleware.Auth(authService)
	mux.Handle("/api/auth/me", authMiddleware(protected))
	mux.Handle("/api/auth/profile", authMiddleware(protected))
	mux.Handle("/api/auth/avatar", authMiddleware(protected))
	mux.Handle("/api/auth/password", authMiddleware(protected))
	mux.Handle("/api/organizations", authMiddleware(protected))
	mux.Handle("/api/organizations/", authMiddleware(protected))
	mux.Handle("/api/teams", authMiddleware(protected))
	mux.Handle("/api/teams/", authMiddleware(protected))
	mux.Handle("/api/boards", authMiddleware(protected))
	mux.Handle("/api/boards/", authMiddleware(protected))
	mux.Handle("/api/columns/", authMiddleware(protected))
	mux.Handle("/api/tasks", authMiddleware(protected))
	mux.Handle("/api/tasks/", authMiddleware(protected))
	mux.Handle("/api/subtasks/", authMiddleware(protected))
	mux.Handle("/api/labels/", authMiddleware(protected))
	mux.Handle("/api/comments/", authMiddleware(protected))
	mux.Handle("/api/attachments/", authMiddleware(protected))
	mux.Handle("/api/invitations", authMiddleware(protected))
	mux.Handle("/api/invitations/", authMiddleware(protected))
	mux.Handle("/api/reports", authMiddleware(protected))
	mux.Handle("/api/reports/", authMiddleware(protected))
	mux.Handle("/api/dashboard/", authMiddleware(protected))
	mux.Handle("/api/notifications", authMiddleware(protected))
	mux.Handle("/api/notifications/", authMiddleware(protected))
	mux.Handle("/ws/", authMiddleware(protected))

	// Apply global middleware chain (outermost first)
	handler := middleware.Chain(
		mux,
		middleware.Logger,
		middleware.Recovery,
		middleware.SecurityHeaders,
		middleware.CORS,
		middleware.SanitizeInput,
		middleware.RateLimit,
		middleware.AuthRateLimit,
	)

	return handler
}
