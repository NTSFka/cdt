
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

### List of supported tools

The list might not be up to date.

- [Bandit](https://github.com/PyCQA/bandit)
- [Black](https://github.com/psf/black)
- [Clang Format](https://clang.llvm.org/docs/ClangFormat.html)
- [Clang Tidy](https://clang.llvm.org/extra/clang-tidy/)
- [CMake](https://cmake.org/)
- [Composer](https://getcomposer.org/)
- [CTest](https://cmake.org/)
- [Docker](https://www.docker.com/)
- [Flake8](https://github.com/PyCQA/flake8)
- [Go](https://go.dev/)
- [Golangci-lint](https://golangci-lint.run/)
- [NilAway](https://github.com/uber-go/nilaway)
- [Mypy](https://www.mypy-lang.org/)
- [Paratest](https://github.com/paratestphp/paratest)
- [PHP](https://www.php.net/)
- [PHP CS Fixer](https://github.com/PHP-CS-Fixer/PHP-CS-Fixer)
- [PHPStan](https://phpstan.org/)
- [PHPUnit](https://phpunit.de/)
- [PIP](https://pip.pypa.io)
- [Pylint](https://www.pylint.org/)
- [Pytest](https://docs.pytest.org)
- [Python](https://www.python.org/)
- [Ruff](https://docs.astral.sh/ruff/)

## Running a tool

The tool can be run with the `tool run` command followed by the tool id and the arguments.

```shell
# run a tool (e.g. go)
cdt tool run go version
```