package settings

import (
	"fmt"

	"github.com/ivanmtzp/go-microservice/config"
)

type Reader interface {
	Grpc() *Grpc
	Log() *Log
	Database() *Database
	Monitoring() *Monitoring
}

type Grpc struct {
	Address string
	GatewayAddress string
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
	port := c.config.GetInt("grpc", "port")
	gatewayPort := c.config.GetInt("grpc", "gateway_port")
	address := fmt.Sprintf("%s:%d", host, port)
	gatewayAddress := fmt.Sprintf("%s:%d", host, gatewayPort)

	return &Grpc{Address: address, GatewayAddress: gatewayAddress}
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


	influxMetricsPusher := InfluxMetricsPusher{}
	enabled := c.config.GetBool("monitoring", "influxdb_pusher", "enabled")
	if enabled {
		influxMetricsPusher.Interval = c.config.GetInt("monitoring", "metrics", "pusher", "interval")
		influxMetricsPusher.Host = c.config.GetString("monitoring", "metrics", "pusher", "host")
		influxMetricsPusher.Port = c.config.GetInt("monitoring", "metrics", "pusher", "port")
		influxMetricsPusher.Database = c.config.GetString("monitoring", "metrics", "pusher", "database")
		influxMetricsPusher.User = c.config.GetString("monitoring", "metrics", "pusher", "user")
		influxMetricsPusher.Password = c.config.GetString("monitoring", "metrics", "password")
	}

	return &Monitoring{Address: address, InfluxMetricsPusher: influxMetricsPusher}

}



