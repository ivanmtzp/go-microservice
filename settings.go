package microservice

import (
	"github.com/ivanmtzp/go-microservice/config"
	"fmt"
)


type GrpcSettings struct {
	address string
	gatewayAddress string
}


func ReadGrpcSettings() *GrpcSettings {
	host := config.GetString("grpc", "host")
	grpcPort := config.GetInt("grpc", "port")
	gatewayPort := config.GetInt("grpc", "gateway_port")
	address := fmt.Sprintf("%s:%d", host, grpcPort)
	gatewayAddress := fmt.Sprintf("%s:%d", host, gatewayPort)
	return &GrpcSettings{address, gatewayAddress}
}


