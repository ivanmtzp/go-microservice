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
	for k, v := range c.config.GetStringMap("grpc", "clients") {
		vMap := v.(map[string]interface{})
		endpoints[k] = fmt.Sprintf("%s:%d", vMap["host"].(string), vMap["port"].(int))
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
				Host:     m["host"].(string),
				Port:     m["port"].(int),
				Database: m["database"].(string),
				User:     m["user"].(string),
				Password: m["password"].(string),
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
	for k, v := range queuesConfig {
		vMap := v.(map[string]interface{})
		queues[k] = &broker.RabbitMqQueueProperties {
			Name: vMap["name"].(string),
			Durable: vMap["durable"].(bool),
			AutoDelete: vMap["auto_delete"].(bool),
			Exclusive: vMap["exclusive"].(bool),
			NoWait:  vMap["no_wait"].(bool),
		}
	}
	consumers :=  make(map[string]*broker.RabbitMqConsumerProperties)
	consumersConfig := c.config.GetStringMap("broker", "rabbitmq", "consumers")
	for k, v := range consumersConfig {
		vMap := v.(map[string]interface{})
		name := ""
		if vMap["name"] != nil {
			name = vMap["name"].(string)
		}
		consumers[k] = &broker.RabbitMqConsumerProperties {
			Name: name,
			QueueName: vMap["queue_name"].(string),
			AutoAck: vMap["auto_ack"].(bool),
			Exclusive: vMap["exclusive"].(bool),
			NoLocal:  vMap["no_local"].(bool),
			NoWait:  vMap["no_wait"].(bool),
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







