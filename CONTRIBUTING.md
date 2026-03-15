# Contributing to Pangolin

Thank you for your interest in contributing to Pangolin!

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported
2. Create a new issue with:
   - Clear title
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details

### Suggesting Features

1. Check existing issues/enhancements
2. Create an issue with:
   - Clear description
   - Use case
   - Proposed solution

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Add tests if applicable
5. Ensure code passes: `go build ./...` and `go vet ./...`
6. Commit with clear messages
7. Push and create PR

## Development Setup

```bash
# Clone
git clone https://github.com/pangolin-cms/staticpress.git
cd staticpress

# Build
go build -o pangolin .

# Run tests (when available)
go test ./...

# Build dashboard
go build -o pangolin-dashboard ./dashboard/
```

## Code Style

- Follow Go standard conventions
- Use meaningful variable names
- Add comments for complex logic
- Keep functions small and focused

## License

By contributing, you agree your code will be licensed under the MIT License.
