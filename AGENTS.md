# Airflow-TUI Onboarding Guide

Welcome to **airflow-tui**, a Terminal User Interface (TUI) for Apache Airflow. This guide will help you understand the project structure, set up your development environment, and get started with building and testing.

## Project Overview

**airflow-tui** is a Go-based terminal user interface application that provides an interactive way to manage and monitor Apache Airflow instances. It enables users to:

- Connect to multiple Airflow servers (v1 and v2 APIs)
- Browse and filter DAGs (Directed Acyclic Graphs)
- Monitor DAG runs and task instances
- View task logs
- Manage authentication with basic auth or token-based auth

### Key Technologies

- **Language**: Go 1.26+
- **UI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) - A Go framework for building terminal UIs
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) - A Go library for styling terminal output
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra) - A Go library for building command-line applications
- **Configuration**: Viper - Configuration management with support for YAML, JSON, and environment variables
- **Testing**: Python pytest for end-to-end testing

## Project Structure

```
airflow-tui/
├── main.go                 # Entry point for the application
├── cmd/                    # Command-line interface commands
│   ├── root.go            # Root command and TUI runner
│   └── config.go          # Configuration-related commands
├── ui/                     # User interface components
│   ├── model.go           # Bubble Tea model and state management
│   ├── theme.go           # UI theming and styling
│   └── ...                # Other UI components
├── airflow/               # Airflow API client code
│   └── ...                # API clients for v1 and v2
├── config/                # Configuration handling
│   └── config.go          # Config file management and paths
├── e2e/                   # End-to-end tests
│   ├── test_tui.py        # Python e2e test suite
│   └── .venv/             # Python virtual environment
├── dags/                  # Sample DAGs for development/testing
├── docker-compose.yaml    # Docker Compose for local Airflow setup
├── Dockerfile             # Docker image for Airflow server
├── go.mod / go.sum        # Go module dependencies
├── Justfile               # Build and test recipes
└── airflow-v2-openapi.yml # OpenAPI spec for Airflow v2
```

## Development Setup

### Prerequisites

- **Go 1.26+**: [Install Go](https://golang.org/doc/install)
- **uv** (Python package manager): [Install uv](https://github.com/astral-sh/uv)
- **Docker & Docker Compose** (optional, for running Airflow locally)
- **just** (optional, for running Justfile recipes): [Install just](https://github.com/casey/just)

### Initial Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/airflow-tui/airflow-tui.git
   cd airflow-tui
   ```

2. **Install Go dependencies**:
   ```bash
   go mod download
   ```

3. **Set up test environment** (for e2e tests):
   ```bash
   just setup-tests
   ```
   Or manually:
   ```bash
   uv venv e2e/.venv --clear
   uv pip install pytest pytest-xdist requests pexpect -p e2e/.venv
   ```

## Building

### Build the binary

```bash
# Using just
just build

# Or directly with Go
go build -o airflow-tui .
```

This creates an executable named `airflow-tui` in the project root.

### Run the application

```bash
# Using just
just run

# Or directly
./airflow-tui
```

## Testing

### Run E2E Tests

```bash
# Run all e2e tests
just e2e-test

# Run a specific test (replaces 'testname' with test name)
just e2e-test testname
```

The e2e tests are located in `e2e/test_tui.py` and are written in Python using pytest.

### Prerequisites for E2E Tests

The tests require a running Airflow instance. Set up a local instance with:

```bash
docker-compose up -d
```

This starts an Airflow server (based on Apache Airflow 2.10.2) with PostgreSQL backend.

## Configuration

### Config File Location

- **Linux/macOS**: `~/.config/airflow-tui/config.yaml`
- **Windows**: `%APPDATA%\airflow-tui\config.yaml`

### Config File Format

```yaml
servers:
  - name: "Local Airflow"
    url: "http://localhost:8080"
    api_version: "v2"
    auth:
      type: "basic"
      username: "airflow"
      password: "airflow"

active_server: "Local Airflow"
poll_interval_ms: 2000
```

### Supported Authentication Types

- **basic**: Username and password authentication
- **token**: Bearer token authentication

## Architecture Overview

### Model-View-Update (MVU) Pattern

The application uses Bubble Tea's MVU pattern for state management:

1. **Model** (`ui/model.go`): Holds the application state
2. **View**: Renders the current state to the terminal
3. **Update**: Processes user input and messages, updating the state

### Panel System

The TUI is organized into panels:

- **PanelConfig**: Server and configuration management
- **PanelDAG**: List of DAGs
- **PanelDAGRun**: DAG runs for a selected DAG
- **PanelTaskInstance**: Tasks within a DAG run
- **PanelLogs**: Task execution logs

### API Client

The `airflow/` package provides clients for both Airflow v1 and v2 APIs. Each environment (server) gets its own client instance that handles API communication and authentication.

## Common Development Tasks

### Adding a New Feature

1. Identify which component needs modification (UI, API client, config, etc.)
2. Make changes to the relevant files
3. Test locally with `./airflow-tui`
4. Add e2e tests in `e2e/test_tui.py` if needed
5. Run full test suite: `just e2e-test`

### Debugging

1. Set up a local Airflow instance: `docker-compose up -d`
2. Configure the TUI to connect to it
3. Run with debug output or add `fmt.Println()` statements in the code
4. Rebuild and test: `just build && just run`

### Adding Dependencies

```bash
# Add a new Go dependency
go get github.com/user/package

# Update the Justfile/e2e tests if adding Python dependencies
```

## Local Airflow Setup

A `docker-compose.yaml` is included for easy local testing:

```bash
# Start Airflow and PostgreSQL
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f airflow-webserver
```

The local Airflow instance will be available at `http://localhost:8080` with:
- Username: `airflow`
- Password: `airflow`

## Key Files for Different Tasks

| Task | Primary Files |
|------|----------------|
| Add UI component | `ui/model.go`, `ui/theme.go` |
| Add Airflow API support | `airflow/*.go` |
| Add configuration option | `config/config.go`, `cmd/config.go` |
| Add authentication method | `airflow/*.go` (auth providers) |
| Add e2e test | `e2e/test_tui.py` |

## Useful Commands

```bash
# Format Go code
go fmt ./...

# Run Go linter
go vet ./...

# Check for dependencies
go mod tidy

# View test coverage
go test -cover ./...

# Run e2e tests with verbose output
e2e/.venv/bin/python -m pytest e2e/test_tui.py -v -s
```

## Troubleshooting

### Build Fails

- Ensure Go 1.26+ is installed: `go version`
- Run `go mod download` to fetch dependencies
- Check for typos in file paths (Go is case-sensitive)

### Tests Fail

- Ensure Airflow is running: `docker-compose ps`
- Check Airflow logs: `docker-compose logs airflow-webserver`
- Verify network connectivity to Airflow server
- Ensure e2e virtual environment is set up: `just setup-tests`

### Config File Issues

- Check file permissions on config file
- Verify YAML syntax (indentation matters!)
- Ensure server URLs are valid and accessible

## Next Steps

1. Read through `cmd/root.go` to understand how the TUI starts
2. Explore `ui/model.go` to see the state management
3. Check `airflow/` to understand API integration
4. Review `e2e/test_tui.py` to see test patterns
5. Run the application and explore the UI: `just run`

## Resources

- [Bubble Tea Documentation](https://pkg.go.dev/github.com/charmbracelet/bubbletea)
- [Apache Airflow API Docs](https://airflow.apache.org/docs/apache-airflow/stable/stable-rest-api.html)
- [Go Documentation](https://golang.org/doc/)
- [Cobra Documentation](https://cobra.dev/)

---

For questions, issues, or contributions, please refer to the project's GitHub repository.
