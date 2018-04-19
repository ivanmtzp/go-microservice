package settings

import (
	"fmt"

	"github.com/ivanmtzp/go-microservice/config"
)

type Reader interface {
	Grpc() *Grpc
	Log() *Log
	Database() *Database
	MetricsPusher() *MetricsPusher
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

type MetricsPusher struct {
	Enabled bool
	Interval int
	Host string
	Port int
	Database string
	User string
	Password string
}

type ConfigSettings struct {
	config *config.Config
}

func NewConfigSettings(c *config.Config) *ConfigSettings {
	return &ConfigSettings{config: c}
}

func (c *ConfigSettings) Grpc() *Grpc {
	host := c.config.GetString("grpc", "host")
	grpcPort := c.config.GetInt("grpc", "port")
	gatewayPort := c.config.GetInt("grpc", "gateway_port")
	address := fmt.Sprintf("%s:%d", host, grpcPort)
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

func (c *ConfigSettings) MetricsPusher() *MetricsPusher {
	enabled := c.config.GetBool("metrics", "pusher", "enabled")
	if !enabled {
		return &MetricsPusher{}
	}
	interval := c.config.GetInt("metrics", "pusher", "interval")
	host := c.config.GetString("metrics",  "pusher", "host")
	port := c.config.GetInt("metrics", "pusher", "port")
	database := c.config.GetString("metrics",  "pusher", "database")
	user := c.config.GetString("metrics",  "pusher", "user")
	password := c.config.GetString("metrics", "password")

	return &MetricsPusher{Enabled: enabled, Interval: interval, Host: host, Port: port, Database: database, User: user, Password: password}
}



