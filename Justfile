# Just recipes for airflow-tui

# Default recipe - show help
default:
    @just --list

# Create virtual environment and install test dependencies
setup-tests:
    uv venv e2e/.venv --clear
    uv pip install pytest pytest-xdist requests pexpect -p e2e/.venv

# Run e2e tests in parallel (builds first)
e2e-test testname="": build
    #!/bin/bash
    set -e
    if [ -z "{{testname}}" ]; then
        e2e/.venv/bin/python -m pytest e2e/test_tui.py -v -n auto
    else
        e2e/.venv/bin/python -m pytest e2e/test_tui.py::test_{{testname}} -v
    fi

# Build the project
build:
    go build -o airflow-tui .

# Run Go tests
test:
    go test ./...

# Run the TUI
run:
    ./airflow-tui

generate-airflow-api:
  npx @openapitools/openapi-generator-cli generate -i airflow-v2-openapi.yml -g go -o airflow/api/generated
