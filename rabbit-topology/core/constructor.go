package rabbit_topology

import (
	"github.com/dullkingsman/go-pkg/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConstructTopology(cluster BrokerCluster) {
	var channels = map[string]*amqp.Channel{}

	for _, exchange := range cluster.Topology.Exchanges {
		var channelFor = cluster.Name + "/" + exchange.Vhost

		channel, ok := channels[channelFor]

		if !ok || channel == nil {
			channel = NewChannel(cluster.Name, exchange.Vhost)
		}

		var err error

		if !exchange.Passive {
			err = channel.ExchangeDeclare(
				exchange.Name,
				exchange.Kind,
				exchange.Durable,
				exchange.AutoDelete,
				exchange.Internal,
				exchange.NoWait,
				exchange.Args,
			)
		} else {
			err = channel.ExchangeDeclarePassive(
				exchange.Name,
				exchange.Kind,
				exchange.Durable,
				exchange.AutoDelete,
				exchange.Internal,
				exchange.NoWait,
				exchange.Args,
			)
		}

		if err != nil {
			utils.LogFatal("rabbit-topology-constructor", "failed to declare exchange "+utils.GreyString(exchange.DisplayName)+" on "+utils.GreyString(channelFor)+": "+err.Error())
		}

		utils.LogSuccess("rabbit-topology-constructor", "declared exchange "+utils.GreyString(exchange.DisplayName)+" on "+utils.GreyString(channelFor))
	}

	for _, queue := range cluster.Topology.Queues {
		var channelFor = cluster.Name + "/" + queue.Vhost

		channel, ok := channels[channelFor]

		if !ok || channel == nil {
			channel = NewChannel(cluster.Name, queue.Vhost)
		}

		var err error

		if !queue.Passive {
			_, err = channel.QueueDeclare(
				queue.Name,
				queue.Durable,
				queue.AutoDelete,
				queue.Exclusive,
				queue.NoWait,
				queue.Args,
			)
		} else {
			_, err = channel.QueueDeclarePassive(
				queue.Name,
				queue.Durable,
				queue.AutoDelete,
				queue.Exclusive,
				queue.NoWait,
				queue.Args,
			)
		}

		if err != nil {
			utils.LogFatal("rabbit-topology-constructor", "failed to declare queue "+utils.GreyString(queue.DisplayName)+" on "+utils.GreyString(channelFor)+": "+err.Error())
		}

		utils.LogSuccess("rabbit-topology-constructor", "declared queue "+utils.GreyString(queue.DisplayName)+" on "+utils.GreyString(channelFor))
	}

	for _, binding := range cluster.Topology.Bindings {
		var channelFor = cluster.Name + "/" + binding.Vhost

		channel, ok := channels[channelFor]

		if !ok || channel == nil {
			channel = NewChannel(cluster.Name, binding.Vhost)
		}

		var err error

		if binding.DestinationType == "queue" {
			err = channel.QueueBind(
				binding.Destination,
				binding.RoutingKey,
				binding.Source,
				binding.NoWait,
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
			utils.LogFatal("rabbit-topology-constructor", "failed to bind "+binding.DestinationType+" "+utils.GreyString(binding.Destination)+" to exchange "+utils.GreyString(binding.Source)+" on "+utils.GreyString(channelFor)+": "+err.Error())
		}

		utils.LogSuccess("rabbit-topology-constructor", "bound queue "+utils.GreyString(binding.Destination)+" to exchange "+utils.GreyString(binding.Source)+" on "+utils.GreyString(channelFor))
	}
}

func ConstructTopologies(topologies BrokerTopologies) {
	for _, cluster := range topologies {
		ConstructTopology(cluster)
	}
}
