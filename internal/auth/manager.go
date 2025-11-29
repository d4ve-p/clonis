package auth

import (
	"os"
	"sync"
	"time"
)

type Manager struct {
	password string
	sessions map[string]time.Time
	mu sync.Mutex
}

func New() *Manager {
	pwd := os.Getenv("CLONIS_PASSWORD")
	if pwd == "" {
		pwd = "clonis_admin"
	}
	
	return &Manager {
		password: pwd,
		sessions: make(map[string]time.Time),
	}
	
}