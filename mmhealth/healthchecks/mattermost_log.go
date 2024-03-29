package healthchecks

import (
	"sort"
	"strings"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

func (p *ProcessPacket) logChecks() (results []CheckResult) {

	checks := map[string]CheckFunc{
		"h003": p.h003,
		"h004": p.h004,
		"h005": p.h005,
	}
	testResults := []CheckResult{}

	for id, check := range checks {
		result := check(p.Checks.MattermostLog)
		result.ID = id
		testResults = append(testResults, result)
	}

	return p.sortResults(testResults)
}

func (p *ProcessPacket) topLogs() []TopLogs {

	topLogs := make(map[string]TopLogs)

	for _, log := range p.packet.Logs {
		if log.Level == "debug" {
			continue
		}
		info, exists := topLogs[log.Msg]

		if exists {
			info.Count++
			topLogs[log.Msg] = info
		} else {
			topLogs[log.Msg] = TopLogs{Caller: log.Caller, Count: 1, Level: log.Level}
		}
	}

	// Convert map to slice
	logsSlice := make([]TopLogs, 0, len(topLogs))
	for msg, info := range topLogs {
		logsSlice = append(logsSlice, TopLogs{Count: info.Count, Caller: info.Caller, Msg: msg, Level: info.Level})
	}

	// Sort slice by count in descending order
	sort.Slice(logsSlice, func(i, j int) bool {
		return logsSlice[i].Count > logsSlice[j].Count
	})

	sliceLength := len(logsSlice)
	if sliceLength > 10 {
		sliceLength = 10
	}

	return logsSlice[:sliceLength]

}

func (p *ProcessPacket) h003(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h003", checks, Pass)

	for _, log := range p.packet.Logs {
		if strings.Contains(log.Msg, "context deadline exceeded") {
			result.Status = Fail
			result.Result = check.Result.Fail
			return result
		}
	}

	return result
}

func (p *ProcessPacket) h004(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h004", checks, Pass)

	for _, log := range p.packet.Logs {
		if strings.Contains(log.Msg, "i/o timeout") {
			result.Status = Fail
			result.Result = check.Result.Fail
			return result
		}
	}
	return result
}

func (p *ProcessPacket) h005(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h005", checks, Pass)

	for _, log := range p.packet.Logs {
		if strings.Contains(log.Msg, "Error while creating session for user access token") {
			result.Status = Fail
			result.Result = check.Result.Fail
			return result
		}
	}

	return result
}
