# Contributing to Kontainer

Thank you for considering contributing to Kontainer! 🎉

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs actual behavior
- **Screenshots** if applicable
- **Environment details**: OS, browser, Docker version, etc.
- **Error logs** from browser console or Docker logs

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Clear title and description**
- **Use case**: Why would this feature be useful?
- **Examples**: How would it work?
- **Alternatives**: What alternatives have you considered?

### Pull Requests

1. **Fork the repository** and create your branch from `master`
2. **Make your changes**:
   - Follow the existing code style
   - Add comments for complex logic
   - Update documentation if needed
3. **Test your changes**:
   - Build and run the application
   - Test affected features thoroughly
   - Check for regressions
4. **Commit your changes**:
   - Use clear, descriptive commit messages
   - Reference issues in commits (e.g., "Fix #123")
5. **Push to your fork** and submit a pull request

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git
- A text editor or IDE

### Building from Source

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/kontainer.git
cd kontainer

# Install dependencies
go mod download

# Build
go build -o kontainer cmd/kontainer/main.go

# Run
./kontainer
```

### Project Structure

```
kontainer/
├── cmd/kontainer/       # Application entry point
├── internal/
│   ├── api/            # HTTP handlers and routing
│   ├── database/       # Database initialization
│   ├── models/         # Data models
│   └── service/        # Business logic
└── web/static/         # Frontend assets (CSS, JS, images)
```

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Keep functions focused and readable
- Add comments for non-obvious code

### Testing

```bash
# Run tests (when available)
go test ./...

# Run with race detector
go test -race ./...
```

## Git Workflow

1. Create a feature branch: `git checkout -b feature/amazing-feature`
2. Make your changes and commit: `git commit -m 'Add amazing feature'`
3. Push to your fork: `git push origin feature/amazing-feature`
4. Open a Pull Request

### Commit Message Guidelines

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Keep first line under 72 characters
- Reference issues and PRs in the body

Example:
```
Add dark mode toggle to settings page

- Implement theme switcher in settings
- Store preference in localStorage
- Apply theme on page load

Fixes #42
```

## Feature Requests

We use GitHub Issues to track feature requests. Before submitting:

1. **Search existing issues** to avoid duplicates
2. **Describe the problem** you're trying to solve
3. **Propose a solution** if you have one in mind
4. **Be open to discussion** about implementation

## Questions?

- Open a [Discussion](https://github.com/yourusername/kontainer/discussions) for general questions
- Check [existing issues](https://github.com/yourusername/kontainer/issues) for known problems
- Review [documentation](README.md) and [technical docs](TECHNICAL-DOCS.md)

## Code of Conduct

### Our Pledge

We are committed to providing a friendly, safe, and welcoming environment for all contributors.

### Our Standards

**Positive behavior includes:**
- Being respectful and inclusive
- Accepting constructive criticism
- Focusing on what's best for the community
- Showing empathy toward others

**Unacceptable behavior includes:**
- Harassment, insults, or derogatory comments
- Trolling or inflammatory comments
- Publishing others' private information
- Other conduct inappropriate in a professional setting

### Enforcement

Project maintainers have the right to remove, edit, or reject comments, commits, code, and other contributions that don't align with this Code of Conduct.

## Attribution

Thank you to all contributors who help make Kontainer better! 🙏

This project follows the [Contributor Covenant](https://www.contributor-covenant.org/) Code of Conduct.
