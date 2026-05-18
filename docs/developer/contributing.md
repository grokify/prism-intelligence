# Contributing

Thank you for your interest in contributing to PRISM. This guide covers the contribution process and code standards.

## Getting Started

### Fork and Clone

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/prism.git
cd prism
git remote add upstream https://github.com/grokify/prism-intelligence.git
```

### Install Dependencies

```bash
go mod download
```

### Verify Setup

```bash
go build ./cmd/prism
go test -v ./...
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Write code following the [Code Style](#code-style) guidelines
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Run linter
golangci-lint run
```

### 4. Commit Your Changes

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
git commit -m "feat(analysis): add initiative recommendation engine"
git commit -m "fix(score): handle nil metrics in calculation"
git commit -m "docs(cli): add export command documentation"
```

**Commit Types:**

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Code formatting (no logic change) |
| `refactor` | Code restructuring |
| `test` | Adding or updating tests |
| `chore` | Maintenance tasks |

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Code Style

### Go Conventions

- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines

### Naming

- Use clear, descriptive names
- Acronyms should be consistent case: `SLO`, `slo`, not `Slo`
- Package names should be lowercase, single words

### Error Handling

Always handle errors explicitly:

```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Avoid
result, _ := doSomething()  // Don't ignore errors
```

### Comments

- Export functions need doc comments
- Use complete sentences
- Explain "why" not "what"

```go
// CalculatePRISMScore computes the overall health score for a document.
// It combines maturity levels, metric performance, and customer awareness
// using configurable weights.
func (d *PRISMDocument) CalculatePRISMScore(
    stageWeights map[string]float64,
    domainWeights map[string]float64,
) PRISMScore {
    // ...
}
```

### Testing

- Table-driven tests are preferred
- Test both success and error cases
- Use meaningful test names

```go
func TestMetric_MeetsSLO(t *testing.T) {
    tests := []struct {
        name     string
        metric   Metric
        expected bool
    }{
        {
            name: "higher is better - meets target",
            metric: Metric{
                Current:        99.9,
                TrendDirection: TrendHigherBetter,
                SLO:            &SLO{Operator: "gte", Value: 99.5},
            },
            expected: true,
        },
        // ... more cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.metric.MeetsSLO()
            if got != tt.expected {
                t.Errorf("MeetsSLO() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

## Adding Features

### New CLI Command

1. Create `cmd/prism/mycommand.go`:

```go
package main

import (
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand [file]",
    Short: "Brief description",
    Long:  `Longer description with examples.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runMyCommand,
}

func init() {
    rootCmd.AddCommand(myCmd)
    myCmd.Flags().StringP("output", "o", "", "Output file")
}

func runMyCommand(cmd *cobra.Command, args []string) error {
    // Implementation
    return nil
}
```

2. Add documentation in `docs/cli/mycommand.md`

3. Update `mkdocs.yml` navigation

### New Export Format

1. Create `export/myformat.go`:

```go
package export

import "github.com/grokify/prism-intelligence"

type MyFormatDocument struct {
    // Define output structure
}

func ConvertToMyFormat(doc *prism.PRISMDocument) *MyFormatDocument {
    // Conversion logic
}
```

2. Add tests in `export/myformat_test.go`

3. Add CLI subcommand under `prism export myformat`

### New Metric Type

1. Add constant in `constants.go`:

```go
const MetricTypeMyType MetricType = "mytype"
```

2. Update `ValidMetricTypes()` function

3. Add handling in `score.go` if needed

4. Update JSON Schema by running `cd schema && go run generate.go`

## Documentation

### Code Documentation

- All exported types and functions need doc comments
- Include examples for complex APIs

### User Documentation

- Located in `docs/` (MkDocs format)
- Preview locally: `mkdocs serve`
- Follow existing structure for new pages

### Changelog

- Update `CHANGELOG.md` for user-facing changes
- Group changes by type (Added, Changed, Fixed, etc.)

## Pull Request Guidelines

### Before Submitting

- [ ] Tests pass (`go test -v ./...`)
- [ ] Linting passes (`golangci-lint run`)
- [ ] Documentation updated if needed
- [ ] Commit messages follow conventional commits
- [ ] No unrelated changes included

### PR Description

Include:

- Summary of changes
- Motivation/context
- Testing performed
- Related issues (use `Fixes #123` to auto-close)

### Review Process

1. Maintainers will review within a few days
2. Address feedback with new commits
3. Squash commits if requested before merge

## Questions?

- Open a [GitHub Discussion](https://github.com/grokify/prism-intelligence/discussions) for questions
- Open an [Issue](https://github.com/grokify/prism-intelligence/issues) for bugs or feature requests

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
