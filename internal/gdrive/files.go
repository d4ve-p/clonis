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

func (s *Service) PruneOldBackups(ctx context.Context, retentionCount int) error {
	if retentionCount < 1 {
		return nil
	}

	folderID, err := s.SetupBackupFolder(ctx)
	if err != nil {
		return err
	}

	client, err := s.GetClient(ctx)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	list, err := client.Files.List().
		Q(query).
		OrderBy("createdTime").
		Fields("files(id, name, createdTime)").
		Do()
	if err != nil {
		return fmt.Errorf("failed to list files for pruning: %w", err)
	}

	totalFiles := len(list.Files)
	if totalFiles <= retentionCount {
		return nil // No cleanup needed
	}

	toDeleteCount := totalFiles - retentionCount
	
	for i := range toDeleteCount {
		file := list.Files[i]
		if err := client.Files.Delete(file.Id).Do(); err != nil {
			fmt.Printf("Failed to delete old backup %s: %v\n", file.Name, err)
			continue
		}
	}

	return nil
}