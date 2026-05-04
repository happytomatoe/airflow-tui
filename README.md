# airflow-tui

A terminal user interface for Apache Airflow built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- Browse DAGs, DAG runs, and task instances
- View task logs with syntax highlighting
- Gantt chart visualization for task durations
- Search and filter DAGs
- Visual selection mode for bulk operations
- Configurable server connections

## Requirements

- Go 1.26+
- Airflow 2.10.2+ (API compatible)

## Installation

```bash
git clone https://github.com/airflow-tui/airflow-tui.git
cd airflow-tui
just build
```

## Configuration

On first run, a config file is created at `~/.config/airflow-tui/config.yaml`:

```yaml
servers: []
active_server: ""
poll_interval_ms: 2000
```

### Add a server

```bash
./airflow-tui config add local http://localhost:8080 -a basic -u airflow -p airflow
./airflow-tui config list
./airflow-tui config remove local
```

## Usage

```bash
./airflow-tui
```

### Key Bindings

| Key | Action |
|-----|--------|
| `q` | Quit |
| `?` | Toggle help |
| `/` | Search/filter |
| `V` | Visual selection mode |
| `1-5` | Switch tabs |
| `←/→` | Switch tabs |
| `g` | Toggle gantt chart |
| `j/k` | Navigate down/up |
| `G` / `gg` | Jump to bottom/top |

## Development

### Build

```bash
just build
```

### Run e2e tests (requires Airflow running)

```bash
docker-compose up -d  # Start Airflow
just setup-tests      # Install test dependencies
just e2e-test         # Run all tests
```

Or with an existing Airflow instance:

```bash
just e2e-test testname  # Run specific test
```

## License

MIT License - see [LICENSE](LICENSE).