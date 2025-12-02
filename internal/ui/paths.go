package ui

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/d4ve-p/clonis/internal/database"
)

func (h *Handler) AddPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathStr := r.FormValue("path")
	pathType := r.FormValue("type")

	// Basic Validation
	if pathStr == "" {
		h.RenderError(w, "Path is required", nil)
		return
	}
	
	if _, err := os.ReadDir(pathStr); err != nil {
		h.RenderError(w, "Path does not exist or permisison denied", err)
		log.Printf("Error adding path: %v", err)
		return
	}

	// Add to Database
	err := h.Store.AddPath(pathStr, pathType)
	if err != nil {
		// check error, eg. unique path constraint
		h.RenderError(w, fmt.Sprintf("Error adding path: %v", err), err)
	}

	h.renderPathsList(w)
}

func (h *Handler) DeletePathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get ID from Query String (?id=1)
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Delete from DB
	if err := h.Store.DeletePath(id); err != nil {
		log.Printf("Error deleting path: %v", err)
	}

	h.renderPathsList(w)
}

func (h *Handler) renderPathsList(w http.ResponseWriter) {
	paths, err := h.Store.GetPaths()
	if err != nil {
		h.RenderError(w, "Failed to fetch paths", err)
		return
	}

	data := struct{ Paths []database.Path }{Paths: paths}

	tmpl := template.Must(template.ParseFiles("templates/fragments/path_list.html"))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Render error: %v", err)
	}
}
