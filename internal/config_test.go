package internal

import (
	"testing"
)

const configPath = "/home/user/test/ahab.json"

var configMin = Configuration{
	AhabVersion:       "0.1",
	Environment:       []string{},
	Hostname:          "",
	ImageURI:          "ubuntu:18.04",
	ManualPermissions: false,
	Name:              "",
	Options:           []string{},
	ShareX11:          false,
	Volumes:           []string{},
	Workdir:           "",
}

var configMax = Configuration{
	AhabVersion: "0.1",
	Environment: []string{
		"SOME=$THINGONE",
		"$OTHER=THINGTWO",
	},
	Hostname:          "myhost",
	ImageURI:          "ubuntu:18.04",
	ManualPermissions: true,
	Name:              "myname",
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
}

func TestContainerPathName(t *testing.T) {
	name := ContainerPathName(configPath)
	expectStrEq("home_user_test", name, t)
}

func TestContainerName(t *testing.T) {
	name := ContainerName(&configMin, configPath)
	expectStrEq("home_user_test", name, t)

	name = ContainerName(&configMax, configPath)
	expectStrEq("myname", name, t)
}

func TestCheckConfigVersion(t *testing.T) {
	Version = "2.0"

	// version older, no error expected
	if err := checkConfigVersion("0.0.1"); err != nil {
		t.Errorf("Expected version error checking ahab version '%s' against test version '0.0.1'", Version)
	}

	// version same, no error expected
	if err := checkConfigVersion("2"); err != nil {
		t.Errorf("Expected version error checking ahab version '%s' against test version '2'", Version)
	}

	// version newer, error expected
	if err := checkConfigVersion("3.2"); err == nil {
		t.Errorf("Unexpected version error checking ahab version '%s' against test version '3.2'", Version)
	}
}

func expectStrEq(expected string, actual string, t *testing.T) {
	if expected != actual {
		t.Errorf("\nExpected: %s\nActual: %s", expected, actual)
	}
}
