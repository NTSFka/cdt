
# Environment and environment provider

The tool that provides the environment must implement the `EnvironmentProvider` interface. The most important function
is `CreateEnvironment` which creates the actual environment in which the tools are run. The created environment must
implement the `Environment` interface.

## Environment

The environment instance must provide a way to start and stop the environment, find an executable in the environment,
and execute a command in the environment.

## Example

It's not simple to provide an example, better to see the implementation for docker, docker-compose or python.
