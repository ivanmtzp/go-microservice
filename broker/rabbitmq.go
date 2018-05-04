package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/ivanmtzp/go-microservice/log"
	"fmt"
)

struct Broker
{

}


func NewRabbitmqBroker() *Broker {
	return &Broker{}
}



func initRabbitmqConnection()(*amqp.Connection, *amqp.Channel, *amqp.Queue, <-chan amqp.Delivery, error)  {

	// rabbitUrl := viper.Get( "rabbitmq_url").(string)
	rabbitUrl := "amqp://localhost:5672/"
	log.Infof("initializing rabbitmq on %s", rabbitUrl)

	conn, err := amqp.Dial(rabbitUrl)
	failOnError(err, fmt.Sprintf("failed to connect to RabbitMQ"))

	ch, err := conn.Channel()
	failOnError(err, "failed to open a rabbitmq channel")

	queue, err := ch.QueueDeclare(
		"vol_calibrator_queue", // name
		true,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to declare rabbitmq queue named vol_calibrator_queue, %s", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to set rabbitmq channel QoS, %s", err)
	}

	msgsChannel, err := ch.Consume(
		queue.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to register a rabbitmq consumer, %s", err)
	}

	return conn, ch, &queue, msgsChannel, nil
}







