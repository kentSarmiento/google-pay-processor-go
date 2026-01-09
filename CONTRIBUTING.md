# Contributing to Google Pay Processor

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and be patient with questions
- Focus on constructive criticism
- Keep discussions professional

## Getting Started

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR-USERNAME/google-pay-processor-go.git
   cd google-pay-processor-go
   ```
3. **Set up development environment**
   ```bash
   go mod download
   ```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test improvements

### 2. Make Your Changes

Follow the coding standards:
- Use `go fmt` to format code
- Run `go vet` to check for issues
- Write tests for new functionality
- Update documentation as needed

### 3. Write Tests

All new code should have tests:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 4. Lint Your Code

```bash
make lint
```

Or manually:
```bash
golangci-lint run ./...
```

### 5. Commit Your Changes

Write clear commit messages:
```bash
git commit -m "feat: add new feature X"
```

Commit message format:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test changes
- `refactor:` - Code refactoring
- `style:` - Code style changes
- `chore:` - Build process or auxiliary tool changes

### 6. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 7. Create a Pull Request

- Go to the original repository on GitHub
- Click "New Pull Request"
- Select your fork and branch
- Fill in the PR template
- Submit for review

## Coding Standards

### Go Style Guide

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

Key points:
- Use `gofmt` for formatting
- Follow effective Go guidelines
- Use meaningful variable names
- Keep functions small and focused
- Write self-documenting code

### Project-Specific Standards

1. **Clean Architecture**: Follow the layer separation
   - Domain layer: Pure business logic
   - Application layer: Use cases
   - Infrastructure layer: External integrations
   - Presentation layer: HTTP handlers

2. **Error Handling**: Wrap errors with context
   ```go
   if err != nil {
       return fmt.Errorf("%w: additional context", domain.ErrType, err)
   }
   ```

3. **Logging**: Use structured logging
   ```go
   logger.Info("message", map[string]interface{}{
       "key": "value",
   })
   ```

4. **Testing**: Write table-driven tests
   ```go
   tests := []struct {
       name string
       input string
       expected string
   }{
       // test cases
   }
   ```

5. **Interfaces**: Keep them small and focused
   ```go
   type Reader interface {
       Read(p []byte) (n int, err error)
   }
   ```

## Testing Guidelines

### Unit Tests

- Test one component in isolation
- Mock external dependencies
- Use table-driven tests
- Test error cases
- Aim for >80% coverage

Example:
```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {name: "valid input", input: "test", want: "result", wantErr: false},
        {name: "empty input", input: "", want: "", wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("MyFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("MyFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests

- Test component interactions
- Use real implementations where possible
- Test with external services
- Clean up resources after tests

## Documentation

### Code Documentation

Use GoDoc format:
```go
// ProcessPayment processes a Google Pay token and executes the payment.
// It validates the token, decrypts it, and forwards to the payment processor.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - token: Google Pay token to process
//   - amount: Payment amount in smallest currency unit
//
// Returns:
//   - PaymentResponse containing the result
//   - Error if processing fails
func ProcessPayment(ctx context.Context, token *GooglePayToken, amount int64) (*PaymentResponse, error) {
    // implementation
}
```

### README Updates

Update README.md if:
- Adding new features
- Changing configuration
- Modifying API endpoints
- Updating dependencies

### Additional Documentation

Add to `docs/` directory for:
- Architecture decisions
- Integration guides
- Deployment instructions
- Troubleshooting guides

## Pull Request Guidelines

### PR Title

Use conventional commit format:
```
feat: add new Google Pay token validation
fix: correct timeout handling in payment processor
docs: update API documentation
```

### PR Description

Include:
- What changed and why
- How to test the changes
- Screenshots (for UI changes)
- Breaking changes (if any)
- Related issues

Template:
```markdown
## Description
Brief description of changes

## Changes Made
- Item 1
- Item 2

## Testing
How to test the changes

## Screenshots
(if applicable)

## Breaking Changes
(if any)

## Related Issues
Closes #123
```

### PR Checklist

Before submitting:
- [ ] Tests pass locally
- [ ] Linter passes
- [ ] Documentation updated
- [ ] Tests added for new functionality
- [ ] No breaking changes (or documented)
- [ ] Commit messages follow convention
- [ ] Branch is up to date with main

## Review Process

1. **Automated Checks**: CI runs tests and linters
2. **Code Review**: Maintainers review code
3. **Feedback**: Address review comments
4. **Approval**: At least one maintainer approves
5. **Merge**: Maintainer merges the PR

### Review Timeline

- Simple fixes: 1-2 days
- Features: 3-5 days
- Major changes: 1-2 weeks

## Release Process

1. Version bump following semantic versioning
2. Update CHANGELOG.md
3. Create release tag
4. Build and push Docker image
5. Deploy to staging
6. Deploy to production

## Getting Help

- Open an issue for bugs or feature requests
- Use discussions for questions
- Check existing issues and docs first
- Be specific and provide context

## License

By contributing, you agree that your contributions will be licensed under the project's license.

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- GitHub contributors page

Thank you for contributing! 🎉
