package types

type MattermostLogEntry struct {
	Timestamp       string `json:"timestamp"`
	Level           string `json:"level"`
	Msg             string `json:"msg"`
	Caller          string `json:"caller"`
	HttpCode        int    `json:"http_code"`
	Status          string `json:"status"`
	StatusCode      string `json:"status_code"`
	SchedulerName   string `json:"scheduler_name"`
	Worker          string `json:"worker"`
	WorkerName      string `json:"worker_name"`
	MissingPluginId string `json:"missing_plugin_id"`
	PluginId        string `json:"plugin_id"`
	Error           string `json:"error"`
	Method          string `json:"method"`
	Action          string `json:"action"`
	JobId           string `json:"job_id"`
}
