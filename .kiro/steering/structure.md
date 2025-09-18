# Project Structure

## File Organization
This is a single-module Go project with a flat structure optimized for simplicity:

```
├── cmd.go              # CLI command definitions and main entry point
├── fileutils.go        # File discovery and I/O utilities
├── diff.go            # Diff logic and table formatting
├── cmd_test.go        # CLI command tests
├── fileutils_test.go  # File utility tests
├── diff_test.go       # Diff functionality tests
├── go.mod             # Go module definition
├── go.sum             # Dependency checksums
└── README.md          # Documentation
```

## Architecture Patterns

### Separation of Concerns
- **cmd.go**: Cobra command setup, argument parsing, and CLI orchestration
- **fileutils.go**: Pure file system operations (finding files, writing output)
- **diff.go**: Comparison logic and table rendering

### Testing Strategy
- **One test file per source file** pattern
- **Functional testing** with temporary directories and files
- **Table-driven tests** where appropriate
- **Integration testing** of CLI commands using cobra's testing utilities

### Code Organization Principles
- **Flat structure** - no subdirectories for this simple tool
- **Single responsibility** - each file has a clear, focused purpose
- **Standard library first** - minimize external dependencies
- **Error handling** - explicit error returns with CLI-appropriate exit codes

### Naming Conventions
- **Files**: lowercase with descriptive names (`fileutils.go`, not `utils.go`)
- **Functions**: camelCase with clear, action-oriented names (`findFiles`, `showDiff`)
- **Tests**: `Test` prefix with descriptive suffixes (`TestCobraCreateCommand`)

## Development Workflow
- Tests are co-located with source files
- No build scripts or complex tooling required
- Standard Go project layout for easy onboarding