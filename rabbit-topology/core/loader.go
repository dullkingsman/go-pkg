package rabbit_topology

import (
	"encoding/json"
	"github.com/dullkingsman/go-pkg/utils"
	"os"
	"path/filepath"
)

func LoadBrokerTopologies(topologiesFilePath string) BrokerTopologies {
	var filePath = filepath.Join(topologiesFilePath)

	fileContent, err := os.ReadFile(filePath)

	if err != nil {
		utils.LogFatal("rabbit-topology-loader", "could not read %s", utils.GreyString(filePath))
	}

	var topologies BrokerTopologies

	if err := json.Unmarshal(fileContent, &topologies); err != nil {
		utils.LogFatal("rabbit-topology-loader", "could not unmarshal the contents of %s", utils.GreyString(filePath))
	}

	utils.LogSuccess("rabbit-topology-loader", "loaded brokers topologies")

	return topologies
}
