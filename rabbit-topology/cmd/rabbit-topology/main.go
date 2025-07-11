package main

import (
	rabbittopology "github.com/dullkingsman/go-pkg/rabbit-topology/core"
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func main() {
	var args = os.Args[1:]

	var (
		goExec             = false
		envFilePath        = ".env"
		topologiesFilePath = "topologies.json"
		currentCommand     = ""
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
			case "--topologies", "-t":
				checkOptionValue("--topologies", index, args)
				topologiesFilePath = args[index+1]
				skipNext = true
			case "--env", "-e":
				checkOptionValue("--env", index, args)
				envFilePath = args[index+1]
				skipNext = true
			}
		}
	}

	if envLoadErr := godotenv.Load(envFilePath); envLoadErr != nil {
		utils.LogFatal("rabbit-topology-connections-loader", "could not load environment file: "+envLoadErr.Error())
	}

	if currentCommand == "" {
		utils.LogError("rabbit-topology", "no valid command was provided")
		printGeneralHelp(1)
	}

	if goExec {
		utils.LogInfo("rabbit-topology", "being run using 'go run'")
	}

	if !filepath.IsAbs(topologiesFilePath) {
		getwd, err := os.Getwd()

		if err != nil {
			utils.LogError("rabbit-topology", "could not get the current working directory: "+err.Error())
			printGeneralHelp(1)
		}

		topologiesFilePath = filepath.Clean(getwd + "/" + topologiesFilePath)
	}

	var topologies = rabbittopology.LoadBrokerTopologies(topologiesFilePath)

	if currentCommand != "generate" {
		rabbittopology.LoadBrokerConnections(topologies, envFilePath)
		defer rabbittopology.CloseConnections()
	}

	switch currentCommand {
	case "generate":
		rabbittopology.GenerateClientModel(topologies, topologiesFilePath)
	case "construct":
		rabbittopology.ConstructTopologies(topologies)
	default:
		rabbittopology.DeconstructTopologies(topologies)
	}

}

func isValidOption(value string) bool {
	return value == "--env" ||
		value == "-e" ||
		value == "--topologies" ||
		value == "-t" ||
		value == "--go-exec" ||
		value == "-g" ||
		value == "--help"
}

func checkOptionValue(option string, index int, args []string) {
	if index >= len(args)-1 {
		utils.LogError("rabbit-topology", "value for "+option+" not set")
		printGeneralHelp(1)
	}

	if !isValidOption(args[index+1]) {
		utils.LogError("rabbit-topology", "invalid command option used as value for "+option)
		printGeneralHelp(1)
	}
}

func printGeneralHelp(exitCode int) {
	var generalHelp = `Usage: go run github.com/dullkingsman/go-pkg/rabbit-topology [COMMAND] [OPTIONS]

Commands:
  generate    [OPTIONS]  Generates go definitions for the topologies
  construct   [OPTIONS]  Constructs the topology in the broker
  deconstruct [OPTIONS]  Cleans up the topology in the broker

Options:
  -e, --env        PATH  Specify the .env file to pull connections urls and other configurations from
  -t, --topologies PATH  Specify the topologies.json file that defines the broker 
  -g, --go-exec   		 Specify whether it is being run using 'go run' or as a pre-generated binary
  --help          		 Show this help message and exit

Examples:
  go run github.com/dullkingsman/go-pkg/rabbit-topology construct -e .env --topologies topologies.json
  go run github.com/dullkingsman/go-pkg/rabbit-topology --help` + "\n"

	utils.PrintHelp(generalHelp, exitCode)
}
