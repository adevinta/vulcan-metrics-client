# vulcan-metrics-client
Metrics client abstraction used in Vulcan components.

Supports metrics types:
- count
- gauge
- histogram
- distribution

Currently supports implementation for DataDog statsd service only.

Metrics are divided in `Metric`, which defaults to a sample rate of 1.0, and `RatedMetric` which allows to specify its sample rate between 0 (everything is sampled) and 1 (no sample); see the [sample rates section](https://docs.datadoghq.com/developers/metrics/dogstatsd_metrics_submission/?tab=go#sample-rates) in DataDog.

Example of pushing a metric:
```
metricsClient, err := metrics.NewClient(metrics.DataDogClient, true, "localhost", 8125)
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
```
metricsClient, err := metrics.NewClient(metrics.DataDogClient, true, "localhost", 8125)
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
