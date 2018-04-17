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

type ConfigFileSettings struct {

}

func (c *ConfigFileSettings) Grpc() *GrpcSettings {
	host := config.GetString("grpc", "host")
	grpcPort := config.GetInt("grpc", "port")
	gatewayPort := config.GetInt("grpc", "gateway_port")
	address := fmt.Sprintf("%s:%d", host, grpcPort)
	gatewayAddress := fmt.Sprintf("%s:%d", host, gatewayPort)

	return &GrpcSettings{address: address, gatewayAddress: gatewayAddress}
}

func (c *ConfigFileSettings) Log() *LogSettings {
	level := config.GetString("log", "level")

	return &LogSettings{level: level}
}

func (c *ConfigFileSettings) Database() *DatabaseSettings {
	dialect := config.GetString("database", "dialect")
	database := config.GetString("database", "name")
	host := config.GetString("database", "host")
	port := config.GetInt("database", "port")
	user := config.GetString("database", "user")
	password := config.GetString("database", "password")
	pool := config.GetInt("database", "pool")

	return &DatabaseSettings{dialect: dialect, database: database, host: host, port: port, user: user, password: password, pool: pool}
}

func (c *ConfigFileSettings) MetricsPusher() *MetricsPusherSettings {
	enabled := config.GetBool("metrics", "pusher", "enabled")
	if !enabled {
		return &MetricsPusherSettings{}
	}
	interval := config.GetInt("metrics", "pusher", "interval")
	host := config.GetString("metrics",  "pusher", "host")
	port := config.GetInt("metrics", "pusher", "port")
	database := config.GetString("metrics",  "pusher", "database")
	user := config.GetString("metrics",  "pusher", "user")
	password := config.GetString("metrics", "password")

	return &MetricsPusherSettings{enabled: enabled, interval: interval, host: host, port: port, database: database, user: user, password: password}
}



