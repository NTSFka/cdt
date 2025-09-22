# Environment

An environment allows running commands and tools in a specific context like in a docker container.
The environment is specified by the tool name and a parameter separated by double colon in
the following way: `<tool-id>:<parameter>`. Some tools might not require a parameter or can use the default one,
so only the tool id can be specified: `<tool-id>`.

```shell
# List available environments
cdt env list

# The environment status
cdt env status

# Start the environment
cdt env start

# Stop the environment
cdt env stop

# Use environment - it overrides detected or configured environment
cdt -e <tool-id>:<parameter> ...
```

The chosen environment is used for all commands, and when it is not already running, it will be used otherwise it
will be started and automatically stopped at the end of the command.

## Execute commands in the environment

When an unsupported tool is required to execute, it's possible to use the `cdt exec` command which invoke command in
selected environment.

```shell
# prints: go version go1.24.5 linux/amd64
cdt -e docker:golang:1.24 exec go version
```

## Supported environments

The following table lists the supported environments:

| Tool ID          | Aliases | Parameter                             |
|------------------|---------|---------------------------------------|
| `system`         | `s`     | -                                     |
| `docker`         | `d`     | image name                            |
| `docker-compose` | `dc`    | service name                          |
| `python`         | `pyenv` | path to virtual environment directory |

## Automatic environment detection

Some environments can be detected automatically from the project directory,
so no parameter or configuration is necessary. When the environment is detected,
the environment is used automatically. For cases when multiple environments are
detected, the first one is used.

| Tool ID  | Detection                               |
|----------|-----------------------------------------|
| `python` | presence of `.venv` or `venv` directory |
