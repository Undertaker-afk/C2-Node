package nostr

import "errors"

var (
	ErrNotConnected = errors.New("nostr manager not connected")
	ErrNoRelays     = errors.New("no relays configured")
	ErrInvalidKey   = errors.New("invalid key format")
	ErrPublishFailed = errors.New("failed to publish event")
)
