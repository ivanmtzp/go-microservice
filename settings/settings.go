package settings

import (
	"fmt"

	"github.com/ivanmtzp/go-microservice/config"
)

type Reader interface {
	GrpcServer() *GrpcServer
	GrpcClient() *GrpcClient
	Log() *Log
	Database() *Database
	Monitoring() *Monitoring
	RabbitMqBroker() *RabbitMqBroker
}

type GrpcServer struct {
	Address string
	GatewayAddress string
}

type GrpcClient struct {
	Endpoints map[string]string
}

type Log struct {
	Level string
}

type Database struct {
	Dialect string
	Database string
	Host string
	Port int
	User string
	Password string
	Pool int
}

type InfluxMetricsPusher struct {
	Enabled bool
	Interval int
	Host string
	Port int
	Database string
	User string
	Password string
}

type Monitoring struct {
	Address string
	InfluxMetricsPusher InfluxMetricsPusher
}

type ConfigSettings struct {
	config *config.Config
}

type RabbitMqQueue struct {
	Name string
	Durable bool
	AutoDelete bool
	Exclusive bool
	NoWait bool
}

type RabbitMqConsumer struct {
	Name string
	QueueName string
	AutoAck bool
	Exclusive bool
	NoLocal bool
	NoWait bool
}

type RabbitMqBroker struct {
	Address string
	PrefetchCount int
	PrefetchSize int
	Queues map[string]*RabbitMqQueue
	Consumers map[string]*RabbitMqConsumer
}


func NewConfigSettings(c *config.Config) *ConfigSettings {
	return &ConfigSettings{config: c}
}

func (c* ConfigSettings) RabbitMqBroker() *RabbitMqBroker {
	host := c.config.GetString("broker", "rabbitmq", "host")
	port := c.config.GetInt("broker", "rabbitmq", "port")
	address := fmt.Sprintf("amqp://%s:%d/", host, port)
	prefetchCount := c.config.GetInt("broker", "rabbitmq", "qos", "prefetch_count")
	prefetchSize := c.config.GetInt("broker", "rabbitmq", "qos", "prefetch_size")

	queues := make(map[string]*RabbitMqQueue)
	queuesConfig := c.config.GetStringMap("broker", "rabbitmq", "queues")
	for k, v := range queuesConfig {
		vMap := v.(map[string]interface{})
		queues[k] = &RabbitMqQueue {
			Name: vMap["name"].(string),
			Durable: vMap["durable"].(bool),
			AutoDelete: vMap["auto_delete"].(bool),
			Exclusive: vMap["exclusive"].(bool),
			NoWait:  vMap["no_wait"].(bool),
		}
	}

	consumers :=  make(map[string]*RabbitMqConsumer)
	consumersConfig := c.config.GetStringMap("broker", "rabbitmq", "consumers")
	for k, v := range consumersConfig {
		vMap := v.(map[string]interface{})
		name := ""
		if vMap["name"] != nil {
			name = vMap["name"].(string)
		}
		consumers[k] = &RabbitMqConsumer {Name: name,
			QueueName: vMap["queue_name"].(string),
			AutoAck: vMap["auto_ack"].(bool),
			Exclusive: vMap["exclusive"].(bool),
			NoLocal:  vMap["no_local"].(bool),
			NoWait:  vMap["no_wait"].(bool),
		}
	}
	return &RabbitMqBroker{
		Address: address,
		PrefetchCount: prefetchCount,
		PrefetchSize: prefetchSize,
		Queues: queues,
		Consumers: consumers,
	}
}


func (c *ConfigSettings) GrpcServer() *GrpcServer {
	host := c.config.GetString("host")
	grpcPort := c.config.GetInt("grpc", "server", "port")
	gatewayPort := c.config.GetInt("grpc", "server", "gateway_port")
	address := fmt.Sprintf("%s:%d", host, grpcPort)
	gatewayAddress := fmt.Sprintf("%s:%d", host, gatewayPort)

	return &GrpcServer{Address: address, GatewayAddress: gatewayAddress}
}

func (c* ConfigSettings) GrpcClient() *GrpcClient {
	endpoints := make(map[string]string)
	clients := c.config.GetStringMap("grpc", "clients")
	for k, v := range clients {
		vMap := v.(map[string]interface{})
		host := vMap["host"].(string)
		port := vMap["port"].(int)
		endpoints[k] = fmt.Sprintf("%s:%d", host, port)
	}
	return &GrpcClient{Endpoints: endpoints}
}

func (c *ConfigSettings) Log() *Log{
	level := c.config.GetString("log", "level")

	return &Log{Level: level}
}

func (c *ConfigSettings) Database() *Database {
	dialect := c.config.GetString("database", "dialect")
	database := c.config.GetString("database", "name")
	host := c.config.GetString("database", "host")
	port := c.config.GetInt("database", "port")
	user := c.config.GetString("database", "user")
	password := c.config.GetString("database", "password")
	pool := c.config.GetInt("database", "pool")

	return &Database{Dialect: dialect, Database: database, Host: host, Port: port, User: user, Password: password, Pool: pool}
}

func (c *ConfigSettings) Monitoring() *Monitoring {
	host := c.config.GetString("host")
	port := c.config.GetInt("monitoring", "port")
	address := fmt.Sprintf("%s:%d", host, port)

	imp := InfluxMetricsPusher{}
	enabled := c.config.GetBool("monitoring", "metrics", "influxdb_pusher", "enabled")
	if enabled {
		imp.Interval = c.config.GetInt("monitoring", "metrics", "influxdb_pusher", "interval")
		imp.Host = c.config.GetString("monitoring", "metrics",  "influxdb_pusher", "host")
		imp.Port = c.config.GetInt("monitoring", "metrics", "influxdb_pusher", "port")
		imp.Database = c.config.GetString("monitoring", "metrics",  "influxdb_pusher", "database")
		imp.User = c.config.GetString("monitoring", "metrics", "influxdb_pusher", "user")
		imp.Password = c.config.GetString( "monitoring", "metrics", "influxdb_pusher", "password")	
	}
	

	return &Monitoring{Address: address, InfluxMetricsPusher: imp}
}



