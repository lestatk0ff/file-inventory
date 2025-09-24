
# file-inventory

![GitHub release (latest by SemVer)](https://img.shields.io/github/v/release/lestatk0ff/file-inventory?sort=semver)

`file-inventory` is a Go command-line tool to:
- List all files in a directory (including subdirectories) and output the list to a text file (create command)
- Compare two inventory files and show the diff in a clean table format (diff command)

## Features

- **Fast file discovery**: Recursively scans directories to find all files
- **Clean output**: Outputs file paths one per line for easy processing
- **Professional diff display**: Shows differences in a formatted table with clear indicators
- **Cross-platform**: Works on Windows, macOS, and Linux
- **Flexible filtering**: Include/exclude files using glob patterns
- **Path options**: Support for relative paths and hidden files
- **Sorted output**: Optional alphabetical sorting of file paths
- **Robust error handling**: Graceful handling of permission errors and invalid inputs
- **No external dependencies** for basic functionality


## Usage

### Create inventory file

```
file-inventory create DIR_PATH [flags]
```

**Flags:**
- `-o, --output string`: Output file name (default: file-inventory.txt)
- `--sort`: Sort file paths alphabetically in output
- `--full`: Use full absolute paths (default: relative paths from scan directory)
- `--hidden`: Include hidden files and directories
- `--include strings`: Include only files matching these glob patterns
- `--exclude strings`: Exclude files matching these glob patterns

**Examples:**
```bash
# Basic usage (uses relative paths by default)
file-inventory create ./mydir -o inventory1.txt

# Sort output and use full absolute paths
file-inventory create ./mydir --sort --full -o inventory1.txt

# Include only text files, including hidden ones
file-inventory create ./mydir --include "*.txt" --hidden -o inventory1.txt

# Exclude all .log and .tmp files
file-inventory create ./mydir --exclude "*.log" --exclude "*.tmp" -o inventory1.txt
```


### Diff two inventory files

```
file-inventory diff FILE1 FILE2
```

Compares two inventory files and displays differences in a clean table format. The output shows:
- **file_path**: The path of files that differ between the inventories
- **FILE1 column**: Shows `+` if file exists only in FILE1, `-` if missing from FILE1
- **FILE2 column**: Shows `+` if file exists only in FILE2, `-` if missing from FILE2

Example:
```
file-inventory diff inventory1.txt inventory2.txt
```

**Sample diff output:**
```
 file_path        │ inventory1.txt │ inventory2.txt
──────────────────┼────────────────┼────────────────
 docs/readme.txt  │ +              │ -
 src/main.go      │ -              │ +
 tests/unit.go    │ +              │ -
```


## Example Output (inventory file)

```
testdir/file1.mp3
testdir/subdir/file2.txt
testdir/subdir/nested/file3.doc
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [tablewriter](https://github.com/olekukonko/tablewriter) - Table formatting for diff output

## Testing

The project includes comprehensive test coverage organized by functionality:

- `cmd_test.go` - Tests for CLI commands and cobra integration
- `fileutils_test.go` - Tests for file discovery and writing utilities
- `diff_test.go` - Tests for diff functionality and table output

### Running Tests

Run all tests:
```
go test -v ./...
```

Run specific test categories:
```bash
# Test CLI commands and integration
go test -v -run "TestRun|TestCobra"

# Test file utilities and configuration
go test -v -run "TestFind|TestWrite|TestConfig"

# Test diff functionality
go test -v -run "TestShowDiff|TestReadFileLines"
```


## Build

Build the executable:
```
go build -o file-inventory .
```

Or install directly:
```
go install
```

## Project Structure

```
├── cmd.go           # CLI command definitions and main entry point
├── fileutils.go     # File discovery and I/O utilities
├── diff.go          # Diff logic and table formatting
├── cmd_test.go      # CLI command tests
├── fileutils_test.go # File utility tests
└── diff_test.go     # Diff functionality tests
```
