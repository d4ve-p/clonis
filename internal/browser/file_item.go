package browser

import (
	"os"
	"path/filepath"
	"sort"
)

type FileItem struct {
	Name string
	Path string
	Type FileType // file or folder
	ParentDir string // for going back
}

func GetPath(path string) ([]FileItem, error) {
	__DEFAULT_ROOT_PATH := os.Getenv("ROOT_FOLDER")
	if __DEFAULT_ROOT_PATH == "" {
		__DEFAULT_ROOT_PATH = "/hostfs"
	}
	
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	
	var items []FileItem
	
	if path == "" {
		path = __DEFAULT_ROOT_PATH
	}
	
	if path != __DEFAULT_ROOT_PATH {
		items = append(items, FileItem{
			Name: "..",
			Path: filepath.Dir(path),
			Type: Folder,
			ParentDir: "",
		})
	}
	
	for _, e := range entries {
		fullPath := filepath.Join(path, e.Name())
		
		itemType := File
		if e.IsDir() {
			itemType = Folder
		} else if e.Type()&os.ModeSymlink != 0 {
			if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
				itemType = Folder
			}
		}
		
		items = append(items, FileItem{
			Name: e.Name(),
			Path: fullPath,
			Type: itemType,
			ParentDir: path,
		})
	}
	
	sort.Slice(items, func(i, j int) bool {
		if items[i].Type != items[j].Type {
			return items[i].Type == Folder
		}
		
		return items[i].Name < items[j].Name
	})
	
	return items, nil
}