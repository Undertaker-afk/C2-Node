# System Architecture Documentation

## Overview

The Go-based Remote Management System is designed as a decentralized, peer-to-peer infrastructure for managing remote devices with an emphasis on security, privacy, and ethical operation.

## Core Architectural Principles

1. **Decentralized**: No central server or dependency on third-party infrastructure
2. **End-to-End Encrypted**: All communications protected with encryption
3. **Privacy-First**: Tor integration ensures anonymity
4. **Serverless**: Direct device-to-device communication
5. **Persistent**: Automatic reconnection and state preservation
6. **Scalable**: Support for fleets of devices
7. **Ethical**: Mandatory disclaimers and audit logging

## Component Architecture

### Panel Component

```
┌─────────────────────────────────────────────────┐
│         Fyne GUI Application                    │
├─────────────────────────────────────────────────┤
│ Splash Screen (20 second disclaimer)            │
├─────────────────────────────────────────────────┤
│ Tab Interface:                                  │
│  ├─ Dashboard (overview/status)                │
│  ├─ Devices (device list/discovery)            │
│  ├─ Files (file browser)                       │
│  ├─ Terminal (remote shell)                    │
│  ├─ Metrics (system graphs)                    │
│  ├─ Scripts (batch execution)                  │
│  ├─ Clipboard (sync manager)                   │
│  ├─ VNC (desktop viewer)                       │
│  └─ Settings (configuration)                   │
└─────────────────────────────────────────────────┘
         │            │             │
         ▼            ▼             ▼
    ┌──────────┬──────────┬──────────┐
    │ P2P      │ Nostr    │ Clipboard│
    │ Manager  │ Manager  │ Manager  │
    └──────────┴──────────┴──────────┘
         │             │
         └─────┬───────┘
               ▼
    ┌──────────────────────┐
    │ libp2p + Tor Stack   │
    └──────────────────────┘
```

### Stub Component

```
┌─────────────────────────────────────────────────┐
│         Stub Agent (Headless)                   │
├─────────────────────────────────────────────────┤
│ Request Handler (TCP Listener)                  │
├─────────────────────────────────────────────────┤
│ Feature Handlers:                               │
│  ├─ File Operations Manager                    │
│  ├─ Script Executor                            │
│  ├─ Metrics Collector                          │
│  ├─ Shell Allocator (PTY)                      │
│  ├─ VNC Server                                 │
│  ├─ Clipboard Manager                          │
│  └─ Keylogger (toggleable)                     │
├─────────────────────────────────────────────────┤
│ Persistence Layer:                              │
│  ├─ P2P Manager                                │
│  ├─ Nostr Manager                              │
│  ├─ Tor Manager                                │
│  └─ Session Logger                             │
└─────────────────────────────────────────────────┘
         │              │              │
         ▼              ▼              ▼
    ┌──────────┬───────────┬──────────┐
    │ System   │ Nostr     │ Tor      │
    │ Calls    │ Network   │ Network  │
    └──────────┴───────────┴──────────┘
```

## Communication Protocols

### Protocol Stack (OSI Model)

```
Layer 7: Application Protocols
├─ VNC (remote desktop)
├─ SSH-like shell
├─ File transfer (custom)
├─ Metrics streaming
├─ Script execution
├─ Clipboard sync
└─ Nostr DMs

Layer 6: Serialization
├─ JSON
└─ Binary frames

Layer 5: Session/Stream Management
├─ libp2p Streams (multiplexed)
└─ Session IDs

Layer 4: Transport
├─ libp2p multiplex protocols
└─ TCP/Tor

Layer 3: Network
├─ libp2p DHT
├─ Peer discovery
└─ libp2p Routing

Layer 2: Anonymization
├─ Tor circuits
├─ Onion services (bine)
└─ Exit nodes

Layer 1: Physical
├─ Tor network
└─ Internet
```

### Message Types

#### P2P Streams (libp2p)

Each feature uses dedicated multiplexed streams:

1. **File Operations Stream** (`/remotefiles/1.0.0`)
   - Request: `{type: "list"|"download"|"upload"|"delete", path: "..."}` 
   - Response: `{type: "...", data: [...], error: null}`

2. **Shell Stream** (`/shell/1.0.0`)
   - Bidirectional: stdin/stdout/stderr multiplexing
   - Control frames for signal handling (Ctrl+C, etc.)

3. **Metrics Stream** (`/metrics/1.0.0`)
   - Unidirectional push from stub to panel
   - Periodic snapshots with timestamps

4. **Script Stream** (`/scripts/1.0.0`)
   - Execute command with args
   - Stream stdout/stderr in real-time
   - Session tracking and logging

5. **VNC Stream** (`/vnc/1.0.0`)
   - VNC protocol frames
   - Compression support

6. **Clipboard Stream** (`/clipboard/1.0.0`)
   - Text content sync
   - Bidirectional

#### Nostr Events (Kind 4 - DMs)

Used for out-of-band signaling:

1. **Status Update**
   ```json
   {
     "event": "status",
     "stub_id": "device1",
     "online": true,
     "onion": "abc123...onion",
     "libp2p_addr": "/ip4/.../p2p/..."
   }
   ```

2. **Anomaly Alert**
   ```json
   {
     "event": "anomaly",
     "type": "temperature|memory|cpu",
     "value": 85.5,
     "threshold": 80.0,
     "severity": "critical|warning|info"
   }
   ```

3. **Configuration Update**
   ```json
   {
     "event": "config",
     "key": "polling_interval",
     "value": 5
   }
   ```

## Data Flow Examples

### File Download Flow

```
Panel                          Network                        Stub
  │                              │                              │
  ├─ Create file request ────────┼──────────────────────────────┤
  │  (path: /etc/config)         │                              │
  │                              │  ┌────────────────────────┐  │
  │                              │  │ Validate path          │  │
  │                              │  │ Check permissions      │  │
  │                              │  │ Open file              │  │
  │                              │  │ Stream chunks          │  │
  │  ◄────────────────────────────┼──────────────────────────┤  │
  │  (file chunks)               │                              │
  │                              │                              │
  ├─ Verify hash                 │                              │
  │                              │                              │
  └─ Save locally                │                              │
```

### Script Execution Flow

```
Panel                          Network                        Stub
  │                              │                              │
  ├─ Submit script ──────────────┼──────────────────────────────┤
  │  (content, args, timeout)    │                              │
  │                              │  ┌────────────────────────┐  │
  │                              │  │ Validate script        │  │
  │                              │  │ Create session ID      │  │
  │  ◄────────────────────────────┼──────────────────────────┤  │
  │  (session_id)                │                              │
  │                              │                              │
  ├─ Open output stream ─────────┼──────────────────────────────┤
  │                              │  ┌────────────────────────┐  │
  │                              │  │ Execute in PTY/shell   │  │
  │                              │  │ Stream stdout/stderr   │  │
  │  ◄────────────────────────────┼──────────────────────────┤  │
  │  (output stream)             │                              │
  │                              │                              │
  ├─ Display output              │                              │
  │  (real-time)                 │                              │
  │                              │  ┌────────────────────────┐  │
  │                              │  │ Process exits          │  │
  │                              │  │ Send exit code         │  │
  │  ◄────────────────────────────┼──────────────────────────┤  │
  │  (exit_code)                 │                              │
  │                              │                              │
  └─ Log session                 │                              │
```

### Metrics Collection Flow

```
Stub                           Network                        Panel
  │                              │                              │
  ├─ Poll system metrics ───────┐│                              │
  │  (CPU, MEM, TEMP, UP)        ││                              │
  │                              ││ ┌──────────────────────────┐
  │  ├─ Stream to panel ────────┬┼─┼─────────────────────────┤ │
  │  │  (timestamp, values)    │││ │ Display in dashboard    │ │
  │  │                         │││ │ Update graph widgets    │ │
  │  │                         │││ │ Check for anomalies     │ │
  │  │                         │││ └──────────────────────────┘
  │  │                         │││
  │  ├─ Check thresholds ──────┐│                              │
  │  │  (local anomaly detect)  ││                              │
  │  │                         ││ ┌──────────────────────────┐
  │  ├─ Alert if anomaly ──────┬┼─┼─ Nostr DM ──────────────┤ │
  │  │  (publish to Nostr)     │││ │ (alert notification)   │ │
  │  │                         │││ └──────────────────────────┘
  │  │                         │││
  │  └─ Wait for interval ─────┘││
  │     (5 seconds default)      ││
  │                              ││
  └──────────────────────────────┘│
                                  │
```

## Security Architecture

### Encryption Layers

```
Application Data
    │
    ▼ (libp2p Noise Protocol)
Encrypted P2P Stream
    │
    ▼ (Tor encryption)
Tor Circuit
    │
    ▼
Internet
```

### Key Management

```
Onion Keys (Tor v3)
  ├─ Generated by oniongen-go (builder phase)
  ├─ Stored in ~/.remote-mgmt/keys/
  ├─ Embedded in stub binary
  └─ Never regenerated (persistent)

libp2p Keys
  ├─ Auto-generated on first run
  ├─ Stored locally
  ├─ Used for peer identification
  └─ Part of libp2p identity

Nostr Keys (for signaling)
  ├─ Auto-generated or user-provided
  ├─ Used for DM encryption
  └─ Stored in config
```

### Access Control

```
File Operations
  ├─ Path validation (sandbox to base dir)
  ├─ Permission checks
  └─ ACL (future expansion)

Script Execution
  ├─ Script validation
  ├─ Timeout enforcement
  ├─ Resource limits
  └─ Sandboxing (OS-dependent)

Keylogger
  ├─ Explicit enable/disable command
  ├─ Requires confirmation
  ├─ Audit logging
  └─ Ethical gate (disclaimer)
```

## State Management

### Panel State
- Connected device list
- Stream handles per device
- UI component state
- Settings/configuration
- Anomaly history

### Stub State
- Peer connectivity status
- Open file handles
- Running scripts/sessions
- Metrics cache
- Configuration overrides

### Persistence

Panel:
- Device list: Memory + config file
- Credentials: Encrypted local storage
- History: SQLite or JSON files

Stub:
- Peer connections: Maintained by libp2p
- Onion keys: Loaded from binary
- Session logs: File system
- Configuration: File system

## Scalability Considerations

### Horizontal Scaling
- **Multiple Panels**: Independent instances can control different stub groups
- **Many Stubs**: Panel handles 100+ devices with list pagination
- **Load Distribution**: Stubs publish status via Nostr for discovery

### Vertical Scaling
- **Panel**: CPU/GPU for graph rendering, UI responsiveness
- **Stub**: Resource constraints on ARM (minimal overhead)

### Network Scaling
- **Tor Bandwidth**: Not optimized for bulk transfers
- **DHT Propagation**: Time to discover new peers (~30s typical)
- **Relay Redundancy**: Multiple Nostr relays for reliability

## Failure Scenarios

### Network Disconnection
- Stub automatically attempts reconnection
- Exponential backoff on retry
- Status updates via Nostr when online
- Panel shows offline indicator

### Device Offline
- Panel detects via timeout on streams
- Updates device status indicator
- Queues commands for later delivery
- Alerts user of anomalies if available

### Tor Unavailable
- System falls back to libp2p direct P2P
- May reduce anonymity but maintains functionality
- Warning in UI/logs

### Nostr Relay Down
- Uses multiple relays (fallback mechanism)
- Degrades to libp2p DHT discovery
- Delayed status updates

## Performance Profile

### Panel
- CPU: <5% at rest, <20% with active interactions
- Memory: 100-300MB typical
- Network: Event-driven, low baseline bandwidth
- Disk I/O: Minimal, logs only

### Stub (Desktop)
- CPU: <2% at rest, <15% under load
- Memory: 50-100MB typical
- Network: 5-50KB/s depending on features
- Disk I/O: Script execution + logging

### Stub (Raspberry Pi)
- CPU: <1% at rest, <10% under load
- Memory: 30-50MB typical
- Network: 2-20KB/s
- Thermal: <60°C typical under load

## Future Enhancements

1. **WebRTC Transport**: Direct device-to-device without Tor for LAN
2. **QUIC Protocol**: Faster, more efficient than TCP
3. **Clustering**: Multi-panel coordination
4. **Kubernetes Integration**: Manage containerized workloads
5. **Mobile App**: iOS/Android companion app
6. **Hardware Acceleration**: GPU computing offload
7. **Database Backend**: Centralized logging (optional)
8. **API Gateway**: RESTful interface for integrations

---

For implementation details, see individual package documentation in `pkg/*/`.
