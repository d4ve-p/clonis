package ui

import (
	"html/template"
	"log"
	"net/http"

	"github.com/d4ve-p/clonis/internal/gdrive"
	uimodel "github.com/d4ve-p/clonis/internal/ui-model"
)

func (h* Handler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	paths, err := h.Store.GetPaths()
	if err != nil {
		h.RenderError(w, "Failed to load paths", err)
		return
	}
	
	logs, err := h.Store.GetRecentLogs(10)
	if err != nil {
		h.RenderError(w, "Failed to get logs", err)
		return
	}
	
	settings, err := h.Store.GetSettings()
	if err != nil {
		h.RenderError(w, "Failed to load settings", err)
		return
	}
	
	driveService := gdrive.Get(h.Store)
	isConnected := driveService.IsConnected()
	
	data := uimodel.DashboardData{
		User: true,
		Paths: paths,
		Logs: logs,
		Settings: settings,
		DriveLinked: isConnected,
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