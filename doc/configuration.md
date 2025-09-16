
# Configuration

The tool behavior can be changed by the configuration file or configuration options.

## Configuration options

```shell
# Use custom configuration file
cds --config cds.yml

# Use specific workflow
cds --workflow go
```

## Configuration file

The configuration file (`cds.yml`) can be placed in the root of the project and the tool will automatically load it.

### Workflow

Workflow configuration allows you to specify the workflow to use. It can be a name of a predefined workflow:

```yaml
project:
  # Use PHP workflow
  workflow: php
```

When a predefined workflow is not enough, you can specify a full workflow specification in the form of commands and
tools to use. The specified tools must support corresponding tool functionality to be really used, otherwise the command
fails.

```yaml
project:
  # Use custom workflow
  workflow:
    configure: make
    build: make
    test: make
    format: clang-format
    lint: cppcheck
    run: make
    dependency: conan
```
