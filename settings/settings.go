package settings

import (
	"fmt"

	"github.com/ivanmtzp/go-microservice/config"
	"github.com/ivanmtzp/go-microservice/database"
	"github.com/ivanmtzp/go-microservice/monitoring"
	"github.com/ivanmtzp/go-microservice/broker"
)

type Reader interface {
	Log() *Log
	Database() *database.Properties
	GrpcServer() *GrpcServer
	GrpcClient() *GrpcClient
	Monitoring() *Monitoring
	RabbitMqBroker() *RabbitMqBroker
}


type Log struct {
	Level string
}

type GrpcServer struct {
	Address string
	GatewayAddress string
}

type GrpcClient struct {
	Endpoints map[string]string
}

type InfluxDbMetricsPusher struct {
	InfluxDbProperties *monitoring.InfluxDbProperties
	Interval int
}

type Monitoring struct {
	Address string
	InfluxDbMetricsPusher *InfluxDbMetricsPusher
}

type RabbitMqBroker struct {
	Address string
	PrefetchCount int
	PrefetchSize int
	Queues map[string]*broker.RabbitMqQueueProperties
	Consumers map[string]*broker.RabbitMqConsumerProperties
}

type ConfigSettings struct {
	config *config.Config
}

func NewConfigSettings(c *config.Config) *ConfigSettings {
	return &ConfigSettings{config: c}
}

func (c *ConfigSettings) Log() *Log{
	return &Log{
		Level: c.config.GetString("log", "level"),
	}
}

func (c *ConfigSettings) Database() *database.Properties {
	m := c.config.GetStringMap("database")
	return &database.Properties{
		Dialect: m.GetString("dialect"),
		Database: m.GetString("name"),
		Host: m.GetString("host"),
		Port: m.GetInt("port"),
		User: m.GetString("user"),
		Password: m.GetString("password"),
		Pool: m.GetInt("pool"),
	}
}

func (c *ConfigSettings) GrpcServer() *GrpcServer {
	host := c.config.GetString("host")
	return &GrpcServer{
		Address: fmt.Sprintf("%s:%d", host, c.config.GetInt("grpc", "server", "port")),
		GatewayAddress: fmt.Sprintf("%s:%d", host, c.config.GetInt("grpc", "server", "gateway_port")),
	}
}

func (c* ConfigSettings) GrpcClient() *GrpcClient {
	endpoints := make(map[string]string)
	for k, _ := range c.config.GetStringMap("grpc", "clients") {
		m := c.config.GetStringMap("grpc", "clients", k)
		endpoints[k] = fmt.Sprintf("%s:%d", m.GetString("host"), m.GetInt("port"))
	}
	return &GrpcClient{Endpoints: endpoints}
}

func (c *ConfigSettings) Monitoring() *Monitoring {
	var imp *InfluxDbMetricsPusher
	m := c.config.GetStringMap("monitoring", "metrics", "influxdb_pusher")
	if m != nil {
		imp = &InfluxDbMetricsPusher{
			Interval: m["interval"].(int),
			InfluxDbProperties: &monitoring.InfluxDbProperties{
				Host:     m.GetString("host"),
				Port:     m.GetInt("port"),
				Database: m.GetString("database"),
				User:     m.GetStringWithDefault("user", ""),
				Password: m.GetStringWithDefault("password", ""),
			},
		}
	}
	return &Monitoring{
		Address: fmt.Sprintf("%s:%d", c.config.GetString("host"), c.config.GetInt("monitoring", "port")),
		InfluxDbMetricsPusher: imp,
	}
}

func (c* ConfigSettings) RabbitMqBroker() *RabbitMqBroker {
	queues := make(map[string]*broker.RabbitMqQueueProperties)
	queuesConfig := c.config.GetStringMap("broker", "rabbitmq", "queues")
	for k, _ := range queuesConfig {
		m := c.config.GetStringMap("broker", "rabbitmq", "queues", k)
		queues[k] = &broker.RabbitMqQueueProperties {
			Name: m.GetString("name"),
			Durable: m.GetBool("durable"),
			AutoDelete: m.GetBool("auto_delete"),
			Exclusive: m.GetBool("exclusive"),
			NoWait: m.GetBool("no_wait"),
		}
	}
	consumers :=  make(map[string]*broker.RabbitMqConsumerProperties)
	consumersConfig := c.config.GetStringMap("broker", "rabbitmq", "consumers")
	for k, _ := range consumersConfig {
		m := c.config.GetStringMap("broker", "rabbitmq", "consumers", k)
		consumers[k] = &broker.RabbitMqConsumerProperties {
			Name: m.GetStringWithDefault("name", ""),
			QueueName: m.GetString("queue_name"),
			AutoAck: m.GetBool("auto_ack"),
			Exclusive: m.GetBool("exclusive"),
			NoLocal: m.GetBool("no_local"),
			NoWait: m.GetBool("no_wait"),
		}
	}
	return &RabbitMqBroker{
		Address: fmt.Sprintf("amqp://%s:%d/",
			c.config.GetString("broker", "rabbitmq", "host"),
			c.config.GetInt("broker", "rabbitmq", "port")	),
		PrefetchCount: c.config.GetInt("broker", "rabbitmq", "qos", "prefetch_count"),
		PrefetchSize: c.config.GetInt("broker", "rabbitmq", "qos", "prefetch_size"),
		Queues: queues,
		Consumers: consumers,
	}
}







