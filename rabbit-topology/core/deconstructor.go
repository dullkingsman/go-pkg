package rabbit_topology

import (
	"github.com/dullkingsman/go-pkg/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeconstructTopology(cluster BrokerCluster) {
	var channels = map[string]*amqp.Channel{}

	for _, binding := range cluster.Topology.Bindings {
		var channelFor = cluster.Name + "/" + binding.Vhost

		channel, ok := channels[channelFor]

		if !ok || channel == nil {
			channel = NewChannel(cluster.Name, binding.Vhost)
		}

		var err error

		if binding.DestinationType == "queue" {
			err = channel.QueueUnbind(
				binding.Destination,
				binding.RoutingKey,
				binding.Source,
				binding.Args,
			)
		} else {
			err = channel.ExchangeBind(
				binding.Destination,
				binding.RoutingKey,
				binding.Source,
				binding.NoWait,
				binding.Args,
			)
		}

		if err != nil {
			utils.LogFatal("rabbit-topology-deconstructor", "failed to unbind "+binding.DestinationType+" "+utils.GreyString(binding.Destination)+" to exchange "+utils.GreyString(binding.Source)+" on "+utils.GreyString(channelFor)+": "+err.Error())
		}

		utils.LogSuccess("rabbit-topology-deconstructor", "unbound queue "+utils.GreyString(binding.Destination)+" to exchange "+utils.GreyString(binding.Source)+" on "+utils.GreyString(channelFor))
	}

	for _, queue := range cluster.Topology.Queues {
		var channelFor = cluster.Name + "/" + queue.Vhost

		channel, ok := channels[channelFor]

		if !ok || channel == nil {
			channel = NewChannel(cluster.Name, queue.Vhost)
		}

		var err error

		_, err = channel.QueueDelete(
			queue.Name,
			false,
			false,
			queue.NoWait,
		)

		if err != nil {
			utils.LogFatal("rabbit-topology-deconstructor", "failed to delete queue "+utils.GreyString(queue.DisplayName)+" on "+utils.GreyString(channelFor)+": "+err.Error())
		}

		utils.LogSuccess("rabbit-topology-deconstructor", "deleted queue "+utils.GreyString(queue.DisplayName)+" on "+utils.GreyString(channelFor))
	}

	for _, exchange := range cluster.Topology.Exchanges {
		var channelFor = cluster.Name + "/" + exchange.Vhost

		channel, ok := channels[channelFor]

		if !ok || channel == nil {
			channel = NewChannel(cluster.Name, exchange.Vhost)
		}

		var err error

		err = channel.ExchangeDelete(
			exchange.Name,
			false,
			exchange.NoWait,
		)

		if err != nil {
			utils.LogFatal("rabbit-topology-deconstructor", "failed to delete exchange "+utils.GreyString(exchange.DisplayName)+" on "+utils.GreyString(channelFor)+": "+err.Error())
		}

		utils.LogSuccess("rabbit-topology-deconstructor", "deleted exchange "+utils.GreyString(exchange.DisplayName)+" on "+utils.GreyString(channelFor))
	}
}

func DeconstructTopologies(topologies BrokerTopologies) {
	for _, cluster := range topologies {
		DeconstructTopology(cluster)
	}
}
