package ui

import (
	"html/template"
	"log"
	"net/http"

	"github.com/d4ve-p/clonis/internal/database"
)

type Handler struct {
	Store *database.Store
}

func New(store *database.Store) *Handler {
	return &Handler {
		Store: store,
	}
}

func (h *Handler) RenderError(w http.ResponseWriter, msg string, err error) {
	if err != nil {
		log.Printf("App Error: %s - %v", msg, err)
	}
	
	w.Header().Set("HX-Retarget", "#toast")
	w.Header().Set("HX-Reswap", "innerHTML")
	
	tmpl, templateError := template.ParseFiles("templates/fragments/error_toast.html")
	if templateError != nil {
		log.Printf("Critical error: Missing error_toast.html template: %v", templateError)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	
	data := struct{Message string}{Message: msg}
	
	if executeErr := tmpl.Execute(w, data); executeErr != nil {
		log.Printf("Error rendering toast template: %v", executeErr)
	}
}
