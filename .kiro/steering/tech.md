# Technology Stack

## Language & Runtime
- **Go 1.24.7** - Modern Go version with latest features
- Cross-platform compatibility (Windows, macOS, Linux)

## Dependencies
- **cobra** (`github.com/spf13/cobra`) - CLI framework for command structure and argument parsing
- **tablewriter** (`github.com/olekukonko/tablewriter`) - Professional table formatting for diff output
- **Standard library only** for core file operations (filepath, os, bufio)

## Build System
Standard Go toolchain with no additional build tools required.

## Common Commands

### Development
```bash
# Run all tests with verbose output
go test -v ./...

# Run specific test categories
go test -v -run "TestCobra"        # CLI command tests
go test -v -run "TestFind|TestInventory|TestWrite"  # File utility tests
go test -v -run "TestShowDiff"     # Diff functionality tests
```

### Building
```bash
# Build executable
go build -o file-inventory .

# Install directly to GOPATH/bin
go install
```

### Usage
```bash
# Basic create inventory
file-inventory create ./mydir -o inventory1.txt

# Advanced create with filtering and options (relative paths by default)
file-inventory create ./mydir --sort --include "*.txt" --exclude "*.log" -o inventory1.txt

# Use full absolute paths
file-inventory create ./mydir --full --sort -o inventory1.txt

# Compare inventories
file-inventory diff inventory1.txt inventory2.txt
```

## Code Style
- Standard Go formatting (`gofmt`)
- Clear, descriptive function names
- Minimal external dependencies philosophy
- Proper error handling with wrapped errors and graceful degradation
- Configuration-driven approach with struct-based options
- Table-driven tests for comprehensive coverage
- Separation of CLI logic from core functionality