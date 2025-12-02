package uimodel

import "github.com/d4ve-p/clonis/internal/database"

type DashboardData struct {
	User bool
	Paths []database.Path
	Logs []database.LogEntry
	Settings map[string]string
	DriveLinked bool
}