package microservice

import (
	"fmt"
	"time"
	"strconv"

	"github.com/gobuffalo/pop"
	"github.com/ivanmtzp/go-microservice/grpc"
	"github.com/ivanmtzp/go-microservice/log"
	"github.com/ivanmtzp/go-microservice/metrics"

)


type MicroService struct {
	name string
	settings SettingsReader
	grpcServer *grpc.Server
	database *pop.Connection
	optionalMetricsPusher bool
}


func New(name string, sr SettingsReader) *MicroService {
	return &MicroService{name: name, settings: sr}
}

func (ms *MicroService) WithGrpcAndGateway(sr grpc.ServiceRegistrator,gsr grpc.HttpGatewayServiceRegistrator) *MicroService {
	ms.grpcServer = grpc.New(ms.settings.Grpc().address, sr).
		WithGateway(ms.settings.Grpc().gatewayAddress, gsr)
	return ms
}

func (ms *MicroService) WithDatabase() (*MicroService, error) {
	dbSettings := ms.settings.Database()
	connectionDetails := &pop.ConnectionDetails{	Dialect: dbSettings.dialect, Database: dbSettings.database,
		Host: dbSettings.host, Port: strconv.Itoa(dbSettings.port), User: dbSettings.user, Password: dbSettings.password,
		Pool: dbSettings.pool, IdlePool: 0}
	connection, err := pop.NewConnection(connectionDetails)
	if err != nil {
		return ms, fmt.Errorf("failed to create the database connection %s", err)
	}
	ms.database = connection
	return ms, nil
}

func (ms *MicroService) WithOptionalMetricsPusher() *MicroService {
	ms.optionalMetricsPusher = true
	return ms
}


func (ms *MicroService) Run() {

	if ms.database != nil {
		log.Info("connecting to database ", ms.database.URL())
		err := ms.database.Open()
		if err != nil {
			log.FailOnError(err, "couldn't open connection to database ")
		}
		defer ms.database.Close()
	}

	// err = handlers.PingDb(db)
	// if err != nil {
	//	log.FailOnError(err, fmt.Sprintf("ping check to database failed, database url: %s", db.URL()))
	// }

	if ms.grpcServer != nil {
		// fire the gRPC server in a goroutine
		go func() {
			log.Infof("starting HTTP/2 gRPC server on %s", ms.settings.Grpc().address)
			err:= ms.grpcServer.Run()
			log.FailOnError(err, fmt.Sprint("failed to start gRPC server " , ms.settings.Grpc().address))
		}()
		// fire the http grpc gateway in a goroutine
		go func() {
			log.Infof("starting HTTP/1.1 REST server on %s", ms.settings.Grpc().gatewayAddress)
			err := ms.grpcServer.RunHttpGateway()
			log.FailOnError(err, fmt.Sprint("failed to start Http gateway server ", ms.settings.Grpc().gatewayAddress))
		}()
	}

	//	fire the metrics pusher
	if ms.optionalMetricsPusher {
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




