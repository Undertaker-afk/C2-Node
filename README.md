# Go-Based Remote Management System

A comprehensive, pure-Go remote management toolkit designed for managing devices such as Raspberry Pis and other hardware over decentralized P2P networks with Tor anonymity and end-to-end encryption.

## System Overview

This project consists of three main components:

### 1. **GUI Panel** (`cmd/panel/`)
A Windows-optimized desktop control interface built with Fyne that allows:
- **Dashboard**: Overview of connected devices and system status
- **Device Management**: Monitor and interact with multiple remote stubs
- **File Browser**: Remote file listing, upload, and download
- **Terminal Emulation**: Interactive shell access to remote devices
- **System Metrics**: Real-time monitoring of CPU, memory, temperature, and uptime
- **Script Execution**: Execute scripts remotely with output streaming
- **Clipboard Sync**: Bidirectional clipboard synchronization
- **Configuration Management**: Bulk device configuration updates

**Mandatory 20-Second Splash Screen**: On startup, displays "THIS IS ONLY FOR EDUCATIONAL AND LEGITIMATE PURPOSES" for exactly 20 seconds as an ethical reminder before proceeding to the main interface.

### 2. **Stub Agent** (`cmd/stub/`)
A lightweight, headless agent that runs on target devices (Windows, Linux, ARM):
- **Cross-Platform Support**: Compiled for Windows (AMD64), Linux (AMD64), and Linux (ARM for Raspberry Pi)
- **Persistent Services**: Runs as background service with auto-start on boot
- **P2P Connectivity**: Joins libp2p swarms with automatic reconnection
- **Tor Integration**: Embeds Tor onion services for anonymous inbound connectivity
- **System Monitoring**: Collects and streams metrics to the panel
- **Anomaly Detection**: Local detection with alerts via Nostr DMs
- **Remote Execution**: Secure script execution in sandboxed environment
- **Session Logging**: Comprehensive audit logs for all operations
- **OTA Updates**: Self-updating capability with secure verification

### 3. **Builder Tool** (`cmd/builder/`)
Cross-compilation utility for generating platform-specific binaries:
- Compiles panel and stubs for multiple target platforms
- Embeds Tor onion keys using oniongen-go for persistent addresses
- Supports vanity onion generation via regex patterns
- Version stamping and security verification
- Command-line configurability

## Key Features

### Networking & Security
- **Decentralized P2P**: Uses libp2p for peer-to-peer connectivity without port forwarding
- **Tor Integration**: All traffic anonymized via Tor overlay network using bine
- **End-to-End Encryption**: Secure communication via libp2p Noise protocol
- **Serverless Operation**: No central server required; direct device-to-device communication
- **Persistence**: Automatic reconnection with connection state preservation

### Device Management
- **Multi-Device Support**: Manage multiple stubs simultaneously from single panel
- **Auto-Discovery**: Devices discovered via Nostr or libp2p DHT
- **Status Indicators**: Visual indicators for online/offline status and anomalies
- **Fleet Management**: Bulk operations across multiple devices

### Monitoring & Diagnostics
- **Real-Time Metrics**: CPU usage, memory, temperature, uptime
- **Anomaly Detection**: Statistical analysis with threshold-based alerts
- **Nostr Alerts**: Anomaly notifications via Nostr direct messages
- **Session Logging**: Complete audit trail of all operations

### Remote Access
- **VNC**: Remote desktop streaming with Tight/ZRLE compression
- **Shell Access**: Interactive terminal with PTY allocation
- **File Operations**: Secure file browsing, upload, and download
- **Clipboard Sync**: Seamless text transfer between panel and stub

### Automation
- **Script Execution**: Run commands/scripts with output streaming
- **Batch Operations**: Execute configurations across device fleet
- **OTA Updates**: Remotely update stub binaries with verification

## Library Stack

| Library | Purpose | Version |
|---------|---------|---------|
| `fyne` | GUI panel desktop application | v2.4.0 |
| `go-libp2p` | P2P networking core | v0.32.0 |
| `bine` | Tor embedding | v0.2.0 |
| `oniongen-go` | Tor v3 onion key generation | v0.2.0 |
| `go-nostr` | Nostr protocol implementation | v0.29.0 |
| `go-vnc` | VNC server/client | Latest |
| `pty` | PTY allocation for shells | v1.1.21 |
| `filebrowser` | Remote file operations | Latest |
| `keylogger` | Keystroke capture | Latest |
| `gopsutil` | System metrics collection | v3.21.11 |
| `gonum` | Statistical anomaly detection | v0.14.0 |
| `clipboard` | Clipboard synchronization | v0.1.4 |

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Panel (GUI)                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Dashboard │ Devices │ Files │ Terminal │ Metrics ...  │ │
│  └────────────────────────────────────────────────────────┘ │
│             │                │                               │
│             └────────────────┴──────────────────────────────┘│
└────────────────┬─────────────────────────────────────────────┘
                 │
        ┌────────┴────────┐
        │                 │
    ┌───────────────────────────────┐
    │  libp2p + Tor Swarm           │
    │  (NAT Traversal, DHT)         │
    └───────────────────────────────┘
        │                 │
    ┌───▼───┐         ┌───▼───┐         ┌───▼───┐
    │Stub 1 │         │Stub 2 │  ...    │Stub N │
    │Linux  │         │RPi    │         │Win    │
    │ARM    │         │       │         │       │
    └───────┘         └───────┘         └───────┘
        │                 │                 │
   ┌────────────────────────────────────────────┐
   │  Nostr Relays (Status & Alerts)            │
   └────────────────────────────────────────────┘
```

## Getting Started

### Prerequisites
- Go 1.21+
- For Windows: Visual Studio Build Tools (for CGO)
- For Linux: GCC/G++ toolchain
- Tor daemon (for full functionality)

### Installation

1. Clone the repository:
```bash
git clone <repo-url>
cd gobased-remote-mgmt
```

2. Download dependencies:
```bash
go mod download
```

### Building

Build all components:
```bash
go run ./cmd/builder -all -out ./bin
```

Build specific components:
```bash
# Panel only
go run ./cmd/builder -panel -out ./bin

# Stub for Linux ARM (Raspberry Pi)
go run ./cmd/builder -stub -out ./bin

# With vanity onion
go run ./cmd/builder -all -vanity "^mydevice" -out ./bin
```

### Running

#### Panel (Windows)
```bash
./bin/panel-windows-amd64.exe
```

#### Stub Agent
```bash
# Linux
./bin/stub-linux-amd64

# Raspberry Pi
./bin/stub-linux-arm

# Windows
./bin/stub-windows-amd64.exe
```

With options:
```bash
# Specify custom port and Nostr relays
./bin/stub-linux-amd64 -port 9000 -nostr-relays "wss://relay.damus.io,wss://relay.nostr.band"

# Enable debug logging
./bin/stub-linux-amd64 -debug

# Custom Tor vanity pattern
./bin/stub-linux-amd64 -vanity "^device[a-z0-9]{2}"
```

## Configuration

### Nostr Relays
Configure in the panel settings tab:
```
wss://relay.damus.io
wss://relay.nostr.band
wss://nos.lol
```

### libp2p Bootstraps
Default public bootstraps used. Can be customized in settings.

### Tor Configuration
- Embedded via bine library
- Onion keys generated/loaded from `~/.remote-mgmt/keys/`
- Configurable vanity patterns

## Ethical Guidelines

This tool is designed exclusively for:
- Educational purposes
- Authorized device management
- Legitimate system administration
- Responsible remote access scenarios

**Users must**:
- Obtain explicit authorization before connecting to devices
- Respect device owner privacy and rights
- Maintain detailed audit logs of all operations
- Comply with applicable laws and regulations
- Use keylogging only with explicit consent and for authorized purposes

The mandatory splash screen serves as a reminder of these responsibilities.

## Security Considerations

### Encryption
- End-to-end via libp2p Noise protocol
- Tor provides anonymity layer
- All credentials should be managed securely

### Key Management
- Tor onion keys stored in encrypted form
- libp2p peer identities managed locally
- Nostr keys should use hardware wallets when possible

### Access Control
- File operations sandboxed to allowed directories
- Script execution with configurable permissions
- Session logging for audit trails

## Performance Notes

### Panel Requirements
- Windows 10/11+
- Minimal CPU/memory overhead
- Responsive UI for up to 100+ connected devices

### Stub Optimization for Raspberry Pi
- Minimal goroutines for ARM efficiency
- Configurable polling intervals (default: 5s)
- Binary size: ~15-20MB
- Memory footprint: <50MB typical
- CPU usage: <5% idle, <20% under load

### Network
- Optimized for high-latency Tor connections
- Automatic reconnection with exponential backoff
- Connection pooling for multiple operations
- Compression for file transfers

## Project Structure

```
.
├── cmd/
│   ├── panel/          # GUI panel (Windows)
│   ├── stub/           # Stub agent (Windows/Linux/ARM)
│   └── builder/        # Build tool
├── pkg/
│   ├── common/         # Shared types
│   ├── tormanager/     # Tor integration
│   ├── network/        # libp2p networking
│   ├── metrics/        # System metrics collection
│   ├── nostr/          # Nostr protocol
│   ├── fileops/        # File operations
│   ├── scripts/        # Script execution
│   └── clipboard/      # Clipboard sync
├── README.md
├── LICENSE
└── go.mod
```

## Development

### Testing
```bash
go test ./...
```

### Building for Different Platforms
```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o stub-win.exe ./cmd/stub

# Linux x86_64
GOOS=linux GOARCH=amd64 go build -o stub-linux ./cmd/stub

# ARM (Raspberry Pi)
GOOS=linux GOARCH=arm GOARM=7 go build -o stub-arm ./cmd/stub
```

## Troubleshooting

### Tor Connection Issues
- Verify Tor daemon is running: `ps aux | grep tor`
- Check port accessibility: `netstat -an | grep 9050`
- Enable debug logging: `-debug` flag

### Metrics Not Collecting
- Verify gopsutil dependencies on Linux: `apt-get install sysstat`
- Check temperature sensors: `cat /sys/class/thermal/*/temp`

### libp2p DHT Timeout
- Add bootstrap nodes manually
- Check firewall rules
- Verify NAT traversal settings

## Contributing

Guidelines:
1. Follow existing code style
2. Add tests for new features
3. Update documentation
4. Ensure cross-platform compatibility
5. Test on ARM (if possible)

## License

[Specify appropriate license]

## References

- [libp2p Documentation](https://docs.libp2p.io/)
- [Fyne Documentation](https://docs.fyne.io/)
- [Nostr Protocol](https://github.com/nostr-protocol/nostr)
- [Tor Documentation](https://www.torproject.org/docs/)
- [oniongen-go](https://github.com/rdkr/oniongen-go)

## Support

For issues, questions, or contributions:
- Open GitHub issues
- Check documentation wiki
- Review example configurations

---

**Disclaimer**: This software is provided for educational and authorized use only. Users are responsible for ensuring they have proper authorization before accessing any systems. Misuse may violate applicable laws.
