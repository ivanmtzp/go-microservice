package microservice

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/ivanmtzp/go-microservice/grpc"
	"github.com/ivanmtzp/go-microservice/log"
	"github.com/ivanmtzp/go-microservice/config"
	"github.com/ivanmtzp/go-microservice/metrics"
)


type MicroService struct {
	name string
	settings SettingsReader
	grpcServer *grpc.Server

}


func New(name string, sr SettingsReader) *MicroService {
	return &MicroService{name: name, settings: sr}
}

func (ms *MicroService) WithGrpcAndGateway(sr grpc.ServiceRegistrator,gsr grpc.HttpGatewayServiceRegistrator) *MicroService {
	ms.grpcServer = grpc.New(ms.settings.GrpcSettings().address, sr).
		WithGateway(ms.settings.GrpcSettings().gatewayAddress, gsr)
	return ms
}


func (ms *MicroService) Run() {

	pop.ConfigName = "microservice.yml"
	db, err := pop.Connect("local")
	if err != nil {
		log.FailOnError(err, "failed to connect to database")
	}
	defer db.Close()

	// err = handlers.PingDb(db)
	// if err != nil {
	//	log.FailOnError(err, fmt.Sprintf("ping check to database failed, database url: %s", db.URL()))
	// }

	if ms.grpcServer != nil {
		// fire the gRPC server in a goroutine
		go func() {
			log.Infof("starting HTTP/2 gRPC server on %s", ms.settings.GrpcSettings().address)
			err:= ms.grpcServer.Run()
			log.FailOnError(err, fmt.Sprint("failed to start gRPC server " , ms.settings.GrpcSettings().address))
		}()
		// fire the http grpc gateway in a goroutine
		go func() {
			log.Infof("starting HTTP/1.1 REST server on %s", ms.settings.GrpcSettings().gatewayAddress)
			err := ms.grpcServer.RunHttpGateway()
			log.FailOnError(err, fmt.Sprint("failed to start Http gateway server ", ms.settings.GrpcSettings().gatewayAddress))
		}()
	}

	//	fire the metrics pusher
	metricsPushEnabled := config.GetBool("metrics", "push_enabled")
	if metricsPushEnabled {
		go func() {
			metricsInterval := config.GetInt("metrics", "interval")
			metricsHostUrl := fmt.Sprintf("http://%s:%d", config.GetString("metrics", "host"),
				config.GetString("metrics", "port"))
			metricsDatabase := config.GetString("metrics", "database")
			metricsUser := config.GetString("metrics", "user")
			metricsPassword := config.GetString("metrics", "password")
			log.Infof("starting InfluxDb Metrics pushing to: %s, database: %s, user: %s,  with interval: %d", metricsHostUrl,
				metricsDatabase, metricsUser, metricsInterval)
			metrics.StartInfluxDbPusher(time.Second * time.Duration(metricsInterval), metricsHostUrl, metricsDatabase, metricsUser, metricsPassword)
		}()
	}

	// infinite loop
	select {}

}




