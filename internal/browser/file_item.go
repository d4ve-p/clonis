package browser

type FileItem struct {
	Name string
	Path string
	Type FileType // file or folder
	ParentDir string // for going back
}