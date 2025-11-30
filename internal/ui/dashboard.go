package ui

import (
	"html/template"
	"log"
	"net/http"

	"github.com/d4ve-p/clonis/internal/database"
)

func (h* Handler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	paths, err := h.Store.GetPaths()
	if err != nil {
		http.Error(w, "Failed to load paths", http.StatusInternalServerError)
		return
	}
	
	logs, err := h.Store.GetRecentLogs(10)
	if err != nil {
		http.Error(w, "Failed to load logs", http.StatusInternalServerError)
		return
	}
	
	settings, err := h.Store.GetSettings()
	if err != nil {
		http.Error(w, "Failed to load settings", http.StatusInternalServerError)
		return
	}
	
	data := struct {
		User bool
		Paths []database.Path
		Logs []database.LogEntry
		Settings map[string]string
	}{
		User: true,
		Paths: paths,
		Logs: logs,
		Settings: settings,
	}
	
	tmpl := template.Must(template.ParseFiles(
		"templates/layout.html",
		"templates/dashboard.html",
		"templates/fragments/path_list.html",
	))
	
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("Template error: %v.", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}