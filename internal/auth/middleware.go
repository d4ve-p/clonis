package auth

import (
	"net/http"
	"time"
)

func (m* Manager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" || r.URL.Path == "/static/style.css" {
			next.ServeHTTP(w, r)
			return
		}
		
		c, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
		m.mu.Lock()
		expiry, exists := m.sessions[c.Value]
		m.mu.Unlock()
		
		if !exists || time.Now().After(expiry) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return 
		}
		
		next.ServeHTTP(w, r)
	})
}