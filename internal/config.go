package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ConfigFileName holds the name of the config file
const ConfigFileName = "ahab.json"

// UserConfigFilePath holds the path of the user's config file, relative to their home dir
const UserConfigFilePath = ".config/ahab/config.json"

// Version holds the build-time ahab version
var Version string

// Configuration contains docker config fields
type Configuration struct {
	AhabVersion       string            `json:"ahab"`
	BuildContext      string            `json:"buildContext"`
	Command           string            `json:"command"`
	Dockerfile        string            `json:"dockerfile"`
	Entrypoint        string            `json:"entrypoint"`
	Environment       []string          `json:"environment"`
	Hostname          string            `json:"hostname"`
	ImageURI          string            `json:"image"`
	Init              []string          `json:"init"`
	Name              string            `json:"name"`
	Options           []string          `json:"options"`
	Permissions       PermConfiguration `json:"permissions"`
	RestartAfterSetup bool              `json:"restartAfterSetup"`
	ShareDisplay      bool              `json:"shareDisplay"`
	User              string            `json:"user"`
	Volumes           []string          `json:"volumes"`
	Workdir           string            `json:"workdir"`
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

// return a non-nil error if config is invalid
func (config *Configuration) validateConfig() (err error) {
	if config.AhabVersion == "" {
		err = fmt.Errorf("Missing required version field 'ahab'")
	} else if config.ImageURI == "" && config.Dockerfile == "" {
		err = fmt.Errorf("Either 'image' or 'dockerfile' must be present")
	} else if config.ImageURI != "" && config.Dockerfile != "" {
		err = fmt.Errorf("'image' and `dockerfile' cannot both be present")
	}
	return
}
