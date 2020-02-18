package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContainerPathName(t *testing.T) {
	name := ContainerPathName(minConfig.path)
	expectStrEq(minConfig.containerName, name, t)
}

func TestContainerName(t *testing.T) {
	name := ContainerName(minConfig.data, minConfig.path)
	expectStrEq(minConfig.containerName, name, t)

	name = ContainerName(maxConfig.data, maxConfig.path)
	expectStrEq(maxConfig.containerName, name, t)
}

func TestProjectConfig(t *testing.T) {
	os.Chdir(noConfDir)
	_, _, err := ProjectConfig()
	if err == nil {
		t.Errorf("Unexpected success finding project config: %s", err)
	}

	os.Chdir(exampleConfDir)
	_, configPath, err := ProjectConfig()
	if err != nil {
		t.Errorf("Error finding project config: %s", err)
	} else if configPath != exampleConfPath {
		t.Errorf("Error finding project config path. Expected %s, Found %s", exampleConfPath, configPath)
	}

	os.Chdir(exampleConfChildDir)
	_, configPath, err = ProjectConfig()
	if err != nil {
		t.Errorf("Error finding project config from child dir: %s", err)
	} else if configPath != exampleConfPath {
		t.Errorf("Error finding project config path from child dir. Expected %s, Found %s", exampleConfPath, configPath)
	}
}

func TestUserConfig(t *testing.T) {
	config, err := UserConfig()
	if err != nil {
		t.Errorf("Unexpected error finding present user config: %s", err)
	} else if config.HideCommands != true {
		t.Errorf("Unexpected value for hideCommands in user config file. Expected true, found %v", config.HideCommands)
	}

	// Find and temporarily rename user config file
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, userConfigFilePath)
	configPathTmp := filepath.Join(filepath.Dir(configPath), "tmp.json")
	os.Rename(configPath, configPathTmp)
	defer os.Rename(configPathTmp, configPath)

	config, err = UserConfig()
	if err != nil {
		t.Errorf("Unexpected error finding nonexistent user config: %s", err)
	} else if config.HideCommands != false {
		t.Errorf("Unexpected value for hideCommands in nonexistent config file. Expected false, found %v", config.HideCommands)
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
	expectedConfPath := filepath.Join(exampleConfDir, configFileName)
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

func TestMissingConfigVars(t *testing.T) {
	var config = Configuration{ImageURI: "ubuntu:18.04"}
	missing := missingConfigVars(&config)
	expectStrEq("ahab", missing, t)

	config = Configuration{AhabVersion: "0.1"}
	missing = missingConfigVars(&config)
	expectStrEq("image", missing, t)

	missing = missingConfigVars(&Configuration{})
	expectStrEq("ahab, image", missing, t)

	missing = missingConfigVars(minConfig.data)
	expectStrEq("", missing, t)
}
