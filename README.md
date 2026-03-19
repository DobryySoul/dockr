[![CI & Release](https://github.com/DobryySoul/dockr/actions/workflows/ci.yaml/badge.svg?event=release)](https://github.com/DobryySoul/dockr/actions/workflows/ci.yaml)
# Dockr 🐳

**Dockr** is a smart CLI utility for safely cleaning up unused Docker resources (images, containers, volumes, and networks). The tool is written in Go and provides a user-friendly command-line interface to keep your host machine clean.

## Features

- **Smart Analysis**: Finds orphaned images, exited containers, and unused volumes/networks.
- **Safe Deletion**: Supports interactive mode (`-i`) to prompt for confirmation before cleaning up.
- **Exceptions**: Ability to protect specific images from deletion by their tags (`-e`).
- **Dry-Run Mode**: Allows you to view a report of what would be deleted without actually making changes to the system (`-d`).
- **Informative**: Colored and structured table output with a calculation of freed disk space.

## Installation

**Linux / macOS (Quick Install):**
```bash
curl -sSfL https://raw.githubusercontent.com/DobryySoul/dockr/main/scripts/install.sh | bash
```

**Using Go:**
```bash
go install github.com/DobryySoul/dockr@latest
```

Alternatively, you can build the project yourself if you have Go installed:

```bash
# Clone the repository
git clone https://github.com/DobryySoul/dockr.git
cd dockr

# Build the project
go build -o dockr main.go
```

## Usage

Simply run the utility from the command line:

```bash
dockr [flags]
```

### Available Flags:
- `-d, --dry-run` — Simulation mode: prints information about resources that would be deleted, without actually removing them.
- `-i, --interactive` — Interactive mode: asks for user confirmation before deleting resources.
- `-e, --exclude-tags` — Exclude specific image tags from deletion (can be specified multiple times, e.g., `-e latest -e prod`).
- `-a, --all` — Delete ALL unused resources (including potentially important ones).
- `-v, --version` — Show the current application version.

## Uninstallation

If you used the installation script (`install.sh`), remove the binary:
```bash
sudo rm /usr/local/bin/dockr
```

If you installed it via `go install`, remove it from your Go bin path:
```bash
rm $(go env GOPATH)/bin/dockr
```

## Testing

The project is covered by unit tests to verify the correctness of the business logic (rules for determining if resources are "unused"). To run the tests, execute the following command:

```bash
go test -v ./...
```

## Project Structure

The codebase is built with clean architecture principles and standard Go project conventions in mind. Purpose of main directories:

```text
.
├── cmd/                # CLI commands (based on Cobra). Initialization and flag setup
│   └── root.go         # Root command 'dockr'
├── internal/           # Internal application business logic (cannot be imported externally)
│   ├── analyzer/       # Analysis logic: determining if a resource is used or can be deleted
│   ├── cleaner/        # Methods for actually deleting objects from Docker
│   ├── docker/         # Docker SDK wrapper, methods for interacting with Docker Daemon
│   └── domain/         # Core data structures and models (e.g., UnusedResources)
├── pkg/                # Public packages (potentially reusable)
│   └── formatter/      # Output formatting utilities (tables, colored text, calculations)
├── scripts/            # Helper bash scripts (e.g., install.sh)
├── Makefile            # Automation commands (build, test, linters, etc.)
└── main.go             # Application entry point
```

## Dependencies

The project relies on reliable open-source solutions:
- [Cobra](https://github.com/spf13/cobra) — A framework for creating powerful CLI applications.
- [Docker Engine API / moby](https://github.com/moby/moby) — The official Go client for interacting with the Docker API.
- [fatih/color](https://github.com/fatih/color) — A handy package for formatting and printing colored text to the console.

