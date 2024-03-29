package healthchecks

import (
	"fmt"
	"sort"
	"time"

	"github.com/Masterminds/semver"
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/mattermost/mattermost/server/public/model"
)

func (p *ProcessPacket) pluginChecks() (results []PluginResults) {

	plugins := p.packet.Plugins
	pluginResults := []PluginResults{}

	for _, plugin := range plugins.Active {
		parsedPlugin := getPluginResults(plugin, p.Config, true)
		pluginResults = append(pluginResults, parsedPlugin)
	}

	for _, plugin := range plugins.Inactive {
		parsedPlugin := getPluginResults(plugin, p.Config, false)
		pluginResults = append(pluginResults, parsedPlugin)
	}

	sort.Slice(pluginResults, func(i, j int) bool {
		return pluginResults[i].LatestReleaseDate > pluginResults[j].LatestReleaseDate
	})

	return pluginResults

}

func getPluginResults(plugin *model.PluginInfo, config types.ConfigFile, isActive bool) PluginResults {

	pluginEntry := getPluginInfoFromConfig(plugin.Id, config)

	parsedPlugin := PluginResults{
		PluginID:         plugin.Id,
		PluginName:       plugin.Name,
		PluginURL:        findAPluginURL(plugin, pluginEntry),
		Active:           isActive,
		LatestVersion:    pluginEntry.Latest,
		InstalledVersion: plugin.Version,
		SupportLevel:     "unknown",
	}

	if pluginEntry.LatestReleaseDate != "" {
		t, err := time.Parse(time.RFC3339, pluginEntry.LatestReleaseDate)
		if err != nil {
			fmt.Printf("Error parsing latest release date for plugin %s: %s\n", plugin.Id, err)
		}
		parsedPlugin.LatestReleaseDate = t.Format("2006-01-02")
	}

	if pluginEntry.SupportLevel != "" {
		parsedPlugin.SupportLevel = pluginEntry.SupportLevel
	}

	isUpdated, err := returnIsUpdated(parsedPlugin.LatestVersion, parsedPlugin.InstalledVersion)

	if err != nil {
		fmt.Printf("Error parsing version for plugin %s: %s\n", plugin.Id, err)
	}

	parsedPlugin.IsUpdated = isUpdated
	return parsedPlugin
}

func returnIsUpdated(latest string, installed string) (bool, error) {
	installedVersion, err := semver.NewVersion(installed)
	if err != nil {
		return false, err
	}
	latestSemVersion, err := semver.NewVersion(latest)
	if err != nil {
		return false, err
	}
	isUpdated := installedVersion.Compare(latestSemVersion) >= 0
	return isUpdated, nil
}

func getPluginInfoFromConfig(pluginID string, config types.ConfigFile) types.PluginEntry {
	for ID, p := range config.Plugins {
		if ID == pluginID {
			return p
		}
	}

	return types.PluginEntry{}
}

func findAPluginURL(plugin *model.PluginInfo, pluginEntry types.PluginEntry) string {

	if pluginEntry.Repo != "" {
		return pluginEntry.Repo
	}

	if plugin.HomepageURL != "" {
		return plugin.HomepageURL
	}
	if plugin.SupportURL != "" {
		return plugin.SupportURL
	}
	if plugin.ReleaseNotesURL != "" {
		return plugin.ReleaseNotesURL
	}

	return ""
}
