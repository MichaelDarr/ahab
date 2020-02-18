package internal

import (
	"path/filepath"
	"testing"
)

const noConfDir = "/mnt/empty"
const exampleConfDir = "/mnt/project"
const exampleConfPath = "/mnt/project/" + ConfigFileName
const exampleConfChildDir = "/mnt/project/src"

func expectStrEq(expected string, actual string, t *testing.T) {
	if expected != actual {
		t.Errorf("\nExpected: %s\nActual: %s", expected, actual)
	}
}

// generate a mini config object at path /test/[testDir]/ahab.json
func miniConfig(testDir string) *Container {
	return &Container{
		FilePath: filepath.Join("/test", testDir, ConfigFileName),
		Fields: &Configuration{
			AhabVersion: "0.1",
			ImageURI:    "golang:1.13.7-buster",
		},
	}
}
