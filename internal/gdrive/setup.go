package gdrive

import (
	"context"
	"errors"
	"fmt"
	"os"

	"google.golang.org/api/drive/v3"
)

func (s *Service) SetupBackupFolder(ctx context.Context) (string, error) {
	client, err := s.GetClient(ctx)
	if err != nil {
		return "", err
	}
	
	clonisID, err := s.getOrCreateFolder(client, "Clonis", "")
	if err != nil {
		return "", fmt.Errorf("Failed to setup root Clonis folder: %w", err)
	}
	
	serverID := os.Getenv("SERVER_NAME")
	if serverID == "" {
		return "", errors.New("SERVER_NAME env should be defined")
	}
	
	serverFolderID, err := s.getOrCreateFolder(client, serverID, clonisID)
	if err != nil {
		return "", fmt.Errorf("Failed to setup server folder: %w", err)
	}
	
	return serverFolderID, nil
}

func (s *Service) getOrCreateFolder(srv *drive.Service, name string, parentID string) (string, error) {
	query := fmt.Sprintf("mimeType = 'application/vnd.google-apps.folder' and name = '%s' and trashed = false", name)
	if parentID != "" {
		query += fmt.Sprintf(" and '%s' in parents", parentID)
	}

	list, err := srv.Files.List().Q(query).Fields("files(id)").Do()
	if err != nil {
		return "", err
	}

	if len(list.Files) > 0 {
		return list.Files[0].Id, nil
	}

	folderMetadata := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
	}
	
	if parentID != "" {
		folderMetadata.Parents = []string{parentID}
	}

	folder, err := srv.Files.Create(folderMetadata).Fields("id").Do()
	if err != nil {
		return "", err
	}

	return folder.Id, nil
}