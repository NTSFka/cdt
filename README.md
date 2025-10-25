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

The tool supports running tools in different environments like docker. More information is available in
the [documentation](doc/environment.md).

```shell
# Build and run main.go in a docker container
cdt -e docker:golang:1.24 run main.go
```

### Configuration

The tool tries to detect the best workflow from the project directory, but in some cases the default ones might not be
suitable for your project. The tool reads configuration file `cdt.yml` from root directory and use it as configuration.
More information is available in the [documentation](doc/configuration.md).

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
cdt --output build/debug configure
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

## Installation

The cdt can be installed using go tool:

```shell
go install github.com/NTSFka/cdt
```

### Auto-completion

The tool supports auto-completion for bash and zsh via [cli.urfave.org/v3](https://cli.urfave.org/v3/examples/completions/shell-completions/#default-auto-completion).

## Supported tools

The list of supported tools types is available in the [documentation](doc/tool.md).

## Supported workflow types

The list of supported workflow types is available in the [documentation](doc/workflow.md).


## Development

The cdt tool can be extended by additional tools or workflow types. More information is available in the
[developer documentation](doc/dev).
