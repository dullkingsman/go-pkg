package rabbit_topology

import (
	"github.com/dullkingsman/go-pkg/utils"
	"os"
	"path/filepath"
)

func GenerateClientModel(topologies BrokerTopologies, topologiesFilePath string) {
	var buffer = "package client\n\n"

	for _, cluster := range topologies {
		utils.LogInfo("rabbit-topology-generator", "constructing topology model for "+cluster.Name)

		var tmp = GetTopologyHierarchy(cluster)

		var typesBuffer = clusterTypeBuffer{
			Cluster:   "",
			VHosts:    map[string]string{},
			Exchanges: map[string]string{},
			Queues:    map[string]string{},
		}

		var clusterName = utils.AnyToPascalCase(cluster.Name)

		typesBuffer.Cluster = "type cluster" + clusterName + " struct {\n\tName string\n"

		for vhost, exchanges := range tmp {
			var vhostName = utils.AnyToPascalCase(vhost)

			typesBuffer.Cluster += "\t" + vhostName + " vhost" + vhostName + " \n"

			typesBuffer.VHosts[vhost] = "type vhost" + vhostName + " struct {\n\tName string\n"

			for exchange, queues := range exchanges {
				var exchangeName = utils.AnyToPascalCase(exchange)

				typesBuffer.VHosts[vhost] += "\t" + exchangeName + " exchange" + exchangeName + " \n"

				typesBuffer.Exchanges[exchange] = "type exchange" + exchangeName + " struct {\n\tName string\n"

				for queue, bindings := range queues {
					var queueName = utils.AnyToPascalCase(queue)

					typesBuffer.Exchanges[exchange] += "\t" + queueName + " queue" + queueName + " \n"

					typesBuffer.Queues[queue] = "type queue" + queueName + " struct {\n\tName string\n"

					for routingKey := range bindings {
						typesBuffer.Queues[queue] += "\t" + utils.AnyToPascalCase(routingKey) + " string \n"
					}

					typesBuffer.Queues[queue] += "\n}\n\n"
				}

				typesBuffer.Exchanges[exchange] += "\n}\n\n"
			}

			typesBuffer.VHosts[vhost] += "\n}\n\n"
		}

		typesBuffer.Cluster += "\n}\n\n"

		var clusterDefBuffer = typesBuffer.Cluster

		for _, vhost := range typesBuffer.VHosts {
			clusterDefBuffer += vhost
		}

		for _, exchange := range typesBuffer.Exchanges {
			clusterDefBuffer += exchange
		}

		for _, queue := range typesBuffer.Queues {
			clusterDefBuffer += queue
		}

		clusterDefBuffer += "\n"

		clusterDefBuffer += "var " + clusterName + " = cluster" + clusterName + "{\n"
		clusterDefBuffer += "\tName: \"" + cluster.Name + "\",\n"

		for vhost, exchanges := range tmp {
			var vhostName = utils.AnyToPascalCase(vhost)

			clusterDefBuffer += "\t" + vhostName + " : vhost" + vhostName + "{\n"
			clusterDefBuffer += "\t\tName: \"" + vhost + "\",\n"

			for exchange, queues := range exchanges {
				var exchangeName = utils.AnyToPascalCase(exchange)

				clusterDefBuffer += "\t\t" + exchangeName + ": exchange" + exchangeName + "{\n"
				clusterDefBuffer += "\t\t\tName: \"" + exchange + "\",\n"

				for queue, bindings := range queues {
					var queueName = utils.AnyToPascalCase(queue)

					clusterDefBuffer += "\t\t\t" + queueName + ": queue" + queueName + "{\n"
					clusterDefBuffer += "\t\t\t\tName: \"" + queue + "\",\n"

					for routingKey := range bindings {
						clusterDefBuffer += "\t\t\t\t" + utils.AnyToPascalCase(routingKey) + ": \"" + routingKey + "\",\n"
					}

					clusterDefBuffer += "\t\t\t},\n"
				}

				clusterDefBuffer += "\t\t},\n"
			}

			clusterDefBuffer += "\t},\n"
		}

		clusterDefBuffer += "}\n"

		buffer += clusterDefBuffer
	}

	formatted, err := utils.FormatAsGoCode(buffer)

	if err != nil {
		utils.LogFatal("rabbit-topology-generator", "could not format code: "+err.Error())
	}

	var topologiesDir = filepath.Dir(topologiesFilePath)

	var filePath = filepath.Join(topologiesDir + "/broker/client/definition.go")

	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		utils.LogFatal("rabbit-topology-generator", "could not create directories: "+err.Error())
	}

	if err := utils.WriteToFile(filePath, formatted); err != nil {
		utils.LogFatal("rabbit-topology-generator", "could not write to file: "+err.Error())
	}
}

func GetTopologyHierarchy(cluster BrokerCluster) VhostHierarchy {
	var tmp = VhostHierarchy{}

	var vhosts = cluster.Topology.VHosts

	for _, vhost := range vhosts {
		tmp[vhost.Name] = ExchangeHierarchy{}
	}

	var exchanges = cluster.Topology.Exchanges

	for _, exchange := range exchanges {
		if _, ok := tmp[exchange.Vhost]; !ok {
			utils.LogFatal("rabbit-topology-generator", "vhost "+exchange.Vhost+" not found")
		}

		tmp[exchange.Vhost][exchange.Name] = QueueHierarchy{}
	}

	var queues = cluster.Topology.Queues

	for _, queue := range queues {
		if _, ok := tmp[queue.Vhost]; !ok {
			utils.LogFatal("rabbit-topology-generator", "vhost "+queue.Vhost+" not found")
		}

		if _, ok := tmp[queue.Vhost][queue.Exchange]; !ok {
			utils.LogFatal("rabbit-topology-generator", "exchange "+queue.Exchange+" not found")
		}

		tmp[queue.Vhost][queue.Exchange][queue.Name] = BindingHierarchy{}
	}

	var bindings = cluster.Topology.Bindings

	for _, binding := range bindings {
		if _, ok := tmp[binding.Vhost]; !ok {
			utils.LogFatal("rabbit-topology-generator", "vhost "+binding.Vhost+" not found")
		}

		if _, ok := tmp[binding.Vhost][binding.Source]; !ok {
			utils.LogFatal("rabbit-topology-generator", "exchange "+binding.Source+" not found")
		}

		if _, ok := tmp[binding.Vhost][binding.Source][binding.Destination]; !ok {
			utils.LogFatal("rabbit-topology-generator", "queue "+binding.Destination+" not found")
		}

		tmp[binding.Vhost][binding.Source][binding.Destination][binding.RoutingKey] = binding.RoutingKey
	}

	return tmp
}

type VhostHierarchy = map[string]ExchangeHierarchy

type ExchangeHierarchy = map[string]QueueHierarchy

type QueueHierarchy = map[string]BindingHierarchy

type BindingHierarchy = map[string]string

type clusterTypeBuffer struct {
	Cluster   string
	VHosts    map[string]string
	Exchanges map[string]string
	Queues    map[string]string
}
