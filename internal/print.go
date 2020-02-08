package internal

import (
	"fmt"
	"os"
)

// PrintErr prints an error to the console (if non-nil)
func PrintErr(err error) {
	if errStr := err.Error(); errStr != "" {
		StylePrint("red", errStr)
	}
}

// PrintErrFatal prints an error to the console and terminates the program (if non-nil)
func PrintErrFatal(err error) {
	if err != nil {
		PrintErr(err)
		os.Exit(1)
	}
}

// PrintWarning prints a warning to the console. Severity: 0-dev, 1-moderate, 2-severe
// TODO: condition warning prints on severity/verbosity, set from config or CLI
func PrintWarning(severity int, warning string) {
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

// surround a string with style codes
func stylize(style string, str string) string {
	textCode, ok := textCodes[style]
	if ok {
		return textCode + str + textCodes["reset"]
	}
	PrintWarning(0, "Unsupported text style: "+style)
	return str
}
