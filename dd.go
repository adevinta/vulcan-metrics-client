package metrics

import (
	"fmt"
	"os"
	"strconv"

	"github.com/DataDog/datadog-go/statsd"
)

const (
	// DDEnabled represents the config env var to enabled DD client.
	DDEnabled = "DOGSTATSD_ENABLED"
	// DDHost represents the config env var to set DD client's statsd host.
	DDHost = "DOGSTATSD_HOST"
	// DDPort represents the config env var to set DD client's statsd port.
	DDPort = "DOGSTATSD_PORT"

	defaultHost = "localhost"
	defaultPort = 8125
)

// ddClient represents a DataDog metrics client.
type ddClient struct {
	enabled   bool
	statsdCli statsd.ClientInterface
}

type pushMetricsFunc func(name string, value, rate float64, tags []string)

// newDDClient creates a new DataDog metrics client.
func newDDClient(enabled bool, statsdHost string, statsdPort int) (Client, error) {
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

	return &ddClient{
		enabled:   enabled,
		statsdCli: statsdCli,
	}, nil
}

// newDDClientFromEnv creates a new DataDog metrics client
// reading its configuration from environment.
func newDDClientFromEnv() (Client, error) {
	enabled, _ := strconv.ParseBool(os.Getenv(DDEnabled))
	host := os.Getenv(DDHost)
	port, _ := strconv.ParseInt(os.Getenv(DDPort), 10, 0)

	return newDDClient(enabled, host, int(port))
}

// Push pushes the specified metric with rate 1.
func (c *ddClient) Push(m Metric) {
	if !c.enabled {
		return
	}
	c.PushWithRate(RatedMetric{Metric: m, Rate: 1})
}

// PushWithRate pushes the input metric with the specified rate.
func (c *ddClient) PushWithRate(m RatedMetric) {
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

func (c *ddClient) count(name string, value, rate float64, tags []string) {
	c.statsdCli.Count(name, int64(value), tags, rate)
}

func (c *ddClient) gauge(name string, value, rate float64, tags []string) {
	c.statsdCli.Gauge(name, value, tags, rate)
}

func (c *ddClient) histogram(name string, value, rate float64, tags []string) {
	c.statsdCli.Histogram(name, value, tags, rate)
}

func (c *ddClient) distribution(name string, value, rate float64, tags []string) {
	c.statsdCli.Distribution(name, value, tags, rate)
}
