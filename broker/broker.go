package broker

import (
	"github.com/streadway/amqp"
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

type RabbitMqBroker struct
{
	address string
	connection *amqp.Connection
	channel *amqp.Channel
	queues map[string]*amqp.Queue
	consumerChannels map[string]<-chan amqp.Delivery
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
	return &RabbitMqBroker{address: address, connection: connection, channel:channel, queues: make(map[string]*amqp.Queue)}, nil
}

func (b *RabbitMqBroker) WithQueue(p *RabbitMqQueueProperties) (*amqp.Queue, error) {
	queue, err := b.channel.QueueDeclare( p.Name, p.Durable, p.AutoDelete, p.Exclusive, p.NoWait, nil)
	if err != nil {
		return nil, err
	}
	b.queues[p.Name] = &queue
	return &queue, nil
}

func (b *RabbitMqBroker) WithConsumerChannel(p *RabbitMqConsumerProperties) (<-chan amqp.Delivery, error) {
	consumerChannel, err := b.channel.Consume(p.QueueName, p.Name, p.AutoAck, p.Exclusive, p.NoLocal, p.NoWait, nil)
	if err != nil {
		return nil, err
	}
	b.consumerChannels[p.QueueName] = consumerChannel
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






