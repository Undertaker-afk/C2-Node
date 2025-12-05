package scripts

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/remotemgmt/gobased-remote-mgmt/pkg/common"
)

// ScriptExecutor executes remote scripts
type ScriptExecutor struct {
	mu              sync.RWMutex
	maxTimeout      time.Duration
	runningScripts  map[string]*ScriptSession
}

// ScriptSession represents a running script session
type ScriptSession struct {
	ID          string
	Command     string
	StartTime   time.Time
	EndTime     *time.Time
	ExitCode    *int
	Stdout      bytes.Buffer
	Stderr      bytes.Buffer
	ctx         context.Context
	cancel      context.CancelFunc
	cmd         *exec.Cmd
}

// NewScriptExecutor creates a new script executor
func NewScriptExecutor() *ScriptExecutor {
	return &ScriptExecutor{
		maxTimeout:     5 * time.Minute,
		runningScripts: make(map[string]*ScriptSession),
	}
}

// ExecuteScript executes a script synchronously
func (se *ScriptExecutor) ExecuteScript(ctx context.Context, req common.ScriptExecutionRequest) (*common.ScriptExecutionResult, error) {
	se.mu.Lock()

	// Check if script already running
	if _, exists := se.runningScripts[req.ScriptID]; exists {
		se.mu.Unlock()
		return nil, fmt.Errorf("script already running: %s", req.ScriptID)
	}

	session := &ScriptSession{
		ID:        req.ScriptID,
		Command:   req.Content,
		StartTime: time.Now(),
	}

	session.ctx, session.cancel = context.WithCancel(ctx)
	se.runningScripts[req.ScriptID] = session
	se.mu.Unlock()

	defer func() {
		se.mu.Lock()
		delete(se.runningScripts, req.ScriptID)
		se.mu.Unlock()
	}()

	// Apply timeout
	timeout := time.Duration(req.Timeout) * time.Second
	if timeout == 0 {
		timeout = se.maxTimeout
	}

	ctx, cancel := context.WithTimeout(session.ctx, timeout)
	defer cancel()

	// Execute command
	var cmd *exec.Cmd
	if len(req.Args) > 0 {
		cmd = exec.CommandContext(ctx, req.Content, req.Args)
	} else {
		cmd = exec.CommandContext(ctx, req.Content)
	}

	cmd.Stdout = &session.Stdout
	cmd.Stderr = &session.Stderr

	err := cmd.Run()
	
	now := time.Now()
	session.EndTime = &now
	
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	session.ExitCode = &exitCode

	result := &common.ScriptExecutionResult{
		ScriptID: req.ScriptID,
		ExitCode: exitCode,
		Stdout:   session.Stdout.String(),
		Stderr:   session.Stderr.String(),
		Duration: int(now.Sub(session.StartTime).Seconds()),
	}

	log.Printf("Script %s completed with exit code %d in %d seconds", req.ScriptID, exitCode, result.Duration)

	return result, nil
}

// GetSessionStatus gets the status of a script session
func (se *ScriptExecutor) GetSessionStatus(scriptID string) *ScriptSession {
	se.mu.RLock()
	defer se.mu.RUnlock()

	return se.runningScripts[scriptID]
}

// CancelScript cancels a running script
func (se *ScriptExecutor) CancelScript(scriptID string) error {
	se.mu.RLock()
	session, exists := se.runningScripts[scriptID]
	se.mu.RUnlock()

	if !exists {
		return fmt.Errorf("script not found: %s", scriptID)
	}

	session.cancel()
	log.Printf("Script %s cancelled", scriptID)

	return nil
}

// GetRunningScripts returns all running scripts
func (se *ScriptExecutor) GetRunningScripts() []string {
	se.mu.RLock()
	defer se.mu.RUnlock()

	ids := make([]string, 0, len(se.runningScripts))
	for id := range se.runningScripts {
		ids = append(ids, id)
	}

	return ids
}

// WaitForScript waits for a script to complete
func (se *ScriptExecutor) WaitForScript(scriptID string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			se.mu.RLock()
			session, exists := se.runningScripts[scriptID]
			se.mu.RUnlock()

			if !exists {
				return nil
			}

			if session.EndTime != nil {
				return nil
			}
		}
	}
}
