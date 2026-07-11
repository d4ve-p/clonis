package backup

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/d4ve-p/clonis/internal/database"
	"github.com/d4ve-p/clonis/internal/gdrive"
	"github.com/d4ve-p/clonis/internal/model"
)

type Engine struct {
	Store     *database.Store
	Drive     *gdrive.Service
	mu        sync.Mutex
	isRunning bool
}

func New(store *database.Store, drive *gdrive.Service) *Engine {
	return &Engine {
		Store: store,
		Drive: drive,
	}
}

func CleanUpTempFiles() error {
	log.Println("Cleaning up temp files from previous sessions...")

	tmpDir := os.TempDir()

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to read temp dir: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "backup_") && strings.HasSuffix(file.Name(), ".zip") {
			filePath := filepath.Join(tmpDir, file.Name())
			log.Printf("Removing orphaned temp file: %s", filePath)
			if err := os.Remove(filePath); err != nil {
				log.Printf("Warning: Failed to remove orphaned temp file %s: %v", filePath, err)
			}
		}
	}

	return nil
}

func (e *Engine) RunNow(ctx context.Context) error {
	e.mu.Lock()
	if e.isRunning {
		e.mu.Unlock()
		return errors.New("backup already in progress")
	}
	e.isRunning = true
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		e.isRunning = false
		e.mu.Unlock()
	}()

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
	
	if !e.Drive.IsConnected() {
		updateLogHelper("SKIPPED", "Google Drive is not connected, Please connect in settings", 0)
		return nil
	}
	
	if _, err := e.Drive.SetupBackupFolder(ctx); err != nil {
		updateLogHelper("FAILED", "Google Drive connection test failed (Revoked?): "+err.Error(), 0)
		return err
	}
	
	paths, err := e.Store.GetPaths()
	if err != nil {
		updateLogHelper("FAILED", "Database error: " + err.Error(), 0)
		return err
	}
	
	if len(paths) == 0 {
		updateLogHelper("CANCELED", "No registered path", 0)
		return errors.New("No registered path to backup")
	}
	
	// Prepare temp directory
	tmpDir := os.TempDir()
	fileName := fmt.Sprintf("backup_%s.zip", time.Now().Format("2006-01-02_15-04-05"))
	localZipPath := fmt.Sprintf("%s/%s", tmpDir, fileName)

	// Ensures that Zip is always cleaned up after function exits
	defer os.Remove(localZipPath); 
	
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
	
	// Retention Policy
	settings, _ := e.Store.GetSettings()
		retentionStr := settings["retention_count"]
		retention, _ := strconv.Atoi(retentionStr)
		
		if retention > 0 {
			log.Printf("Enforcing retention policy (Keep %d)...", retention)
			if err := e.Drive.PruneOldBackups(ctx, retention); err != nil {
				log.Printf("Warning: Failed to prune old backups: %v", err)
			}
		}
	
	updateLogHelper("SUCCESS", "Backup uploaded successfully!", size)
	log.Println("Backup job successfully completed")
	return nil
}

func (e *Engine) IsRunning() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.isRunning
}