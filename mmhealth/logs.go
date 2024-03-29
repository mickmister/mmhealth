package mmhealth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

// ParseLogs determines the log format and parses accordingly
func ParseLogs(logData []byte) ([]types.MattermostLogEntry, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(logData))
	var parsedLogs []types.MattermostLogEntry
	var err error

	if scanner.Scan() {
		firstLine := scanner.Bytes()
		var entry types.MattermostLogEntry

		// Try to unmarshal the first line as JSON
		if json.Unmarshal(firstLine, &entry) == nil {
			parsedLogs, err = ReadJSONLog(logData)
		} else {
			parsedLogs, err = ReadTextLog(logData)
		}
	}

	if err != nil {
		return nil, err
	}

	return parsedLogs, nil
}

// ReadJSONLog reads and parses JSON formatted logs
func ReadJSONLog(logData []byte) ([]types.MattermostLogEntry, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(logData))
	var logEntries []types.MattermostLogEntry

	for scanner.Scan() {
		line := scanner.Bytes()
		var entry types.MattermostLogEntry
		err := json.Unmarshal(line, &entry)
		if err != nil {
			return nil, err
		}
		logEntries = append(logEntries, entry)
	}

	return logEntries, nil
}

func ReadTextLog(logData []byte) ([]types.MattermostLogEntry, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(logData))
	var logEntries []types.MattermostLogEntry

	setters := createSetters()
	var lastEntry *types.MattermostLogEntry

	for scanner.Scan() {
		line := scanner.Text()
		entry, isComplete, err := parseLogLine(line, setters)
		if err != nil {
			return nil, err
		}

		if isComplete {
			logEntries = append(logEntries, entry)
			lastEntry = &logEntries[len(logEntries)-1]
		} else if lastEntry != nil {
			lastEntry.Msg += "\n" + line
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logEntries, nil
}

func parseLogLine(line string, setters map[string]func(*types.MattermostLogEntry, string)) (types.MattermostLogEntry, bool, error) {
	logRegex := regexp.MustCompile(`^(\w+)\s+\[([^\]]+)\]\s+(.*?)\s+(\w+=.*)?$`)
	matches := logRegex.FindStringSubmatch(line)

	if len(matches) < 4 {
		return types.MattermostLogEntry{}, false, nil
	}

	entry := types.MattermostLogEntry{
		Level:     matches[1],
		Timestamp: matches[2],
		Msg:       matches[3],
	}

	if len(matches) == 5 && matches[4] != "" {
		err := parseKeyValuePairs(matches[4], setters, &entry)
		if err != nil {
			return types.MattermostLogEntry{}, false, err
		}
	}

	return entry, true, nil
}

func parseKeyValuePairs(kvPairs string, setters map[string]func(*types.MattermostLogEntry, string), entry *types.MattermostLogEntry) error {
	partsRegex := regexp.MustCompile(`(\w+)="([^"]*)"`)
	parts := partsRegex.FindAllStringSubmatch(kvPairs, -1)

	for _, part := range parts {
		key := part[1]
		value := part[2]
		if setter, ok := setters[key]; ok {
			setter(entry, value)
		}
	}

	return nil
}

func createSetters() map[string]func(*types.MattermostLogEntry, string) {
	return map[string]func(*types.MattermostLogEntry, string){
		"caller": func(e *types.MattermostLogEntry, v string) { e.Caller = v },
		"http_code": func(e *types.MattermostLogEntry, v string) {
			if val, err := strconv.Atoi(v); err == nil {
				e.HttpCode = val
			}
		},
		"status":            func(e *types.MattermostLogEntry, v string) { e.Status = v },
		"status_code":       func(e *types.MattermostLogEntry, v string) { e.StatusCode = v },
		"scheduler_name":    func(e *types.MattermostLogEntry, v string) { e.SchedulerName = v },
		"worker":            func(e *types.MattermostLogEntry, v string) { e.Worker = v },
		"worker_name":       func(e *types.MattermostLogEntry, v string) { e.WorkerName = v },
		"missing_plugin_id": func(e *types.MattermostLogEntry, v string) { e.MissingPluginId = v },
		"plugin_id":         func(e *types.MattermostLogEntry, v string) { e.PluginId = v },
		"error":             func(e *types.MattermostLogEntry, v string) { e.Error = v },
		"method":            func(e *types.MattermostLogEntry, v string) { e.Method = v },
		"action":            func(e *types.MattermostLogEntry, v string) { e.Action = v },
		"job_id":            func(e *types.MattermostLogEntry, v string) { e.JobId = v },
	}
}
