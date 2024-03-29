package healthchecks

import (
	"testing"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

func TestH001(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		siteURL        string
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name:           "h001 - SiteURL not set",
			siteURL:        "",
			expectedStatus: Fail,
			expectedResult: "Not set",
		},
		{
			name:           "h001 - SiteURL set",
			siteURL:        "http://localhost",
			expectedStatus: Pass,
			expectedResult: "Set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ServiceSettings.SiteURL = &tc.siteURL
			checkStatus(t, p.h001, tc.expectedStatus, tc.expectedResult)
		})
	}
}
func TestA001(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		linkPreviews   bool
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name:           "a001 - Enable link previews is false",
			linkPreviews:   false,
			expectedStatus: Warn,
			expectedResult: "Not enabled",
		},
		{
			name:           "a001 - Enable link previews is true",
			linkPreviews:   true,
			expectedStatus: Pass,
			expectedResult: "Enabled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ServiceSettings.EnableLinkPreviews = &tc.linkPreviews
			checkStatus(t, p.a001, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestA002(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		sessionLength  bool
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name:           "a002 - Extend session length with activity is false",
			sessionLength:  false,
			expectedStatus: Warn,
			expectedResult: "Not enabled",
		},
		{
			name:           "a002 - Extend session length with activity is true",
			sessionLength:  true,
			expectedStatus: Pass,
			expectedResult: "Enabled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ServiceSettings.ExtendSessionLengthWithActivity = &tc.sessionLength
			checkStatus(t, p.a002, tc.expectedStatus, tc.expectedResult)
		})
	}
}
func TestP002(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name             string
		notificationType string
		expectedStatus   types.CheckStatus
		expectedResult   string
	}{
		{
			name:             "p002 - not using ID notifications",
			notificationType: "full",
			expectedStatus:   Warn,
			expectedResult:   "Set to `full`",
		},
		{
			name:             "p002 - using ID notifications",
			notificationType: "id_loaded",
			expectedStatus:   Pass,
			expectedResult:   "Set to `id_loaded`",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.EmailSettings.PushNotificationContents = &tc.notificationType
			checkStatus(t, p.p002, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH002(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		enableIndexing bool
		liveIndexing   int
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name:           "h002 - ES not in use",
			enableIndexing: false,
			liveIndexing:   10,
			expectedStatus: Ignore,
			expectedResult: "Elasticsearch disabled",
		},
		{
			name:           "h002 - ES in use and live indexing is set to default",
			enableIndexing: true,
			liveIndexing:   1,
			expectedStatus: Fail,
			expectedResult: "Uses default value",
		},
		{
			name:           "h002 - ES in use and live indexing is configured",
			enableIndexing: true,
			liveIndexing:   10,
			expectedStatus: Pass,
			expectedResult: "Modified to a value greater than default of 1 - `10`",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ElasticsearchSettings.EnableIndexing = &tc.enableIndexing
			p.packet.Config.ElasticsearchSettings.LiveIndexingBatchSize = &tc.liveIndexing
			checkStatus(t, p.h002, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH010(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name                  string
		enableIndexing        bool
		enableSearching       bool
		enableAutocomplete    bool
		disableDatabaseSearch bool
		expectedStatus        types.CheckStatus
		expectedResult        string
	}{
		{
			name:                  "h010 - using Database search",
			enableIndexing:        false,
			enableSearching:       false,
			enableAutocomplete:    false,
			disableDatabaseSearch: false,
			expectedStatus:        Pass,
			expectedResult:        "Database",
		},
		{
			name:                  "h010 - using Elasticsearch",
			enableIndexing:        true,
			enableSearching:       true,
			enableAutocomplete:    true,
			disableDatabaseSearch: false,
			expectedStatus:        Pass,
			expectedResult:        "Elasticsearch",
		},
		{
			name:                  "h010 - No search enabled",
			enableIndexing:        false,
			enableSearching:       false,
			enableAutocomplete:    false,
			disableDatabaseSearch: true,
			expectedStatus:        Fail,
			expectedResult:        "No search enabled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ElasticsearchSettings.EnableIndexing = &tc.enableIndexing
			p.packet.Config.ElasticsearchSettings.EnableSearching = &tc.enableSearching
			p.packet.Config.ElasticsearchSettings.EnableAutocomplete = &tc.enableAutocomplete
			p.packet.Config.SqlSettings.DisableDatabaseSearch = &tc.disableDatabaseSearch
			checkStatus(t, p.h010, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestP003(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		IdAttribute    string
		ldapEnabled    bool
		emailAttribute string
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name:           "p003 - not using LDAP",
			IdAttribute:    "",
			emailAttribute: "",
			ldapEnabled:    false,
			expectedStatus: Ignore,
			expectedResult: "LDAP disabled",
		},
		{
			name:           "p003 - LDAP enabled and ID attribute set",
			IdAttribute:    "uniqueID",
			emailAttribute: "email",
			ldapEnabled:    true,
			expectedStatus: Pass,
			expectedResult: `ID attribute set to \newline
\textbf{\seqsplit{uniqueID}}.`,
		},
		{
			name:           "p003 - LDAP enabled and ID attribute set the same",
			IdAttribute:    "myEmail",
			emailAttribute: "myEmail",
			ldapEnabled:    true,
			expectedStatus: Warn,
			expectedResult: "Using email",
		},
		{
			name:           "p003 - LDAP enabled no ID set",
			IdAttribute:    "",
			emailAttribute: "",
			ldapEnabled:    true,
			expectedStatus: Warn,
			expectedResult: "Using email",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.LdapSettings.Enable = &tc.ldapEnabled
			p.packet.Config.LdapSettings.IdAttribute = &tc.IdAttribute
			p.packet.Config.LdapSettings.EmailAttribute = &tc.emailAttribute

			checkStatus(t, p.p003, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestP004(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		idAttribute    string
		samlEnabled    bool
		emailAttribute string
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name:           "p004 - not using SAML",
			idAttribute:    "",
			emailAttribute: "",
			samlEnabled:    false,
			expectedStatus: Ignore,
			expectedResult: "SAML disabled",
		},
		{
			name:           "p004 - SAML enabled and ID attribute set",
			idAttribute:    "uniqueID",
			emailAttribute: "email",
			samlEnabled:    true,
			expectedStatus: Pass,
			expectedResult: `ID attribute set to \newline
\textbf{\seqsplit{uniqueID}}.`,
		},
		{
			name:           "p004 - SAML enabled and ID attribute set the same",
			idAttribute:    "myEmail",
			emailAttribute: "myEmail",
			samlEnabled:    true,
			expectedStatus: Warn,
			expectedResult: "Using email",
		},
		{
			name:           "p004 - SAML enabled no ID set",
			idAttribute:    "",
			emailAttribute: "",
			samlEnabled:    true,
			expectedStatus: Warn,
			expectedResult: "Using email",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.SamlSettings.Enable = &tc.samlEnabled
			p.packet.Config.SamlSettings.IdAttribute = &tc.idAttribute
			p.packet.Config.SamlSettings.EmailAttribute = &tc.emailAttribute

			checkStatus(t, p.p004, tc.expectedStatus, tc.expectedResult)
		})
	}
}
