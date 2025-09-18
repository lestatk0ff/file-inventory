# Product Overview

`file-inventory` is a cross-platform Go CLI tool for file system inventory management. It provides two core functions:

1. **Create**: Recursively scan directories and generate text-based file inventories
2. **Diff**: Compare two inventory files and display differences in a clean table format

## Key Features
- Fast recursive file discovery using Go's standard library
- Clean, line-by-line text output for easy processing
- Professional table-formatted diff output with clear visual indicators
- Cross-platform compatibility (Windows, macOS, Linux)
- Flexible filtering with glob patterns (include/exclude)
- Path options (relative paths, hidden files)
- Sorted output capability
- Robust error handling and input validation
- Minimal external dependencies

## Use Cases
- File system auditing and change tracking
- Backup verification
- Directory synchronization monitoring
- Build artifact comparison