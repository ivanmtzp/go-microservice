package microservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/ivanmtzp/go-microservice/metrics"
	"github.com/ivanmtzp/go-microservice/log"
	"github.com/ivanmtzp/go-microservice/config"
	"github.com/ivanmtzp/go-microservice/grpc"
)


var grpcServer *grpc.Server

func Configure (name, envPrefix string) {
	log.SetAppName(name)

	if err := config.Read( strings.ToLower(envPrefix), "./config", "microservice", config.Yaml); err != nil {
		log.FailOnError(err, "error reading configuration file: ./config/microservice.yaml")
	}

	logLevel := config.GetString("log", "level")
	if logLevel != "" {
		err := log.SetLevel(logLevel)
		if err != nil {
			log.FailOnError(err, fmt.Sprint("configuration error, invalid log level: ", logLevel))
		}
	}
	log.Info("Log level set to ", log.Level())
	log.Environment(strings.ToUpper(envPrefix))
}

func RegisterGrpcServiceWithHttpGateway(grpcSettings *GrpcSettings, ) {
	grpcServer = grpc.New(grpcSettings.address)

}



func Run() {
	log.Info("starting microservice")

	pop.ConfigName = "microservice.yml"
	db, err := pop.Connect("database")
	if err != nil {
		log.FailOnError(err, "failed to connect to database")
	}
	defer db.Close()

	// err = handlers.PingDb(db)
	// if err != nil {
	//	log.FailOnError(err, fmt.Sprintf("ping check to database failed, database url: %s", db.URL()))
	// }

	grpcSettings := ReadGrpcSettings()
	if grpcSettings.address != "" {
		// fire the gRPC server in a goroutine
		go func() {
			log.Infof("starting HTTP/2 gRPC server on %s", grpcSettings.address)
			grpcServer.Run()
			log.FailOnError(err, fmt.Sprint("Failed to start gRPC server " , grpcSettings.address))
		}()
		// fire the http grpc gateway in a goroutine
		go func() {
			log.Infof("starting HTTP/1.1 REST server on %s", restAddress)
			err := grpc.RunRestGrpcGatewayServer(restAddress, grpcAddress)
			log.FailOnError(err, fmt.Sprintf("Failed to start Http REST server. Host: %s, Port: %d", host, restPort))
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




