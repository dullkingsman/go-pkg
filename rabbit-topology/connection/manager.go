package connection

import (
	"context"
	"errors"
	"github.com/dullkingsman/go-pkg/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
	"time"
)

func (p *ChannelPool) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	var channel, err = p.GetChannel()

	if err != nil {
		return err
	}

	return channel.Publish(exchange, key, mandatory, immediate, msg)
}

func (p *ChannelPool) PublishWithContent(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	var channel, err = p.GetChannel()

	if err != nil {
		return err
	}

	return channel.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
}

func (p *ChannelPool) PublishWithDeferredConfirm(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) (*amqp.DeferredConfirmation, error) {
	var channel, err = p.GetChannel()

	if err != nil {
		return nil, err
	}

	return channel.PublishWithDeferredConfirm(exchange, key, mandatory, immediate, msg)
}

func (p *ChannelPool) PublishWithDeferredConfirmWithContent(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) (*amqp.DeferredConfirmation, error) {
	var channel, err = p.GetChannel()

	if err != nil {
		return nil, err
	}

	return channel.PublishWithDeferredConfirmWithContext(ctx, exchange, key, mandatory, immediate, msg)
}

func (p *ChannelPool) Get(queue string, autoAck bool) (msg amqp.Delivery, ok bool, err error) {
	channel, err := p.GetChannel()

	if err != nil {
		return amqp.Delivery{}, false, err
	}

	return channel.Get(queue, autoAck)
}

func (p *ChannelPool) GlobalQos(prefetchCount, prefetchSize int) error {
	channel, err := p.GetChannel()

	if err != nil {
		return err
	}

	return channel.Qos(prefetchSize, prefetchCount, true)
}

func (p *ChannelPool) Consume(
	queue string,
	consumer string,
	autoAck bool,
	exclusive bool,
	noLocal bool,
	noWait bool,
	args amqp.Table,
	prefetchCount,
	prefetchSize int,
	handler func(message amqp.Delivery),
	exchange ...string,
) error {
	if p == nil {
		return errors.New("channel pool is not initialized")
	}

	var channel, err = p.GetChannel()

	if err != nil || channel == nil {
		if err == nil {
			err = errors.New("obtained channel is nil")
		}

		return err
	}

	if queue == "" && (len(exchange) <= 0 || exchange[0] == "") {
		return errors.New("queue not specified")
	}

	if queue == "" {
		q, err := channel.QueueDeclare(
			"",
			false,
			true,
			true,
			false,
			nil,
		)

		if err != nil {
			return err
		}

		err = channel.QueueBind(
			q.Name,
			"",
			exchange[0],
			false,
			nil,
		)

		if err != nil {
			return err
		}

		queue = q.Name
	}

	if prefetchCount > 0 || prefetchSize > 0 {
		err = channel.Qos(prefetchCount, prefetchSize, false)

		if err != nil {
			return err
		}
	}

	messages, err := channel.Consume(
		queue,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)

	if err != nil {
		return err
	}

	utils.LogInfo("consumer("+consumer+")", "listening for messages...")

	go func() {
		defer func() {
			err = p.ReturnChannel(channel)

			if err != nil {
				utils.LogError("consumer("+consumer+")", "could not return channel")
			}
		}()

		for message := range messages {
			handler(message)
		}
	}()

	return nil
}

// GetChannel fetches a channel from the pool, creating one if necessary.
func (p *ChannelPool) GetChannel(_timeout ...time.Duration) (*Channel, error) {
	if p.conn.IsClosed() {
		utils.LogInfo("channel-pool(fetcher)", "connection is closed: reconnecting...")

		var err = ConnectToBroker(p.url, p.cluster, p.vhost, p.dependents, p)

		if err != nil {
			return nil, err
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	var timeout = 5 * time.Second

	if len(_timeout) > 0 {
		timeout = _timeout[0]
	}

	// First, try to get an available channel immediately
	for {
		select {
		case ch := <-p.available:
			var err = ch.Usable()

			if err != nil {
				return nil, err
			}

			if ch.IsClosed() {
				continue // Skip closed channels and retry
			}

			p.inUse[ch] = struct{}{}

			return ch, nil
		default:
			// If the pool is not full, create a new channel
			if len(p.inUse)+len(p.available) < p.maxSize {
				return p.newChannelLocked()
			}
		}

		// Unlock before waiting to avoid deadlocks
		p.mu.Unlock()

		select {
		case ch, ok := <-p.available:
			var err = ch.Usable()

			if err != nil {
				return nil, err
			}

			if !ok {
				return nil, errors.New("channel pool closed")
			}

			p.mu.Lock() // Reacquire lock before modifying state

			if ch.IsClosed() {
				continue // Skip closed channels and retry
			}

			p.inUse[ch] = struct{}{}

			return ch, nil

		case <-time.After(timeout): // Timeout reached
			return nil, errors.New("timeout waiting for available channel")
		}
	}
}

// ReturnChannel returns a channel to the pool.
func (p *ChannelPool) ReturnChannel(ch *Channel) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.inUse[ch]; !ok {
		return errors.New("channel does not belong to the pool")
	}

	delete(p.inUse, ch)

	if ch.IsClosed() {
		return nil
	}

	select {
	case p.available <- ch:
		return nil
	default:
		_ = ch.Close() // If pool is full, close the channel
		return nil
	}
}

// Close closes the pool and all its channels.
func (p *ChannelPool) Close(reconnection bool) {
	if !reconnection {
		close(p.closeSignal)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.available != nil {
		close(p.available)
		p.available = nil
	}

	for ch := range p.inUse {
		if err := ch.Close(); err != nil {
			utils.LogError("channel-pool(cleaner)", "could not close channel: "+err.Error())
		}
	}

	p.inUse = make(map[*Channel]struct{})

	if reconnection == false {
		var err = p.conn.Close()

		if err != nil {
			utils.LogError("channel-pool(cleaner)", "could not close broker connection: "+err.Error())
		}

		utils.LogInfo("channel-pool(cleaner)", "closed connection with broker "+utils.GreyString(p.cluster.Name))
	}
}

// newChannelLocked creates a new channel.
func (p *ChannelPool) newChannelLocked() (*Channel, error) {
	if p.conn.IsClosed() {
		return nil, errors.New("connection is closed")
	}

	ch, err := p.conn.Channel()

	if err != nil {
		return nil, err
	}

	if ch == nil {
		return nil, errors.New("channel is nil")
	}

	var channel = &Channel{ch}

	p.inUse[channel] = struct{}{}

	return channel, nil
}

// cleanupIdleChannels periodically closes idle channels.
func (p *ChannelPool) cleanupIdleChannels() {
	var ticker = time.NewTicker(p.idleTimeout)

	defer ticker.Stop()

	for {
		select {
		case <-p.closeSignal:
			utils.LogInfo("channel-pool(cleaner)", "closing idle channel cleaner...")
			return
		case <-ticker.C:
			if p.conn.IsClosed() {
				utils.LogInfo("channel-pool(cleaner)", "connection is closed: reconnecting...")
				_ = ConnectToBroker(p.url, p.cluster, p.vhost, p.dependents, p)
			}

			p.mu.Lock()

			// process idle channels
			for i := 0; i < len(p.available); i++ {
				select {
				case ch := <-p.available:
					if ch.IsClosed() {
						_ = ch.Close()
					} else {
						select {
						case p.available <- ch: // return to pool if still valid
						default:
							_ = ch.Close() // close if the pool is full
						}
					}
				default:
					break
				}
			}

			p.mu.Unlock()
		}
	}
}

func (c *Channel) Usable(err ...error) error {
	if len(err) > 0 && err[0] != nil {
		return err[0]
	}

	if c.Channel == nil {
		return errors.New("channel is nil")
	}

	return nil
}

var BrokerChannelPools = map[string]*ChannelPool{}

func GetChannelPool(clusterName string, vhostName string) *ChannelPool {
	return BrokerChannelPools[clusterName+"/"+vhostName]
}

func CloseConnections() {
	for _, pool := range BrokerChannelPools {
		pool.Close(false)
	}
}

func ConnectToBroker(
	url string,
	cluster *BrokerCluster,
	vhost *BrokerVhost,
	dependents []func(reconnection bool),
	oldChannelPool ...*ChannelPool,
) error {
	if len(oldChannelPool) > 0 {
		oldChannelPool[0].Close(true)
	}

	var conn, err = NewConnection(url, vhost)

	if err != nil {
		return err
	}

	pool, err := NewChannelPool(conn, int(conn.Config.ChannelMax), 90*time.Second, oldChannelPool...)

	if err != nil || pool == nil {
		if err == nil {
			err = errors.New("channel pool is not initialized")
		}

		return err
	}

	pool.cluster = cluster
	pool.vhost = vhost
	pool.url = url
	pool.dependents = dependents

	utils.LogSuccess("channel-pool(connection-manager)", "connected to broker "+utils.GreyString(url)+" at "+utils.GreyString(vhost.Name))

	utils.LogInfo("channel-pool(connection-manager)", "maximum allowed channels on this connection are "+utils.GreyString(strconv.Itoa(int(pool.conn.Config.ChannelMax))))

	BrokerChannelPools[cluster.Name+"/"+vhost.Name] = pool

	for _, dependent := range dependents {
		dependent(len(oldChannelPool) > 0)
	}

	return nil
}

func NewConnection(url string, vhost *BrokerVhost) (conn *amqp.Connection, err error) {
	var completeUrl = url + "/" + vhost.Name

	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(completeUrl)

		if err != nil || conn == nil {
			if err == nil {
				err = errors.New("connection is nil")
			}

			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	return
}

// NewChannelPool creates a new pool with the specified max size and idle timeout.
func NewChannelPool(
	conn *amqp.Connection,
	maxSize int,
	idleTimeout time.Duration,
	oldChannelPool ...*ChannelPool,
) (*ChannelPool, error) {
	if maxSize <= 0 {
		return nil, errors.New("channel pool max size must be greater than 0")
	}

	var pool = &ChannelPool{}

	if len(oldChannelPool) > 0 {
		pool = oldChannelPool[0]
	}

	pool.conn = conn
	pool.maxSize = maxSize
	pool.idleTimeout = idleTimeout
	pool.closeSignal = make(chan struct{})
	pool.inUse = make(map[*Channel]struct{})
	pool.available = make(chan *Channel, maxSize)

	if len(oldChannelPool) == 0 {
		go pool.cleanupIdleChannels()
	}

	return pool, nil
}
