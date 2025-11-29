package database

import "time"

type Path struct {
	ID int
	Path string
	Type string
	AddedAt time.Time
}

func (s* Store) GetPaths() ([]Path, error) {
	rows, err := s.Db.Query("SELECT id, path, type, added_at FROM paths ORDER BY path ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var paths []Path
	for rows.Next() {
		var p Path
		if err := rows.Scan(&p.ID, &p.Path, &p.Type, &p.AddedAt); err != nil {
			return nil, err
		}
		paths = append(paths, p);
	}
	
	return paths, nil
}

func (s *Store) AddPath(pathStr, pathType string) error {
	_, err := s.Db.Exec("INSERT INTO paths (path, type) VALUES (?, ?)", pathStr, pathType)
	return err
}

func (s *Store) DeletePath(id int) error {
	_, err := s.Db.Exec("DELETE FROM paths WHERE id = ?", id)
	return err
}