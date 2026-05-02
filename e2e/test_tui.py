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

if __name__ == '__main__':
    pytest.main([__file__, "-v", "-n", "auto"])