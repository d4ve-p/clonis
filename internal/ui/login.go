package ui

import (
	"html/template"
	"net/http"
)

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/layout.html",
		"templates/login.html",	
	))
	
	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}