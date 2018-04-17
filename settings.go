package microservice

import (
	"fmt"

	"github.com/ivanmtzp/go-microservice/config"
)

type SettingsReader interface {
	GrpcSettings() *GrpcSettings
	LogSettings() *LogSettings
}

type GrpcSettings struct {
	address string
	gatewayAddress string
}

type LogSettings struct {
	level string
}

type ConfigFileSettings struct {

}

func (c *ConfigFileSettings) GrpcSettings() *GrpcSettings {
	host := config.GetString("grpc", "host")
	grpcPort := config.GetInt("grpc", "port")
	gatewayPort := config.GetInt("grpc", "gateway_port")
	address := fmt.Sprintf("%s:%d", host, grpcPort)
	gatewayAddress := fmt.Sprintf("%s:%d", host, gatewayPort)

	return &GrpcSettings{address: address, gatewayAddress: gatewayAddress}
}

func (c *ConfigFileSettings) LogSettings() *LogSettings {
	level := config.GetString("log", "level")

	return &LogSettings{level: level}
}


