package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Configuration contains docker config fields
type Configuration struct {
	AhabVersion string            `json:"ahab"`
	Command     string            `json:"command"`
	Entrypoint  string            `json:"entrypoint"`
	Environment []string          `json:"environment"`
	Hostname    string            `json:"hostname"`
	ImageURI    string            `json:"image"`
	Name        string            `json:"name"`
	Options     []string          `json:"options"`
	Permissions PermConfiguration `json:"permissions"`
	ShareX11    bool              `json:"shareX11"`
	User        string            `json:"user"`
	Volumes     []string          `json:"volumes"`
	Workdir     string            `json:"workdir"`
}

// PermConfiguration contains information regarding container user permissions setup
type PermConfiguration struct {
	CmdSet            string   `json:"cmdSet"`
	Disable           bool     `json:"disable"`
	Groups            []string `json:"groups"`
	RestartAfterSetup bool     `json:"restartAfterSetup"`
}

// UserConfiguration contains global user config fields
type UserConfiguration struct {
	Environment  []string `json:"environment"`
	Options      []string `json:"options"`
	HideCommands bool     `json:"hideCommands"`
	Volumes      []string `json:"volumes"`
}

// Version is the build-time dcfg version
var Version string

// configFileName holds the name of the config file
const configFileName = "ahab.json"

// userConfigFilePath holds the path of the user's config file, relative to their home dir
const userConfigFilePath = ".config/ahab/config.json"

// ContainerPathName returns a filepath-based container name for a given config file
func ContainerPathName(configPath string) string {
	name := filepath.Dir(configPath)
	name = strings.TrimPrefix(name, "/")
	return strings.ReplaceAll(name, "/", "_")
}

// ContainerName returns a consistent container name for config file
func ContainerName(config *Configuration, configPath string) string {
	if config.Name == "" {
		return ContainerPathName(configPath)
	}
	return config.Name
}

// ProjectConfig finds and parses the docker config file relative to the working directory
func ProjectConfig() (config *Configuration, configPath string, err error) {
	curDir, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("Failed to get working directory: %s", err)
		return
	}

	configPath, err = findConfigPath(curDir)
	if err != nil {
		err = fmt.Errorf("Failed to find config file: %s", err)
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

	err = checkConfigVersion(config.AhabVersion)
	return
}

// UserConfig finds and parses the user's docker config file
func UserConfig() (userConfig *UserConfiguration, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("Failed to get user home directory: %s", err)
		return
	}

	configPath := filepath.Join(homeDir, userConfigFilePath)
	configFile, err := os.Open(configPath)
	if err != nil && os.IsNotExist(err) {
		var blankConfig UserConfiguration
		return &blankConfig, nil
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&userConfig); err != nil {
		err = fmt.Errorf("Failed to parse user config file: %s", err)
	}
	return
}

// checkConfigVersion returns a non-nil err if the passed version is newer the active dcfg version
func checkConfigVersion(configVersion string) error {
	configVersionOrd, err := versionOrdinal(configVersion)
	if err != nil {
		return err
	}
	selfVersionOrd, err := versionOrdinal(Version)
	if err != nil {
		return err
	}

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
	if config.AhabVersion == "" {
		missingVars = appendToStrList(missingVars, "ahab")
	}
	if config.ImageURI == "" {
		missingVars = appendToStrList(missingVars, "image")
	}
	return
}

// versionOrdinal standardizes version strings for easy comparison
// see https://stackoverflow.com/a/18411978
func versionOrdinal(version string) (string, error) {
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
			return "", fmt.Errorf("versionOrdinal: Invalid version '%s'", version)
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo), nil
}
