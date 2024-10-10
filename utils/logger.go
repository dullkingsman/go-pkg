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

func ColorString(string string, color string, disableColor ...bool) string {
	if len(disableColor) > 0 {
		if disableColor[0] {
			return string
		}
	}

	return color + string + Reset
}

func RedString(string string, disableColor ...bool) string {
	return ColorString(string, Red, disableColor...)
}

func BrightRedString(string string, disableColor ...bool) string {
	return ColorString(string, BrightRed, disableColor...)
}

func GreenString(string string, disableColor ...bool) string {
	return ColorString(string, Green, disableColor...)
}

func YellowString(string string, disableColor ...bool) string {
	return ColorString(string, Yellow, disableColor...)
}

func MagentaString(string string, disableColor ...bool) string {
	return ColorString(string, Magenta, disableColor...)
}

func CyanString(string string, disableColor ...bool) string {
	return ColorString(string, Cyan, disableColor...)
}

func GreyString(string string, disableColor ...bool) string {
	return ColorString(string, Grey, disableColor...)
}

func LightBlueString(string string, disableColor ...bool) string {
	return ColorString(string, LightBlue, disableColor...)
}

func LogDebug(tag string, format string, v ...any) {
	log.Printf(GreyString(tag)+" "+format+"\n", v...)
}

func LogFatal(tag string, format string, v ...any) {
	log.Fatalf(BrightRedString(tag)+" "+format+"\n", v...)
}

func LogPanic(tag string, format string, v ...any) {
	log.Panicf(BrightRedString(tag)+" "+format+"\n", v...)
}

func LogError(tag string, format string, v ...any) {
	log.Printf(BrightRedString(tag)+" "+format+"\n", v...)
}

func LogErrorWithStack(tag string, format string, v ...any) {
	log.Printf(BrightRedString(tag)+" "+format+"\n", v...)
	log.Printf("%s", debug.Stack())
}

func LogInfo(tag string, format string, v ...any) {
	log.Printf(LightBlueString(tag)+" "+format+"\n", v...)
}

func LogWarn(tag string, format string, v ...any) {
	log.Printf(YellowString(tag)+" "+format+"\n", v...)
}

func LogSuccess(tag string, format string, v ...any) {
	log.Printf(GreenString(tag)+" "+format+"\n", v...)
}
