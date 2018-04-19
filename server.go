package microservice

import (
	"fmt"
	"time"
	"os"
	"strings"
	"strconv"

	"github.com/gobuffalo/pop"
	"github.com/ivanmtzp/go-microservice/grpc"
	"github.com/ivanmtzp/go-microservice/log"
	"github.com/ivanmtzp/go-microservice/metrics"
	"github.com/ivanmtzp/go-microservice/config"
	"github.com/ivanmtzp/go-microservice/database"
)


type MicroService struct {
	name string
	settings SettingsReader
	grpcServer *grpc.Server
	httpGatewayServer *grpc.HttpGatewayServer
	database *database.Database
	optionalMetricsPusher bool
}

func (ms *MicroService) Database() *database.Database {
	return ms.database
}

func New(name string, sr SettingsReader) *MicroService {
	return &MicroService{name: name, settings: sr}
}

func NewFromSettingsFile(name, envPrefix string) (*MicroService, error) {
	conf := config.New()
	if err:= conf.Read( envPrefix, "./config", "microservice", config.Yaml); err != nil {
		return nil, fmt.Errorf("error reading configuration file: ./config/microservice.yaml, %s", err)
	}
	configSettings := NewConfigSettings(conf)
	logLevel := configSettings.Log().level
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


func (ms *MicroService) WithGrpcAndGateway(sr grpc.ServiceRegistrator, gsr grpc.HttpGatewayServiceRegistrator) *MicroService {
	grpcSettings := ms.settings.Grpc()
	ms.grpcServer = grpc.New(grpcSettings.address, sr)
	ms.httpGatewayServer = grpc.NewHttpGateway(grpcSettings.gatewayAddress, grpcSettings.address, gsr)
	return ms
}

func (ms *MicroService) WithDatabase(healthCheck func (connection *pop.Connection)error) (*MicroService, error) {
	dbs := ms.settings.Database()
	connectionDetails := &pop.ConnectionDetails{ Dialect: dbs.dialect, Database: dbs.database,
		Host: dbs.host, Port: strconv.Itoa(dbs.port), User: dbs.user, Password: dbs.password,
		Pool: dbs.pool, IdlePool: 0}

	db, err := database.New(connectionDetails, healthCheck)
	if err != nil {
		return nil, err
	}
	ms.database = db
	return ms, nil
}

func (ms *MicroService) WithOptionalMetricsPusher() *MicroService {
	ms.optionalMetricsPusher = true
	return ms
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

	if ms.optionalMetricsPusher {
		//	fire the metrics pusher
		go func() {
			mps := ms.settings.MetricsPusher()
			hostUrl := fmt.Sprintf("http://%s:%d", mps.host, mps.port)
			log.Infof("starting InfluxDb Metrics pushing to: %s, database: %s, user: %s,  with interval: %d", hostUrl,
				mps.database, mps.user, mps.interval)
			metrics.StartInfluxDbPusher(time.Second * time.Duration(mps.interval), hostUrl, mps.database, mps.user, mps.password)
		}()
	}

	// infinite loop
	select {}

}




