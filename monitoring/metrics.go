package monitoring

import (
	"time"
	"github.com/rcrowley/go-metrics"
	"github.com/vrischmann/go-metrics-influxdb"
	"io"
)

func StartInfluxDbPusher(interval time.Duration, hostUrl, database, user, password string) {
	influxdb.InfluxDB(metricsRegistry,
		interval,
		hostUrl,
		database,
		user,
		password,
	)
}

func RegisterTimer(names ...string) {
	for _, name := range names {
		metrics.GetOrRegisterTimer(name, metricsRegistry)
	}
}

func UpdateTimerSince(name string, ts time.Time, unit time.Duration) {
	metrics.GetOrRegisterTimer(name, metricsRegistry).Update(time.Since(ts) / unit)
}

func WriteJsonMetrics(w io.Writer) {
	metrics.WriteJSONOnce(metricsRegistry, w)
}

var metricsRegistry metrics.Registry


func init() {
	metricsRegistry = metrics.NewRegistry()
}




