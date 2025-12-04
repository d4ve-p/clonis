package scheduler

import (
	"context"
	"log"
)

func (s* Service) scheduleJob() {
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
	
	id, err := s.Cron.AddFunc(interval, func() {
		log.Println("⏰ Cron Triggered: Starting Backup...")
		if err := s.Engine.RunNow(context.Background()); err != nil {
			log.Printf("⏰ Cron Job Failed: %v", err)
		}
	})
	
	if err != nil {
		log.Printf("CRITICAL: Invalid cron schedule '%s': %v", interval, err)
		return		
	}
	
	s.entryID = id
}