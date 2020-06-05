package metrics

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

const (
	defaultHost = "localhost"
	defaultPort = 8125
)

// DDClient represents a DataDog metrics client.
type DDClient struct {
	enabled   bool
	statsdCli statsd.ClientInterface
}

type pushMetricsFunc func(name string, value, rate float64, tags []string)

// NewDDClient creates a new DataDog metrics client.
func NewDDClient(enabled bool, statsdHost string, statsdPort int) (*DDClient, error) {
	if statsdHost == "" {
		statsdHost = defaultHost
	}
	if statsdPort == 0 {
		statsdPort = defaultPort
	}

	statsdAddrs := fmt.Sprintf("%s:%d", statsdHost, statsdPort)
	statsdCli, err := statsd.New(statsdAddrs)
	if err != nil {
		return nil, err
	}

	return &DDClient{
		enabled:   enabled,
		statsdCli: statsdCli,
	}, nil
}

// Push pushes the specified metric with rate 1.
func (c *DDClient) Push(m Metric) {
	if !c.enabled {
		return
	}
	c.PushWithRate(RatedMetric{Metric: m, Rate: 1})
}

// PushWithRate pushes the input metric with the specified rate.
func (c *DDClient) PushWithRate(m RatedMetric) {
	if !c.enabled {
		return
	}

	var pushFunc pushMetricsFunc

	switch m.Typ {
	case Count:
		pushFunc = c.count
	case Gauge:
		pushFunc = c.gauge
	case Histogram:
		pushFunc = c.histogram
	case Distribution:
		pushFunc = c.distribution
	default:
		return
	}

	pushFunc(m.Name, m.Value, m.Rate, m.Tags)
}

func (c *DDClient) count(name string, value, rate float64, tags []string) {
	c.statsdCli.Count(name, int64(value), tags, rate)
}

func (c *DDClient) gauge(name string, value, rate float64, tags []string) {
	c.statsdCli.Gauge(name, value, tags, rate)
}

func (c *DDClient) histogram(name string, value, rate float64, tags []string) {
	c.statsdCli.Histogram(name, value, tags, rate)
}

func (c *DDClient) distribution(name string, value, rate float64, tags []string) {
	c.statsdCli.Distribution(name, value, tags, rate)
}
