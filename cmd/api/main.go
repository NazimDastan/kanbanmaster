package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"kanbanmaster/cmd/api/routes"
	"kanbanmaster/cmd/config"
	"kanbanmaster/cmd/services"
)

func main() {
	cfg := config.Load()
	db := config.ConnectDB(cfg.DatabaseURL)
	defer db.Close()

	// Start deadline scheduler
	notifService := services.NewNotificationService(db)
	scheduler := services.NewScheduler(db, notifService)
	scheduler.Start()

	router := routes.SetupRouter(cfg, db)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("KanbanMaster API starting on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Printf("Server failed: %v", err)
		os.Exit(1)
	}
}
