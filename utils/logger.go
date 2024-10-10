package utils

import (
	"log"
	"runtime/debug"
)

const (
	Reset     = "\033[0;0m"
	Red       = "\033[0;31m"
	BrightRed = "\033[1;31m"
	Green     = "\033[0;32m"
	Yellow    = "\033[0;33m"
	Magenta   = "\033[0;35m"
	Cyan      = "\033[0;36m"
	Grey      = "\033[0;90m"
	LightBlue = "\033[0;94m"
)

// ColorString returns the specified string with the specified color.
func ColorString(string string, color string, disableColor ...bool) string {
	if len(disableColor) > 0 {
		if disableColor[0] {
			return string
		}
	}

	return color + string + Reset
}

// RedString returns a red colored string.
func RedString(string string, disableColor ...bool) string {
	return ColorString(string, Red, disableColor...)
}

// BrightRedString returns a bright red colored string.
func BrightRedString(string string, disableColor ...bool) string {
	return ColorString(string, BrightRed, disableColor...)
}

// GreenString returns a green colored string.
func GreenString(string string, disableColor ...bool) string {
	return ColorString(string, Green, disableColor...)
}

// YellowString returns a yellow colored string.
func YellowString(string string, disableColor ...bool) string {
	return ColorString(string, Yellow, disableColor...)
}

// MagentaString returns a magenta colored string.
func MagentaString(string string, disableColor ...bool) string {
	return ColorString(string, Magenta, disableColor...)
}

// CyanString returns a cyan colored string.
func CyanString(string string, disableColor ...bool) string {
	return ColorString(string, Cyan, disableColor...)
}

// GreyString returns a grey colored string.
func GreyString(string string, disableColor ...bool) string {
	return ColorString(string, Grey, disableColor...)
}

// LightBlueString returns a light blue colored string.
func LightBlueString(string string, disableColor ...bool) string {
	return ColorString(string, LightBlue, disableColor...)
}

// LogDebug logs a debug message.
func LogDebug(tag string, format string, v ...any) {
	log.Printf(GreyString(tag)+" "+format+"\n", v...)
}

// LogFatal logs an error message and exits the program with exit code 1.
func LogFatal(tag string, format string, v ...any) {
	log.Fatalf(BrightRedString(tag)+" "+format+"\n", v...)
}

// LogPanic logs an error message and exits the program with panic.
func LogPanic(tag string, format string, v ...any) {
	log.Panicf(BrightRedString(tag)+" "+format+"\n", v...)
}

// LogError logs an error message.
func LogError(tag string, format string, v ...any) {
	log.Printf(BrightRedString(tag)+" "+format+"\n", v...)
}

// LogErrorWithStack logs an error message with stack trace.
func LogErrorWithStack(tag string, format string, v ...any) {
	log.Printf(BrightRedString(tag)+" "+format+"\n", v...)
	log.Printf("%s", debug.Stack())
}

// LogInfo logs an info message.
func LogInfo(tag string, format string, v ...any) {
	log.Printf(LightBlueString(tag)+" "+format+"\n", v...)
}

// LogWarn logs a warning message.
func LogWarn(tag string, format string, v ...any) {
	log.Printf(YellowString(tag)+" "+format+"\n", v...)
}

// LogSuccess logs a success message.
func LogSuccess(tag string, format string, v ...any) {
	log.Printf(GreenString(tag)+" "+format+"\n", v...)
}
