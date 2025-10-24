# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Research Phase
- Comprehensive platform research for Schoology and PowerSchool APIs
- Architecture design documentation
- Security and privacy guidelines
- Initial project structure

## [0.1.0] - 2025-10-23

### Added
- Project initialization
- Go backend skeleton with Gin framework
- SQLite database schema and migrations
- Docker containerization setup
- Configuration management with Viper
- API route structure (stubs)
- Documentation:
  - Research findings (RESEARCH.md)
  - Architecture design (ARCHITECTURE.md)
  - Security and privacy guidelines (SECURITY_AND_PRIVACY.md)
  - README with quick start guide
- Development tools:
  - Makefile for common tasks
  - .gitignore for Go and Docker
  - .env.example for configuration
  - config.example.yaml

### In Progress
- Schoology OAuth implementation
- PowerSchool OAuth implementation
- API client for Schoology
- API client for PowerSchool
- React frontend

### Planned
- Background sync service
- Cache layer implementation
- Assignment aggregation
- Grade calculation and GPA
- Calendar/events widget
- Message inbox
- Multi-student support
- Dashboard UI

---

## Version History

### Phase 1: MVP (Weeks 1-2)
- [x] Research and documentation
- [x] Project structure
- [ ] Schoology integration
- [ ] Basic dashboard
- [ ] Docker deployment

### Phase 2: Full Integration (Weeks 3-4)
- [ ] PowerSchool integration
- [ ] Multiple student support
- [ ] Background sync
- [ ] Full dashboard features

### Phase 3: Polish (Weeks 5-6)
- [ ] GPA calculation
- [ ] Real-time updates
- [ ] Mobile-responsive UI
- [ ] Notification system

---

[Unreleased]: https://github.com/leftathome/trunchbull/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/leftathome/trunchbull/releases/tag/v0.1.0
