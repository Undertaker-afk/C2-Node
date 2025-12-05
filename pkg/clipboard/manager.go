package clipboard

import (
	"context"
	"log"
	"sync"
	"time"
)

// ClipboardManager manages clipboard synchronization
type ClipboardManager struct {
	mu              sync.RWMutex
	ctx             context.Context
	cancel          context.CancelFunc
	currentContent  string
	updateCallback  UpdateCallback
	lastSyncTime    time.Time
	autoSyncEnabled bool
	syncInterval    time.Duration
}

// UpdateCallback is called when clipboard content changes
type UpdateCallback func(ctx context.Context, content string) error

// NewClipboardManager creates a new clipboard manager
func NewClipboardManager(ctx context.Context) *ClipboardManager {
	childCtx, cancel := context.WithCancel(ctx)
	return &ClipboardManager{
		ctx:             childCtx,
		cancel:          cancel,
		currentContent:  "",
		lastSyncTime:    time.Now(),
		autoSyncEnabled: false,
		syncInterval:    2 * time.Second,
	}
}

// SetUpdateCallback sets the callback for content updates
func (cm *ClipboardManager) SetUpdateCallback(callback UpdateCallback) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.updateCallback = callback
}

// SetContent sets the local clipboard content
func (cm *ClipboardManager) SetContent(content string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// TODO: Use atotto/clipboard library to set system clipboard
	cm.currentContent = content
	cm.lastSyncTime = time.Now()

	log.Printf("Clipboard content updated (%d bytes)", len(content))

	return nil
}

// GetContent gets the local clipboard content
func (cm *ClipboardManager) GetContent() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// TODO: Use atotto/clipboard library to get system clipboard
	return cm.currentContent
}

// SyncWithRemote synchronizes with remote clipboard
func (cm *ClipboardManager) SyncWithRemote(remoteContent string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if remoteContent == cm.currentContent {
		return nil
	}

	cm.currentContent = remoteContent
	cm.lastSyncTime = time.Now()

	log.Printf("Clipboard synchronized with remote (%d bytes)", len(remoteContent))

	// Call update callback if set
	if cm.updateCallback != nil {
		go func() {
			err := cm.updateCallback(cm.ctx, remoteContent)
			if err != nil {
				log.Printf("Error updating clipboard: %v", err)
			}
		}()
	}

	return nil
}

// StartAutoSync starts automatic clipboard synchronization
func (cm *ClipboardManager) StartAutoSync(callback UpdateCallback) error {
	cm.mu.Lock()
	cm.autoSyncEnabled = true
	cm.updateCallback = callback
	cm.mu.Unlock()

	go func() {
		ticker := time.NewTicker(cm.syncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-cm.ctx.Done():
				return
			case <-ticker.C:
				cm.checkAndSync()
			}
		}
	}()

	log.Println("Clipboard auto-sync started")
	return nil
}

// checkAndSync checks for clipboard changes and syncs
func (cm *ClipboardManager) checkAndSync() {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// TODO: Check if local clipboard changed and trigger sync
	// This would involve using atotto/clipboard to poll the system clipboard
}

// StopAutoSync stops automatic synchronization
func (cm *ClipboardManager) StopAutoSync() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.autoSyncEnabled = false
	log.Println("Clipboard auto-sync stopped")

	return nil
}

// GetLastSyncTime returns the last sync time
func (cm *ClipboardManager) GetLastSyncTime() time.Time {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.lastSyncTime
}

// IsAutoSyncEnabled returns whether auto-sync is enabled
func (cm *ClipboardManager) IsAutoSyncEnabled() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.autoSyncEnabled
}

// Close closes the clipboard manager
func (cm *ClipboardManager) Close() error {
	cm.cancel()
	log.Println("Clipboard Manager closed")
	return nil
}
