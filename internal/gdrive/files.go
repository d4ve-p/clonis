package gdrive

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/drive/v3"
)

func (s* Service) UploadFile(ctx context.Context, localPath string, filename string) (*drive.File, error) {
	// Get target server folder
	folderID, err := s.SetupBackupFolder(ctx)
	if err != nil {
		return nil, err
	}
	
	client, err := s.GetClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to authenticate: %w", err)
	}
	
	// Open local file
	f, err := os.Open(localPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	// Drive metadata preparation
	fileMetadata := &drive.File {
		Name: filename,
		Parents: []string{folderID},
	}
	
	driveFile, err := client.Files.Create(fileMetadata).Media(f).Do()
	if err != nil {
		return nil, err
	}
	
	return driveFile, nil
}