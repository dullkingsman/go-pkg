package main

import (
	prizzle "github.com/dullkingsman/go-pkg/prizzle/core"
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	var args = os.Args[1:]

	var (
		driver         = "postgres"
		goExec         = false
		envFilePath    = ".env"
		schemaFilePath = "model.prisma"
		currentCommand = ""
	)

	var argsLength = len(args)

	if argsLength > 0 {
		var skipNext = false

		for index, arg := range args {
			if skipNext {
				skipNext = false
				continue
			}

			switch arg {
			case "generate", "construct", "deconstruct":
				currentCommand = arg
			case "--help":
				printGeneralHelp(0)
			case "--go-exec", "-g":
				goExec = true
			case "--driver", "-d":
				checkOptionValue("--driver", index, args)
				driver = args[index+1]
				skipNext = true
			case "--schema", "-s":
				checkOptionValue("--schema", index, args)
				schemaFilePath = args[index+1]
				skipNext = true
			case "--env", "-e":
				checkOptionValue("--env", index, args)
				envFilePath = args[index+1]
				skipNext = true
			}
		}
	}

	if envLoadErr := godotenv.Load(envFilePath); envLoadErr != nil {
		utils.LogFatal("prizzle-connections-loader", "could not load environment file: "+envLoadErr.Error())
	}

	if currentCommand == "" {
		utils.LogError("prizzle", "no valid command was provided")
		printGeneralHelp(1)
	}

	if goExec {
		utils.LogInfo("prizzle", "being run using 'go run'")
	}

	if !isSupportedDriver(driver) {
		utils.LogError("prizzle", "unsupported driver: "+driver)
		printGeneralHelp(1)
	}

	var client = prizzle.LoadDatabaseClusterNode(driver, prizzle.ClusterNodeConfig{
		Url: os.Getenv("DATABASE_URL"),
	})

	defer prizzle.CloseDatabaseClusterNode(client)

	switch currentCommand {
	case "generate":
		prizzle.GenerateClientModel(driver, client, schemaFilePath)
	}
}

func isSupportedDriver(driver string) bool {
	return driver == "postgres" || driver == "sqlite3"
}

func isValidOption(value string) bool {
	return value == "--env" ||
		value == "-e" ||
		value == "--driver" ||
		value == "-d" ||
		value == "--schema" ||
		value == "-s" ||
		value == "--go-exec" ||
		value == "-g" ||
		value == "--help"
}

func checkOptionValue(option string, index int, args []string) {
	if index >= len(args)-1 {
		utils.LogError("prizzle", "value for "+option+" not set")
		printGeneralHelp(1)
	}

	if !isValidOption(args[index]) {
		utils.LogError("prizzle", "invalid command option used as value for "+option)
		printGeneralHelp(1)
	}
}

func printGeneralHelp(exitCode int) {
	var generalHelp = `Usage: go run github.com/dullkingsman/go-pkg/prizzle [COMMAND] [OPTIONS]

Commands:
  generate [OPTIONS]  Generates go definitions for the database schema

Options:
  -e, --env    	PATH  Specify the .env file to pull connections urls and other configurations from
  -d, --driver  DRIVER  Specify the database driver to use (postgres, sqlite3)
  -s, --schema  PATH  Specify the topologies.json file that defines the broker 
  -g, --go-exec   	  Specify whether it is being run using 'go run' or as a pre-generated binary
  --help          	  Show this help message and exit

Examples:
  go run github.com/dullkingsman/go-pkg/prizzle construct -e .env --schema model.prisma
  go run github.com/dullkingsman/go-pkg/prizzle --help` + "\n"

	utils.PrintHelp(generalHelp, exitCode)
}
