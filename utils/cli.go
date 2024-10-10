package utils

import "os"

// PrintHelp prints the given help text and exits with the given exit code.
func PrintHelp(text string, exitCode int) {
	if exitCode != 0 {
		text = "\n" + text
	}

	print(text)

	os.Exit(exitCode)
}
