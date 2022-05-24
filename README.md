# vulcan-metrics-client

Metrics client abstraction used in Vulcan components.

## Configuration and behaviour

Internally the client abstraction works as a pool of clients for the supported implementations. Each client is enabled/disabled and configured through environment variables, so for each input metric it will be pushed through every client in the pool.

Current supported clients and its configurations are:
DataDog

- DOGSTATSD_ENABLED
- DOGSTATSD_HOST
- DOGSTATSD_PORT

Metrics are divided in `Metric`, which defaults to a sample rate of 1.0, and `RatedMetric` which allows to specify its sample rate between 0 (everything is sampled) and 1 (no sample); see the [sample rates section](https://docs.datadoghq.com/developers/metrics/dogstatsd_metrics_submission/?tab=go#sample-rates) in DataDog.

Current supported metric types are:

- count
- gauge
- histogram
- distribution

Example of pushing a metric:

```go
metricsClient, err := metrics.NewClient()
if err != nil {
    return err
}

metricsClient.Push(metrics.Metric{
    Name: "vulcan.scan.count",
    Typ: metrics.Count,
    Value: 1.0,
    Tags: []string{"team:purple"},
})
```

Example of pushing a rated metric:

```go
metricsClient, err := metrics.NewClient()
if err != nil {
    return err
}

metricsClient.PushWithRate(metrics.RatedMetric{
    Metric: Metric{
        Name:  "vulcan.requests.count",
        Typ:   metrics.Count,
        Value: 1,
        Tags:  []string{"team:purple"},
    },
    Rate: 0.5,
})
```
