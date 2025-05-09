# CDT — Common Developer Tool

A front-end tool that aims to simplify developer workflow. It unifies the interface for working
with different types of projects and programming languages.

## Usage

```shell
# Run main target - also do everything that is needed to run main target
cdt run main
```

**Some functions might not be available if they are not supported by the project or the corresponding tool is not
available**.

### Testing

```shell
# Run all tests
cdt test

# Run specific subset of tests
cdt test test1
```

### Code formatting

Tool allows reformatting the whole project or specific file by formatting tool defined by project workflow.

```shell
# Reformat whole project
cdt format

# Reformat file or multiple files
cdt format file1 file2
```

It's possible to check the code format without changing files. It can be used in CI to
verify if code is properly formatted.

```shell
cdt format --check
```

### Code linting

```shell
# Lint whole project
cdt lint

# Lint file or multiple files
cdt lint file1 file2
```

### Project configuration

Some project needs to be configured (e.g., CMake). In most cases, it's handled by the tool, but if a specific
configuration is required, it's possible to invoke configuration manually.

```shell
cdt configure

# Custom build directory
cdt --build build/debug configure
```

### Project compilation

Targets can be compiled by `build` command. Some other command might invoke this command automatically if needed.

```shell
# Build whole project
cdt build

# Build specific target
cdt build main
```

### Project information

For getting information about the current project the `project` command can be used.

```shell
# Obtain available targets
cdt project targets

# List all files grouped by targets
cdt project files
```

## Supported project types

- CMake
    - Clang Format
    - Clang Tidy
    - CppCheck
