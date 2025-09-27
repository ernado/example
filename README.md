# example

Example project with best practices.

## Installation

### atlas

```bash
curl -sSf https://atlasgo.sh | sh
```

## Structure


| File/Directory       | Description                                                         |
|----------------------|---------------------------------------------------------------------|
| `.github/`           | GitHub workflows and configurations                                 |
| `cmd/`               | Main applications                                                   |
| `internal/`          | Private application and library code. Most code **SHOULD** be here. |
| `pkg/`               | Directory that **MUST NOT** exist                                   |
| `.codecov.yml`       | Codecov configuration                                               |
| `.golangci.yml`      | GolangCI-Lint configuration                                         |
| `.editorconfig`      | Editor configuration                                                |
| `Dockerfile`         | Dockerfile for building the application                             |
| `go.mod`             | Go module definition. Tools are defined here.                       |
| `go.sum`             | Go module checksums                                                 |
| `Makefile`           | Makefile with common commands                                       |
| `README.md`          | This file                                                           |
| `LICENSE`            | License file                                                        |
| `go.coverage.sh`     | Script to generate coverage report                                  |
| `_oas`               | OpenAPI specifications                                              |
| `go.test.sh`         | Script to run tests                                                 |
| `generate.go`        | Code generation entrypoint                                          |
| `migrate.Dockerfile` | Docker file for ent migrations                                      |

### .github

#### Dependencies files

1. Dependabot configuration files with groups for otel and golang dependencies.
2. Dependency

#### Workflows

- Commit linting
- Dependency checks
- Linting
- Tests

##  _oas

OpenAPI specifications.

## generate.go

Code generation entrypoint.

## go.mod

Note that tools are defined here.
Example:

```
tool github.com/ogen-go/ogen/cmd/ogen
```

## internal

Most code SHOULD be here.

### ent

Ent ORM code.

Note `entc.go` and `generate.go` files.

#### schema

Ent schemas.
