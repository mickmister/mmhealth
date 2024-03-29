package healthchecks

import (
	"runtime"
	"strconv"
	"testing"

	"github.com/coltoneshaw/mmhealth/mmhealth"
	"github.com/coltoneshaw/mmhealth/mmhealth/files"
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/mattermost/mattermost/server/public/model"
)

func mockProcessPacket(t *testing.T) (*ProcessPacket, error) {
	p := &ProcessPacket{}
	p.log = mmhealth.HandleError

	p.packet.Config = model.Config{}
	p.packet.Config.SetDefaults()

	checks, err := files.ReadChecksFile()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	p.Checks = checks

	config, err := files.ReadConfigFile()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	p.Config = config

	return p, nil
}

func setupTest(t *testing.T, checkType string) (
	*ProcessPacket,
	func(t *testing.T, testFunc func(checks map[string]types.Check) CheckResult, expectedStatus types.CheckStatus, expectedResult string)) {
	p, err := mockProcessPacket(t)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	checks := map[string]types.Check{}

	if checkType == "" || (checkType != "config" && checkType != "environment" && checkType != "mattermostLog" && checkType != "packet") {
		t.Fatalf("checkType is incorrect")
	}

	switch checkType {
	case "config":
		checks = p.Checks.Config
	case "environment":
		checks = p.Checks.Environment
	case "mattermostLog":
		checks = p.Checks.MattermostLog
	case "packet":
		checks = p.Checks.Packet

	}

	checkStatus := func(t *testing.T, testFunc func(checks map[string]types.Check) CheckResult, expectedStatus types.CheckStatus, expectedResult string) {
		test := testFunc(checks)
		_, file, line, _ := runtime.Caller(1)

		if test.Status != expectedStatus {
			t.Errorf("Failed at %s:%d: Expected status '%v', got '%v'.", file, line, expectedStatus, test.Status)
		}

		if test.Result != expectedResult {
			t.Errorf("Failed at %s:%d: Expected result '%v', got '%v'.", file, line, expectedResult, test.Result)
		}
	}

	return p, checkStatus
}

func TestSortResults(t *testing.T) {
	p := &ProcessPacket{}
	mockResults := []CheckResult{
		{
			ID:       "2",
			Status:   Fail,
			Severity: types.High,
		},
		{
			ID:       "6",
			Status:   Pass,
			Severity: types.High,
		},
		{
			ID:       "5",
			Status:   Pass,
			Severity: types.Urgent,
		},
		{
			ID:       "1",
			Status:   Fail,
			Severity: types.Urgent,
		},
		{
			ID:       "7",
			Status:   Pass,
			Severity: types.Medium,
		},
		{
			ID:       "3",
			Status:   Fail,
			Severity: types.Medium,
		},
		{
			ID:       "8",
			Status:   Pass,
			Severity: types.Low,
		},
		{
			ID:       "4",
			Status:   Fail,
			Severity: types.Low,
		},
		{
			ID:       "0",
			Status:   Error,
			Severity: types.Low,
		},
	}

	testResults := p.sortResults(mockResults)

	for id, test := range testResults {
		if strconv.Itoa(id) != test.ID {
			t.Errorf("Failed at %d: Expected id '%v', got '%v'.", id, id, test.ID)
		}
	}
}
