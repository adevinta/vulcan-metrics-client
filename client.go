package metrics

import (
	"errors"
	"os"
	"strconv"
)

const (
	// Config env variables.
	envClientType = "METRICS_CLIENT_TYPE"
	envEnabled    = "METRICS_ENABLED"
	envHost       = "METRICS_HOST"
	envPort       = "METRICS_PORT"

	// DataDogClient represents the DataDog
	// implementation for metrics client.
	DataDogClient = "dd"

	// Count represents a count metric type.
	Count = iota
	// Gauge represents a gauge metric type.
	Gauge = iota
	// Histogram represents a histogram metric type.
	Histogram = iota
	// Distribution represents a distribution metric type.
	Distribution = iota
)

var (
	// ErrUnsupportedClientType indicates that the requested client type is not supported.
	ErrUnsupportedClientType = errors.New("Error: unsupported client type")
)

// Type represents the type of metric to push.
// Supports Count, Gauge, Histogram and Distribution types.
type Type int

// ClientType represents the requested metrics
// client implementation.
type ClientType string

// Metric represents a metric.
type Metric struct {
	Name  string
	Typ   Type
	Value float64
	Tags  []string
}

// RatedMetric represents a metric with rate.
type RatedMetric struct {
	Metric
	Rate float64
}

// Client represents a metrics service client.
type Client interface {
	Push(metric Metric)
	PushWithRate(ratedMetric RatedMetric)
}

// NewClient creates a new metrics client based on given input.
func NewClient(typ ClientType, enabled bool, host string, port int) (Client, error) {
	switch typ {
	case DataDogClient:
		return newDDClient(enabled, host, port)
	default:
		return nil, ErrUnsupportedClientType
	}
}

// NewClientFromEnv creates a new metrics client based on environment variables.
func NewClientFromEnv() (Client, error) {
	clientType := ClientType(os.Getenv(envClientType))
	enabled, _ := strconv.ParseBool(os.Getenv(envEnabled))
	host := os.Getenv(envHost)
	port, _ := strconv.ParseInt(os.Getenv(envPort), 10, 0)

	return NewClient(clientType, enabled, host, int(port))
}
