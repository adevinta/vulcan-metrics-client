package metrics

const (
	// Count represents a count metric type.
	Count = iota
	// Gauge represents a gauge metric type.
	Gauge = iota
	// Histogram represents a histogram metric type.
	Histogram = iota
	// Distribution represents a distribution metric type.
	Distribution = iota
)

// Type represents the type of metric to push.
// Supports Count, Gauge, Histogram and Distribution types.
type Type int

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

// NewClient creates a new metrics client based on environment variables.
func NewClient() (Client, error) {
	return newClientPool()
}
