package dfg

import (
	"encoding/json"
	"os"
)

// DockerConfig contains all docker config fields
type DockerConfig struct {
	Image string `json:"image"`
}

// LoadDockerConfig reads data from a docker json config file
func LoadDockerConfig(path string) (parsedConfig DockerConfig, configErr error) {
	configFile, configErr := os.Open(path)

	if configErr != nil {
		return
	}

	configRaw := make([]byte, 200)
	bytesRead, configErr := configFile.Read(configRaw)

	if configErr != nil {
		return
	}

	configErr = json.Unmarshal(configRaw[:bytesRead], &parsedConfig)

	if configErr != nil {
		return
	}

	return
}
