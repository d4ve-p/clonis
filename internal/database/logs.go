package database

import (
	"database/sql"

	"github.com/d4ve-p/clonis/internal/model"
)

func (s *Store) GetRecentLogs(limit int) ([]model.LogEntry, error) {
	rows, err := s.Db.Query(`
		SELECT id, status, message, total_size_bytes, started_at, completed_at 
		FROM logs 
		ORDER BY id DESC 
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.LogEntry
	for rows.Next() {
		var l model.LogEntry
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

func (s* Store) CreateLog(log model.LogEntry) (int, error) {
	res, err := s.Db.Exec(`INSERT INTO Logs (status, message, total_size_bytes, started_at, completed_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, ?)`, log.Status, log.Message, log.TotalSizeBytes, log.CompletedAt)
	if err != nil {
		return -1, err
	}
	
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	
	return int(id), nil
}

func (s* Store) UpdateLog(id int, log model.LogEntry) error {
	_, err := s.Db.Exec(`UPDATE Logs SET status=?, message=?, total_size_bytes=?, started_at=?, completed_at=CURRENT_TIMESTAMP WHERE id=?`, log.Status, log.Message, log.TotalSizeBytes, id)
	return err
}