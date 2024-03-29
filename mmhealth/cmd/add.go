package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	mmhealth "github.com/coltoneshaw/mmhealth/mmhealth"
	"github.com/coltoneshaw/mmhealth/mmhealth/files"
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a check to the yaml file",
	Long:  "Interactive dialog to add a check to the yaml file before building the check. ",
	RunE:  addCmdF,
}

func init() {

	AddCmd.Hidden = mmhealth.BuildVersion != "(devel)"

	RootCmd.AddCommand(
		AddCmd,
	)
}

var qs = []*survey.Question{
	{
		Name:      "name",
		Prompt:    &survey.Input{Message: "What is the name of this check?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "group",
		Prompt: &survey.Select{
			Message: "Choose a check group:",
			Options: []string{"environment", "config", "packet", "mattermostLog", "notificationLog", "plugins"},
		},
		Validate: survey.Required,
	},
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Choose a check type:",
			Options: []string{"proactive", "health", "adoption"},
		},
		Validate: survey.Required,
	},
	{
		Name: "severity",
		Prompt: &survey.Select{
			Message: "Choose a check severity:",
			Options: []string{"urgent", "high", "medium", "low"},
			Default: "medium",
		},
		Validate: survey.Required,
	},
	{
		Name:     "description",
		Prompt:   &survey.Input{Message: "What is the description of this check?"},
		Validate: survey.Required,
	},
	{
		Name:     "pass",
		Prompt:   &survey.Input{Message: "What is the pass message?"},
		Validate: survey.Required,
	},
	{
		Name:     "fail",
		Prompt:   &survey.Input{Message: "What is the fail message?"},
		Validate: survey.Required,
	},
	{
		Name: "ignore",
		Prompt: &survey.Input{
			Message: "What is the ignore message? (Optional)",
			Help:    "If you don't want to show anything, just press enter",
		},
	},
}

func addCmdF(cmd *cobra.Command, args []string) error {
	answers := struct {
		Name        string
		Type        string
		Group       string
		Severity    string
		Description string
		Pass        string
		Fail        string
		Ignore      string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return errors.Wrap(err, "Failed to ask questions")
	}

	checks, err := files.ReadChecksFile()
	if err != nil {
		return errors.Wrap(err, "Failed to read checks file")
	}

	newKey := generateCheckKey(answers.Type, checks)

	newCheck := types.Check{
		Name:        answers.Name,
		Result:      types.Result{Pass: answers.Pass, Fail: answers.Fail, Ignore: answers.Ignore},
		Description: answers.Description,
		Severity:    types.CheckSeverity(answers.Severity),
		Type:        types.CheckType(answers.Type),
	}

	switch answers.Group {
	case "config":
		checks.Config[newKey] = newCheck
		checks.Config = sortGroup(checks.Config)
	case "packet":
		checks.Packet[newKey] = newCheck
		checks.Packet = sortGroup(checks.Packet)
	case "mattermostLog":
		checks.MattermostLog[newKey] = newCheck
		checks.MattermostLog = sortGroup(checks.MattermostLog)
	case "notificationLog":
		checks.NotificationLog[newKey] = newCheck
		checks.NotificationLog = sortGroup(checks.NotificationLog)
	case "plugins":
		checks.Plugins[newKey] = newCheck
		checks.Plugins = sortGroup(checks.Plugins)
	case "environment":
		checks.Environment[newKey] = newCheck
		checks.Environment = sortGroup(checks.Environment)
	}

	// Marshal the Config struct back into YAML
	err = storeChecksFile(checks)

	if err != nil {
		return errors.Wrap(err, "Failed to write checks file")
	}

	switch answers.Group {
	case "config":
		fmt.Printf("Check %s added successfully. Edit ./mmhealth/healthchecks/config.go to build the check.", newKey)
	case "packet":
		fmt.Printf("Check %s added successfully. Edit ./mmhealth/healthchecks/packet.go to build the check.", newKey)
	case "mattermostLog":
		fmt.Printf("Check %s added successfully. Edit ./mmhealth/healthchecks/mattermostLog.go to build the check.", newKey)
	case "notificationLog":
		fmt.Printf("Check %s added successfully. Edit ./mmhealth/healthchecks/notificationLog.go to build the check.", newKey)
	case "plugins":
		fmt.Printf("Check %s added successfully. Edit ./mmhealth/healthchecks/plugins.go to build the check.", newKey)
	case "environment":
		fmt.Printf("Check %s added successfully. Edit ./mmhealth/healthchecks/environment.go to build the check.", newKey)
	}
	return nil

}

// parses the existing yaml file and finds the highest existing value and returns the next value
func generateCheckKey(checkType string, checks types.ChecksFile) string {
	prefix := string(checkType[0])
	highest := 0

	for key := range checks.Environment {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	for key := range checks.Config {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	for key := range checks.MattermostLog {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	for key := range checks.NotificationLog {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	for key := range checks.Plugins {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	for key := range checks.Packet {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	return fmt.Sprintf("%s%03d", prefix, highest+1)
}

func sortGroup(checks map[string]types.Check) map[string]types.Check {
	var keys []string
	for k := range checks {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})

	// Create a new sorted map
	sortedChecks := make(map[string]types.Check)
	for _, k := range keys {
		sortedChecks[k] = checks[k]
	}

	// Replace the 'config' group with the sorted map
	return sortedChecks
}

func storeChecksFile(checks types.ChecksFile) error {
	data, err := yaml.Marshal(&checks)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal checks file")

	}
	return os.WriteFile(filepath.Join("./mmhealth/files", "checks.yaml"), data, 0644)
}
