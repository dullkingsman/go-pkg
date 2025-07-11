package rabbit_topology

import (
	"context"
	"fmt"
	"github.com/dullkingsman/go-pkg/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"time"
)

var BrokersConnections = map[string]*amqp.Connection{}

func LoadBrokerConnections(topologies BrokerTopologies, envFilePath string) {
	for _, cluster := range topologies {
		var url = os.Getenv(cluster.UrlKey)

		if url == "" {
			if !cluster.Optional {
				utils.LogFatal("rabbit-topology-connections-loader", "%s not set in environment file", cluster.UrlKey)
			}

			continue
		}

		for _, vhost := range cluster.Topology.VHosts {
			var completeUrl = url + "/" + vhost.Name

			utils.LogInfo("rabbit-topology-brokers-connections-loader", "connecting to broker "+utils.GreyString(completeUrl)+"...")

			var conn, err = amqp.Dial(completeUrl)

			if err != nil || conn == nil {
				var _err = completeUrl

				if err != nil {
					_err = err.Error()
				}

				var errorString = fmt.Sprintf("could not connect to broker " + utils.GreyString(url) + " at " + utils.GreyString(vhost.Name) + ": " + _err)

				utils.LogFatal("rabbit-topology-brokers-connections-loader", errorString)
			}

			utils.LogSuccess("rabbit-topology-brokers-connections-loader", "connected to broker "+utils.GreyString(url)+" at "+utils.GreyString(vhost.Name))

			BrokersConnections[cluster.Name+"/"+vhost.Name] = conn
		}
	}
}

func GetConnection(clusterName string, vhostName string) *amqp.Connection {
	return BrokersConnections[clusterName+"/"+vhostName]
}

func NewChannel(clusterName string, vhostName string) *amqp.Channel {
	var conn = GetConnection(clusterName, vhostName)

	var channel, err = conn.Channel()

	if err != nil {
		utils.LogError("rabbit-topology-channel-retriever", "could not create channel on "+utils.GreyString(clusterName)+": "+err.Error())
		return nil
	}

	return channel
}

func NewChannelWithContext(clusterName string, vhostName string, timeout ...time.Duration) (*amqp.Channel, context.Context, context.CancelFunc) {
	var _timeout time.Duration = 5

	if len(timeout) > 0 {
		_timeout = timeout[0]
	}

	var channel = NewChannel(clusterName, vhostName)

	if channel == nil {
		return nil, nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), _timeout*time.Second)

	return channel, ctx, cancel
}

func CloseConnections() {
	for clusterName, conn := range BrokersConnections {
		err := conn.Close()

		if err != nil {
			utils.LogError("rabbit-topology-connections-cleaner", "could not close broker connection: "+err.Error())
		}

		utils.LogInfo("rabbit-topology-connections-cleaner", "closed connection with broker "+utils.GreyString(clusterName))
	}
}
