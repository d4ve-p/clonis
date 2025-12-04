package ui

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/d4ve-p/clonis/internal/scheduler"
	"github.com/robfig/cron/v3"
)

type SettingsHandler struct {
	UI *Handler
	Scheduler *scheduler.Service
}

func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.UI.RenderError(w, "Method not allowed", nil)
		return
	}
	
	interval := strings.TrimSpace(r.FormValue("backup_interval"))
	retention := strings.TrimSpace(r.FormValue("retention"))
	
	// Validate retention is a number
	retention_int, err := strconv.Atoi(retention)
	if err != nil {
		h.UI.RenderError(w, "Retention must be a number", err)
		return
	}
	
	// Validate retention >= 1
	if retention_int < 1 {
		h.UI.RenderError(w, "Retention count must at least be 1", nil)
		return
	}
	
	// Cron syntax validation
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, err := parser.Parse(interval); err != nil {
		h.UI.RenderError(w, "Invalid Cron format. Try '0 3 * * *'" ,err)
		return
	}
	
	// Updating database
	if err := h.UI.Store.UpdateSetting("backup_interval", interval); err != nil {
		h.UI.RenderError(w, "Failed to save interval to database", err)
		return
	}
	if err := h.UI.Store.UpdateSetting("retention", retention); err != nil {
		h.UI.RenderError(w, "Failed to save retention to database", err)
		return
	}
	
	// Restart to apply immediately
	h.Scheduler.Restart()
	
	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}