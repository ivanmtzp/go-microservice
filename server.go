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


type MicroService struct
{
	name string
}

func NewMicroService (name, envPrefix string) *MicroService {
	log.AppName(name)

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

	return &MicroService{name: name}
}


func (m *MicroService) Run() {
	log.Info("starting microservice")

	host := config.GetString("host")
	grpcPort := config.GetInt("grpc_port")
	restPort := config.GetInt("rest_port")

	log.Infof("environment: %s, Host: %s, Grpc port: %d, Rest Port: %d", config.Environment(), host, grpcPort, restPort)

	grpcAddress := fmt.Sprintf("%s:%d", host, grpcPort)
	restAddress := fmt.Sprintf("%s:%d", host, restPort)

	db, err := pop.Connect(config.Environment())
	if err != nil {
		log.FailOnError(err, "failed to connect to database")
	}
	defer db.Close()

	// err = handlers.PingDb(db)
	// if err != nil {
	//	log.FailOnError(err, fmt.Sprintf("ping check to database failed, database url: %s", db.URL()))
	// }

	// fire the gRPC server in a goroutine
	go func() {
		log.Infof("starting HTTP/2 gRPC server on %s", grpcAddress)
		err := grpc.RunGrpcServer(grpcAddress)
		log.FailOnError(err, fmt.Sprintf("Failed to start gRPC server. Host: %s, Port: %d", host, grpcPort))
	}()

	// fire the REST server in a goroutine
	go func() {
		log.Infof("starting HTTP/1.1 REST server on %s", restAddress)
		err := grpc.RunRestGrpcGatewayServer(restAddress, grpcAddress)
		log.FailOnError(err, fmt.Sprintf("Failed to start Http REST server. Host: %s, Port: %d", host, restPort))
	}()

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




