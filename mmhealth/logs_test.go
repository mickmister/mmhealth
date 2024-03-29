package mmhealth

import (
	"fmt"
	"testing"
)

func TestParseLogs(t *testing.T) {
	// Test data
	jsonLog := []byte(`{"timestamp":"2024-01-01T00:00:00Z","level":"debug","msg":"test message","caller":"main.go:1"}`)
	textLog := []byte(`debug [2024-01-01 00:00:00.000 Z] test message                   caller="main.go:1"`)

	// Test JSON log
	logs, err := ParseLogs(jsonLog)
	if err != nil {
		t.Errorf("ParseLogs failed with error: %v", err)
	}
	if len(logs) != 1 || logs[0].Msg != "test message" {
		t.Errorf("ParseLogs did not correctly parse JSON log")
	}

	// Test text log
	logs, err = ParseLogs(textLog)
	if err != nil {
		t.Errorf("ParseLogs failed with error: %v", err)
	}
	if len(logs) != 1 || logs[0].Msg != "test message" {
		t.Errorf("ParseLogs did not correctly parse text log")
	}
}

func TestReadJSONLog(t *testing.T) {
	logData := []byte(`{"timestamp":"2024-01-01T00:00:00Z","level":"debug","msg":"test message","caller":"main.go:1"}`)

	logs, err := ReadJSONLog(logData)
	if err != nil {
		t.Errorf("ReadJSONLog failed with error: %v", err)
	}
	if len(logs) != 1 || logs[0].Msg != "test message" {
		t.Errorf("ReadJSONLog did not correctly parse log")
	}
}

func TestReadTextLog(t *testing.T) {
	// Test data
	logData := []byte(`debug [2024-01-01 17:46:19.649 Z] test message                   caller="main.go:1"`)

	logs, err := ReadTextLog(logData)
	if err != nil {
		t.Errorf("ReadTextLog failed with error: %v", err)
	}

	fmt.Println(logs)
	if len(logs) != 1 || logs[0].Msg != "test message" {
		t.Errorf("ReadTextLog did not correctly parse log")
	}
}
