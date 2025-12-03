package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/d4ve-p/clonis/internal/database"
)

// CreateArchive zips the targets into a single file at destPath
func CreateArchive(destPath string, targets []database.Path) error {
	// Create the Zip File
	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()

	// Iterate over every path in our manifest
	for _, target := range targets {
		// We walk the path (works for both file and folder)
		err := filepath.Walk(target.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// If one file fails (permission?), log it but don't stop the whole backup
				// For now, we return nil to continue walking
				return nil 
			}

			// Create the Header for the Zip Entry
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// 4. Sanitize the Name
			// Turn "/hostfs/var/www/index.html" -> "var/www/index.html"
			// This makes the backup portable and clean.
			
			root_prefix := os.Getenv("ROOT_FOLDER")
			relPath := strings.TrimPrefix(path, root_prefix)
			relPath = strings.TrimPrefix(relPath, "/") // Remove leading slash
			header.Name = relPath

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate // Compress files
			}

			// 5. Write Header
			writer, err := w.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			// 6. Write File Content
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}