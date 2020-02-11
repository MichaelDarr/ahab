package internal

import (
	"testing"
)

var configMinl = Configuration{
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
	name := ContainerPathName("/home/user/test/ahab.json")
	expectStrEq("home_user_test", name, t)
}

func expectStrEq(expected string, actual string, t *testing.T) {
	if expected != actual {
		t.Errorf("\nExpected: %s\nActual: %s", expected, actual)
	}
}
