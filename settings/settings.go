package microservice

import (
	"fmt"

	"github.com/ivanmtzp/go-microservice/config"
)

type SettingsReader interface {
	Grpc() *GrpcSettings
	Log() *LogSettings
	Database() *DatabaseSettings
	MetricsPusher() *MetricsPusherSettings
}

type GrpcSettings struct {
	address string
	gatewayAddress string
}

type LogSettings struct {
	level string
}

type DatabaseSettings struct {
	dialect string
	database string
	host string
	port int
	user string
	password string
	pool int
}

type MetricsPusherSettings struct {
	enabled bool
	interval int
	host string
	port int
	database string
	user string
	password string
}

type ConfigSettings struct {
	config *config.Config
}

func NewConfigSettings(c *config.Config) *ConfigSettings{
	return &ConfigSettings{config: c}
}

func (c *ConfigSettings) Grpc() *GrpcSettings {
	host := c.config.GetString("grpc", "host")
	grpcPort := c.config.GetInt("grpc", "port")
	gatewayPort := c.config.GetInt("grpc", "gateway_port")
	address := fmt.Sprintf("%s:%d", host, grpcPort)
	gatewayAddress := fmt.Sprintf("%s:%d", host, gatewayPort)

	return &GrpcSettings{address: address, gatewayAddress: gatewayAddress}
}

func (c *ConfigSettings) Log() *LogSettings {
	level := c.config.GetString("log", "level")

	return &LogSettings{level: level}
}

func (c *ConfigSettings) Database() *DatabaseSettings {
	dialect := c.config.GetString("database", "dialect")
	database := c.config.GetString("database", "name")
	host := c.config.GetString("database", "host")
	port := c.config.GetInt("database", "port")
	user := c.config.GetString("database", "user")
	password := c.config.GetString("database", "password")
	pool := c.config.GetInt("database", "pool")

	return &DatabaseSettings{dialect: dialect, database: database, host: host, port: port, user: user, password: password, pool: pool}
}

func (c *ConfigSettings) MetricsPusher() *MetricsPusherSettings {
	enabled := c.config.GetBool("metrics", "pusher", "enabled")
	if !enabled {
		return &MetricsPusherSettings{}
	}
	interval := c.config.GetInt("metrics", "pusher", "interval")
	host := c.config.GetString("metrics",  "pusher", "host")
	port := c.config.GetInt("metrics", "pusher", "port")
	database := c.config.GetString("metrics",  "pusher", "database")
	user := c.config.GetString("metrics",  "pusher", "user")
	password := c.config.GetString("metrics", "password")

	return &MetricsPusherSettings{enabled: enabled, interval: interval, host: host, port: port, database: database, user: user, password: password}
}



