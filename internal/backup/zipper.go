package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/d4ve-p/clonis/internal/database"
)

func CreateArchive(destPath string, targets []database.Path) error {
	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()
	
	root_folder := os.Getenv("ROOT_FOLDER")

	for _, target := range targets {
		walkRoot, err := filepath.EvalSymlinks(target.Path)
		if err != nil {
			fmt.Printf("Warning: could not resolve symblink: %v", err)
			continue
		}
		
		err = filepath.Walk(walkRoot, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Log permission errors but don't abort the backup
				return nil 
			}

			// Sanitize the Name
			relPath := strings.TrimPrefix(path, root_folder)

			// Create the Header
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = relPath

			// --- Handle Symlinks ---
			if info.Mode()&os.ModeSymlink != 0 {
				linkTarget, err := os.Readlink(path)
				if err != nil {
					fmt.Printf("Error reading symlink: %v", err)
					return nil
				}
				
				// Symlinks are stored as uncompressed text pointing to the target
				header.Method = zip.Store
				writer, err := w.CreateHeader(header)
				if err != nil {
					return err
				}
				_, err = writer.Write([]byte(linkTarget))
				return err
			}
			// -----------------------------

			if info.IsDir() {
				header.Name += "/"
				_, err = w.CreateHeader(header)
				return err
			}

			// Regular Files
			header.Method = zip.Deflate
			writer, err := w.CreateHeader(header)
			if err != nil {
				return err
			}

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