package database

import (
	"database/sql"
	"time"
)

type LogEntry struct {
	ID             int
	Status         string
	Message        string
	TotalSizeBytes int64
	StartedAt      time.Time
	CompletedAt    *time.Time
}

func (s *Store) GetRecentLogs(limit int) ([]LogEntry, error) {
	rows, err := s.Db.Query(`
		SELECT id, status, message, total_size_bytes, started_at, completed_at 
		FROM logs 
		ORDER BY id DESC 
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var l LogEntry
		var completedAt sql.NullTime
		
		if err := rows.Scan(&l.ID, &l.Status, &l.Message, &l.TotalSizeBytes, &l.StartedAt, &completedAt); err != nil {
			return nil, err
		}
		
		if completedAt.Valid {
			t := completedAt.Time
			l.CompletedAt = &t
		}
		logs = append(logs, l)
	}
	return logs, nil
}