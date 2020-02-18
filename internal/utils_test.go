package internal

import (
	"os"
	"testing"
)

func TestExpandEnvs(t *testing.T) {
	os.Setenv("WHERE", "WORLD")
	os.Setenv("WHAT", "FOO")
	testEnvStrings := []string{"HELLO=$WHERE", "$WHAT=BAR"}
	expected := []string{"HELLO=WORLD", "FOO=BAR"}
	expandedStrings := expandEnvs(&testEnvStrings)
	expectStrsEq(&expandedStrings, &expected, t)
}
