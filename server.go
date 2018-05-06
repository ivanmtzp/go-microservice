package microservice

import (
	"fmt"
	"os"
	"strings"
	"github.com/ivanmtzp/go-microservice/broker"
	"github.com/ivanmtzp/go-microservice/grpc"
	"github.com/ivanmtzp/go-microservice/log"
	"github.com/ivanmtzp/go-microservice/config"
	"github.com/ivanmtzp/go-microservice/database"
	"github.com/ivanmtzp/go-microservice/settings"
	"github.com/ivanmtzp/go-microservice/monitoring"
)

type GrpcClientsMap map[string]*grpc.Client

type MicroService struct {
	name string
	settings settings.Reader
	statusServer *monitoring.StatusServer
	grpcServer *grpc.Server
	grpcClients GrpcClientsMap
	httpGatewayServer *grpc.HttpGatewayServer
	database *database.Database
	rabbitMqBroker *broker.RabbitMqBroker
}


func New(name string, sr settings.Reader) *MicroService {
	return &MicroService{name: name, settings: sr, statusServer: monitoring.NewStatusServer(), grpcClients: make(map[string]*grpc.Client)}
}

func NewWithSettingsFile(name, envPrefix, filename string) (*MicroService, error) {
	conf := config.New()
	if err:= conf.Read( envPrefix, filename); err != nil {
		return nil, fmt.Errorf("error reading configuration file %s, %s", filename, err)
	}

	configSettings := settings.NewConfigSettings(conf)
	logLevel := configSettings.Log().Level
	if logLevel != "" {
		if err := log.SetLevel(logLevel); err != nil {
			return nil, fmt.Errorf("configuration error, invalid log level: %s, ", err)
		}
	}
	log.Infof("log level set to %s", log.Level())
	log.Debug("environment variables: ")
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, envPrefix) {
			log.Debug(e)
		}
	}
	return New(name, configSettings), nil
}


func (ms *MicroService) WithGrpcAndGatewayServer(sr grpc.ServerServiceRegistrationFunc, gsr grpc.GatewayServerServiceRegistrationFunc, gatewayhealthCheckEndpoint string) (*grpc.Server, *grpc.HttpGatewayServer, error) {
	grpcSettings := ms.settings.GrpcServer()
	grpcServer := grpc.NewServer(grpcSettings.Address, sr)
	gatewayServer, err := grpc.NewHttpGatewayServer(grpcSettings.GatewayAddress, grpcSettings.Address, gsr, gatewayhealthCheckEndpoint)
	if err != nil {
		return nil, nil, err
	}
	ms.grpcServer = grpcServer
	ms.httpGatewayServer = gatewayServer
	ms.statusServer.RegisterHealthCheck("grpc_gateway", ms.httpGatewayServer)
	return grpcServer, gatewayServer, nil
}

func (ms *MicroService) WithGrpcClient(name string, serviceCreator grpc.CreateClientServiceFunc) (*grpc.Client, error) {
	settings := ms.settings.GrpcClient()
	address, ok := settings.Endpoints[name]
	if !ok {
		return nil, fmt.Errorf("grpc client address not found in settings: %s", name)
	}
	client, err := grpc.NewClient(address, serviceCreator)
	if err != nil {
		return nil, err
	}
	ms.grpcClients[name] = client
	return client, nil
}


func (ms *MicroService) WithGrpcClients(clients map[string]grpc.CreateClientServiceFunc) (GrpcClientsMap, error) {
	endpoints := ms.settings.GrpcClient().Endpoints
	for name, sc := range clients {
		address, ok := endpoints[name]
		if !ok {
			for _, client := range ms.grpcClients {
				client.Close()
			}
			return nil, fmt.Errorf("grpc client address not found in settings: %s", name)
		}
		client, err := grpc.NewClient(address, sc)
		if err != nil {
			for _, client := range ms.grpcClients {
				client.Close()
			}
			return nil, err
		}
		ms.grpcClients[name] = client
	}
	return ms.grpcClients, nil
}

func (m GrpcClientsMap) Close(){
	for _, client := range m {
		client.Close()
	}
}

func (ms *MicroService) WithDatabase(healthCheckQuery string) (*database.Database, error) {
	dbs := ms.settings.Database()

	db, err := database.New(dbs, healthCheckQuery)
	if err != nil {
		return nil, err
	}
	ms.database = db
	ms.statusServer.RegisterHealthCheck("database", ms.database)
	return ms.database, nil
}

func (ms *MicroService) WithMonitoring() {
	monSettings := ms.settings.Monitoring()
	ms.statusServer.Enable(monSettings.Address)
}

func (ms *MicroService) WithRabbitMqBroker() (*broker.RabbitMqBroker, error) {
	settings := ms.settings.RabbitMqBroker()
	rabbitmq, err := broker.NewRabbitMqBroker(settings.Address, settings.PrefetchCount, settings.PrefetchSize)
	if err != nil {
		return nil, err
	}
	return rabbitmq, nil
}

func (ms *MicroService) Run() {

	if ms.grpcServer != nil {
		go func() {
			log.Infof("starting HTTP/2 gRPC server on %s", ms.grpcServer.Address())
			err:= ms.grpcServer.Run()
			log.FailOnError(err, fmt.Sprint("failed to start gRPC server " , ms.grpcServer.Address()))
		}()
		if ms.httpGatewayServer != nil {
			go func() {
				log.Infof("starting HTTP/1.1 gateway server on %s for grpc server endpoint %s", ms.httpGatewayServer.Address(), ms.httpGatewayServer.GrpcEndpointAddress() )
				err := ms.httpGatewayServer.Run()
				log.FailOnError(err, fmt.Sprint("failed to start Http gateway server ", ms.httpGatewayServer.Address()))
			}()
		}
	}

	if ms.statusServer.Enabled() {
		go func() {
			log.Infof("starting HTTP/1.1 monitoring server on %s", ms.settings.Monitoring().Address)
			ms.statusServer.Run()
		}()
		go func() {
			mps := ms.settings.Monitoring().InfluxDbMetricsPusher
			if mps != nil {
				log.Infof("starting InfluxDb Metrics pushing to host: %s, port: %d, database: %s, user: %s,  with interval: %d",
					mps.InfluxDbProperties.Host,
					mps.InfluxDbProperties.Port,
					mps.InfluxDbProperties.Database,
					mps.InfluxDbProperties.User,
					mps.Interval)
				monitoring.RunInfluxDbMetricsPusher(mps.InfluxDbProperties, mps.Interval)
			}
		}()
	} else {
		log.Warning("monitoring server is disabled")
	}

	select {}

}




