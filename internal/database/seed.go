package database

func (s *Store) seedDefaults() {
	defaults := map[string]string {
		"backup_interval": "0 3 * * *",
		"retention_count": "5",
	}
	
	for k, v := range defaults {
		s.Db.Exec("INSERT OR IGNORE INTO settings (key, value) VALUES (?, ?)", k, v)
	}
}