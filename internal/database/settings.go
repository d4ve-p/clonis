package database

func (s* Store) GetSettings() (map[string]string, error) {
	rows, err := s.Db.Query("SELECT key, value from settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	settings := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		settings[k] = v
	}
	return settings, nil
}

func (s *Store) UpdateSetting(key, value string) error {
	_, err := s.Db.Exec("INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?", key, value, value)
	return err
}