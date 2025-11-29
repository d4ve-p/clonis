package auth

import (
	"net/http"
	"time"
)

func (m* Manager) LoginHandler(w http.ResponseWriter, r* http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	inputPwd := r.FormValue("password")
	
	if inputPwd != m.password {
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}
	
	token := generateToken()
	
	m.mu.Lock()
	m.sessions[token] = time.Now().Add(24 * time.Hour)
	m.mu.Unlock()
	
	http.SetCookie(w, &http.Cookie{
		Name:	"session_token",
		Value:	token,
		Expires:time.Now().Add(24 * time.Hour),
		HttpOnly:true,
		Path:	"/",
	})
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m* Manager) LogoutHandler (w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err == nil {
		m.mu.Lock()
		delete(m.sessions, c.Value)
		m.mu.Unlock()
	}
	
	http.SetCookie(w, &http.Cookie{
		Name: "session_token",
		Value: "",
		Expires: time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path: "/",
	})
	
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}