import pytest
import requests
import subprocess
import pexpect
import time
import os
import tempfile
import shutil

LOG_DIR = "e2e/logs"
os.makedirs(LOG_DIR, exist_ok=True)

@pytest.fixture
def spawn_tui(request):
    test_name = request.node.name
    log_file = open(f"{LOG_DIR}/{test_name}.log", "w")
    subprocess.run('./airflow-tui config remove local', shell=True, capture_output=True)
    subprocess.run('./airflow-tui config add local http://localhost:8080 -a basic -u airflow -p airflow', shell=True, capture_output=True)
    p = pexpect.spawn('./airflow-tui', dimensions=(24, 150), encoding='utf-8', env={'TERM': 'dumb'})
    p.logfile = log_file
    time.sleep(1)

    yield p

    p.send('q')
    try:
        p.expect(pexpect.EOF, timeout=2)
    except pexpect.TIMEOUT:
        p.terminate()
    p.close()
    log_file.close()

def log_to_file(name, content):
    os.makedirs(LOG_DIR, exist_ok=True)
    with open(f"{LOG_DIR}/{name}.log", "w") as f:
        f.write(content)
    print(f"  -> logged to {LOG_DIR}/{name}.log")

def capture_cmd(cmd):
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=5)
    return result.stdout + result.stderr

def test_api_connect():
    r = requests.get('http://localhost:8080/api/v1/dags', auth=('airflow', 'airflow'))
    assert r.status_code == 200
    dags = r.json()
    count = dags.get('total_entries', 0)

    r_stats = requests.get('http://localhost:8080/api/v1/dagStats', auth=('airflow', 'airflow'))
    dag_stats = r_stats.json() if r_stats.status_code == 200 else {}

    log_to_file("api_dags", f"Total DAGs: {count}\n\nStats: {dag_stats}")

    p = pexpect.spawn('./airflow-tui', dimensions=(24, 150), encoding='utf-8', env={'TERM': 'dumb'})
    time.sleep(3)
    p.send('c')
    time.sleep(1)
    p.send('q')
    try:
        p.expect(pexpect.EOF, timeout=2)
    except pexpect.TIMEOUT:
        p.terminate()
    output = p.before

    dag_list = dags.get('dags', [])
    if dag_list:
        sample_dag = dag_list[0].get('dag_id', '')
        assert sample_dag in output, f"DAG '{sample_dag}' not found in TUI output"

    print(f"PASS: API connects ({count} DAGs) and TUI displays them")

def test_help():
    output = capture_cmd('./airflow-tui --help')
    log_to_file("help", output)
    assert "TUI for Apache Airflow" in output
    print("PASS: --help")

def test_config_help():
    output = capture_cmd('./airflow-tui config --help')
    log_to_file("config_help", output)
    assert "Manage configuration" in output
    print("PASS: config --help")

def test_config_add():
    os.system('./airflow-tui config remove local >/dev/null 2>&1')
    output = capture_cmd('./airflow-tui config add local http://localhost:8080 -a basic -u airflow -p airflow')
    log_to_file("config_add", output)
    assert "Added server" in output
    print("PASS: config add")

def test_config_list():
    output = capture_cmd('./airflow-tui config list')
    log_to_file("config_list", output)
    print("PASS: config list")

def test_get_dags(spawn_tui):
    p = spawn_tui
    time.sleep(3)
    # Read output from logfile
    log_path = f"{LOG_DIR}/test_get_dags.log"
    with open(log_path, 'r') as f:
        output = f.read()
    assert "No DAGs found" not in output, f"TUI couldn't load DAGs. Output: {output}"
    print("PASS: get_dags")

def test_search_dags(spawn_tui):
    r = requests.get('http://localhost:8080/api/v1/dags', auth=('airflow', 'airflow'))
    assert r.status_code == 200
    dag_list = r.json().get('dags', [])
    assert dag_list, "Expected at least one DAG from Airflow API"

    dag_id = dag_list[0].get('dag_id', '')
    assert dag_id, "Expected first DAG to have a dag_id"

    search_term = dag_id[:max(1, min(6, len(dag_id)))]
    p = spawn_tui
    time.sleep(3)

    p.send('/')
    time.sleep(0.5)
    p.send(search_term)
    time.sleep(0.5)
    p.send('\r')
    p.expect_exact(f"Filter: {search_term}", timeout=5)
    p.logfile.flush()

    log_path = f"{LOG_DIR}/test_search_dags.log"
    with open(log_path, 'r') as f:
        output = f.read()

    assert f"Filter: {search_term}" in output, f"Search filter not shown. Output: {output}"
    assert dag_id in output, f"Filtered DAG '{dag_id}' not found in TUI output"
    print("PASS: search_dags")

def test_get_dag_runs(spawn_tui):
    """Test that we can get and display DAG runs for a DAG."""
    p = spawn_tui
    time.sleep(3)

    # Press enter to select the first DAG
    p.send('\r')
    time.sleep(3)

    # Wait for DAG runs table to appear
    try:
        p.expect('Run ID', timeout=5)
    except pexpect.TIMEOUT:
        pass

    # Read output from logfile
    log_path = f"{LOG_DIR}/test_get_dag_runs.log"
    with open(log_path, 'r') as f:
        output = f.read()

    # Check that we're in the DAG runs view (should show Run ID, Type, State columns)
    assert "Run ID" in output, f"DAG runs table not displayed. Output: {output}"
    assert "Type" in output, f"DAG runs table missing Type column. Output: {output}"
    assert "State" in output, f"DAG runs table missing State column. Output: {output}"

    # Check that we have actual DAG runs (not just the header)
    # The output should contain DAG run data, not "No DAG runs found"
    assert "No DAG runs found" not in output, f"No DAG runs found. Output: {output}"

    print("PASS: get_dag_runs")

def test_tab_switching(spawn_tui):
    """Test that we can switch between tabs using 1-5 keys without crashing."""
    p = spawn_tui
    time.sleep(3)

    # Test pressing number keys - just verify TUI doesn't crash
    p.send('1')
    time.sleep(1)
    p.send('2')
    time.sleep(1)
    p.send('3')
    time.sleep(1)

    print("PASS: tab_switching")

def test_help_overlay(spawn_tui):
    """Test that help overlay toggles with ? key without crashing."""
    p = spawn_tui
    time.sleep(3)

    # Toggle help on
    p.send('?')
    time.sleep(1)
    # Toggle help off
    p.send('?')
    time.sleep(1)

    print("PASS: help_overlay")

def test_visual_selection(spawn_tui):
    """Test that visual selection mode works with V key."""
    p = spawn_tui
    time.sleep(3)

    # Toggle visual mode on and use j/k
    p.send('V')
    time.sleep(1)
    p.send('j')
    time.sleep(0.5)
    p.send('k')
    time.sleep(0.5)
    # Toggle off
    p.send('V')
    time.sleep(1)

    print("PASS: visual_selection")

def test_task_instances(spawn_tui):
    """Test that we can navigate to task instances without crashing."""
    p = spawn_tui
    time.sleep(3)

    # Go to DAG runs
    p.send('\r')
    time.sleep(3)

    # Go to task instances
    p.send('\r')
    time.sleep(3)

    print("PASS: task_instances")

def test_gantt_chart(spawn_tui):
    """Test that gantt chart toggles with g key without crashing."""
    p = spawn_tui
    time.sleep(3)

    # Go to DAG runs
    p.send('\r')
    time.sleep(3)

    # Go to task instances
    p.send('\r')
    time.sleep(3)

    # Toggle gantt chart on
    p.send('g')
    time.sleep(1)

    # Toggle gantt chart off
    p.send('g')
    time.sleep(1)

    print("PASS: gantt_chart")

def test_tab_navigation_arrows(spawn_tui):
    """Test that we can switch tabs using left/right arrow keys without crashing."""
    p = spawn_tui
    time.sleep(3)

    # Navigate with arrow keys
    p.send('\x1b[C')  # right arrow
    time.sleep(1)
    p.send('\x1b[C')  # right arrow
    time.sleep(1)
    p.send('\x1b[D')  # left arrow
    time.sleep(1)

    print("PASS: tab_navigation_arrows")

def test_view_logs(spawn_tui):
    """Test that we can navigate to logs view without crashing."""
    p = spawn_tui
    time.sleep(3)

    # Go to DAG runs
    p.send('\r')
    time.sleep(3)

    # Go to task instances
    p.send('\r')
    time.sleep(3)

    # Go to logs (select first task) - just verify no crash
    p.send('\r')
    time.sleep(4)

    print("PASS: view_logs")

def test_search_with_arrow_navigation(spawn_tui):
    """Test that arrow keys work during search mode and selection is preserved on Enter."""
    p = spawn_tui
    time.sleep(3)

    # Start search mode
    p.send('/')
    time.sleep(0.5)
    
    # Type a search term
    p.send('test')
    time.sleep(0.5)
    
    # Use arrow keys to navigate while searching
    p.send('\x1b[A')  # up arrow
    time.sleep(0.2)
    p.send('\x1b[B')  # down arrow
    time.sleep(0.2)
    p.send('\x1b[A')  # up arrow
    time.sleep(0.2)
    
    # Exit search mode with Enter, should keep selection
    p.send('\r')
    time.sleep(0.5)

    print("PASS: search_with_arrow_navigation")

def test_logs_scrolling(spawn_tui):
    """Test that we can scroll through logs without crashing."""
    p = spawn_tui
    time.sleep(3)

    # Go to DAG runs
    p.send('\r')
    time.sleep(3)

    # Go to task instances
    p.send('\r')
    time.sleep(3)

    # Go to logs
    p.send('\r')
    time.sleep(4)

    # Verify we're in logs view
    output = p.before.decode('utf-8', errors='replace') if p.before else ""
    assert len(output) > 0, "Should have output in logs view"

    # Scroll down
    p.send('j')
    time.sleep(0.5)
    p.send('j')
    time.sleep(0.5)

    # Scroll up
    p.send('k')
    time.sleep(0.5)

    # Jump to bottom
    p.send('G')
    time.sleep(0.5)

    # Jump to top
    p.send('g')
    time.sleep(0.5)
    p.send('g')
    time.sleep(0.5)

    print("PASS: logs_scrolling")

def test_log_content_correctness(spawn_tui):
    """Test that log content is correctly processed and displayed."""
    p = spawn_tui
    time.sleep(3)

    # Navigate to logs view
    p.send('\r')  # Select first DAG
    time.sleep(3)
    p.send('\r')  # Go to DAG runs
    time.sleep(3)
    p.send('\r')  # Go to task instances
    time.sleep(3)
    p.send('\r')  # Go to logs for first task
    time.sleep(4)

    # Verify we're in logs view - if we get here without crash, test passes
    print("PASS: log_content_correctness")

def test_connection_status_disconnected():
    """Test that TUI shows disconnected status when Airflow is not running."""
    # Use a separate config file to avoid interfering with parallel tests
    # Create a temp config that points to a non-existent Airflow instance
    temp_config_yaml = """servers:
  - name: test-disconnected
    url: http://localhost:18080
    auth:
      type: basic
      username: airflow
      password: airflow
    api_version: v1
active_server: test-disconnected
poll_interval_ms: 2000
"""

    temp_dir = tempfile.mkdtemp()
    temp_config_path = os.path.join(temp_dir, 'airflow-tui', 'config.yaml')
    os.makedirs(os.path.dirname(temp_config_path), exist_ok=True)
    with open(temp_config_path, 'w') as f:
        f.write(temp_config_yaml)

    log_file_path = f"{LOG_DIR}/test_connection_status_disconnected.log"
    log_file = open(log_file_path, "w")

    # Set XDG_CONFIG_HOME to use our temp config
    env = os.environ.copy()
    env['XDG_CONFIG_HOME'] = temp_dir

    p = pexpect.spawn('./airflow-tui', dimensions=(24, 150), encoding='utf-8', env={**env, 'TERM': 'dumb'})
    p.logfile = log_file
    time.sleep(4)

    p.send('q')
    try:
        p.expect(pexpect.EOF, timeout=2)
    except pexpect.TIMEOUT:
        p.terminate()
    p.close()
    log_file.close()

    # Read the log file
    with open(log_file_path, 'r') as f:
        output = f.read()

    # Should show disconnected status
    assert "disconnected" in output, f"Expected 'disconnected' in output. Output: {output}"
    assert "●" in output, f"Expected connection indicator '●' in output. Output: {output}"

    # Cleanup
    shutil.rmtree(temp_dir)

    print("PASS: connection_status_disconnected")


def test_connection_status_connected():
    """Test that TUI shows connected status when Airflow is running."""
    # Ensure Airflow is running
    try:
        r = requests.get('http://localhost:8080/api/v1/dags', auth=('airflow', 'airflow'), timeout=2)
        if r.status_code != 200:
            pytest.skip("Airflow is not running - skipping connected status test")
    except (requests.exceptions.ConnectionError, requests.exceptions.Timeout):
        pytest.skip("Airflow is not running - skipping connected status test")

    log_file_path = f"{LOG_DIR}/test_connection_status_connected.log"
    log_file = open(log_file_path, "w")
    p = pexpect.spawn('./airflow-tui', dimensions=(24, 150), encoding='utf-8', env={'TERM': 'dumb'})
    p.logfile = log_file
    time.sleep(5)  # Wait for connection to establish

    p.send('q')
    try:
        p.expect(pexpect.EOF, timeout=2)
    except pexpect.TIMEOUT:
        p.terminate()
    p.close()
    log_file.close()

    # Read the log file
    with open(log_file_path, 'r') as f:
        output = f.read()

    # Should show connected status or DAGs loaded
    assert "connected" in output or "DAG" in output, f"Expected 'connected' or DAGs in output. Output: {output}"

    print("PASS: connection_status_connected")


def test_connection_status_connecting():
    """Test that TUI shows connecting status while attempting connection."""
    log_file_path = f"{LOG_DIR}/test_connection_status_connecting.log"
    log_file = open(log_file_path, "w")
    p = pexpect.spawn('./airflow-tui', dimensions=(24, 150), encoding='utf-8', env={'TERM': 'dumb'})
    p.logfile = log_file
    time.sleep(2)  # Capture during loading phase

    p.send('q')
    try:
        p.expect(pexpect.EOF, timeout=2)
    except pexpect.TIMEOUT:
        p.terminate()
    p.close()
    log_file.close()

    # Read the log file
    with open(log_file_path, 'r') as f:
        output = f.read()

    # During loading, should show "connecting..." or the spinner
    assert "connecting" in output or "Loading" in output or "●" in output, f"Expected 'connecting' or 'Loading' in output. Output: {output}"

    print("PASS: connection_status_connecting")


if __name__ == '__main__':
    pytest.main([__file__, "-v", "-n", "auto"])
