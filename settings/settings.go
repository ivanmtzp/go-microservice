package settings

import (
	"fmt"

	"github.com/ivanmtzp/go-microservice/config"
)

type Reader interface {
	Grpc() *Grpc
	GrpcClient() *GrpcClient
	Log() *Log
	Database() *Database
	Monitoring() *Monitoring
}

type Grpc struct {
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

func NewConfigSettings(c *config.Config) *ConfigSettings {
	return &ConfigSettings{config: c}
}

func (c *ConfigSettings) Grpc() *Grpc {
	host := c.config.GetString("host")
	grpcPort := c.config.GetInt("grpc", "server", "port")
	gatewayPort := c.config.GetInt("grpc", "server", "gateway_port")
	address := fmt.Sprintf("%s:%d", host, grpcPort)
	gatewayAddress := fmt.Sprintf("%s:%d", host, gatewayPort)

	return &Grpc{Address: address, GatewayAddress: gatewayAddress}
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



