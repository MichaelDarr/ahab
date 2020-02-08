package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Configuration contains all docker config fields
type Configuration struct {
	DcfgVersion string   `json:"dcfg-version"`
	ImageURI    string   `json:"image"`
	Options     []string `json:"options"`
}

// configFileName holds the name of the config file
const configFileName string = "dcfg.json"

// ContainerName returns a filepath-based container name for a given config file
func ContainerName(configPath string) (name string, err error) {
	name = filepath.Dir(configPath)
	name = strings.TrimPrefix(name, "/") // we don't want every name to begin with _
	name = strings.ReplaceAll(name, "/", "_")
	return
}

// Config finds and parses the docker config file relative to the working directory
func Config() (config *Configuration, configPath string, err error) {
	curDir, _ := os.Getwd()
	configPath, err = findConfigPath(curDir)
	if err != nil {
		return
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		err = fmt.Errorf("Failed to open config file '%s': %s", configPath, err)
		return
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&config); err != nil {
		err = fmt.Errorf("Failed to parse config file '%s': %s", configPath, err)
		return
	}

	missingVars := missingConfigVars(config)
	if missingVars != "" {
		err = fmt.Errorf("Config file '%s' missing required fields: %s", configPath, missingVars)
		return
	}

	err = checkConfigVersion(config.DcfgVersion)
	return
}

// checkConfigVersion returns a non-nil err if the passed version is newer the active dcfg version
func checkConfigVersion(configVersion string) error {
	configVersionOrd, selfVersionOrd := versionOrdinal(configVersion), versionOrdinal(Version)
	if configVersionOrd > selfVersionOrd {
		return fmt.Errorf("Config file requires dcfg >= %s (your version: %s)", configVersion, Version)
	}
	return nil
}

// findConfigPath recursively searches for a config file starting at topDir, ending at fs root
func findConfigPath(topDir string) (configPath string, err error) {
	configTestPath := filepath.Join(topDir, configFileName)
	_, err = os.Stat(configTestPath)
	if err != nil && os.IsNotExist(err) && filepath.Clean(topDir) != "/" {
		configPath, err = findConfigPath(filepath.Join(topDir, ".."))
	} else if err != nil && os.IsNotExist(err) {
		err = fmt.Errorf("No config file '%s' found in current or parent directories", configFileName)
	} else if err != nil {
		err = fmt.Errorf("Failed to find config file '%s': %s", configFileName, err)
	} else {
		configPath = configTestPath
	}

	return
}

// missingConfigVars returns a comma-separated string of missing required config fields
func missingConfigVars(config *Configuration) (missingVars string) {
	if config.ImageURI == "" {
		missingVars = appendToStrList(missingVars, "image")
	}
	if config.DcfgVersion == "" {
		missingVars = appendToStrList(missingVars, "version")
	}

	return
}

// see https://stackoverflow.com/a/18411978
func versionOrdinal(version string) string {
	// ISO/IEC 14651:2011
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			panic("VersionOrdinal: invalid version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo)
}
