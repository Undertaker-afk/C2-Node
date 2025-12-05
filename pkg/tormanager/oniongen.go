package tormanager

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"runtime"
	"sync"
)

// OnionKeyGenerator handles generation of Tor v3 onion keys
type OnionKeyGenerator struct {
	pattern string
	workers int
}

// NewOnionKeyGenerator creates a new onion key generator
func NewOnionKeyGenerator(pattern string) *OnionKeyGenerator {
	workers := runtime.NumCPU()
	return &OnionKeyGenerator{
		pattern: pattern,
		workers: workers,
	}
}

// GenerateVanityOnion generates a v3 onion address matching the pattern
func (g *OnionKeyGenerator) GenerateVanityOnion() (onionAddr string, privateKey ed25519.PrivateKey, err error) {
	if g.pattern == "" {
		return g.generateRandomOnion()
	}

	re, err := regexp.Compile(g.pattern)
	if err != nil {
		return "", nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	resultChan := make(chan struct {
		addr string
		key  ed25519.PrivateKey
	}, 1)

	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	for i := 0; i < g.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stopChan:
					return
				default:
					pub, priv, err := ed25519.GenerateKey(rand.Reader)
					if err != nil {
						continue
					}

					onionAddr := g.deriveOnionAddress(pub)
					if re.MatchString(onionAddr) {
						select {
						case resultChan <- struct {
							addr string
							key  ed25519.PrivateKey
						}{onionAddr, priv}:
						default:
						}
						return
					}
				}
			}
		}()
	}

	result := <-resultChan
	close(stopChan)
	wg.Wait()

	return result.addr, result.key, nil
}

// generateRandomOnion generates a random v3 onion address
func (g *OnionKeyGenerator) generateRandomOnion() (string, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", nil, err
	}

	onionAddr := g.deriveOnionAddress(pub)
	return onionAddr, priv, nil
}

// deriveOnionAddress derives a v3 onion address from a public key
// This is a simplified version - real implementation would use Tor's proper v3 checksum algorithm
func (g *OnionKeyGenerator) deriveOnionAddress(pub ed25519.PublicKey) string {
	// For now, return a hex representation with .onion suffix
	// Real implementation would use proper v3 address derivation with checksum
	addr := hex.EncodeToString(pub[:16]) + ".onion"
	return addr
}

// ExportPrivateKey exports the private key in Tor format
func ExportPrivateKey(privKey ed25519.PrivateKey) ([]byte, error) {
	// Tor v3 format requires specific encoding
	return privKey, nil
}

// ImportPrivateKey imports a Tor format private key
func ImportPrivateKey(data []byte) (ed25519.PrivateKey, error) {
	if len(data) != 64 {
		return nil, fmt.Errorf("invalid private key length: %d", len(data))
	}
	return ed25519.PrivateKey(data), nil
}
