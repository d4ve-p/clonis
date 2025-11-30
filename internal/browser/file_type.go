package browser

import "strings"

type FileType int

const (
	Folder FileType = iota
	File FileType = iota
)

var fileTypeNames = map[FileType]string {
	Folder: "folder",
	File: "file",
}

var fileTypeValeus = map[string]FileType {
	"folder": Folder,
	"file": File,
}

func (f FileType) String() string {
	if name, ok := fileTypeNames[f]; ok {
		return name
	}
	
	return "unknown"
}

func ParseFileType(s string) FileType {
	normalized := strings.ToLower(strings.TrimSpace(s))
	
	if val, ok := fileTypeValeus[normalized]; ok {
		return val
	}
	
	return File
}