# Workflow

Workflow captures called programs for a project. The tool has some predefined workflows, or a custom workflow for
the project can be defined via a [configuration file](configuration.md).

## Predefined workflows

Named workflows that are predefined by the tool. They can be used via configuration file or as CLI parameter: `--workflow <ID>`.

| ID     | Name   |
|--------|--------|
| cmake  | CMake  |
| go     | Go     |
| php    | PHP    |
| python | Python |

The full list can be obtained via `workflow list` command.
