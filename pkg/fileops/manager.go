package fileops

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/remotemgmt/gobased-remote-mgmt/pkg/common"
)

// FileOperationManager manages remote file operations
type FileOperationManager struct {
	mu              sync.RWMutex
	basePath        string
	maxFileSize     int64
	allowedPaths    map[string]bool
}

// NewFileOperationManager creates a new file operation manager
func NewFileOperationManager(basePath string) (*FileOperationManager, error) {
	info, err := os.Stat(basePath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("base path must be a directory: %s", basePath)
	}

	return &FileOperationManager{
		basePath:     basePath,
		maxFileSize:  1024 * 1024 * 100, // 100MB
		allowedPaths: make(map[string]bool),
	}, nil
}

// AddAllowedPath adds a path that can be accessed
func (fm *FileOperationManager) AddAllowedPath(path string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	absPath, err := filepath.Abs(filepath.Join(fm.basePath, path))
	if err != nil {
		return err
	}

	if !filepath.HasPrefix(absPath, fm.basePath) {
		return fmt.Errorf("path outside base directory: %s", path)
	}

	fm.allowedPaths[absPath] = true
	return nil
}

// ListFiles lists files in a directory
func (fm *FileOperationManager) ListFiles(path string) ([]common.FileInfo, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	fullPath := filepath.Join(fm.basePath, path)
	if err := fm.validatePath(fullPath); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		log.Printf("Error reading directory %s: %v", fullPath, err)
		return nil, err
	}

	files := make([]common.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, common.FileInfo{
			Name:    entry.Name(),
			Size:    info.Size(),
			Mode:    info.Mode().String(),
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime(),
		})
	}

	return files, nil
}

// DownloadFile downloads a file
func (fm *FileOperationManager) DownloadFile(path string) ([]byte, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	fullPath := filepath.Join(fm.basePath, path)
	if err := fm.validatePath(fullPath); err != nil {
		return nil, err
	}

	stat, err := os.Stat(fullPath)
	if err != nil {
		log.Printf("Error stat file %s: %v", fullPath, err)
		return nil, err
	}

	if stat.Size() > fm.maxFileSize {
		return nil, fmt.Errorf("file too large: %d > %d", stat.Size(), fm.maxFileSize)
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("cannot download directory: %s", path)
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		log.Printf("Error reading file %s: %v", fullPath, err)
		return nil, err
	}

	return data, nil
}

// UploadFile uploads a file
func (fm *FileOperationManager) UploadFile(path string, data []byte) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if int64(len(data)) > fm.maxFileSize {
		return fmt.Errorf("file too large: %d > %d", len(data), fm.maxFileSize)
	}

	fullPath := filepath.Join(fm.basePath, path)
	if err := fm.validatePath(fullPath); err != nil {
		return err
	}

	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Error creating directory %s: %v", dir, err)
		return err
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		log.Printf("Error writing file %s: %v", fullPath, err)
		return err
	}

	return nil
}

// DeleteFile deletes a file
func (fm *FileOperationManager) DeleteFile(path string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fullPath := filepath.Join(fm.basePath, path)
	if err := fm.validatePath(fullPath); err != nil {
		return err
	}

	if err := os.Remove(fullPath); err != nil {
		log.Printf("Error deleting file %s: %v", fullPath, err)
		return err
	}

	return nil
}

// GetFileInfo gets information about a file
func (fm *FileOperationManager) GetFileInfo(path string) (*common.FileInfo, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	fullPath := filepath.Join(fm.basePath, path)
	if err := fm.validatePath(fullPath); err != nil {
		return nil, err
	}

	stat, err := os.Stat(fullPath)
	if err != nil {
		log.Printf("Error stat file %s: %v", fullPath, err)
		return nil, err
	}

	return &common.FileInfo{
		Name:    stat.Name(),
		Size:    stat.Size(),
		Mode:    stat.Mode().String(),
		IsDir:   stat.IsDir(),
		ModTime: stat.ModTime(),
	}, nil
}

// validatePath ensures the path is within the base directory
func (fm *FileOperationManager) validatePath(fullPath string) error {
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return err
	}

	absBase, err := filepath.Abs(fm.basePath)
	if err != nil {
		return err
	}

	if !filepath.HasPrefix(absPath, absBase) {
		return fmt.Errorf("path outside base directory: %s", fullPath)
	}

	return nil
}

// StreamDownload returns a reader for streaming downloads
func (fm *FileOperationManager) StreamDownload(path string) (io.ReadCloser, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	fullPath := filepath.Join(fm.basePath, path)
	if err := fm.validatePath(fullPath); err != nil {
		return nil, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		log.Printf("Error opening file %s: %v", fullPath, err)
		return nil, err
	}

	return file, nil
}

// StreamUpload writes from a reader to a file
func (fm *FileOperationManager) StreamUpload(path string, reader io.Reader) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fullPath := filepath.Join(fm.basePath, path)
	if err := fm.validatePath(fullPath); err != nil {
		return err
	}

	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Error creating directory %s: %v", dir, err)
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		log.Printf("Error creating file %s: %v", fullPath, err)
		return err
	}
	defer file.Close()

	_, err = io.CopyN(file, reader, fm.maxFileSize)
	if err != nil && err != io.EOF {
		log.Printf("Error writing file %s: %v", fullPath, err)
		return err
	}

	return nil
}
