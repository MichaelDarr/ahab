package internal

import "testing"

const noConfDir = "/mnt/empty"
const exampleConfDir = "/mnt/project"
const exampleConfPath = "/mnt/project/" + configFileName
const exampleConfChildDir = "/mnt/project/src"

type loadedConfig struct {
	containerName string
	path          string
	data          *Configuration
}

var minConfig = &loadedConfig{
	containerName: "ahab_test_min",
	path:          "/ahab/test/min/" + configFileName,
	data: &Configuration{
		AhabVersion: "0.1",
		ImageURI:    "ubuntu:18.04",
	},
}

var maxConfig = &loadedConfig{
	containerName: "maxConfig",
	path:          "/ahab/test/max/" + configFileName,
	data: &Configuration{
		AhabVersion: "0.1",
		Environment: []string{
			"SOME=$THINGONE",
			"$OTHER=THINGTWO",
		},
		Hostname: "myhost",
		ImageURI: "ubuntu:18.04",
		Name:     "maxConfig",
		Options: []string{
			"--gpus",
			"$ALL",
		},
		ShareX11: true,
		Volumes: []string{
			"./:/mnt/cur",
			"~/:/mnt/home",
			"/mnt:/mnt/nest",
		},
		Workdir: "/mnt/cur",
	},
}

func expectStrEq(expected string, actual string, t *testing.T) {
	if expected != actual {
		t.Errorf("\nExpected: %s\nActual: %s", expected, actual)
	}
}
