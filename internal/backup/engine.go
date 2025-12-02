package backup

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/d4ve-p/clonis/internal/database"
	"github.com/d4ve-p/clonis/internal/gdrive"
	"github.com/d4ve-p/clonis/internal/model"
)

type Engine struct {
	Store *database.Store
	Drive *gdrive.Service
}

func New(store *database.Store, drive *gdrive.Service) *Engine {
	return &Engine {
		Store: store,
		Drive: drive,
	}
}

func (e *Engine) RunNow(ctx context.Context) error {
	// Create initial log entry
	logID, err := e.Store.CreateLog(model.LogEntry{
		ID: 0,
		Status: "IN_PROGRESS",
		Message: "Backup Started",
		TotalSizeBytes: 0,
	})
	if err != nil {
		return fmt.Errorf("Failed to create log entry: %w", err)
	}
	
	// Helper to update log on exit
	updateLogHelper := func(status, msg string, size int64) {
		if err := e.Store.UpdateLog(model.LogEntry{
			ID: logID,
			Status: status,
			Message: msg,
			TotalSizeBytes: size,
		}); err != nil {
			log.Printf("Failed to update log status: %v", err)
		}
	}
	
	paths, err := e.Store.GetPaths()
	if err != nil {
		updateLogHelper("FAILED", "Database error: " + err.Error(), 0)
		return err
	}
	
	// Prepare temp directory
	tmpDir := os.TempDir()
	fileName := fmt.Sprintf("backup_%s.zip", time.Now().Format("2006-01-02_15-04-05"))
	localZipPath := fmt.Sprintf("%s/%s", tmpDir, fileName)
	
	// Creating an archive
	log.Printf("Zipping files to %s...", localZipPath)
	if err := CreateArchive(localZipPath, paths); err != nil {
		updateLogHelper("FAILED", "Zipping failed: " + err.Error(), 0)
		return err
	}
	
	// Check size
	fi, err := os.Stat(localZipPath)
	if err != nil {
		updateLogHelper("FAILED", "File stat failed", 0)
		return err
	}
	size := fi.Size()
	
	// Upload to drive
	log.Printf("Uploading %s (%d bytes)...", fileName, size)
	if _, err := e.Drive.UploadFile(ctx, localZipPath, fileName); err != nil {
		updateLogHelper("FAILED", "Fail uploading zip to drive: " + err.Error(), size)
		return err
	}
	
	os.Remove(localZipPath)
	
	updateLogHelper("SUCCESS", "Backup uploaded successfully!", size)
	log.Println("Backup job successfully completed")
	return nil
}