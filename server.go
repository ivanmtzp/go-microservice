package microservice

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"time"
	"github.com/gobuffalo/pop"
	"github.com/ivanmtzp/go-microservice/grpc"
	"github.com/ivanmtzp/go-microservice/log"
	"github.com/ivanmtzp/go-microservice/config"
	"github.com/ivanmtzp/go-microservice/database"
	"github.com/ivanmtzp/go-microservice/settings"
	"github.com/ivanmtzp/go-microservice/monitoring"
)


type MicroService struct {
	name string
	settings settings.Reader
	statusServer *monitoring.StatusServer
	grpcServer *grpc.Server
	grpcClient map[string]*grpc.Client
	httpGatewayServer *grpc.HttpGatewayServer
	database *database.Database
}


func New(name string, sr settings.Reader) *MicroService {
	return &MicroService{name: name, settings: sr, statusServer: monitoring.NewStatusServer(), grpcClient: make(map[string]*grpc.Client)}
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


func (ms *MicroService) WithGrpcAndGateway(sr grpc.ServerServiceRegistrationFunc, gsr grpc.GatewayServerServiceRegistrationFunc, gatewayhealthCheckEndpoint string) (*grpc.Server, *grpc.HttpGatewayServer, error) {
	grpcSettings := ms.settings.Grpc()
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
	ms.grpcClient[name] = client
	return client, nil
}

func (ms *MicroService) WithDatabase(healthCheckQuery string) (*database.Database, error) {
	dbs := ms.settings.Database()
	connectionDetails := &pop.ConnectionDetails{ Dialect: dbs.Dialect, Database: dbs.Database,
		Host: dbs.Host, Port: strconv.Itoa(dbs.Port), User: dbs.User, Password: dbs.Password,
		Pool: dbs.Pool, IdlePool: 0}

	db, err := database.New(connectionDetails, healthCheckQuery)
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

func (ms *MicroService) Run() {

	if ms.database != nil {
		// connect to database
		log.Info("connecting to database ", ms.database.Connection().URL())
		if err := ms.database.Open() ; err != nil {
			log.FailOnError(err, "couldn't open connection to database ")
		}
		defer ms.database.Close()
	}

	if ms.grpcServer != nil {
		// fire the gRPC server in a goroutine
		go func() {
			log.Infof("starting HTTP/2 gRPC server on %s", ms.grpcServer.Address())
			err:= ms.grpcServer.Run()
			log.FailOnError(err, fmt.Sprint("failed to start gRPC server " , ms.grpcServer.Address()))
		}()
		if ms.httpGatewayServer != nil {
			// fire the http grpc gateway in a goroutine
			go func() {
				log.Infof("starting HTTP/1.1 gateway server on %s for grpc server endpoint %s", ms.httpGatewayServer.Address(), ms.httpGatewayServer.GrpcEndpointAddress() )
				err := ms.httpGatewayServer.Run()
				log.FailOnError(err, fmt.Sprint("failed to start Http gateway server ", ms.httpGatewayServer.Address()))
			}()
		}
	}

	for _, client := range ms.grpcClient {
		defer client.Close()
	}


	if ms.statusServer.Enabled() {
		// fire the status server
		go func() {
			log.Infof("starting HTTP/1.1 monitoring server on %s", ms.settings.Monitoring().Address)
			ms.statusServer.Run()
		}()
		//	fire the metrics pusher
		go func() {
			mps := &ms.settings.Monitoring().InfluxMetricsPusher
			if mps.Enabled {
				hostUrl := fmt.Sprintf("http://%s:%d", mps.Host, mps.Port)
				log.Infof("starting InfluxDb Metrics pushing to: %s, database: %s, user: %s,  with interval: %d", hostUrl,
					mps.Database, mps.User, mps.Interval)
				monitoring.StartInfluxDbPusher(time.Second*time.Duration(mps.Interval), hostUrl, mps.Database, mps.User, mps.Password)
			}
		}()
	} else {
		log.Warning("monitoring server is disabled")
	}

	// infinite loop
	select {}

}




