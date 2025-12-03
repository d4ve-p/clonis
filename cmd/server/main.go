package main

import (
	"log"
	"net/http"
	"os"

	"github.com/d4ve-p/clonis/internal/auth"
	"github.com/d4ve-p/clonis/internal/backup"
	"github.com/d4ve-p/clonis/internal/database"
	"github.com/d4ve-p/clonis/internal/gdrive"
	"github.com/d4ve-p/clonis/internal/scheduler"
	"github.com/d4ve-p/clonis/internal/ui"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Clonis Backup Manager Starting...")
	
	godotenv.Load()
	
	// DB setup
	dbStore, err := database.GetDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	// Drive service
	driveService := gdrive.Get(dbStore)
	backupEngine := backup.New(dbStore, driveService)
	
	// Cron service
	cronService := scheduler.New(dbStore, backupEngine)
	cronService.Start()
	defer cronService.Stop()
	
	// UI setup
	uiHandler := ui.New(dbStore)
	
	// Drive handler wrapper
	driveHandlers := &ui.DriveHandler{
		UI: uiHandler,
		Service: driveService,
	}
	
	// Backup handler wrapper
	backupHandler := &ui.BackupHandler{
		UI: uiHandler,
		Engine: backupEngine,
	}
	
	// Routes setup
	mux := http.NewServeMux()
	
	// Auth setup
	authManager := auth.New()
	
	// Serve Static Files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// Auth Routes
	loginPageHandler := http.HandlerFunc(uiHandler.ShowLogin)
	
	mux.HandleFunc("GET /login", loginPageHandler)
	mux.HandleFunc("POST /login", authManager.LoginHandler)
	mux.HandleFunc("GET /logout", authManager.LogoutHandler)
	
	// Dashboard
	dashboardHandler := http.HandlerFunc(uiHandler.DashboardHandler)
	mux.Handle("/", authManager.Middleware(dashboardHandler))
	// Dashboard - Browser
	browserHandler := http.HandlerFunc(uiHandler.BrowseHandler)
	mux.Handle("/browse", authManager.Middleware(browserHandler))
	
	// Google Drive Routes
	mux.HandleFunc("/drive/connect", driveHandlers.Connect)
	mux.HandleFunc("/drive/callback", driveHandlers.Callback)
	
	// Backup Routes
	runBackupHandler := http.HandlerFunc(backupHandler.Run)
	mux.Handle("/backup/run", authManager.Middleware(runBackupHandler))
	
	// Path Management
	addPathHandler := http.HandlerFunc(uiHandler.AddPathHandler)
	deletePathHandler := http.HandlerFunc(uiHandler.DeletePathHandler)
	mux.Handle("/add-path", authManager.Middleware(addPathHandler))
	mux.Handle("/delete-path", authManager.Middleware(deletePathHandler))
	
	port := os.Getenv("PORT")
	log.Printf("Server running on :%v\n", port)
	
	err = http.ListenAndServe(":" + port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
