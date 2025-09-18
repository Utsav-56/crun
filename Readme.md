# CRUN - C/C++ Compile and Run Tool

A fast and simple command-line tool to compile and run C/C++ files with automatic compiler detection, caching, and flexible configuration options.

## Features

- **Automatic Compiler Detection**: Supports clang, gcc, zig, cl (MSVC), and bytes
- **Smart Caching**: Only recompiles when source files are modified
- **Flexible Flag System**: Both short (`-v`) and long (`--verbose`) flags supported
- **Cross-Platform**: Works on Windows, Linux, and macOS
- **Extension Auto-Detection**: Automatically finds source files even without extensions
- **Custom Output Control**: Specify output names and directories
- **Runtime Arguments**: Pass arguments directly to compiled binaries

## Installation

1. Clone or download the repository
2. Build the tool:
    ```bash
    go build -o crun .
    ```
3. Add to your PATH for global access

## Quick Start

```bash
# Simple compile and run
crun main.c

# Compile with specific compiler
crun -c gcc main.cpp

# Always recompile with verbose output
crun --verbose --recompile main.c

# Custom output name and directory
crun -o myprogram -d ./bin main.c
```

## Command Syntax

```bash
crun [flags] <filename>
```

## Flags Reference

### Core Flags

| Short | Long          | Description                           | Example          |
| ----- | ------------- | ------------------------------------- | ---------------- |
| `-v`  | `--verbose`   | Verbose mode - don't clear log output | `crun -v main.c` |
| `-n`  | `--recompile` | Always recompile source file          | `crun -n main.c` |
| `-h`  | `--help`      | Show help message                     | `crun -h`        |

### Compiler Control

| Short | Long         | Description                     | Example                      |
| ----- | ------------ | ------------------------------- | ---------------------------- |
| `-c`  | `--compiler` | Manually choose compiler        | `crun -c clang main.c`       |
| `-e`  | `--extra`    | Extra flags to pass to compiler | `crun -e "-O2 -Wall" main.c` |

### Output Control

| Short | Long          | Description                   | Example                |
| ----- | ------------- | ----------------------------- | ---------------------- |
| `-o`  | `--output`    | Output binary name            | `crun -o myapp main.c` |
| `-d`  | `--directory` | Directory to store the binary | `crun -d ./bin main.c` |

### Runtime Control

| Short | Long         | Description                     | Example                      |
| ----- | ------------ | ------------------------------- | ---------------------------- |
| `-r`  | `--run-args` | Arguments to pass to the binary | `crun -r "arg1 arg2" main.c` |

## Supported Compilers

The tool automatically detects and uses the first available compiler in this order:

1. **clang** - LLVM C/C++ compiler
2. **gcc** - GNU Compiler Collection
3. **zig** - Zig compiler (using `zig cc`)
4. **cl** - Microsoft Visual C++ compiler
5. **bytes** - Custom compiler (if available)

### Manual Compiler Selection

```bash
# Use specific compiler
crun --compiler gcc main.c
crun -c clang main.cpp
crun -c zig main.c
```

## File Extension Handling

CRUN automatically detects source files even without extensions:

```bash
# These are equivalent if main.c exists
crun main.c
crun main
```

**Supported Extensions**: `.c`, `.cpp`, `.cc`, `.cxx`, `.h`, `.hpp`, `.hh`, `.hxx`

## Caching System

CRUN implements smart caching to avoid unnecessary recompilation:

- Compares source file modification time with binary modification time
- Only recompiles if source is newer or binary doesn't exist
- Cache stored in `.crun` directory in current working directory
- Use `-n` or `--recompile` to force recompilation

```bash
# Force recompilation
crun -n main.c
crun --recompile main.c
```

## Advanced Usage Examples

### Development Workflow

```bash
# Debug build with all warnings
crun -c gcc -e "-g -Wall -Wextra -fsanitize=address" main.c

# Optimized release build
crun -c clang -e "-O3 -DNDEBUG" -o release -d ./bin main.c

# Cross-platform build with zig
crun -c zig -e "-target x86_64-windows" main.c
```

### Project Organization

```bash
# Organize binaries in separate directory
crun -d ./build -o debug_version main.c

# Custom naming convention
crun -o "myapp_v1.2" -d ./releases main.c
```

### Testing and Arguments

```bash
# Run with test arguments
crun -r "--input test.txt --output result.txt" main.c

# Verbose compilation with runtime args
crun -v -r "file1.txt file2.txt" processor.c
```

### Combining Flags

```bash
# Multiple short flags
crun -vn -c gcc -e "-O2" main.c

# Mixed short and long flags
crun --verbose -n --compiler clang main.cpp

# Complex build configuration
crun --compiler gcc --extra "-std=c11 -Wall -O2" --output optimized --directory ./bin --run-args "input.txt" main.c
```

## Compiler-Specific Behaviors

### GCC/Clang

```bash
# Standard compilation
args: ["-o", "output.exe", "source.c"]

# With extra flags
args: ["-o", "output.exe", "-O2", "-Wall", "source.c"]
```

### Microsoft Visual C++ (cl)

```bash
# Standard compilation
args: ["/Fe:output.exe", "source.c"]

# With extra flags
args: ["/O2", "/W4", "/Fe:output.exe", "source.c"]
```

### Zig Compiler

```bash
# Uses zig cc mode
args: ["cc", "-o", "output.exe", "source.c"]
```

## Directory Structure

```
your-project/
├── main.c
├── .crun/              # Auto-created cache directory
│   ├── main.exe        # Compiled binaries
│   └── other_file.exe
└── bin/                # Custom output directory (if specified)
    └── release.exe
```

## Error Handling

### Common Errors and Solutions

1. **"No supported compiler found"**

    ```bash
    # Install a supported compiler
    # Windows: Install Visual Studio or mingw
    # Linux: sudo apt install gcc
    # macOS: xcode-select --install
    ```

2. **"Specified compiler 'xyz' not found"**

    ```bash
    # Check available compilers
    crun -h  # Shows supported compilers

    # Use auto-detection instead
    crun main.c
    ```

3. **"Failed to compile the source file"**

    ```bash
    # Use verbose mode to see full output
    crun -v main.c

    # Check syntax and includes
    crun -e "-Wall" main.c
    ```

## Environment Variables

CRUN respects standard environment variables:

- `PATH` - For finding compilers
- `TEMP`/`TMP` - For temporary files (fallback)

## Performance Tips

1. **Use caching**: Let CRUN detect when recompilation is needed
2. **Organize outputs**: Use `-d` to keep binaries organized
3. **Batch operations**: Use scripts for multiple files
4. **Verbose mode**: Use `-v` only when debugging

## Integration Examples

### With Make

```makefile
debug:
    crun -c gcc -e "-g -DDEBUG" -o debug_app -d ./build main.c

release:
    crun -c clang -e "-O3 -DNDEBUG" -o release_app -d ./dist main.c
```

### With Scripts

```bash
#!/bin/bash
# build.sh
crun --compiler gcc --extra "-std=c11 -Wall -Wextra" --output myapp --directory ./bin src/main.c
```

### With IDEs

Most IDEs can be configured to use CRUN as a build tool:

**VS Code tasks.json**:

```json
{
	"label": "crun build",
	"type": "shell",
	"command": "crun",
	"args": ["-v", "-c", "gcc", "${file}"],
	"group": "build"
}
```

## Troubleshooting

### Windows-Specific Issues

1. **Path separators**: CRUN handles Windows paths automatically
2. **MSVC environment**: Ensure Visual Studio tools are in PATH
3. **Permissions**: Run as administrator if needed for system directories

### Linux/macOS Issues

1. **Compiler installation**: Use package managers (apt, yum, brew)
2. **Permissions**: Check execute permissions on CRUN binary
3. **Library paths**: Use `-e` flag for custom library paths

## Contributing

To add new compilers or features:

1. Add compiler to `SUPPORTED_COMPAILERS` slice
2. Add compilation logic in `compileSourceFile()`
3. Update help text and documentation
4. Test on target platforms

### Adding New Flags

```go
// Add to flagAliases map
"--new-flag": "-x",

// Add to Flags struct
type Flags struct {
    // ... existing fields
    newFlag bool
}

// Add to parseFlags() switch statement
case "-x":
    flags.newFlag = true
```

## License

[Add your license information here]

## Changelog

### Version 1.0.0

- Initial release with basic compilation and caching
- Support for major C/C++ compilers
- Flag system with short and long options
- Cross-platform compatibility
- Automatic file extension detection
- Smart recompilation based on modification times

---

**Made with ❤️ for C/C++ developers who want fast, simple compilation workflows.**
