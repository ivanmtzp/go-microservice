package monitoring

import (
	"time"
	"fmt"

	"github.com/influxdata/influxdb/client"
	"github.com/rcrowley/go-metrics"
	"net/url"
)


type InfluxDbPusher struct {
	address string
	database string
	username string
	password string
	tags     map[string]string
	interval time.Duration
	client *client.Client
}

func NewInfluxDbPusher (address, username, password, database string, tags map[string]string, interval time.Duration) (*InfluxDbPusher, error) {
	url, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	client, err := client.NewClient(client.Config{
		URL:      *url,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return &InfluxDbPusher{
		address: address,
		database: database,
		username: username,
		password: password,
		tags: tags,
		interval: interval,
		client: client,
	}, nil
}

func (p *InfluxDbPusher) Address() string {
	return p.address
}

func (p *InfluxDbPusher) Database() string {
	return p.database
}

func (p *InfluxDbPusher) Run() error {
	ticker := time.Tick(p.interval)
	for range ticker {
		if err := p.send(); err != nil {
			return err
		}
	}
	return nil
}

func (p *InfluxDbPusher) send() error {
	var pts []client.Point

	metricsRegistry.Each(func(name string, i interface{}) {
		now := time.Now()

		switch metric := i.(type) {
		case metrics.Counter:
			ms := metric.Snapshot()
			pts = append(pts, client.Point{
				Measurement: fmt.Sprintf("%s.count", name),
				Tags:        p.tags,
				Fields: map[string]interface{}{
					"value": ms.Count(),
				},
				Time: now,
			})
		case metrics.Gauge:
			ms := metric.Snapshot()
			pts = append(pts, client.Point{
				Measurement: fmt.Sprintf("%s.gauge", name),
				Tags:        p.tags,
				Fields: map[string]interface{}{
					"value": ms.Value(),
				},
				Time: now,
			})
		case metrics.GaugeFloat64:
			ms := metric.Snapshot()
			pts = append(pts, client.Point{
				Measurement: fmt.Sprintf("%s.gauge", name),
				Tags:        p.tags,
				Fields: map[string]interface{}{
					"value": ms.Value(),
				},
				Time: now,
			})
		case metrics.Histogram:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
			pts = append(pts, client.Point{
				Measurement: fmt.Sprintf("%s.histogram", name),
				Tags:        p.tags,
				Fields: map[string]interface{}{
					"count":    ms.Count(),
					"max":      ms.Max(),
					"mean":     ms.Mean(),
					"min":      ms.Min(),
					"stddev":   ms.StdDev(),
					"variance": ms.Variance(),
					"p50":      ps[0],
					"p75":      ps[1],
					"p95":      ps[2],
					"p99":      ps[3],
					"p999":     ps[4],
					"p9999":    ps[5],
				},
				Time: now,
			})
		case metrics.Meter:
			ms := metric.Snapshot()
			pts = append(pts, client.Point{
				Measurement: fmt.Sprintf("%s.meter", name),
				Tags:        p.tags,
				Fields: map[string]interface{}{
					"count": ms.Count(),
					"m1":    ms.Rate1(),
					"m5":    ms.Rate5(),
					"m15":   ms.Rate15(),
					"mean":  ms.RateMean(),
				},
				Time: now,
			})
		case metrics.Timer:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
			pts = append(pts, client.Point{
				Measurement: fmt.Sprintf("%s.timer", name),
				Tags:        p.tags,
				Fields: map[string]interface{}{
					"count":    ms.Count(),
					"max":      ms.Max(),
					"mean":     ms.Mean(),
					"min":      ms.Min(),
					"stddev":   ms.StdDev(),
					"variance": ms.Variance(),
					"p50":      ps[0],
					"p75":      ps[1],
					"p95":      ps[2],
					"p99":      ps[3],
					"p999":     ps[4],
					"p9999":    ps[5],
					"m1":       ms.Rate1(),
					"m5":       ms.Rate5(),
					"m15":      ms.Rate15(),
					"meanrate": ms.RateMean(),
				},
				Time: now,
			})
		}
	})

	bps := client.BatchPoints{
		Points:   pts,
		Database: p.database,
	}

	_, err := p.client.Write(bps)
	return err
}

