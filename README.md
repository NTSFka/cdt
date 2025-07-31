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

### Environment

The tool supports running tools in different environments like docker.

```shell
# Build and run main.go in docker container
cdt --environment docker:golang:1.24 run main.go
```

### Configuration

The tool tries to detect the best workflow from the project directory, but in some cases the default ones might not be
suitable for your project. The tool reads configuration file `cdt.yml` from root directory and use it as configuration.

```yaml
project:
  workflow:
    build: make
    lint: clang-tidy
```

### Testing

```shell
# Run all tests
cdt test

# Run specific subset of tests
cdt test test1

# Use specific supported tool
cdt test --tool phpunit test1
```

### Code formatting

Tool allows reformatting the whole project or specific file by formatting the tool defined by the project workflow.

```shell
# Reformat whole project
cdt format

# Reformat file or multiple files
cdt format file1 file2

# Use specific supported tool
cdt format --tool clang-format test1
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

# Use specific supported tool
cdt lint --tool golangci-lint test1
```

### Project configuration

Some projects need to be configured (e.g., CMake). In most cases, it's handled by the tool, but if a specific
configuration is required, it's possible to invoke the configuration manually.

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

# Use specific supported tool
cdt build --tool make test1
```

### Project information

For getting information about the current project the `project` command can be used.

```shell
# Obtain available targets
cdt project targets

# List all files grouped by targets
cdt project files
```

## Supported tools

- [Clang Format](https://clang.llvm.org/docs/ClangFormat.html)
- [Clang Tidy](https://clang.llvm.org/extra/clang-tidy/)
- [CMake](https://cmake.org/)
- [CTest](https://cmake.org/)
- [Docker](https://www.docker.com/)
- [Go](https://go.dev/)
- [Golangci-lint](https://golangci-lint.run/)
- [Paratest](https://github.com/paratestphp/paratest)
- [PHP](https://www.php.net/)
- [PHP CS Fixer](https://github.com/PHP-CS-Fixer/PHP-CS-Fixer)
- [PHPStan](https://phpstan.org/)
- [PHPUnit](https://phpunit.de/)

## Supported project types

- CMake
- Go
- PHP

