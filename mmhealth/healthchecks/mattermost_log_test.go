package healthchecks

import (
	"testing"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

func TestH003(t *testing.T) {
	p, checkStatus := setupTest(t, "mattermostLog")

	testCases := []struct {
		name           string
		logs           []types.MattermostLogEntry
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name: "h003 - logs contain 'context deadline exceeded'",
			logs: []types.MattermostLogEntry{
				{Msg: "context deadline exceeded"},
			},
			expectedStatus: Fail,
			expectedResult: "Found",
		},
		{
			name: "h003 - logs do not contain 'context deadline exceeded'",
			logs: []types.MattermostLogEntry{
				{Msg: ""},
			},
			expectedStatus: Pass,
			expectedResult: "Not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Logs = tc.logs
			checkStatus(t, p.h003, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH004(t *testing.T) {
	p, checkStatus := setupTest(t, "mattermostLog")

	testCases := []struct {
		name           string
		logs           []types.MattermostLogEntry
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name: "h004 - logs contain 'i/o timeout'",
			logs: []types.MattermostLogEntry{
				{Msg: "i/o timeout"},
			},
			expectedStatus: Fail,
			expectedResult: "Found",
		},
		{
			name: "h004 - logs do not contain 'i/o timeout'",
			logs: []types.MattermostLogEntry{
				{Msg: ""},
			},
			expectedStatus: Pass,
			expectedResult: "Not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Logs = tc.logs
			checkStatus(t, p.h004, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH005(t *testing.T) {
	p, checkStatus := setupTest(t, "mattermostLog")

	testCases := []struct {
		name           string
		logs           []types.MattermostLogEntry
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name: "h005 - logs contain 'Error while creating session for user access token'",
			logs: []types.MattermostLogEntry{
				{Msg: "Error while creating session for user access token"},
			},
			expectedStatus: Fail,
			expectedResult: "Found",
		},
		{
			name: "h005 - logs do not contain 'Error while creating session for user access token'",
			logs: []types.MattermostLogEntry{
				{Msg: ""},
			},
			expectedStatus: Pass,
			expectedResult: "Not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Logs = tc.logs
			checkStatus(t, p.h005, tc.expectedStatus, tc.expectedResult)
		})
	}
}
