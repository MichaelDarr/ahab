package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Configuration contains all docker config fields
type Configuration struct {
	ImageURL string `json:"image"`
}

// configFileName holds the name of the config file
const configFileName string = "dfg.json"

var config *Configuration

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
		return fmt.Errorf("Failed to read config file '%s': %s", configPath, err)
	}

	missingVars := missingConfigVars()
	if missingVars != "" {
		return fmt.Errorf("Config file '%s' missing required fields: %s", configPath, missingVars)
	}

	return nil
}

// helper for creating human-readable comma-separated lists
func appendToStrList(list string, newEl string) (finalStr string) {
	if list == "" {
		return newEl
	}
	return list + ", " + newEl
}

// recursively search for a config file starting at topDir, ending at fs root
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

// return a comma-separated string of required fields not present in config file
func missingConfigVars() (missingVars string) {
	if config.ImageURL == "" {
		missingVars = appendToStrList(missingVars, "image")
	}

	return
}
