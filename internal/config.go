package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Version is the build-time dcfg version
var Version string

// ConfigFileName holds the name of the config file
const ConfigFileName = "ahab.json"

// UserConfigFilePath holds the path of the user's config file, relative to their home dir
const UserConfigFilePath = ".config/ahab/config.json"

// Configuration contains docker config fields
type Configuration struct {
	AhabVersion       string             `json:"ahab"`
	Command           string             `json:"command"`
	Entrypoint        string             `json:"entrypoint"`
	Environment       []string           `json:"environment"`
	Hostname          string             `json:"hostname"`
	ImageURI          string             `json:"image"`
	Init              []string           `json:"init"`
	Name              string             `json:"name"`
	Options           []string           `json:"options"`
	Permissions       *PermConfiguration `json:"permissions"`
	RestartAfterSetup bool               `json:"restartAfterSetup"`
	ShareX11          bool               `json:"shareX11"`
	User              string             `json:"user"`
	Volumes           []string           `json:"volumes"`
	Workdir           string             `json:"workdir"`
}

// UserConfiguration contains global user config fields
type UserConfiguration struct {
	Environment  []string `json:"environment"`
	Options      []string `json:"options"`
	HideCommands bool     `json:"hideCommands"`
	Volumes      []string `json:"volumes"`
}

// PermConfiguration contains information regarding container user permissions setup
type PermConfiguration struct {
	CmdSet  string   `json:"cmdSet"`
	Disable bool     `json:"disable"`
	Groups  []string `json:"groups"`
}

// a comma-separated list of required fields missing from a configuration
func (config *Configuration) missingFields() (missingVars string) {
	if config.AhabVersion == "" {
		missingVars = appendToStrList(missingVars, "ahab")
	}
	if config.ImageURI == "" {
		missingVars = appendToStrList(missingVars, "image")
	}
	return
}

// GetConfig finds and parses the config file relative to the working directory
func GetConfig() (config *Configuration, configPath string, err error) {
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

	if missingFields := config.missingFields(); missingFields != "" {
		err = fmt.Errorf("Config file '%s' missing required fields: %s", configPath, missingFields)
		return
	}

	err = checkConfigVersion(config.AhabVersion)
	return
}

// UserConfig finds and parses the user's docker config file
// If the user's config is not found, an empty userConfig is returned
func UserConfig() (userConfig *UserConfiguration, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("Failed to get user home directory: %s", err)
		return
	}

	configPath := filepath.Join(homeDir, UserConfigFilePath)
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
		return fmt.Errorf("Config file requires ahab >= %s (your version: %s)", configVersion, Version)
	}
	return nil
}

// findConfigPath recursively searches for a config file starting at topDir, ending at fs root
func findConfigPath(topDir string) (configPath string, err error) {
	configTestPath := filepath.Join(topDir, ConfigFileName)
	_, err = os.Stat(configTestPath)
	if err != nil && os.IsNotExist(err) && filepath.Clean(topDir) != "/" {
		configPath, err = findConfigPath(filepath.Join(topDir, ".."))
	} else if err != nil && os.IsNotExist(err) {
		err = fmt.Errorf("No config file '%s' found in current or parent directories", ConfigFileName)
	} else if err != nil {
		err = fmt.Errorf("Failed to find config file '%s': %s", ConfigFileName, err)
	} else {
		configPath = configTestPath
	}
	return
}
