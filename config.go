package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Configuration contains all docker config fields
type Configuration struct {
	DcfgVersion string `json:"dcfg-version"`
	ImageURL    string `json:"image"`
}

// config holds the configuration loaded by InitConfig
var config *Configuration

// configFileName holds the name of the config file
const configFileName string = "dcfg.json"

// version holds the build-time dcfg version (set from build command)
var version string

// InitConfig finds and parses the docker config file relative to the working directory
func InitConfig() error {
	curDir, _ := os.Getwd()
	configPath, err := findConfigPath(curDir)
	if err != nil {
		return err
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("Failed to open config file '%s': %s", configPath, err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&config); err != nil {
		return fmt.Errorf("Failed to parse config file '%s': %s", configPath, err)
	}

	missingVars := missingConfigVars()
	if missingVars != "" {
		return fmt.Errorf("Config file '%s' missing required fields: %s", configPath, missingVars)
	}

	err = checkConfigVersion(config.DcfgVersion)
	return err
}

// appendToStrList is a helper for creating human-readable comma-separated lists
func appendToStrList(list string, newEl string) (finalStr string) {
	if list == "" {
		return newEl
	}
	return list + ", " + newEl
}

// checkConfigVersion returns a non-nil err if the passed version is newer the active dcfg version
func checkConfigVersion(configVersion string) error {
	configVersionOrd, selfVersionOrd := versionOrdinal(configVersion), versionOrdinal(version)
	if configVersionOrd > selfVersionOrd {
		return fmt.Errorf("Config file requires dcfg >= %s (your version: %s)", configVersion, version)
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
func missingConfigVars() (missingVars string) {
	if config.ImageURL == "" {
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
