package ui

import (
	"net/http"

	"github.com/d4ve-p/clonis/internal/backup"
)

type BackupHandler struct {
	UI *Handler
	Engine *backup.Engine
}

func (h *BackupHandler) Run(w http.ResponseWriter, r *http.Request) {
	// Method check
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	if err := h.Engine.RunNow(r.Context()); err != nil {
		h.UI.RenderError(w, "Backup Error: " + err.Error(), err)
		return
	}
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}