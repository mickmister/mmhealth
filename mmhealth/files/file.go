package files

import (
	"embed"
	"io/fs"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"gopkg.in/yaml.v3"
)

//go:embed checks.yaml config.yaml
var Files embed.FS

func ReadChecksFile() (types.ChecksFile, error) {
	var checks types.ChecksFile
	data, err := fs.ReadFile(Files, "checks.yaml")
	if err != nil {
		return checks, err
	}

	err = yaml.Unmarshal(data, &checks)
	if err != nil {
		return checks, err
	}

	return checks, nil
}

func ReadConfigFile() (types.ConfigFile, error) {
	var config types.ConfigFile
	data, err := fs.ReadFile(Files, "config.yaml")
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
