# Contributing to Trunchbull

Thank you for your interest in contributing to Trunchbull! This document provides guidelines for contributing to the project.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Security](#security)

---

## Code of Conduct

This project is dedicated to providing a welcoming and inclusive experience for everyone. We expect all contributors to:
- Be respectful and considerate
- Accept constructive criticism gracefully
- Focus on what's best for the community
- Show empathy towards others

---

## Getting Started

### Prerequisites
- Go 1.23 or higher
- Docker and Docker Compose
- Git
- Node.js 18+ (for frontend development)

### Setup Development Environment

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/yourusername/trunchbull.git
   cd trunchbull
   ```

2. **Run setup**
   ```bash
   make setup
   ```

3. **Install dependencies**
   ```bash
   make deps
   ```

4. **Configure environment**
   - Edit `.env` with your API credentials
   - Edit `config/config.yaml` as needed

5. **Run the application**
   ```bash
   make run
   ```

---

## How to Contribute

### Reporting Bugs
- Check if the bug has already been reported in [Issues](https://github.com/leftathome/trunchbull/issues)
- If not, create a new issue with:
  - Clear, descriptive title
  - Steps to reproduce
  - Expected behavior
  - Actual behavior
  - Environment details (OS, Go version, etc.)
  - Logs or screenshots if applicable

### Suggesting Features
- Check [Issues](https://github.com/leftathome/trunchbull/issues) for existing feature requests
- Create a new issue with:
  - Clear description of the feature
  - Use case and benefits
  - Potential implementation approach
  - Any relevant examples or mockups

### Contributing Code
1. Check [Issues](https://github.com/leftathome/trunchbull/issues) for open tasks
2. Comment on the issue to claim it
3. Fork the repository
4. Create a feature branch
5. Make your changes
6. Submit a pull request

---

## Development Workflow

### Branch Naming Convention
- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring
- `test/description` - Test additions or fixes

### Commit Message Guidelines
Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Maintenance tasks

**Examples:**
```
feat(api): add endpoint for fetching student assignments

fix(db): resolve foreign key constraint issue in grades table

docs(readme): update installation instructions
```

### Pull Request Process

1. **Before submitting:**
   - Update documentation if needed
   - Add tests for new functionality
   - Ensure all tests pass
   - Run linters: `make lint`
   - Update CHANGELOG.md

2. **Submit PR with:**
   - Clear title describing the change
   - Description of what changed and why
   - Reference to related issue(s)
   - Screenshots if UI changes

3. **PR Review:**
   - Maintainers will review your PR
   - Address any feedback
   - Once approved, it will be merged

---

## Coding Standards

### Go Code Style
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `go fmt` for formatting
- Use `go vet` for static analysis
- Keep functions small and focused
- Write clear, descriptive names
- Add comments for exported functions

**Example:**
```go
// GetStudentAssignments retrieves all assignments for a given student
// from both Schoology and PowerSchool platforms.
func GetStudentAssignments(studentID string) ([]Assignment, error) {
    // Implementation
}
```

### Project Structure
```
trunchbull/
├── cmd/
│   └── server/          # Main application entry point
├── internal/            # Private application code
│   ├── api/            # HTTP handlers
│   ├── auth/           # OAuth and authentication
│   ├── client/         # API clients (Schoology, PowerSchool)
│   ├── config/         # Configuration management
│   ├── db/             # Database layer
│   ├── models/         # Data models
│   ├── service/        # Business logic
│   └── sync/           # Background sync jobs
├── frontend/           # React frontend
├── docs/               # Documentation
└── deployments/        # Deployment configs
```

### Error Handling
- Always handle errors explicitly
- Use meaningful error messages
- Wrap errors with context: `fmt.Errorf("failed to fetch assignments: %w", err)`
- Log errors appropriately

### Security Best Practices
- Never commit secrets or credentials
- Validate all user input
- Use parameterized queries
- Log security events
- Follow OWASP guidelines

---

## Testing

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/api/...
```

### Writing Tests
- Write unit tests for all business logic
- Write integration tests for API endpoints
- Use table-driven tests where appropriate
- Mock external dependencies

**Example:**
```go
func TestGetStudentAssignments(t *testing.T) {
    tests := []struct {
        name      string
        studentID string
        want      []Assignment
        wantErr   bool
    }{
        {
            name:      "valid student",
            studentID: "student-1",
            want:      []Assignment{/* ... */},
            wantErr:   false,
        },
        {
            name:      "invalid student",
            studentID: "invalid",
            want:      nil,
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := GetStudentAssignments(tt.studentID)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            // Additional assertions
        })
    }
}
```

---

## Documentation

### What to Document
- Public APIs and functions
- Configuration options
- Architecture decisions
- Setup and deployment procedures
- API endpoints

### Documentation Standards
- Use clear, concise language
- Include code examples
- Keep documentation up-to-date with code changes
- Use Markdown for formatting

### Where to Add Documentation
- Code comments for functions and packages
- `docs/` directory for comprehensive guides
- `README.md` for quick start and overview
- API documentation in `docs/API.md`

---

## Security

### Reporting Security Vulnerabilities
**DO NOT open public issues for security vulnerabilities.**

Instead:
1. Email security concerns to: [security contact - TBD]
2. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### Security Guidelines
- Never commit secrets, API keys, or credentials
- Use environment variables for sensitive configuration
- Validate and sanitize all input
- Follow the principle of least privilege
- Keep dependencies updated
- Review [SECURITY_AND_PRIVACY.md](docs/SECURITY_AND_PRIVACY.md)

---

## Additional Resources

### Useful Links
- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [React Documentation](https://react.dev/)
- [Docker Documentation](https://docs.docker.com/)
- [Schoology API Docs](https://developers.schoology.com/)

### Community
- [GitHub Discussions](https://github.com/leftathome/trunchbull/discussions)
- [Issue Tracker](https://github.com/leftathome/trunchbull/issues)

---

## Questions?

If you have questions about contributing:
1. Check existing documentation
2. Search [GitHub Issues](https://github.com/leftathome/trunchbull/issues)
3. Ask in [GitHub Discussions](https://github.com/leftathome/trunchbull/discussions)
4. Reach out to maintainers

---

Thank you for contributing to Trunchbull! Your efforts help parents stay engaged with their children's education.
