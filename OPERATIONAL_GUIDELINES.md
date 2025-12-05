# Operational Guidelines and Ethical Safeguards

## Introduction

This document outlines the ethical framework, operational guidelines, and mandatory safeguards for operating the Go-based Remote Management System. This toolkit is designed exclusively for authorized, legitimate, and educational use.

## Ethical Principles

### Core Values
1. **Authorization**: Never access systems without explicit, written authorization
2. **Transparency**: Maintain clear audit logs of all operations
3. **Privacy**: Respect device owner privacy and data confidentiality
4. **Accountability**: Take responsibility for all actions performed
5. **Compliance**: Adhere to all applicable laws and regulations
6. **Disclosure**: Inform users of data collection and monitoring

### Educational Purpose
This system is designed for:
- **System Administration**: Legitimate remote device management
- **Security Research**: Authorized penetration testing and security assessments
- **IoT Management**: Managing Raspberry Pi and similar device fleets
- **Educational Demonstrations**: Teaching distributed systems and security concepts
- **Disaster Recovery**: Remote system recovery and maintenance

## Operational Security

### Access Control Requirements

**Before operating the stub agent on any device:**

1. ✓ Obtain explicit written authorization from system owner
2. ✓ Document the authorization (date, scope, duration)
3. ✓ Inform all stakeholders of monitoring/management activities
4. ✓ Establish communication channels for reporting issues
5. ✓ Define scope of allowed operations

### Mandatory Splash Screen

**The panel displays a 20-second splash screen on startup** with the message:
```
"THIS IS ONLY FOR EDUCATIONAL AND LEGITIMATE PURPOSES"
```

This serves as:
- A legal reminder of authorized use only
- An acknowledgment of ethical responsibilities
- A pause point for reconsidering actions
- Documentation of intentional operation

**Users cannot bypass or skip this screen.** If you see this message and cannot confirm authorized use, **stop immediately and uninstall the software.**

### Credential Management

#### Tor Onion Keys
- **Storage**: Embedded in binary at build time via oniongen-go
- **Persistence**: Never regenerated (ensures consistent .onion addresses)
- **Protection**: Store in secure location with restricted permissions
- **Backup**: Keep secure backups of onion keys
- **Compromise**: If compromised, rebuild stub with new keys

#### libp2p Keys
- **Generation**: Auto-generated on first stub run
- **Storage**: `~/.remote-mgmt/identity/`
- **Permissions**: `chmod 600` (user read/write only)
- **Rotation**: Rotate keys periodically (monthly recommended)
- **Revocation**: No built-in revocation; depends on trust model

#### Nostr Keys
- **Generation**: Can be user-provided or auto-generated
- **Storage**: `~/.remote-mgmt/nostr/`
- **Protection**: Consider hardware wallet for critical deployments
- **Best Practice**: Use dedicated keys per panel instance

### Network Security

#### Tor Configuration
```bash
# Verify Tor is running
ps aux | grep tor

# Check port accessibility
netstat -an | grep 9050  # SOCKS
netstat -an | grep 9051  # Control

# Monitor Tor logs
journalctl -u tor -f
```

#### Firewall Rules
```bash
# Allow only necessary inbound connections
# (Tor handles P2P traffic)

# Monitor suspicious connections
netstat -an | grep ESTABLISHED
ss -tulpn | grep -i listen
```

#### VPN Considerations
- Use VPN for additional privacy layer if needed
- Tor provides anonymity; VPN adds redundancy
- Disable DNS leaks to prevent IP exposure

## Operational Procedures

### Initial Setup Checklist

- [ ] Verify Go 1.21+ is installed
- [ ] Clone repository from secure source
- [ ] Verify repository signature (if available)
- [ ] Review source code for understanding
- [ ] Obtain authorization documentation
- [ ] Configure Nostr relay addresses
- [ ] Set up secure key storage
- [ ] Test in isolated lab environment first
- [ ] Document deployment plan
- [ ] Brief all stakeholders

### Deployment Procedure

1. **Build Phase**
   ```bash
   # Generate custom onion patterns if desired
   go run ./cmd/builder -all -vanity "^authorized" -version "1.0.0"
   
   # Verify binaries
   ls -lah ./bin/
   ```

2. **Pre-Deployment Testing**
   - Test on representative hardware
   - Verify Tor connectivity
   - Confirm libp2p DHT integration
   - Test all features locally
   - Monitor resource usage

3. **Deployment**
   - Document deployment time/date
   - Verify stub starts correctly
   - Check initial Nostr status publication
   - Monitor for errors in logs

4. **Post-Deployment**
   - Verify connectivity in panel
   - Test each feature independently
   - Monitor metrics for baseline
   - Document any issues

### Monitoring Operations

#### Daily Checks
- [ ] All stubs showing online status
- [ ] No anomalies in metrics
- [ ] No errors in logs
- [ ] Confirm backup systems operational

#### Weekly Reviews
- [ ] Anomaly trending analysis
- [ ] Keylogger usage audit (if enabled)
- [ ] Configuration drift detection
- [ ] Performance trending

#### Monthly Tasks
- [ ] Update stub binaries to latest version
- [ ] Review and rotate keys if needed
- [ ] Audit all operations performed
- [ ] Verify authorization still valid
- [ ] Test disaster recovery procedures

### Incident Response

#### If Unauthorized Access Suspected
1. **IMMEDIATELY**
   - Isolate affected device(s)
   - Disable Tor (if possible)
   - Preserve logs for forensics
   - Alert system owners

2. **WITHIN 1 HOUR**
   - Secure all credentials
   - Rotate Nostr/libp2p keys
   - Rebuild stubs with new onion keys
   - Begin investigation

3. **DOCUMENTATION**
   - Record timeline of events
   - Collect all relevant logs
   - Document findings
   - Report to stakeholders

#### If Data Breach Occurs
- Notify affected parties immediately
- Cooperate with law enforcement if required
- Conduct post-mortem analysis
- Implement preventative measures
- Document lessons learned

## Audit and Logging

### What Gets Logged

#### Panel Logs (`~/.remote-mgmt/panel.log`)
```json
{
  "timestamp": "2024-01-15T10:30:45Z",
  "event": "connect",
  "stub_id": "device1",
  "user": "admin",
  "source_ip": "192.168.1.100"
}
```

#### Stub Logs (`/var/log/remote-mgmt/`)
```json
{
  "timestamp": "2024-01-15T10:30:45Z",
  "event": "file_download",
  "path": "/etc/config",
  "size": 1024,
  "user": "192.168.1.100",
  "status": "success"
}
```

#### Session Logs (per operation)
- Start/end timestamps
- User/peer identification
- Operations performed
- Data transferred
- Exit status/errors

### Log Retention Policy

- **Active Operations**: Full detail for 30 days
- **Archive**: Compressed after 30 days
- **Long-Term**: Keep for 1 year minimum
- **Compliance**: Adjust for regulatory requirements
- **Rotation**: Implement log rotation to prevent disk fill

### Audit Trail Analysis

```bash
# View recent operations
grep -r "file_download" /var/log/remote-mgmt/ | tail -20

# Find all keylogger activity
grep "keylogger" /var/log/remote-mgmt/* | jq '.event'

# Identify anomalies
grep "status.*error" /var/log/remote-mgmt/*

# Generate usage report
awk -F'"' '/file_download|script_execute/ {print $4}' /var/log/remote-mgmt/* | sort | uniq -c
```

## Feature-Specific Guidelines

### File Operations

#### Safe Practices
- ✓ Define allowed directories (avoid sensitive areas)
- ✓ Verify file permissions before downloading
- ✓ Check file sizes before transfer
- ✓ Scan downloaded files for malware
- ✓ Log all file operations

#### Prohibited Actions
- ✗ Accessing files outside allowed paths
- ✗ Modifying critical system files
- ✗ Transferring without appropriate permissions
- ✗ Accessing others' private files

### Script Execution

#### Safe Practices
- ✓ Review scripts before execution
- ✓ Use timeout limits (prevent hanging)
- ✓ Run with minimal necessary privileges
- ✓ Avoid untrusted script sources
- ✓ Enable execution logging

#### Prohibited Actions
- ✗ Executing unvetted/untrusted scripts
- ✗ Running with root/admin unnecessarily
- ✗ Disabling security controls
- ✗ Accessing unauthorized resources

### Keylogging

#### Strict Requirements
- ✓ **Explicit written authorization required**
- ✓ User must be informed of keylogging
- ✓ Use only for legitimate security purposes
- ✓ Enable only when actively monitoring
- ✓ Complete audit trail maintained
- ✓ Automatic disable after session/timeout

#### Prohibited Use Cases
- ✗ Unauthorized monitoring
- ✗ Password theft
- ✗ Spying on personal communications
- ✗ Circumventing user privacy

#### Compliance
- Comply with laws on wiretapping/interception
- Respect employee privacy rights
- Consider union agreements and works councils
- Document informed consent
- Provide transparency to affected parties

### Remote Shell Access

#### Safe Practices
- ✓ Log all commands executed
- ✓ Restrict to necessary users
- ✓ Use read-only mode when possible
- ✓ Monitor for suspicious commands
- ✓ Timeout inactive sessions

#### Prohibited Actions
- ✗ Unauthorized command execution
- ✗ Modifying audit logs
- ✗ Elevating privileges inappropriately
- ✗ Accessing other users' sessions

## Compliance Considerations

### Regulatory Requirements

Depending on jurisdiction and context:

#### GDPR (EU)
- Data minimization: Collect only necessary information
- Data subject rights: Enable user access to their data
- Right to erasure: Delete collected data upon request
- Data Protection Impact Assessment (DPIA): Required for monitoring

#### HIPAA (Healthcare - US)
- Business Associate Agreement required
- Encryption mandatory
- Access controls and audit logging required
- Incident reporting obligations

#### SOX (Financial - US)
- Complete audit trail required
- Segregation of duties
- Change management processes
- Regular access reviews

#### CCPA (California)
- Transparency to California residents
- Opt-out rights for data collection
- Breach notification requirements

### Legal Review

Before operational deployment:
1. **Consult Legal**: Have legal review authorization
2. **Compliance Check**: Verify regulatory adherence
3. **Privacy Policy**: Update if needed
4. **User Agreements**: Ensure users informed
5. **Incident Plan**: Have response procedures ready

## Security Best Practices

### Defense in Depth

```
Layer 1: Prevention
├─ Authorization requirement
├─ Authentication (keys)
└─ Encryption (Noise + Tor)

Layer 2: Detection
├─ Audit logging
├─ Anomaly detection
└─ Alerts (Nostr DMs)

Layer 3: Response
├─ Incident procedures
├─ Key rotation capability
└─ Emergency shutdown

Layer 4: Recovery
├─ Backup procedures
├─ Restore processes
└─ Post-incident analysis
```

### Update Management

```bash
# Check for updates
go get -u ./...

# Pin versions (production)
go mod edit -require=module@version

# Test updates in staging
# Deploy with version tracking
```

### Vulnerability Management

1. Monitor security advisories:
   - Go security bulletins
   - Dependency advisories
   - Known CVEs

2. Patch procedure:
   - Evaluate impact
   - Test in staging
   - Deploy with notification
   - Verify in production

## Training and Awareness

### Required Training
- [ ] Ethical use principles
- [ ] Authorization documentation
- [ ] Operational procedures
- [ ] Security controls
- [ ] Incident response
- [ ] Audit and compliance

### Annual Recertification
- Confirm authorization status
- Review policy changes
- Acknowledge ethical guidelines
- Document in personnel files

### Security Awareness
- Regular briefings on threats
- Social engineering exercises (controlled)
- Phishing simulations
- Security culture emphasis

## Incident Checklist

### Upon Discovering Unauthorized Access
- [ ] Preserve all evidence
- [ ] Immediately isolate system
- [ ] Notify stakeholders
- [ ] Engage legal/compliance
- [ ] Begin forensics
- [ ] Contact law enforcement (if required)
- [ ] Document timeline
- [ ] Implement fixes
- [ ] Post-mortem analysis

### Communication Template
```
INCIDENT REPORT: [Date/Time]

What Happened:
[Describe incident clearly]

When Discovered:
[Date/time of discovery]

Scope:
[Systems/data affected]

Actions Taken:
[Immediate response]

Next Steps:
[Remediation plan]

Contact:
[Responsible party/contact info]
```

## Decommissioning

When discontinuing use:

1. **Revoke Authorization**
   - Notify all stakeholders
   - Document end date
   - Collect sign-offs

2. **Data Cleanup**
   - Delete logs after retention period
   - Securely wipe keys
   - Remove binaries
   - Audit-log the decommissioning

3. **Final Audit**
   - Confirm all stubs offline
   - Verify all data destroyed
   - Archive compliance records
   - Document decommissioning

## Support and Escalation

### Issues/Questions
- Email: [support contact]
- Documentation: https://github.com/remotemgmt/...
- Security Issues: security@[domain]

### Escalation Path
1. First: Review documentation
2. Second: Contact support team
3. Third: Escalate to management
4. Fourth: Legal/compliance review if needed

---

**Remember**: This system grants significant access and control. Use with the utmost responsibility and integrity. When in doubt, ask for authorization before proceeding.

Last Updated: [Date]
Version: 1.0
