package ui

import (
	"context"
	"errors"
	"log"
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
	
	if h.Engine.IsRunning() {
		h.UI.RenderError(w, "Backup is already in progress", errors.New("concurrent backup request ignored"))
		return
	}
	
	go func() {
		if err := h.Engine.RunNow(context.Background()); err != nil {
			log.Printf("Manual backup failed: %v", err)
		}
	}()
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}