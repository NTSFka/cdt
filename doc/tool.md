
# Tools

The tools are the most important part of the cdt tool.

## Supported tools

The full list of supported tools can be obtained by running the `tool list` command. By default, only the tools that
are available in the current environment are listed. It is possible to list a subset of the tools by using the `-a`
or `--tag` flag followed by the tag name, e.g. `go` only for tools for Go language.

```shell
# list tools that are available
cdt tool list

# list all supported tools
cdt tool list -a

# list all tools with tags (e.g. go)
cdt tool --tag go list
```

## Running a tool

The tool can be run with the `tool run` command followed by the tool id and the arguments.

```shell
# run a tool (e.g. go)
cdt tool run go version
```