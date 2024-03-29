package healthchecks

import (
	"fmt"
	"strings"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

func (p *ProcessPacket) configChecks() (results []CheckResult) {

	checks := map[string]CheckFunc{
		"h001": p.h001,
		"h002": p.h002,
		"h010": p.h010,
		"p002": p.p002,
		"p003": p.p003,
		"p004": p.p004,
		"a001": p.a001,
		"a002": p.a002,
	}
	testResults := []CheckResult{}

	for id, check := range checks {
		result := check(p.Checks.Config)
		result.ID = id
		testResults = append(testResults, result)
	}

	return p.sortResults(testResults)
}

//
//
// SERVICE SETTINGS
//
//

func (p *ProcessPacket) h001(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h001", checks, Fail)

	if *p.packet.Config.ServiceSettings.SiteURL != "" {
		result.Result = check.Result.Pass
		result.Status = Pass
	}
	return result
}

func (p *ProcessPacket) a001(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("a001", checks, Fail)

	if *p.packet.Config.ServiceSettings.EnableLinkPreviews {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

func (p *ProcessPacket) a002(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("a002", checks, Fail)

	if *p.packet.Config.ServiceSettings.ExtendSessionLengthWithActivity {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

//
//
// EMAIL SETTINGS
//
//

func (p *ProcessPacket) p002(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("p002", checks, Fail)

	result.Result = fmt.Sprintf(check.Result.Fail, *p.packet.Config.EmailSettings.PushNotificationContents)

	if *p.packet.Config.EmailSettings.PushNotificationContents == "id_loaded" {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

//
//
// ELASTICSEARCH SETTINGS
//
//

func (p *ProcessPacket) h002(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h002", checks, Fail)
	if *p.packet.Config.ElasticsearchSettings.EnableIndexing {
		if *p.packet.Config.ElasticsearchSettings.LiveIndexingBatchSize > 1 {
			result.Result = fmt.Sprintf(check.Result.Pass, *p.packet.Config.ElasticsearchSettings.LiveIndexingBatchSize)
			result.Status = Pass
		}
		return result
	}

	result.Status = Ignore
	result.Result = check.Result.Ignore

	return result
}

// checks to make sure Elasticsearch is enabled OR database search is NOT disabled.

func (p *ProcessPacket) h010(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h010", checks, Fail)
	config := p.packet.Config

	if *config.ElasticsearchSettings.EnableIndexing && *config.ElasticsearchSettings.EnableSearching && *config.ElasticsearchSettings.EnableAutocomplete {
		result.Result = fmt.Sprintf(check.Result.Pass, "Elasticsearch")
		result.Status = Pass
	} else if *config.SqlSettings.DisableDatabaseSearch {
		result.Result = fmt.Sprintf(check.Result.Fail, "No search enabled")
		result.Status = Fail
	} else {
		result.Status = Pass
		result.Result = "Database"
	}

	return result
}

func (p *ProcessPacket) p003(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("p003", checks, Fail)
	config := p.packet.Config

	if !*config.LdapSettings.Enable {
		result.Result = check.Result.Ignore
		result.Status = Ignore
		return result
	}

	if *config.LdapSettings.IdAttribute != "" {
		if !strings.EqualFold(*config.LdapSettings.EmailAttribute, *config.LdapSettings.IdAttribute) &&
			!strings.EqualFold(*config.LdapSettings.IdAttribute, "email") {
			result.Result = fmt.Sprintf(check.Result.Pass, *config.LdapSettings.IdAttribute)
			result.Status = Pass
			return result
		}
	}
	return result
}

func (p *ProcessPacket) p004(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("p004", checks, Fail)
	config := p.packet.Config

	if !*config.SamlSettings.Enable {
		result.Result = check.Result.Ignore
		result.Status = Ignore
		return result
	}

	if *config.SamlSettings.IdAttribute != "" {
		if !strings.EqualFold(*config.SamlSettings.EmailAttribute, *config.SamlSettings.IdAttribute) &&
			!strings.EqualFold(*config.SamlSettings.IdAttribute, "email") {
			result.Result = fmt.Sprintf(check.Result.Pass, *config.SamlSettings.IdAttribute)
			result.Status = Pass
			return result
		}
	}
	return result
}
