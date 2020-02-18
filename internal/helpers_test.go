package internal

import (
	"path/filepath"
	"testing"
)

const noConfDir = "/mnt/empty"
const exampleConfDir = "/mnt/project"
const exampleConfPath = "/mnt/project/" + ConfigFileName
const exampleConfChildDir = "/mnt/project/src"

func expectContainerStatus(wantedStatus int, container *Container, t *testing.T) {
	foundStatus, err := container.Status()
	if err != nil {
		t.Errorf("Error checking status: %s", err)
	} else if foundStatus != wantedStatus {
		t.Errorf("Unexpected container status %s (expected %s)", ParseStatus(foundStatus), ParseStatus(wantedStatus))
	}
}

func prohibitContainerStatus(prohibitStatus int, container *Container, t *testing.T) {
	foundStatus, err := container.Status()
	if err != nil {
		t.Errorf("Error checking status: %s", err)
	} else if foundStatus == prohibitStatus {
		t.Errorf("Observed a prohibitted container status: %s", ParseStatus(foundStatus))
	}
}

func expectStrEq(expected string, actual string, t *testing.T) {
	if expected != actual {
		t.Errorf("\nExpected: %s\nActual: %s", expected, actual)
	}
}

// generate a mini container object with a config at path /test/[testDir]/ahab.json
func miniContainer(testDir string) *Container {
	return &Container{
		FilePath: filepath.Join("/test", testDir, ConfigFileName),
		Fields: &Configuration{
			AhabVersion: "0.1",
			ImageURI:    "golang:1.13.7-buster",
		},
	}
}
