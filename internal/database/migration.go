package database

import "log"

func (s *Store) migrate() error {
	query := `
	-- 1. Paths: The "Backup Schema"
		-- This list defines the 'Manifest' of what gets zipped up.
		CREATE TABLE IF NOT EXISTS paths (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL UNIQUE,           -- e.g. /hostfs/var/www
			type TEXT DEFAULT 'folder',          -- 'file' or 'folder' (for UI icons)
			added_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- 2. Settings: Global configuration
		-- Stores: 'backup_interval' (cron), 'gdrive_token' (json), 'retention_count' (int)
		CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT
		);

		-- 3. Logs: History of the backup executions
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			status TEXT,                         -- 'SUCCESS', 'FAILED'
			message TEXT,                        -- Summary or Error details
			total_size_bytes INTEGER DEFAULT 0,  -- Size of the final zip
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			completed_at DATETIME
		);
	`

	_, err := s.Db.Exec(query)
	if err != nil {
		log.Printf("Error running migration: %v", err)
		return err
	}

	s.seedDefaults()

	log.Println("Database initialized with Simplified Schema")
	return nil
}
