package internal

import (
	"os"
	"testing"
)

func TestExpandEnvs(t *testing.T) {
	os.Setenv("WHERE", "WORLD")
	os.Setenv("WHAT", "FOO")
	testEnvStrings := []string{"HELLO=$WHERE", "$WHAT=BAR"}
	expectedEnvStrings := []string{"HELLO=WORLD", "FOO=BAR"}
	expandedStrings := expandEnvs(&testEnvStrings)
	expectStrsEq(&expandedStrings, &expectedEnvStrings, t)
}

func TestPrepVolumeString(t *testing.T) {
	os.Setenv("WHERE", "here")
	configPath := "/home/ahab/src/ahab.json"
	testVolStrings := []string{
		"~/.config:/home/.config",
		"./build:/build",
		"/tmp:/tmp",
		"~/$WHERE:/mnt/$WHERE",
	}
	expectedVolStrings := []string{
		"/home/ahab/.config:/home/.config",
		"/home/ahab/src/build:/build",
		"/tmp:/tmp",
		"/home/ahab/here:/mnt/here",
	}
	for i, testString := range testVolStrings {
		expandedVol, err := prepVolumeString(testString, configPath)
		if err != nil {
			t.Errorf("Error preparing volume string: %s", err)
		}
		expectStrEq(expectedVolStrings[i], expandedVol, t)
	}
}

func TestSplitGroups(t *testing.T) {
	testGroups := []string{
		"!docker",
		"http",
		"log",
		"!sudo",
		"users",
		"wheel",
	}
	expectedExistingGroups := []string{
		"http",
		"log",
		"users",
		"wheel",
	}
	expectedNewGroups := []string{
		"docker",
		"sudo",
	}
	existingGroups, newGroups := splitGroups(&testGroups)
	expectStrsEq(&expectedExistingGroups, &existingGroups, t)
	expectStrsEq(&expectedNewGroups, &newGroups, t)
}
