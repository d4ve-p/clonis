package uimodel

import (
	"github.com/d4ve-p/clonis/internal/database"
	"github.com/d4ve-p/clonis/internal/model"
)

type DashboardData struct {
	User bool
	Paths []database.Path
	Logs []model.LogEntry
	Settings map[string]string
	DriveLinked bool
}