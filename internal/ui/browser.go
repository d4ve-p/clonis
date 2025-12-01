package ui

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/d4ve-p/clonis/internal/browser"
)

func (h *Handler) BrowseHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	
	path = filepath.Clean(path)
	
	defaultPath := os.Getenv("ROOT_FOLDER")
	if defaultPath == "" {
		defaultPath = "/hostfs"
	}
	
	if !strings.HasPrefix(path, defaultPath) {
		path = defaultPath
	}
	
	entries, err := browser.GetPath(path)
	if err != nil {
		http.Error(w, "Permission denied or invalid path", http.StatusForbidden)
		return
	}
	
	tmpl := template.Must(template.ParseFiles("templates/fragments/browser.html"))
	
	data := struct {
		Current string
		Items []browser.FileItem
	} {
		Current: path,
		Items: entries,
	}
	
	tmpl.Execute(w, data)
}