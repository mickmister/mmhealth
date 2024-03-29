package mmhealth

import "runtime/debug"

var GitCommit string
var GitVersion string

var BuildCommit = func() string {
	if GitCommit != "" {
		return GitCommit
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return ""
}()

var BuildVersion = func() string {
	if GitVersion != "" {
		return GitVersion
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}
	return ""
}()
