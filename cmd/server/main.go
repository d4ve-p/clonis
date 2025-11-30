package main

import (
	"log"
	"net/http"
	"os"
	
	"github.com/joho/godotenv"
	"github.com/d4ve-p/clonis/internal/auth"
	"github.com/d4ve-p/clonis/internal/database"
	"github.com/d4ve-p/clonis/internal/ui"
)

func main() {
	log.Println("Clonis Backup Manager Starting...")
	
	godotenv.Load()
	
	// DB setup
	dbStore, err := database.GetDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	// UI setup
	uiHandler := ui.New(dbStore)
	
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
	
	port := os.Getenv("PORT")
	log.Printf("Server running on :%v\n", port)
	err = http.ListenAndServe(":" + port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
