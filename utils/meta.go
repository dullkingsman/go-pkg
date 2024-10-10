package utils

import (
	"go/format"
	"os"
	"path/filepath"
	"runtime"
)

func FormatGoCode(src string) (string, error) {
	formatted, err := format.Source([]byte(src))
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}

func FindGoConfigDir(compiled bool, startDir ...string) (string, string) {
	var execPath, err = os.Executable()

	if err != nil {
		LogFatal("go-config-finder", "could not get current file path: %s", err.Error())
	}

	if !compiled {
		var _, file, _, ok = runtime.Caller(0)

		if !ok {
			LogFatal("go-config-finder", "could not get current file path")
		}

		execPath = file
	}

	var _startDir = filepath.Dir(execPath)

	if len(startDir) > 0 {
		_startDir = startDir[0]
	}

	dir := _startDir

	var goConfigPath = "go.mod"

	if !compiled {
		goConfigPath = "go.work"
	}

	for {
		var goModPath = filepath.Join(dir, goConfigPath)

		if _, err := os.Stat(goModPath); !os.IsNotExist(err) {
			return _startDir, dir
		}

		var parentDir = filepath.Dir(dir)

		if parentDir == dir {
			LogFatal("go-config-finder", "could not get current file path")
		}

		dir = parentDir
	}
}
