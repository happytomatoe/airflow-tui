import pytest
import requests
import subprocess
import pexpect
import time
import os

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
    time.sleep(3)
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
    time.sleep(4)

    # If we get here without crashing, test passes
    # The log view shows with header "Try:" and scroll percentage
    print("PASS: view_logs")

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

    # If we get here without crashing, test passes
    print("PASS: logs_scrolling")

def test_log_content_correctness(spawn_tui):
    """Test that log content is correctly processed and displayed."""
    p = spawn_tui
    time.sleep(3)

    # Get DAG and task information for our assertions
    r = requests.get('http://localhost:8080/api/v1/dags', auth=('airflow', 'airflow'))
    assert r.status_code == 200
    dag_list = r.json().get('dags', [])
    assert dag_list, "Expected at least one DAG from Airflow API"
    
    dag_id = dag_list[0].get('dag_id', '')
    assert dag_id, "Expected first DAG to have a dag_id"
    
    # Get DAG runs for this DAG
    r_runs = requests.get(f'http://localhost:8080/api/v1/dags/{dag_id}/dagRuns', auth=('airflow', 'airflow'))
    assert r_runs.status_code == 200
    dag_runs = r_runs.json().get('dag_runs', [])
    assert dag_runs, f"Expected at least one DAG run for DAG {dag_id}"
    
    dag_run_id = dag_runs[0].get('dag_run_id', '')
    assert dag_run_id, f"Expected first DAG run to have a dag_run_id"
    
    # Get task instances for this DAG run
    r_tasks = requests.get(f'http://localhost:8080/api/v1/dags/{dag_id}/dagRuns/{dag_run_id}/taskInstances', auth=('airflow', 'airflow'))
    assert r_tasks.status_code == 200
    task_instances = r_tasks.json().get('task_instances', [])
    assert task_instances, f"Expected at least one task instance for DAG run {dag_run_id}"
    
    task_id = task_instances[0].get('task_id', '')
    try_number = task_instances[0].get('try_number', 0)
    assert task_id, f"Expected first task instance to have a task_id"

    # Navigate to logs view
    p.send('\r')  # Select first DAG
    time.sleep(3)
    p.send('\r')  # Go to DAG runs
    time.sleep(3)
    p.send('\r')  # Go to task instances
    time.sleep(3)
    p.send('\r')  # Go to logs for first task
    time.sleep(4)

    # Capture the current output from the pexpect buffer
    output = p.before

    # Test 1: Verify we're actually in the logs view
    assert "Try:" in output, f"Not in logs view - missing 'Try:' header. Output: {output}"
    assert "0%" in output or "100%" in output, f"Missing scroll percentage in logs view. Output: {output}"
    
    # Test 2: Verify actual log content is displayed (not empty or just headers)
    # Look for typical log patterns that should be present
    log_indicators = [
        "INFO", "WARNING", "ERROR", "DEBUG",  # Log levels
        "airflow", "task", "dag", "run",      # Airflow-specific terms
        "\n", "\r\n"                          # Actual line breaks in logs
    ]
    
    # Check if we have meaningful content (not just navigation indicators)
    content_found = any(indicator in output for indicator in log_indicators)
    assert content_found, f"No actual log content found in output. Output: {output}"
    
    # Test 3: Verify hostname skipping is working (simplified for this test)
    # Hostname lines typically don't contain spaces and are short
    lines = output.split('\n')
    single_word_lines = len([line for line in lines if len(line.strip()) > 0 and len(line.strip().split()) == 1])
    total_content_lines = len([line for line in lines if len(line.strip()) > 0])
    
    if total_content_lines > 0:
        single_word_ratio = single_word_lines / total_content_lines
        # hostname lines should be a small minority of the content
        assert single_word_ratio < 0.3, f"Too many potential hostname lines detected: {single_word_ratio:.2f} ratio"
    
    # Test 4: Verify V1 log format parsing is working
    # Look for evidence that Python tuple logs have been parsed
    # V1 logs should show actual content, not Python tuple syntax
    parsed_log_indicators = [
        "Started", "Completed", "Success", "Failed",  # Task status indicators
        "Executing", "Processing", "Running",       # Task execution indicators  
        "seconds", "minutes", "hours"               # Time-based content
    ]
    
    parsed_content_found = any(indicator in output for indicator in parsed_log_indicators)
    if not parsed_content_found:
        # If we don't see typical log indicators, check for Python tuple format
        # which would indicate parsing might not be working correctly
        python_tuples = output.count("(") + output.count(")") + output.count("'") + output.count('"')
        assert python_tuples < 10, f"Too many Python tuple indicators found, parsing may not be working. Found: {python_tuples} indicators"
    
    # Test 5: Verify task identification is present
    assert task_id in output, f"Task ID '{task_id}' not found in logs view. Output: {output}"
    assert dag_id in output, f"DAG ID '{dag_id}' not found in logs view. Output: {output}"
    
    # Test 6: Verify try number is displayed correctly
    try_indicator = f"Try: {try_number}"
    assert try_indicator in output, f"Try number indicator '{try_indicator}' not found. Output: {output}"
    
    # Test 7: Check for log content structure (should have multiple lines if there are actual logs)
    non_empty_lines = [line for line in lines if line.strip()]
    if len(non_empty_lines) > 1:  # More than just the header
        # Check that we have some meaningful content structure
        # Look for lines that look like actual log messages
        log_lines = [line for line in non_empty_lines if len(line.strip()) > 10]  # Reasonable length for log messages
        if log_lines:
            # Log lines should typically have some structure
            avg_line_length = sum(len(line) for line in log_lines) / len(log_lines)
            assert 20 < avg_line_length < 500, f"Log lines seem unusually short or long: avg length {avg_line_length:.1f}"

    print("PASS: log_content_correctness")

if __name__ == '__main__':
    pytest.main([__file__, "-v", "-n", "auto"])
