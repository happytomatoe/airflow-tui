package nav

// Navigation holds the current navigation state for the TUI.
// Fields are exported to allow access from both main ui/ and screen packages.
type Navigation struct {
	Server  string // Airflow server name
	Dag     string // Selected DAG ID
	DagRun  string // Selected DAG Run ID
	Task    string // Selected Task ID
	LogType string // Log type (stdout, stderr, etc.)
}
