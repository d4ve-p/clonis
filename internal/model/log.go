package model

import "time"

type LogEntry struct {
	ID             int
	Status         string
	Message        string
	TotalSizeBytes int64
	StartedAt      time.Time
	CompletedAt    *time.Time
}