package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUserConfig(t *testing.T) {
	// this test runs in a different group than others as it changes the user's config file
	t.Parallel()
	config, err := UserConfig()
	if err != nil {
		t.Errorf("Unexpected error finding present user config: %s", err)
	} else if !config.HideCommands {
		t.Error("hidecommands key is true in user config but not in parsed config object")
	}

	// Find and temporarily rename user config file
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, UserConfigFilePath)
	configPathTmp := filepath.Join(filepath.Dir(configPath), "tmp.json")
	os.Rename(configPath, configPathTmp)
	defer os.Rename(configPathTmp, configPath)

	config, err = UserConfig()
	if err != nil {
		t.Errorf("Unexpected error finding nonexistent user config: %s", err)
	} else if config.HideCommands {
		t.Error("hidecommands read as true from user config, which should be nonexistent")
	}
}

func TestCheckConfigVersion(t *testing.T) {
	Version = "2.0"

	// version older, no error expected
	if err := checkConfigVersion("0.0.1"); err != nil {
		t.Errorf("Error checking ahab version '%s' against test version '0.0.1'", Version)
	}

	// version same, no error expected
	if err := checkConfigVersion("2"); err != nil {
		t.Errorf("Error checking ahab version '%s' against test version '2'", Version)
	}

	// version newer, error expected
	if err := checkConfigVersion("3.2"); err == nil {
		t.Errorf("Unexpected success checking ahab version '%s' against test version '3.2'", Version)
	}
}

func TestFindConfigPath(t *testing.T) {
	// look for config where there isn't one, error expected
	if _, err := findConfigPath(noConfDir); err == nil {
		t.Errorf("Unexpected success finding a config path where there is none: %s", noConfDir)
	}

	// look for config where it is present
	configPath, err := findConfigPath(exampleConfDir)
	expectedConfPath := filepath.Join(exampleConfDir, ConfigFileName)
	if err != nil {
		t.Errorf("Error finding a config file in %s", exampleConfDir)
	} else if configPath != expectedConfPath {
		t.Errorf("Unexpected config path found for %s: '%s' (Expected %s)", exampleConfDir, configPath, expectedConfPath)
	}

	// look for config where it is present in parent dir
	configPath, err = findConfigPath(exampleConfChildDir)
	if err != nil {
		t.Errorf("Error finding a config path in %s", exampleConfChildDir)
	} else if configPath != expectedConfPath {
		t.Errorf("Unexpected config path found for %s: '%s' (Expected %s)", exampleConfDir, configPath, expectedConfPath)
	}
}
