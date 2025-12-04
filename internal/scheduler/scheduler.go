package scheduler

import (
	"log"

	"github.com/d4ve-p/clonis/internal/backup"
	"github.com/d4ve-p/clonis/internal/database"
	"github.com/robfig/cron/v3"
)

type Service struct {
	Cron   *cron.Cron
	Store  *database.Store
	Engine *backup.Engine
	entryID cron.EntryID
}

func New(store *database.Store, engine *backup.Engine) *Service {
	return &Service {
		Cron: cron.New(),
		Store: store,
		Engine: engine,
	}
}

func (s* Service) Start() {
	s.scheduleJob()
	s.Cron.Start()
}

func (s* Service) Restart() {
	if s.entryID != 0 {
		s.Cron.Remove(s.entryID)
	}
	
	log.Println("♻️  Scheduler Restarting with new settings...")
	s.scheduleJob()
}

func (s *Service) Stop() {
	s.Cron.Stop()
}