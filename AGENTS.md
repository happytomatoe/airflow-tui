# Airflow-TUI AGENTS.md — high-signal dev guidance

## Build & run
- `just build` → `go build -o airflow-tui .` (binary output: `./airflow-tui` in repo root)
- `just run` → runs built binary; no auto-rebuild on code change
- `go run .` works but config handling differs — TUI creates `~/.config/flowrs/config.yaml` (note: `flowrs/` not `airflow-tui/`) if missing on first run
- `just e2e-test` builds first, then runs all Python e2e tests in parallel. Requires Airflow running.

## Test environment
- E2E tests: `e2e/test_tui.py` (Python/pytest) — **not Go tests**. No unit/integration tests in Go.
- Local Airflow for e2e: `docker-compose up -d` (Airflow 2.10.2 + PostgreSQL). Available at `http://localhost:8080`, user `airflow` / pass `airflow`.
- E2E venv: `just setup-tests` → `uv venv e2e/.venv && uv pip install pytest pytest-xdist requests pexpect -p e2e/.venv`. Tests use `pexpect` to drive the TUI; avoid changing terminal I/O behavior in UI code without updating tests.
- Run single e2e test: `just e2e-test testname` (matches `test_*` function name).

## Lint / formatting
- `go fmt ./...` (no CI script provided). No `go vet`/`staticcheck` in Justfile; run manually.
- Generated code (`airflow/generated/`) is excluded from manual edits — regenerate from `airflow-v2-openapi.yml` if API changes.