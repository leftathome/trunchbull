# Trunchbull - Student Dashboard

> A self-hosted dashboard for parents to monitor their children's academic progress across multiple learning platforms.

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)

---

## Overview

Trunchbull aggregates student data from multiple learning platforms and presents it in a unified, easy-to-read dashboard. Perfect for parents who want to stay on top of their children's academic progress without logging into multiple platforms.

### Current Platform Support
- **Schoology** (Learning Management System)
- **PowerSchool / The Source** (Student Information System)

### Features
- Outstanding assignments with due dates
- Current grades by course
- GPA calculation
- School calendar events and days off
- Teacher messages
- Multi-student support
- Automatic background syncing
- Self-hosted and secure

---

## 🎯 Important Update: No API Approval Required!

**Good news!** You DON'T need district-level API approval to use Trunchbull.

Our research discovered that both Schoology and PowerSchool can be accessed using your regular parent/student login credentials, just like the official mobile apps do. This eliminates bureaucratic barriers and makes Trunchbull immediately usable.

**Two simple options:**
1. **Session Tokens** (recommended for MVP): Log in via browser, extract session cookies, provide to Trunchbull
2. **Credential Storage**: Provide username/password (encrypted), Trunchbull handles login automatically

**See [Alternative Authentication Guide](docs/ALTERNATIVE_AUTH.md) for complete details.**

This approach has been proven by multiple successful projects, including "The Source: SPS" app which served 17,000+ users for 3 years.

---

## Screenshots

> Coming soon

---

## Prerequisites

- Docker and Docker Compose
- Parent account on Schoology and PowerSchool/The Source
- ~~API access credentials~~ **NOT NEEDED!** Just use your regular login credentials or session tokens

---

## Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/trunchbull.git
cd trunchbull
```

### 2. Configure Environment Variables
```bash
cp .env.example .env
# Edit .env with your API credentials
```

### 3. Start the Application
```bash
docker-compose up -d
```

### 4. Access the Dashboard
Open your browser to http://localhost:8080

---

## Configuration

See the [Configuration Guide](docs/CONFIGURATION.md) for detailed setup instructions.

### Minimum Required Configuration

```bash
# .env file
SCHOOLOGY_CONSUMER_KEY=your_key_here
SCHOOLOGY_CONSUMER_SECRET=your_secret_here
POWERSCHOOL_CLIENT_ID=your_client_id_here
POWERSCHOOL_CLIENT_SECRET=your_client_secret_here
```

---

## Authentication Setup

**Two options - choose what works best for you:**

### Option 1: Session Tokens (Recommended for MVP)

**Schoology:**
1. Log into Schoology in your browser
2. Open Chrome DevTools (F12) → Application → Cookies
3. Copy these values:
   - `SESS<hash>` (session ID)
   - `CSRF_TOKEN`
   - `CSRF_KEY`
   - Your user ID (UID)
4. Add to configuration file

**PowerSchool:**
1. Log into PowerSchool/The Source
2. Open Chrome DevTools → Application → Cookies
3. Copy `pstoken` and session cookies
4. Add to configuration file

**Advantages:** Simple, secure, you control authentication
**Disadvantage:** Must refresh cookies every 1-2 weeks

### Option 2: Credential Storage (Coming in Phase 2)

Provide your username and password (encrypted at rest), and Trunchbull will handle login automatically.

**Advantages:** Fully automated, no manual cookie management
**Disadvantage:** Must store credentials (though encrypted)

**See [Alternative Authentication Guide](docs/ALTERNATIVE_AUTH.md) for detailed instructions with screenshots.**

---

## Security and Privacy

**CRITICAL: READ BEFORE USING**

This application handles sensitive student data protected by FERPA and other privacy laws. You are responsible for:
- Running this application on infrastructure you control
- Implementing appropriate security measures
- Protecting student data from unauthorized access
- Complying with all applicable laws and school policies

**See [docs/SECURITY_AND_PRIVACY.md](docs/SECURITY_AND_PRIVACY.md) for complete details.**

### Key Security Recommendations
- Run behind a firewall, not exposed to the public internet
- Use HTTPS (self-signed certificate acceptable for local use)
- Enable database encryption
- Restrict access to authorized family members only
- Keep Docker and host OS updated
- Implement regular backups

---

## Architecture

Trunchbull consists of:
- **Backend**: Go API server handling OAuth, API clients, data aggregation
- **Frontend**: React SPA for the dashboard UI
- **Database**: SQLite for local data persistence
- **Cache**: In-memory cache for API response optimization

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture documentation.

```
┌─────────────────┐
│   React SPA     │  Dashboard UI
└────────┬────────┘
         │ HTTPS
         ↓
┌─────────────────┐
│   Go Backend    │  API Gateway + Sync
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│  SQLite DB      │  Local Storage
└─────────────────┘

External APIs:
├─ Schoology
└─ PowerSchool
```

---

## Documentation

- [Research Findings](docs/RESEARCH.md) - Platform capabilities and API details
- [Architecture](docs/ARCHITECTURE.md) - Technical architecture and design
- [Security & Privacy](docs/SECURITY_AND_PRIVACY.md) - Critical security information
- [Configuration Guide](docs/CONFIGURATION.md) - Detailed setup instructions
- [API Access](docs/API_ACCESS.md) - How to get API credentials
- [Development](docs/DEVELOPMENT.md) - For contributors

---

## Development

### Local Development Setup

```bash
# Backend
cd cmd/server
go run main.go

# Frontend
cd frontend
npm install
npm start
```

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed development instructions.

### Running Tests

```bash
# Backend tests
go test ./...

# Frontend tests
cd frontend
npm test
```

---

## Roadmap

### Phase 1: MVP (Current)
- [x] Research platform APIs
- [x] Architecture design
- [x] Security documentation
- [ ] Schoology integration
- [ ] Basic dashboard
- [ ] Docker containerization

### Phase 2: Full Integration
- [ ] PowerSchool integration
- [ ] Multiple student support
- [ ] Background sync
- [ ] Calendar widget
- [ ] Message inbox

### Phase 3: Polish
- [ ] GPA calculation
- [ ] Real-time updates
- [ ] Mobile-responsive design
- [ ] Notification system
- [ ] Data export

See [Issues](https://github.com/yourusername/trunchbull/issues) for detailed task tracking.

---

## Contributing

Contributions are welcome! Please:
1. Read [CONTRIBUTING.md](CONTRIBUTING.md)
2. Fork the repository
3. Create a feature branch
4. Submit a pull request

### Security Vulnerabilities

**DO NOT open public issues for security vulnerabilities.**

Email security concerns to: [your-email@example.com]

---

## FAQ

### Why "Trunchbull"?
Named after Miss Trunchbull from Roald Dahl's "Matilda" - a headmistress known for closely monitoring students (though hopefully you'll use this tool more kindly!).

### Is this legal?
Yes, as long as you:
- Have the legal right to access your child's educational records
- Comply with school district policies
- Don't violate platform Terms of Service
- Protect the data appropriately

See [docs/SECURITY_AND_PRIVACY.md](docs/SECURITY_AND_PRIVACY.md) for details.

### Can I use this for multiple families?
The current version is designed for single-family use. Multi-family support may come in a future version, but would require additional security considerations.

### Does this send data to the cloud?
No. All data stays on your local infrastructure. The only external communication is with your school's official platforms.

### What if my school doesn't use Schoology or PowerSchool?
Currently, only these two platforms are supported. Support for Canvas, Google Classroom, and others may be added in the future. Contributions welcome!

### Is there a mobile app?
Not yet, but the dashboard is mobile-responsive. A native mobile app may be developed in Phase 4.

---

## License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

### Important Legal Notes
- This software is provided "as is" without warranty
- Authors are not responsible for misuse or policy violations
- Users are responsible for compliance with all applicable laws
- This is personal, non-commercial software

---

## Acknowledgments

- Thanks to the Schoology and PowerSchool teams for providing APIs
- Inspired by parents who want to stay engaged with their children's education
- Built with Go, React, and lots of coffee

---

## Support

- [Documentation](docs/)
- [GitHub Issues](https://github.com/yourusername/trunchbull/issues)
- [Discussions](https://github.com/yourusername/trunchbull/discussions)

**Note**: This is a community project. Support is provided on a best-effort basis.

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

---

**Made with care for parents who care about education**
