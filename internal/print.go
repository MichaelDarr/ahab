package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PrintCmd prints a single command to the console
func PrintCmd(cmd *exec.Cmd) {
	userConfig, err := UserConfig()
	if err == nil && !userConfig.HideCommands {
		StylePrint("cyan", "$ "+strings.Join(cmd.Args, " "))
	}
}

// PrintDockerHelp parses args for a help flag, printing a help menu and running corresponsing docker help command if requested
func PrintDockerHelp(cmdArgs *[]string, dockerCmd string, helpString string) (helpRequested bool, err error) {
	for _, arg := range *cmdArgs {
		if arg == "-h" || arg == "--help" {
			helpRequested = true
			fmt.Println(helpString)
			helpArgs := []string{dockerCmd, "--help"}
			err = DockerCmd(&helpArgs)
		}
	}
	return
}

// PrintErr prints an error to the console (if non-nil)
func PrintErr(err error) {
	if errStr := err.Error(); errStr != "" {
		StylePrint("red", errStr)
	}
}

// PrintErrStr prints an error string to the console
func PrintErrStr(errStr string) {
	StylePrint("red", errStr)
}

// PrintErrFatal prints an error to the console and terminates the program (if non-nil)
func PrintErrFatal(err error) {
	if err != nil {
		PrintErr(err)
		os.Exit(1)
	}
}

// PrintIndentedPair prints a key/val pair in a human-readable format
func PrintIndentedPair(key string, val string) {
	fmt.Printf("  %-12s%s\n", key, val)
}

// PrintWarning prints a warning string to the console
func PrintWarning(warning string) {
	StylePrint("yellow", warning)
}

// StylePrint prints a string after surrounding it with appropriate style tags
func StylePrint(style string, str string) {
	stylizedStr := stylize(style, str)
	fmt.Println(stylizedStr)
}

// console style code map
var textCodes = map[string]string{
	"blue":    "\x1b[34m",
	"bold":    "\x1b[1m",
	"cyan":    "\x1b[36m",
	"green":   "\x1b[32m",
	"magenta": "\x1b[35m",
	"red":     "\x1b[31m",
	"reset":   "\x1b[0m",
	"yellow":  "\x1b[33m",
}

// appendToStrList is a helper for creating human-readable comma-separated lists
func appendToStrList(list string, newEl string) (finalStr string) {
	if list == "" {
		return newEl
	}
	return list + ", " + newEl
}

// stylize surrounds a string with style codes
func stylize(style string, str string) string {
	textCode, ok := textCodes[style]
	if ok {
		return textCode + str + textCodes["reset"]
	}
	PrintWarning("Unsupported text style: " + style)
	return str
}
