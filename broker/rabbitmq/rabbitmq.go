package broker

import (
	"github.com/streadway/amqp"
)

type RabbitMqBroker struct
{
	address string
	connection *amqp.Connection
	channel *amqp.Channel
	queues map[string]*amqp.Queue
	consumerChannels map[string]<-chan amqp.Delivery
}

func NewRabbitMqBroker(address string, prefetchCount, prefetchSize int) (*Broker, error) {
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

func (b *RabbitMqBroker) WithQueue(name string, durable, autoDelete, exclusive, noWait bool) (*amqp.Queue, error) {
	queue, err := b.channel.QueueDeclare( name, durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		return nil, err
	}
	b.queues[name] = &queue
	return &queue, nil
}

func (b *RabbitMqBroker) WithConsumerChannel(queueName, consumer string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
	consumerChannel, err := b.channel.Consume(queueName, consumer, autoAck, exclusive, noLocal, noWait, nil)
	if err != nil {
		return nil, err
	}
	b.consumerChannels[queueName] = consumerChannel
	return consumerChannel, nil
}







