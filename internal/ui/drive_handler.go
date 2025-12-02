package ui

import (
	"net/http"

	"github.com/d4ve-p/clonis/internal/gdrive"
)

type DriveHandler struct {
	UI *Handler
	Service *gdrive.Service
}

func (h *DriveHandler) Connect(w http.ResponseWriter, r *http.Request) {
	url := h.Service.GetAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *DriveHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.UI.RenderError(w, "No authorization code returned from Google.", nil)
		return
	}

	// Exchange and Save
	err := h.Service.HandleCallback(r.Context(), code)
	if err != nil {
		h.UI.RenderError(w, "Failed to connect Google Drive.", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}