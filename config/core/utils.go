package core

import (
	"github.com/dullkingsman/go-pkg/utils"
	"os"
	"path/filepath"
)

func GetExecutionPaths(goExec bool) (string, string) {
	var startDir, sourceRoot = FindGoConfigDir(goExec)

	utils.LogInfo("server", "StartDir: "+utils.GreyString(startDir))
	utils.LogInfo("server", "SourceRoot: "+utils.GreyString(sourceRoot))

	utils.LogSuccess("server", "loaded paths")

	return startDir, sourceRoot
}

func FindGoConfigDir(workspace bool, startDir ...string) (string, string) {
	var _startDir, err = os.Getwd()

	if err != nil {
		utils.LogFatal("go-config-finder", "could not get current working directory: "+err.Error())
	}

	if len(startDir) > 0 {
		_startDir = startDir[0]
	}

	var dir = _startDir

	var goConfigPath = "go.mod"

	if workspace {
		goConfigPath = "go.work"
	}

	for {
		var goModPath = filepath.Join(dir, goConfigPath)

		if _, err := os.Stat(goModPath); !os.IsNotExist(err) {
			return _startDir, dir
		}

		var parentDir = filepath.Dir(dir)

		if parentDir == dir {
			utils.LogWarn("go-config-finder", "could not get go config file")
			return _startDir, _startDir
		}

		dir = parentDir
	}
}
