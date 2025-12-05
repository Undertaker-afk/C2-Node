package network

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// P2PManager handles libp2p networking
type P2PManager struct {
	mu              sync.RWMutex
	ctx             context.Context
	cancel          context.CancelFunc
	peerID          string
	multiaddrs      []string
	messageHandlers map[string]MessageHandler
	isConnected     bool
}

// MessageHandler is a function that handles incoming messages
type MessageHandler func(ctx context.Context, peerID string, data []byte) error

// NewP2PManager creates a new P2P manager
func NewP2PManager(ctx context.Context) *P2PManager {
	childCtx, cancel := context.WithCancel(ctx)
	return &P2PManager{
		ctx:             childCtx,
		cancel:          cancel,
		messageHandlers: make(map[string]MessageHandler),
		isConnected:     false,
	}
}

// Initialize sets up P2P connectivity
func (pm *P2PManager) Initialize(peerID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.peerID = peerID
	log.Printf("P2P Manager initialized with peer ID: %s", peerID)

	// This would initialize libp2p host and protocols
	// For now, this is a stub implementation
	pm.isConnected = true
	return nil
}

// Connect connects to a remote peer
func (pm *P2PManager) Connect(ctx context.Context, remoteAddr string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	log.Printf("Connecting to peer at %s", remoteAddr)
	// This would use libp2p to connect
	return nil
}

// RegisterMessageHandler registers a handler for a specific protocol
func (pm *P2PManager) RegisterMessageHandler(protocol string, handler MessageHandler) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.messageHandlers[protocol] = handler
	log.Printf("Registered handler for protocol: %s", protocol)
}

// SendMessage sends a message to a peer
func (pm *P2PManager) SendMessage(ctx context.Context, peerID string, protocol string, data []byte) error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if !pm.isConnected {
		return fmt.Errorf("not connected to peer: %s", peerID)
	}

	log.Printf("Sending message to peer %s via protocol %s (%d bytes)", peerID, protocol, len(data))
	// This would use libp2p streams to send
	return nil
}

// GetPeerID returns the local peer ID
func (pm *P2PManager) GetPeerID() string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.peerID
}

// IsConnected returns the connection status
func (pm *P2PManager) IsConnected() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.isConnected
}

// AddMultiaddr adds a multiaddr for this peer
func (pm *P2PManager) AddMultiaddr(addr string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.multiaddrs = append(pm.multiaddrs, addr)
}

// GetMultiaddrs returns all multiaddrs for this peer
func (pm *P2PManager) GetMultiaddrs() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.multiaddrs
}

// Close closes the P2P manager
func (pm *P2PManager) Close() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.cancel()
	pm.isConnected = false
	log.Println("P2P Manager closed")
	return nil
}
