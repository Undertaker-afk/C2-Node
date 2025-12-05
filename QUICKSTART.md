# Quick Start Guide

## 5-Minute Setup

### Step 1: Prerequisites
```bash
# Install Go 1.21+
# Verify installation
go version
# Should output: go version go1.21.x

# For Linux, ensure build tools
sudo apt-get update
sudo apt-get install -y build-essential libssl-dev
```

### Step 2: Clone & Prepare
```bash
cd /path/to/project
go mod download
```

### Step 3: Build Binaries
```bash
# Build for your platform
go run ./cmd/builder -all -out ./bin

# Verify builds
ls -lh ./bin/
```

### Step 4: Configure Nostr (Optional)
Default relays work, but you can customize in settings tab.

### Step 5: Run Panel
```bash
# On Windows
./bin/panel-windows-amd64.exe

# On Linux (compiled with GOOS=linux)
./bin/panel-linux-amd64
```

**The splash screen will appear for 20 seconds** - this is normal and required.

### Step 6: Deploy Stub to Device
```bash
# Copy stub to target device
scp ./bin/stub-linux-amd64 user@device:/tmp/

# SSH into device and run
ssh user@device
sudo /tmp/stub-linux-amd64 -port 9000

# For Raspberry Pi
scp ./bin/stub-linux-arm user@pi:/tmp/
ssh user@pi
/tmp/stub-linux-arm -port 9000
```

### Step 7: View in Panel
Once running, your device should appear in the Devices tab.

---

## Common Commands

### Build Specific Components
```bash
# Panel only
go run ./cmd/builder -panel -out ./bin

# Stub only (all platforms)
go run ./cmd/builder -stub -out ./bin

# With vanity onion
go run ./cmd/builder -all -vanity "^mydevice" -out ./bin
```

### Run Stub with Options
```bash
# Custom port
./bin/stub-linux-amd64 -port 9001

# Custom Nostr relays
./bin/stub-linux-amd64 -nostr-relays "wss://relay.damus.io,wss://nos.lol"

# Debug logging
./bin/stub-linux-amd64 -debug

# Custom stub ID
./bin/stub-linux-amd64 -id "living-room-pi"
```

### Install as Service (Linux)

Create `/etc/systemd/system/remote-mgmt-stub.service`:
```ini
[Unit]
Description=Remote Management Stub
After=network.target

[Service]
Type=simple
User=pi
WorkingDirectory=/home/pi/remote-mgmt
ExecStart=/home/pi/remote-mgmt/stub-linux-arm -port 9000
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable remote-mgmt-stub
sudo systemctl start remote-mgmt-stub
sudo systemctl status remote-mgmt-stub
```

### View Logs
```bash
# Panel logs
tail -f ~/.remote-mgmt/panel.log

# Stub logs (if systemd)
sudo journalctl -u remote-mgmt-stub -f

# Or direct
tail -f /var/log/remote-mgmt/stub.log
```

---

## Features Quick Tour

### Dashboard
Shows overview of connected devices and recent anomalies.

### Devices Tab
- Lists all connected stubs
- Shows online/offline status
- Basic metrics per device
- Click to manage individual device

### File Browser
1. Select device from dropdown
2. Navigate directories
3. Right-click to download
4. Drag-drop to upload
5. Double-click to enter folder

### Remote Terminal
1. Select device
2. Type commands
3. Output appears in real-time
4. Ctrl+C to cancel command
5. Type `exit` to close

### System Metrics
- Real-time CPU, memory, temperature graphs
- Scrollable history
- Anomaly highlights
- Export to CSV

### Script Execution
1. Write/paste script
2. Set timeout (0 = default 5min)
3. Click Execute
4. Monitor output in real-time
5. View exit code when complete

### Settings
- Nostr relay URLs
- libp2p bootstrap nodes
- Polling intervals
- Log retention policy
- Export configuration

---

## Troubleshooting

### Stub Won't Start
```bash
# Check port availability
lsof -i :9000

# Try different port
./bin/stub-linux-amd64 -port 9001

# Enable debug logging
./bin/stub-linux-amd64 -debug

# Check for permission issues
ls -l /tmp/
```

### Device Not Appearing in Panel
1. Check stub is running: `ps aux | grep stub`
2. Verify network connectivity: `ping google.com`
3. Check Tor: `ps aux | grep tor`
4. Verify Nostr relay: `curl wss://relay.damus.io` (using websocat)
5. Wait 30s (DHT propagation)
6. Check logs for errors

### High CPU Usage
- Reduce metrics polling interval: `polling_interval: 10` (in settings)
- Disable unneeded features
- Check for runaway scripts
- Monitor with: `top -p $(pgrep -f stub)`

### Connection Drops
- Normal for Tor - it's higher latency
- Check network stability: `ping -c 100 8.8.8.8`
- Verify Tor is running
- Check for firewalls blocking P2P ports
- Increase timeout values

### Out of Disk Space
```bash
# Check usage
df -h

# Clear old logs
rm -rf ~/.remote-mgmt/logs/*

# Or configure shorter retention
# (in settings panel)
```

---

## Next Steps

1. **Read OPERATIONAL_GUIDELINES.md** - Important for authorized use
2. **Review ARCHITECTURE.md** - Understand system design
3. **Check README.md** - Full documentation
4. **Join Community** - GitHub discussions
5. **Report Issues** - GitHub issues with details

---

## Platform-Specific Notes

### Raspberry Pi
```bash
# SSH setup for easy access
ssh-keygen -t ed25519
ssh-copy-id pi@192.168.1.100

# Install and run
scp ./bin/stub-linux-arm pi@raspberrypi:/home/pi/
ssh pi@raspberrypi
./stub-linux-arm -id "kitchen-pi" -port 9000 &
```

### Windows
```powershell
# Run in PowerShell as Administrator
.\bin\panel-windows-amd64.exe

# For stub on Windows
.\bin\stub-windows-amd64.exe -port 9000
```

### macOS (Build Only)
```bash
# macOS not officially supported for stubruntime,
# but can build tools
GOOS=darwin GOARCH=amd64 go build -o panel-mac ./cmd/panel
```

---

## Security Reminders

⚠️ **Before you start:**
- [ ] Do you have authorization to access the target device?
- [ ] Have you informed the device owner of monitoring?
- [ ] Do you understand the ethical implications?
- [ ] Are you compliant with applicable laws?

✓ The splash screen will remind you of these requirements.

---

## Getting Help

### Documentation
- `README.md` - Full reference
- `ARCHITECTURE.md` - Design details
- `OPERATIONAL_GUIDELINES.md` - Ethics & procedures
- Inline code comments - Implementation details

### Debugging
```bash
# Enable verbose logging everywhere
go run ./cmd/stub -debug
go run ./cmd/panel -debug
go run ./cmd/builder -v -all

# Check system resources
free -h          # Memory
df -h            # Disk
top -b -n 1      # Processes
```

### Community Resources
- GitHub Issues: Report bugs & request features
- GitHub Discussions: Ask questions & share ideas
- Email: [support contact from repo]
- Security Issues: [security contact from repo]

---

**Ready to go! 🚀**

Start with a simple test: Run stub on localhost, connect with panel, browse a file. Then expand to actual devices.

Last Updated: 2024
