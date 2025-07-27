# Contributing to RimPay

Thank you for your interest in contributing to RimPay! This document outlines the process and guidelines for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Style Guidelines](#style-guidelines)
- [Release Process](#release-process)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please be respectful and professional in all interactions.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/rim-pay.git
   cd rim-pay
   ```
3. Add the original repository as upstream:
   ```bash
   git remote add upstream https://github.com/CatoSystems/rim-pay.git
   ```

## Development Setup

### Prerequisites

- Go 1.22 or later
- Make (for build commands)
- Git

### Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run tests to ensure everything works:
   ```bash
   make test
   ```

3. Build the project:
   ```bash
   make build
   ```

## Making Changes

### Workflow

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the [style guidelines](#style-guidelines)

3. Add or update tests as needed

4. Ensure all tests pass:
   ```bash
   make test
   ```

5. Format your code:
   ```bash
   make fmt
   ```

6. Check for issues:
   ```bash
   make vet
   ```

7. Commit your changes with a clear message:
   ```bash
   git commit -m "Add feature: brief description of changes"
   ```

### Commit Message Format

Use clear, concise commit messages:

- `feat: add new feature`
- `fix: resolve bug in payment processing`
- `docs: update README with examples`
- `test: add tests for phone validation`
- `refactor: improve error handling`

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests for specific packages
go test ./pkg/...
go test ./internal/...

# Run tests with coverage
go test -cover ./pkg/...
```

### Writing Tests

- Write tests for all new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Mock external dependencies
- Test both success and error cases

Example test structure:
```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewFeature(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFeature() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("NewFeature() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Pull Request Process

1. Update documentation if needed
2. Update CHANGELOG.md with your changes
3. Ensure all tests pass and code follows style guidelines
4. Create a pull request with:
   - Clear title describing the change
   - Detailed description of what was changed and why
   - Reference any related issues
   - Screenshots or examples if applicable

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] Added tests for new functionality
- [ ] Updated documentation

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
```

## Style Guidelines

### Go Code Style

Follow standard Go conventions:

- Use `gofmt` for formatting
- Follow effective Go guidelines
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors appropriately

### Documentation

- Update README.md for API changes
- Add godoc comments for public APIs
- Include examples in documentation
- Update CHANGELOG.md for all changes

### Example Documentation:
```go
// ProcessBPayPayment processes a payment using the B-PAY provider.
// It validates the request, authenticates with B-PAY, and submits the payment.
//
// Example:
//   request := &BPayPaymentRequest{
//       Amount: money.New(decimal.NewFromInt(10000), "MRU"),
//       PhoneNumber: phone,
//       Reference: "ORDER-123",
//       Passcode: "1234",
//   }
//   response, err := client.ProcessBPayPayment(ctx, request)
func (c *Client) ProcessBPayPayment(ctx context.Context, request *BPayPaymentRequest) (*PaymentResponse, error) {
    // Implementation...
}
```

## Release Process

1. Update version in relevant files
2. Update CHANGELOG.md with release notes
3. Create a release PR
4. After merge, tag the release:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```
5. Create GitHub release with changelog

## Questions?

If you have questions about contributing:

1. Check existing issues and discussions
2. Create a new issue with the `question` label
3. Reach out to maintainers

Thank you for contributing to RimPay!
