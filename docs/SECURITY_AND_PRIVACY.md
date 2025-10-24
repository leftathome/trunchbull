# Security and Privacy Guidelines

## IMPORTANT: READ THIS FIRST

**Trunchbull** is a self-hosted application designed to aggregate and display your children's academic data from their school's learning platforms. This application handles **sensitive educational records protected by FERPA** (Family Educational Rights and Privacy Act) and potentially other privacy regulations.

By using this software, **you accept full responsibility** for the security, privacy, and proper handling of all student data accessed, stored, and displayed by this application.

---

## Overview

This document outlines:
1. What data is collected and how it's used
2. Your responsibilities as a user
3. Security recommendations
4. Privacy considerations
5. Compliance with relevant regulations
6. Incident response procedures

---

## Data We Handle

### Types of Student Data

This application accesses and stores the following types of data:

#### Personal Identifiable Information (PII)
- Student names
- Student ID numbers
- Grade levels
- Course enrollment information

#### Educational Records (FERPA-Protected)
- Assignment details and due dates
- Grades and scores
- GPA calculations
- Attendance records (if available via API)
- Assessment scores
- Teacher communications
- Report cards and progress reports

#### Platform Authentication Data
- OAuth tokens for Schoology
- OAuth tokens for PowerSchool/The Source
- API credentials
- Session information

### Data NOT Collected
- This application does NOT send any student data to external servers (except the school's official platforms)
- No analytics or telemetry
- No third-party tracking
- No cloud storage
- No sharing with other users or families

---

## Your Responsibilities

### As a User of This Software, You Are Responsible For:

#### 1. Legal Compliance
- Ensuring you have the legal right to access the student data
- Complying with FERPA and other applicable privacy laws
- Following your school district's acceptable use policies
- Obtaining proper authorization to use API access (if required)

#### 2. Secure Infrastructure
- Running this application on infrastructure YOU control
- Implementing appropriate network security measures
- Restricting physical and network access to the system
- Keeping the host system patched and updated
- Using strong passwords and authentication

#### 3. Credential Management
- Protecting your school platform login credentials
- Securing API keys and OAuth tokens
- Not sharing credentials with unauthorized persons
- Rotating credentials if compromise is suspected

#### 4. Data Protection
- Implementing backups of the application and data
- Encrypting the database file (strongly recommended)
- Securing the network where the application runs
- Properly disposing of data when no longer needed
- Monitoring for unauthorized access

#### 5. Access Control
- Limiting who can access the dashboard
- Not leaving the dashboard open on shared devices
- Logging out when finished
- Implementing network-level access controls (firewall, VPN)

---

## Security Recommendations

### Deployment Best Practices

#### Network Security

**DO:**
- Run behind a firewall
- Use a VPN if accessing remotely
- Consider running on an isolated VLAN
- Use HTTPS with a valid certificate (self-signed acceptable for local use)
- Restrict access to localhost or specific IP ranges

**DON'T:**
- Expose to the public internet without strong authentication
- Use default passwords or credentials
- Allow unauthenticated access
- Run on unsecured WiFi networks

#### Host Security

**DO:**
- Keep Docker and the host OS updated
- Use a minimal Linux distribution (Alpine, Ubuntu minimal)
- Enable automatic security updates
- Configure host firewall (iptables, ufw)
- Review Docker security best practices
- Run containers with minimal privileges (non-root user)
- Use Docker secrets for sensitive configuration

**DON'T:**
- Run as root inside containers unnecessarily
- Disable SELinux or AppArmor
- Ignore security updates
- Use outdated base images

#### Application Security

**DO:**
- Change default configuration values
- Use strong, unique passwords for any authentication
- Enable database encryption (see Configuration section)
- Rotate OAuth tokens periodically
- Review application logs regularly
- Set appropriate file permissions on data directory (700 or 750)
- Use environment variables for secrets (not hardcoded)

**DON'T:**
- Store credentials in configuration files (use .env)
- Commit secrets to version control
- Share OAuth tokens
- Disable HTTPS in production

#### Data Security

**DO:**
- Encrypt the SQLite database file
- Back up data regularly
- Store backups securely (encrypted)
- Delete old data you no longer need
- Use full disk encryption on the host
- Implement a data retention policy

**DON'T:**
- Store unencrypted backups on cloud services
- Keep data indefinitely without review
- Use world-readable file permissions
- Store data on removable media without encryption

---

## Privacy Considerations

### Data Minimization

This application should only collect data that is necessary for its functionality. Consider:
- Do you need to sync all historical data or just current assignments?
- How long should old assignments be retained?
- Is message history necessary or just unread messages?

**Recommendation**: Configure data retention policies to automatically delete data older than 90 days unless you have a specific need.

### Access Logging

The application logs all API requests and data access. This is for:
- Troubleshooting technical issues
- Detecting unauthorized access
- Understanding sync patterns

**Important**: Log files contain sensitive data. Protect them with the same security measures as the database.

### Third-Party Access

**NEVER:**
- Share your dashboard with anyone outside your immediate family
- Provide API access to third parties
- Use this application for commercial purposes
- Aggregate data from multiple families
- Share data with other parents or students

### School District Policies

Before using this application:
1. Review your school district's acceptable use policy
2. Verify that personal API access is permitted
3. Understand any restrictions on automated access
4. Contact the district IT department if uncertain

**Note**: Some school districts prohibit or restrict third-party applications accessing student data. Seattle Public Schools has previously taken action against unauthorized apps.

---

## Compliance

### FERPA (Family Educational Rights and Privacy Act)

FERPA is a federal law that protects the privacy of student education records. As a parent or guardian, you have rights under FERPA, but you also have responsibilities.

#### Your FERPA Rights
- Access your child's educational records
- Request corrections to inaccurate records
- Control disclosure of records to third parties

#### Your FERPA Responsibilities
- Protect the confidentiality of your child's records
- Not disclose records without proper authorization
- Follow school policies regarding record access

#### How This Application Affects FERPA
- You are accessing records you already have the right to access
- The application is a tool YOU control (not a third party)
- You must ensure no unauthorized access occurs
- You are responsible for protecting the data

### State Privacy Laws

Depending on your location, additional privacy laws may apply:
- California: SOPIPA (Student Online Personal Information Protection Act)
- Other states may have similar laws

**Recommendation**: Familiarize yourself with your state's student data privacy laws.

### School District Policies

Your school district may have additional policies regarding:
- API access by parents
- Automated data collection
- Third-party applications
- Data security requirements

**Action Required**: Review and comply with all applicable district policies.

---

## Data Retention and Deletion

### Retention Policy

**Default Behavior:**
- Assignment data: Retained until 90 days after due date
- Grades: Retained for current and previous grading period
- Messages: Retained until read, then 30 days
- Events: Retained until event date passes, then 30 days
- Audit logs: Retained for 90 days

**Customization:**
You can adjust retention periods in the configuration file.

### Data Deletion

**When to Delete Data:**
- End of school year
- Student graduates or changes schools
- You stop using the application
- Data is no longer needed

**How to Delete Data:**
```bash
# Option 1: Delete specific student data
docker exec trunchbull ./trunchbull delete-student <student-id>

# Option 2: Delete all data (factory reset)
docker-compose down -v
rm -rf ./data

# Option 3: Securely wipe database
shred -vfz -n 3 ./data/trunchbull.db
```

### Decommissioning

When you're done with the application:
1. Delete all student data
2. Revoke OAuth tokens on school platforms
3. Delete Docker images and containers
4. Securely wipe the database file
5. Delete any backups

---

## Incident Response

### If You Suspect a Security Breach

#### Immediate Actions
1. **Shut down the application**
   ```bash
   docker-compose down
   ```

2. **Disconnect from network**
   ```bash
   # Disable network interface
   sudo ip link set <interface> down
   ```

3. **Document what happened**
   - When did you notice the issue?
   - What unusual activity did you observe?
   - Who might have had access?

4. **Assess the scope**
   - What data was potentially accessed?
   - How many students affected?
   - Was data exfiltrated?

#### Next Steps

1. **Change all credentials**
   - School platform passwords
   - API credentials
   - Host system passwords

2. **Revoke OAuth tokens**
   - Log into Schoology and revoke access
   - Log into PowerSchool and revoke access

3. **Review logs**
   ```bash
   # Check application logs
   docker logs trunchbull

   # Check system logs
   sudo journalctl -u docker

   # Check network connections
   sudo ss -tulpn
   ```

4. **Notify affected parties** (if required)
   - School district (if credentials compromised)
   - Other family members
   - Depending on severity, may need to notify school

5. **Restore from clean backup** (if available)
   - Verify backup integrity
   - Restore to known-good state
   - Apply security updates

6. **Implement additional security measures**
   - Add firewall rules
   - Enable additional monitoring
   - Review access controls

### If You Suspect Data Disclosure

If student data was accessed by unauthorized parties:
1. Document the disclosure
2. Determine what data was exposed
3. Notify the school district
4. Take steps to prevent further disclosure
5. Consider legal and regulatory requirements

---

## Vulnerability Reporting

### Found a Security Issue in Trunchbull?

This is an open-source project. If you discover a security vulnerability:

1. **DO NOT open a public GitHub issue**
2. **DO NOT disclose publicly before a fix is available**
3. **DO** email the maintainers privately (contact info in README)
4. **DO** provide details:
   - Description of vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will:
- Acknowledge receipt within 48 hours
- Investigate and develop a fix
- Release a security update
- Credit you for the discovery (if desired)

---

## Configuration for Enhanced Security

### Enable Database Encryption

#### Option 1: SQLite Encryption Extension (SEE)
```bash
# Commercial extension from SQLite
# Requires license: https://www.hwaci.com/sw/sqlite/see.html
```

#### Option 2: SQLCipher (Open Source)
```bash
# Use SQLCipher instead of standard SQLite
# Set encryption key in environment variable
DATABASE_ENCRYPTION_KEY=your-very-strong-key-here
```

#### Option 3: Full Disk Encryption
```bash
# Use LUKS on Linux
sudo cryptsetup luksFormat /dev/sdX
sudo cryptsetup open /dev/sdX encrypted_disk
```

### Use HTTPS with Self-Signed Certificate

```bash
# Generate self-signed certificate
openssl req -x509 -newkey rsa:4096 -nodes \
  -keyout key.pem -out cert.pem -days 365 \
  -subj "/CN=trunchbull.local"

# Configure application to use HTTPS
docker-compose.yml:
  environment:
    - TLS_CERT=/config/cert.pem
    - TLS_KEY=/config/key.pem
```

### Implement Network Access Control

```bash
# Restrict to localhost only
docker-compose.yml:
  ports:
    - "127.0.0.1:8080:8080"

# Or restrict to specific subnet
iptables -A INPUT -p tcp --dport 8080 -s 192.168.1.0/24 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j DROP
```

### Add Basic Authentication (Future Feature)

```yaml
# config.yaml
auth:
  enabled: true
  username: admin
  password_hash: $2a$10$... # bcrypt hash
```

---

## Security Checklist

Before deploying Trunchbull, ensure you have:

### Infrastructure
- [ ] Running on dedicated or isolated system
- [ ] Host OS is up-to-date
- [ ] Docker is up-to-date
- [ ] Firewall is configured
- [ ] Network is secured (not public WiFi)
- [ ] Physical access is restricted

### Application
- [ ] Changed all default configurations
- [ ] Secrets stored in environment variables (not config files)
- [ ] Database encryption enabled (or full disk encryption)
- [ ] HTTPS configured
- [ ] Access restricted to authorized users only
- [ ] Logging enabled

### Operations
- [ ] Backup strategy implemented
- [ ] Backup testing performed
- [ ] Data retention policy defined
- [ ] Incident response plan documented
- [ ] Regular security reviews scheduled
- [ ] Update process established

### Compliance
- [ ] Read and understood FERPA
- [ ] Reviewed school district policies
- [ ] Have authorization for API access
- [ ] Understand your responsibilities
- [ ] Know how to report incidents

---

## Additional Resources

### Privacy Laws
- FERPA: https://www2.ed.gov/policy/gen/guid/fpco/ferpa/index.html
- COPPA: https://www.ftc.gov/enforcement/rules/rulemaking-regulatory-reform-proceedings/childrens-online-privacy-protection-rule
- State Laws: Contact your state education department

### Security Best Practices
- OWASP Top 10: https://owasp.org/www-project-top-ten/
- Docker Security: https://docs.docker.com/engine/security/
- OAuth Security: https://oauth.net/2/security-best-practices/

### School District Resources
- Seattle Public Schools IT: Contact via district website
- PowerSchool Support: https://help.powerschool.com/
- Schoology Support: https://help.schoology.com/

---

## Disclaimer

**THIS SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND.**

The authors and contributors:
- Make no guarantees about security or privacy
- Are not responsible for data breaches or unauthorized access
- Are not responsible for violations of school policies
- Are not responsible for legal or regulatory violations
- Provide this software for personal, non-commercial use only

**You use this software at your own risk.**

By using Trunchbull, you acknowledge that:
1. You have read and understood this document
2. You accept full responsibility for data security and privacy
3. You will comply with all applicable laws and policies
4. You will implement appropriate security measures
5. You understand the risks involved

---

## Questions?

If you have questions about security and privacy:
1. Review the documentation thoroughly
2. Consult with your school district IT department
3. Seek legal advice if needed (we are not lawyers)
4. Open a GitHub discussion (for general questions, not sensitive info)

**Never share sensitive information in public forums.**

---

**Document Version**: 1.0
**Last Updated**: 2025-10-23
**Reviewed By**: Project Maintainers

---

## Acknowledgments

We take security and privacy seriously. This document was created with input from:
- FERPA guidelines
- OWASP security best practices
- Common Sense Media's privacy recommendations
- Real-world incidents and lessons learned

If you have suggestions for improving this document, please submit a pull request or open an issue.
