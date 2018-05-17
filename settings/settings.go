package settings

import (
	"fmt"

	"github.com/ivanmtzp/go-microservice/config"
	"github.com/ivanmtzp/go-microservice/database"
	"github.com/ivanmtzp/go-microservice/broker"
	"time"
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

type InfluxDbProperties struct {
	Address string
	Database string
	User string
	Password string
}

type InfluxDbMetricsPusher struct {
	InfluxDbProperties *InfluxDbProperties
	Interval time.Duration
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
	return &database.Properties{
		Dialect: c.config.GetString("database", "dialect"),
		Database: c.config.GetString("database", "name"),
		Host: c.config.GetString("database", "host"),
		Port: c.config.GetInt("database", "port"),
		User: c.config.GetString("database", "user"),
		Password: c.config.GetString("database", "password"),
		Pool: c.config.GetInt("database", "pool"),
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
		endpoints[k] = fmt.Sprintf("%s:%d", c.config.GetString("grpc", "clients", k, "host"),
			c.config.GetInt("grpc", "clients", k, "port"))
	}
	return &GrpcClient{Endpoints: endpoints}
}

func (c *ConfigSettings) Monitoring() *Monitoring {
	var imp *InfluxDbMetricsPusher
	_, ok := c.config.HasKey("monitoring", "metrics", "influxdb_pusher")
	if ok {
		imp = &InfluxDbMetricsPusher{
			Interval: time.Second * time.Duration(c.config.GetInt("monitoring", "metrics", "influxdb_pusher", "interval")),
			InfluxDbProperties: &InfluxDbProperties{
				Address: fmt.Sprintf("http://%s:%d",  c.config.GetString("monitoring", "metrics", "influxdb_pusher", "host"),
					c.config.GetInt("monitoring", "metrics", "influxdb_pusher", "port")),
				Database: c.config.GetString("monitoring", "metrics", "influxdb_pusher", "database"),
				User:     c.config.GetString("monitoring", "metrics", "influxdb_pusher", "user"),
				Password: c.config.GetString("monitoring", "metrics", "influxdb_pusher", "password"),
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
		queues[k] = &broker.RabbitMqQueueProperties {
			Name: c.config.GetString("broker", "rabbitmq", "queues", k, "name"),
			Durable: c.config.GetBool("broker", "rabbitmq", "queues", k, "durable"),
			AutoDelete: c.config.GetBool("broker", "rabbitmq", "queues", k, "auto_delete"),
			Exclusive: c.config.GetBool("broker", "rabbitmq", "queues", k, "exclusive"),
			NoWait: c.config.GetBool("broker", "rabbitmq", "queues", k, "no_wait"),
		}
	}
	consumers :=  make(map[string]*broker.RabbitMqConsumerProperties)
	consumersConfig := c.config.GetStringMap("broker", "rabbitmq", "consumers")
	for k, _ := range consumersConfig {
		consumers[k] = &broker.RabbitMqConsumerProperties {
			Name: c.config.GetString("broker", "rabbitmq", "consumers", k, "name"),
			QueueName: c.config.GetString("broker", "rabbitmq", "consumers", k, "queue_name"),
			AutoAck: c.config.GetBool("broker", "rabbitmq", "consumers", k, "auto_ack"),
			Exclusive: c.config.GetBool("broker", "rabbitmq", "consumers", k, "exclusive"),
			NoLocal: c.config.GetBool("broker", "rabbitmq", "consumers", k, "no_local"),
			NoWait: c.config.GetBool("broker", "rabbitmq", "consumers", k, "no_wait"),
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







