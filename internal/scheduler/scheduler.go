package scheduler

import (
	"context"
	"log"

	"github.com/d4ve-p/clonis/internal/backup"
	"github.com/d4ve-p/clonis/internal/database"
	"github.com/robfig/cron/v3"
)

type Service struct {
	Cron   *cron.Cron
	Store  *database.Store
	Engine *backup.Engine
}

func New(store *database.Store, engine *backup.Engine) *Service {
	return &Service {
		Cron: cron.New(),
		Store: store,
		Engine: engine,
	}
}

func (s* Service) Start() {
	log.Printf("Scheduler Service Starting...")
	
	settings, err := s.Store.GetSettings()
	if err != nil {
		log.Printf("Scheduler Error: Could not load settings: %v", err)
		return
	}
	
	interval := settings["backup_interval"]
	if interval == "" {
		interval = "0 3 * * *"
		// interval = "@every 20s" // Debug
	}
	
	_, err = s.Cron.AddFunc(interval, func() {
		log.Println("⏰ Cron Triggered: Starting Backup...")
		if err := s.Engine.RunNow(context.Background()); err != nil {
			log.Printf("⏰ Cron Job Failed: %v", err)
		}
	})
	
	if err != nil {
		log.Printf("CRITICAL: Invalid cron schedule '%s': %v", interval, err)
		return		
	}
	
	s.Cron.Start()
}

func (s *Service) Stop() {
	s.Cron.Stop()
}