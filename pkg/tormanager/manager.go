package tormanager

import (
	"crypto/ed25519"
	"fmt"
	"net"
)

// TorManager handles Tor connectivity and onion service management
type TorManager struct {
	enabled      bool
	dataDir      string
	controlAddr  string
	socksAddr    string
	privateKey   ed25519.PrivateKey
	onionAddress string
	listener     net.Listener
}

// NewTorManager creates a new Tor manager
func NewTorManager(dataDir string) (*TorManager, error) {
	return &TorManager{
		enabled:     false,
		dataDir:     dataDir,
		controlAddr: "127.0.0.1:9051",
		socksAddr:   "127.0.0.1:9050",
	}, nil
}

// Initialize sets up Tor connectivity
func (tm *TorManager) Initialize() error {
	// This would initialize Tor using bine library
	// For now, this is a stub implementation
	tm.enabled = true
	return nil
}

// LoadPrivateKey loads a previously generated Tor private key
func (tm *TorManager) LoadPrivateKey(keyData []byte) error {
	privKey, err := ImportPrivateKey(keyData)
	if err != nil {
		return err
	}
	tm.privateKey = privKey
	return nil
}

// GenerateOnionService creates a new onion service
func (tm *TorManager) GenerateOnionService(pattern string) (string, error) {
	gen := NewOnionKeyGenerator(pattern)
	addr, privKey, err := gen.GenerateVanityOnion()
	if err != nil {
		return "", err
	}

	tm.privateKey = privKey
	tm.onionAddress = addr
	return addr, nil
}

// GetOnionAddress returns the current onion address
func (tm *TorManager) GetOnionAddress() string {
	return tm.onionAddress
}

// SetOnionAddress sets the onion address
func (tm *TorManager) SetOnionAddress(addr string) {
	tm.onionAddress = addr
}

// ListenOnion creates a Tor onion listener on the given port
func (tm *TorManager) ListenOnion(port int) (net.Listener, error) {
	if !tm.enabled {
		return nil, fmt.Errorf("Tor not initialized")
	}

	// This would use bine to create a proper Tor listener
	// For now, return a TCP listener on localhost
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}

	tm.listener = listener
	return listener, nil
}

// Close closes the Tor manager
func (tm *TorManager) Close() error {
	if tm.listener != nil {
		return tm.listener.Close()
	}
	return nil
}
