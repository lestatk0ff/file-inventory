
# file-inventory

![GitHub release (latest by SemVer)](https://img.shields.io/github/v/release/lestatk0ff/file-inventory?sort=semver)


`file-inventory` is a Go command-line tool to:
- List all files in a directory (including subdirectories) and output the list to a text file (create command)
- Compare two inventory files and show the diff (diff command)


## Usage

### Create inventory file

```
file-inventory create DIR_PATH [-o OUTPUT_FILE]
```

- `DIR_PATH`: Directory to scan
- `-o OUTPUT_FILE`: (Optional) Output file name (default: file-inventory.txt)

Example:
```
file-inventory create ./mydir -o inventory1
```

### Diff two inventory files

```
file-inventory diff FILE1 FILE2
```

Shows the difference between two inventory files (lines present in one but not the other).

Example:
```
file-inventory diff inventory1 inventory2
```


## Example Output (inventory file)

```
testdir/file1.mp3
testdir/subdir/file2.txt
```

## Running Tests

To run the tests:

```
go test
```


## Build

```
go build -o file-inventory .
```
