package broker

import (
	"github.com/streadway/amqp"
	"sync"
	"strings"
	"fmt"
)

type RabbitMqQueueProperties struct {
	Name string
	Durable bool
	AutoDelete bool
	Exclusive bool
	NoWait bool
}

type RabbitMqConsumerProperties struct {
	Name string
	QueueName string
	AutoAck bool
	Exclusive bool
	NoLocal bool
	NoWait bool
}

type ConsumerHandlerFunc func(channel *amqp.Channel, delivery *amqp.Delivery)

type ConsumerChannel struct {
	channel <-chan amqp.Delivery
	handler ConsumerHandlerFunc
}

type RabbitMqBroker struct {
	address string
	connection *amqp.Connection
	channel *amqp.Channel
	queues map[string]*amqp.Queue
	consumers map[string]*ConsumerChannel
}


func NewRabbitMqBroker(address string, prefetchCount, prefetchSize int) (*RabbitMqBroker, error) {
	connection, err := amqp.Dial(address)
	if err != nil {
		return nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	err = channel.Qos(prefetchCount, prefetchSize, false)
	if err != nil {
		return nil, err
	}
	return &RabbitMqBroker{
		address: address,
		connection: connection,
		channel:channel,
		queues: make(map[string]*amqp.Queue),
		consumers: make(map[string]*ConsumerChannel)}, nil
}

func (b *RabbitMqBroker) Address() string {
	return b.address
}

func (b *RabbitMqBroker) Channel() *amqp.Channel {
	return b.channel
}

func (b *RabbitMqBroker) Queue(id string) (*amqp.Queue, bool) {
	q, ok := b.queues[id]
	return q, ok
}

func (b *RabbitMqBroker) WithQueue(id string, p *RabbitMqQueueProperties) (*amqp.Queue, error) {
	queue, err := b.channel.QueueDeclare( p.Name, p.Durable, p.AutoDelete, p.Exclusive, p.NoWait, nil)
	if err != nil {
		return nil, err
	}
	b.queues[id] = &queue
	return &queue, nil
}

func (b *RabbitMqBroker) WithConsumerChannel(id string, handler ConsumerHandlerFunc, p *RabbitMqConsumerProperties) (<-chan amqp.Delivery, error) {
	queueName := p.QueueName
	if strings.HasPrefix(queueName, "nameFromQueue(") && strings.HasSuffix(queueName, ")") {
		queueId := strings.TrimLeft(strings.TrimRight(queueName, ")"), "nameFromQueue(")
		q, ok := b.queues[queueId]
		if !ok {
			return nil, fmt.Errorf("unable to find queue name for queue id %s", queueId)
		}
		queueName = q.Name
	}
	consumerChannel, err := b.channel.Consume(queueName, p.Name, p.AutoAck, p.Exclusive, p.NoLocal, p.NoWait, nil)
	if err != nil {
		return nil, err
	}
	b.consumers[id] = &ConsumerChannel{
		channel: consumerChannel,
		handler: handler,
	}
	return consumerChannel, nil
}

func (b *RabbitMqBroker) Close() {
	if b.channel != nil {
		b.channel.Close()
	}
	if b.connection != nil {
		b.connection.Close()
	}
}

func (b *RabbitMqBroker) Run() {
	var wg sync.WaitGroup
	wg.Add(len(b.consumers))
	for _, v := range b.consumers {
		go func() {
			defer wg.Done()
			for d := range v.channel {
				v.handler(b.channel, &d)
			}
		}()
	}
	wg.Wait()
}





