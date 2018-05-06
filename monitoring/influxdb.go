package monitoring

import (
	"time"
	"fmt"

	"github.com/vrischmann/go-metrics-influxdb"
)

type InfluxDbProperties struct {
	Host string
	Port int
	Database string
	User string
	Password string
}

func RunInfluxDbMetricsPusher(p *InfluxDbProperties, interval int) {
	influxdb.InfluxDB(
		metricsRegistry,
		time.Second * time.Duration(interval),
		fmt.Sprintf("http://%s:%d", p.Host, p.Port),
		p.Database,
		p.User,
		p.Password,
	)
}
