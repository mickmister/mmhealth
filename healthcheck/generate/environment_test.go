package processpacket

import (
	"testing"
)

func TestH006(t *testing.T) {
	p, checkStatus := setupTest(t, "environment")

	p.Config.Versions.Supported = []string{"6.0.0", "6.1.0", "6.2.0"}
	p.Config.Versions.ESR = "5.7.0"

	testCases := []struct {
		name           string
		serverVersion  string
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "h006 - server version if out of date",
			serverVersion:  "5.0.0",
			expectedStatus: Fail,
			expectedResult: "Unsupported version: 5.0.0",
		},
		{
			name:           "h006 - server version in support",
			serverVersion:  p.Config.Versions.Supported[0],
			expectedStatus: Pass,
			expectedResult: "Supported version: 6.0.0",
		},
		{
			name:           "h006 - server version in support and ESR",
			serverVersion:  p.Config.Versions.ESR,
			expectedStatus: Warn,
			expectedResult: "Supported version: 5.7.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Packet.ServerVersion = tc.serverVersion
			checkStatus(t, p.h006, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH007(t *testing.T) {
	p, checkStatus := setupTest(t, "environment")

	testCases := []struct {
		name           string
		databaseType   string
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "h007 - Database type is not postgres",
			databaseType:   "mysql",
			expectedStatus: Fail,
			expectedResult: "mysql",
		},
		{
			name:           "h007 - Database type is postgres",
			databaseType:   "postgres",
			expectedStatus: Pass,
			expectedResult: "Postgres",
		},
		{
			name:           "h007 - Database type is maria",
			databaseType:   "maria",
			expectedStatus: Fail,
			expectedResult: "maria",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Packet.DatabaseType = tc.databaseType
			checkStatus(t, p.h007, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH008(t *testing.T) {
	p, checkStatus := setupTest(t, "environment")

	testCases := []struct {
		name           string
		license        string
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "h008 - Server os not linux",
			license:        "darwin",
			expectedStatus: Fail,
			expectedResult: "darwin",
		},
		{
			name:           "h008 - server is linux",
			license:        "linux",
			expectedStatus: Pass,
			expectedResult: "Linux",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Packet.ServerOS = tc.license
			checkStatus(t, p.h008, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH009(t *testing.T) {
	p, checkStatus := setupTest(t, "environment")

	testCases := []struct {
		name               string
		totalPosts         int
		enableIndexing     bool
		enableSearching    bool
		enableAutoComplete bool
		expectedStatus     CheckStatus
		expectedResult     string
	}{
		{
			name:               "h009 - total posts is less than 2.5 million",
			totalPosts:         2000000,
			enableIndexing:     false,
			enableSearching:    false,
			enableAutoComplete: false,
			expectedStatus:     Ignore,
			expectedResult:     "<2.5M posts, No Elasticsearch",
		},
		{
			name:               "h009 - Elasticsearch indexing is not enabled",
			totalPosts:         3000000,
			enableIndexing:     false,
			enableSearching:    false,
			enableAutoComplete: false,
			expectedStatus:     Fail,
			expectedResult:     ">2.5M posts, No Elasticsearch",
		},
		{
			name:               "h009 - total posts is greater than 2.5 million and Elasticsearch indexing is enabled",
			totalPosts:         3000000,
			enableIndexing:     true,
			enableSearching:    true,
			enableAutoComplete: true,
			expectedStatus:     Pass,
			expectedResult:     "Elasticsearch Enabled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Packet.TotalPosts = tc.totalPosts
			p.packet.Config.ElasticsearchSettings.EnableIndexing = &tc.enableIndexing
			p.packet.Config.ElasticsearchSettings.EnableSearching = &tc.enableSearching
			p.packet.Config.ElasticsearchSettings.EnableAutocomplete = &tc.enableAutoComplete
			checkStatus(t, p.h009, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}
