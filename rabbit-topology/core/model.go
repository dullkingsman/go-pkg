package rabbit_topology

type BrokerTopologies = []BrokerCluster

type BrokerCluster struct {
	Name     string `json:"name"`
	UrlKey   string `json:"url_key"`
	Topology BrokerTopology
	Optional bool `json:"optional"`
}

type BrokerTopology struct {
	Exchanges []BrokerExchange `json:"exchanges,omitempty"`
	Queues    []BrokerQueue    `json:"queues,omitempty"`
	Bindings  []BrokerBinding  `json:"bindings,omitempty"`
	VHosts    []BrokerVhost    `json:"vhosts,omitempty"`
}

type BrokerVhost struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

type BrokerExchangeArgs struct {
	AlternateExchange *string `json:"alternate-exchange,omitempty"`
}

type BrokerExchange struct {
	Vhost       string                 `json:"vhost,omitempty"`
	Passive     bool                   `json:"passive,omitempty"`
	DisplayName string                 `json:"display_name,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Kind        string                 `json:"type,omitempty"`
	Durable     bool                   `json:"durable,omitempty"`
	AutoDelete  bool                   `json:"auto_delete,omitempty"`
	Internal    bool                   `json:"internal,omitempty"`
	NoWait      bool                   `json:"no_wait,omitempty"`
	Args        map[string]interface{} `json:"arguments,omitempty"`
}

type BrokerQueueArgs struct {
	MessageTtl           *int    `json:"x-message-ttl,omitempty"`
	Expires              *int    `json:"x-expires,omitempty"`
	Overflow             *string `json:"x-overflow,omitempty"` // Values can be drop-head, reject-publish or reject-publish-dlx
	SingleActiveConsumer *bool   `json:"x-single-active-consumer,omitempty"`
	DeadLetterExchange   *string `json:"x-dead-letter-exchange,omitempty"`
	DeadLetterRoutingKey *string `json:"x-dead-letter-routing-key,omitempty"`
	MaxLength            *int    `json:"x-max-length,omitempty"`
	MaxLengthBytes       *int    `json:"x-max-length-bytes,omitempty"`
	MaxPriority          *int    `json:"x-max-priority,omitempty"`
	QueueType            *string `json:"x-queue-type,omitempty"`
	QueueMode            *string `json:"x-queue-mode,omitempty"`
	QueueMasterLocator   *string `json:"x-queue-master-locator,omitempty"`
}

type BrokerQueue struct {
	Vhost       string                 `json:"vhost,omitempty"`
	Exchange    string                 `json:"exchange,omitempty"`
	Passive     bool                   `json:"passive,omitempty"`
	DisplayName string                 `json:"display_name,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Durable     bool                   `json:"durable,omitempty"`
	AutoDelete  bool                   `json:"auto_delete,omitempty"`
	Exclusive   bool                   `json:"exclusive,omitempty"`
	NoWait      bool                   `json:"no_wait,omitempty"`
	Args        map[string]interface{} `json:"arguments,omitempty"`
}

type BrokerBinding struct {
	Vhost           string                 `json:"vhost,omitempty"`
	DestinationType string                 `json:"destination_type,omitempty"` // Values can be queue or exchange
	Source          string                 `json:"source,omitempty"`
	Destination     string                 `json:"destination,omitempty"`
	RoutingKey      string                 `json:"routing_key,omitempty"`
	NoWait          bool                   `json:"no_wait,omitempty"`
	Args            map[string]interface{} `json:"arguments,omitempty"`
}
