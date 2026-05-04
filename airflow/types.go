package airflow

import generated "github.com/airflow-tui/airflow-tui/airflow/generated"

type DAG = generated.DAG
type DAGCollection = generated.DAGCollection
type DAGRun = generated.DAGRun
type DAGRunCollection = generated.DAGRunCollection
type TaskInstance = generated.TaskInstance
type TaskInstanceCollection = generated.TaskInstanceCollection
type DagStatsCollectionItem = generated.DagStatsCollectionItem
type DagStatsCollection = struct {
	Dags *[]DagStatsCollectionItem `json:"dags,omitempty"`
}
type Connection = generated.Connection
type ConnectionCollection = generated.ConnectionCollection
type Pool = generated.Pool
type PoolCollection = generated.PoolCollection
type Variable = generated.Variable
type VariableCollection = generated.VariableCollection
type EventLog = generated.EventLog
type EventLogCollection = generated.EventLogCollection
type ImportError = generated.ImportError
type ImportErrorCollection = generated.ImportErrorCollection
type Dataset = generated.Dataset
type DatasetCollection = generated.DatasetCollection
type DatasetEvent = generated.DatasetEvent
type DatasetEventCollection = generated.DatasetEventCollection
type DAGRunRunType = generated.DAGRunRunType
type DagState = generated.DagState
type TaskState = generated.TaskState
type HealthStatus = generated.HealthStatus
type TriggerRule = generated.TriggerRule
