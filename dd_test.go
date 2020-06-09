package metrics

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DataDog/datadog-go/statsd"
)

type statsdCountMetric struct {
	name  string
	value int64
	rate  float64
	tags  []string
}

type statsdDefaultMetric struct {
	name  string
	value float64
	rate  float64
	tags  []string
}

type mockStatsdClient struct {
	statsd.ClientInterface

	countMetrics        []statsdCountMetric
	gaugeMetrics        []statsdDefaultMetric
	histogramMetrics    []statsdDefaultMetric
	distributionMetrics []statsdDefaultMetric

	expectedCountMetrics        []statsdCountMetric
	expectedGaugeMetrics        []statsdDefaultMetric
	expectedHistogramMetrics    []statsdDefaultMetric
	expectedDistributionMetrics []statsdDefaultMetric
}

func (c *mockStatsdClient) Count(name string, value int64, tags []string, rate float64) error {
	c.countMetrics = append(c.countMetrics, statsdCountMetric{
		name, value, rate, tags,
	})
	return nil
}

func (c *mockStatsdClient) Gauge(name string, value float64, tags []string, rate float64) error {
	c.gaugeMetrics = append(c.gaugeMetrics, statsdDefaultMetric{
		name, value, rate, tags,
	})
	return nil
}

func (c *mockStatsdClient) Histogram(name string, value float64, tags []string, rate float64) error {
	c.histogramMetrics = append(c.histogramMetrics, statsdDefaultMetric{
		name, value, rate, tags,
	})
	return nil
}

func (c *mockStatsdClient) Distribution(name string, value float64, tags []string, rate float64) error {
	c.distributionMetrics = append(c.distributionMetrics, statsdDefaultMetric{
		name, value, rate, tags,
	})
	return nil
}

func (c *mockStatsdClient) Verify() error {
	// Verify length mismatch
	var countMismatch bool
	var mismatchMetric string
	var expectedCount, actualCount int

	if len(c.countMetrics) != len(c.expectedCountMetrics) {
		countMismatch = true
		mismatchMetric = "count"
		actualCount = len(c.countMetrics)
		expectedCount = len(c.expectedCountMetrics)
	} else if len(c.gaugeMetrics) != len(c.expectedGaugeMetrics) {
		countMismatch = true
		mismatchMetric = "gauge"
		actualCount = len(c.gaugeMetrics)
		expectedCount = len(c.expectedGaugeMetrics)
	} else if len(c.histogramMetrics) != len(c.expectedHistogramMetrics) {
		countMismatch = true
		mismatchMetric = "histogram"
		actualCount = len(c.histogramMetrics)
		expectedCount = len(c.expectedHistogramMetrics)
	} else if len(c.distributionMetrics) != len(c.expectedDistributionMetrics) {
		countMismatch = true
		mismatchMetric = "distribution"
		actualCount = len(c.distributionMetrics)
		expectedCount = len(c.expectedDistributionMetrics)
	}

	if countMismatch {
		return fmt.Errorf(
			"Error, count mismatch for %s metrics. Expected %d, but got %d",
			mismatchMetric, expectedCount, actualCount)
	}

	// Verify content mismatch
	var contentMismatch bool
	var expectedContent, actualContent interface{}

	// count
	for _, cm := range c.countMetrics {
		var found bool
		for _, ecm := range c.expectedCountMetrics {
			if reflect.DeepEqual(cm, ecm) {
				found = true
				break
			}
		}
		if !found {
			contentMismatch = true
			mismatchMetric = "count"
			expectedContent = c.expectedCountMetrics
			actualContent = c.countMetrics
		}
	}

	if contentMismatch {
		return verificationError(mismatchMetric, expectedContent, actualContent)
	}

	// gauge
	for _, gm := range c.gaugeMetrics {
		var found bool
		for _, egm := range c.expectedGaugeMetrics {
			if reflect.DeepEqual(gm, egm) {
				found = true
				break
			}
		}
		if !found {
			contentMismatch = true
			mismatchMetric = "gauge"
			expectedContent = c.expectedGaugeMetrics
			actualContent = c.gaugeMetrics
		}
	}

	if contentMismatch {
		return verificationError(mismatchMetric, expectedContent, actualContent)
	}

	// histogram
	for _, hm := range c.histogramMetrics {
		var found bool
		for _, ehm := range c.expectedHistogramMetrics {
			if reflect.DeepEqual(hm, ehm) {
				found = true
				break
			}
		}
		if !found {
			contentMismatch = true
			mismatchMetric = "histogram"
			expectedContent = c.expectedHistogramMetrics
			actualContent = c.histogramMetrics
		}
	}

	if contentMismatch {
		return verificationError(mismatchMetric, expectedContent, actualContent)
	}

	// distribution
	for _, dm := range c.distributionMetrics {
		var found bool
		for _, edm := range c.expectedDistributionMetrics {
			if reflect.DeepEqual(dm, edm) {
				found = true
				break
			}
		}
		if !found {
			contentMismatch = true
			mismatchMetric = "distribution"
			expectedContent = c.expectedDistributionMetrics
			actualContent = c.distributionMetrics
		}
	}

	if contentMismatch {
		return verificationError(mismatchMetric, expectedContent, actualContent)
	}

	return nil
}

func verificationError(metric string, expected, actual interface{}) error {
	return fmt.Errorf(
		"Error, content mismatch for %s metrics. Expected %v, but got %v",
		metric, expected, actual)
}

func TestDDPush(t *testing.T) {
	testCases := []struct {
		metricsEnabled              bool
		inputMetrics                []Metric
		expectedCountMetrics        []statsdCountMetric
		expectedGaugeMetrics        []statsdDefaultMetric
		expectedHistogramMetrics    []statsdDefaultMetric
		expectedDistributionMetrics []statsdDefaultMetric
	}{
		{
			metricsEnabled: true,
			inputMetrics: []Metric{
				{
					Name:  "countMetricA",
					Typ:   Count,
					Value: 1,
					Tags:  []string{"tag:mytag"},
				},
			},
			expectedCountMetrics: []statsdCountMetric{
				{
					name:  "countMetricA",
					value: 1,
					rate:  1,
					tags:  []string{"tag:mytag"},
				},
			},
		},
		{
			metricsEnabled: true,
			inputMetrics: []Metric{
				{
					Name:  "countMetricA",
					Typ:   Count,
					Value: 1,
					Tags:  []string{"tag:mytag"},
				},
				{
					Name:  "countMetricB",
					Typ:   Count,
					Value: 2,
					Tags:  []string{"tag:mytagB"},
				},
				{
					Name:  "histogramMetricA",
					Typ:   Histogram,
					Value: 4.2,
					Tags:  []string{"tag:mytagC"},
				},
			},
			expectedCountMetrics: []statsdCountMetric{
				{
					name:  "countMetricA",
					value: 1,
					rate:  1,
					tags:  []string{"tag:mytag"},
				},
				{
					name:  "countMetricB",
					value: 2,
					rate:  1,
					tags:  []string{"tag:mytagB"},
				},
			},
			expectedHistogramMetrics: []statsdDefaultMetric{
				{
					name:  "histogramMetricA",
					value: 4.2,
					rate:  1,
					tags:  []string{"tag:mytagC"},
				},
			},
		},
		{
			metricsEnabled: true,
			inputMetrics: []Metric{
				{
					Name:  "histogramMetricA",
					Typ:   Histogram,
					Value: 4.2,
					Tags:  []string{"tag:mytagC"},
				},
				{
					Name:  "histogramMetricB",
					Typ:   Histogram,
					Value: 5.2,
					Tags:  []string{"tag:mytagC"},
				},
				{
					Name:  "gaugeMetricA",
					Typ:   Gauge,
					Value: 7.1,
					Tags:  []string{"tag:mytagD"},
				},
				{
					Name:  "distributionMetricA",
					Typ:   Distribution,
					Value: 1.3,
					Tags:  []string{"tag:mytagE"},
				},
			},
			expectedHistogramMetrics: []statsdDefaultMetric{
				{
					name:  "histogramMetricA",
					value: 4.2,
					rate:  1,
					tags:  []string{"tag:mytagC"},
				},
				{
					name:  "histogramMetricB",
					value: 5.2,
					rate:  1,
					tags:  []string{"tag:mytagC"},
				},
			},
			expectedGaugeMetrics: []statsdDefaultMetric{
				{
					name:  "gaugeMetricA",
					value: 7.1,
					rate:  1,
					tags:  []string{"tag:mytagD"},
				},
			},
			expectedDistributionMetrics: []statsdDefaultMetric{
				{
					name:  "distributionMetricA",
					value: 1.3,
					rate:  1,
					tags:  []string{"tag:mytagE"},
				},
			},
		},
		{
			metricsEnabled: false,
			inputMetrics: []Metric{
				{
					Name:  "countMetricA",
					Typ:   Count,
					Value: 1,
					Tags:  []string{"tag:mytag"},
				},
				{
					Name:  "countMetricB",
					Typ:   Count,
					Value: 2,
					Tags:  []string{"tag:mytagB"},
				},
				{
					Name:  "histogramMetricA",
					Typ:   Histogram,
					Value: 4.2,
					Tags:  []string{"tag:mytagC"},
				},
			},
			expectedCountMetrics:        []statsdCountMetric{},
			expectedHistogramMetrics:    []statsdDefaultMetric{},
			expectedGaugeMetrics:        []statsdDefaultMetric{},
			expectedDistributionMetrics: []statsdDefaultMetric{},
		},
	}

	for _, tc := range testCases {
		mockStatsdClient := &mockStatsdClient{
			expectedCountMetrics:        tc.expectedCountMetrics,
			expectedGaugeMetrics:        tc.expectedGaugeMetrics,
			expectedHistogramMetrics:    tc.expectedHistogramMetrics,
			expectedDistributionMetrics: tc.expectedDistributionMetrics,
		}

		ddClient := &ddClient{
			enabled:   tc.metricsEnabled,
			statsdCli: mockStatsdClient,
		}

		for _, im := range tc.inputMetrics {
			ddClient.Push(im)
		}

		if err := mockStatsdClient.Verify(); err != nil {
			t.Fatalf("Error pushing metrics: %v", err)
		}
	}
}

func TestDDPushWithRate(t *testing.T) {
	testCases := []struct {
		metricsEnabled              bool
		inputMetrics                []RatedMetric
		expectedCountMetrics        []statsdCountMetric
		expectedGaugeMetrics        []statsdDefaultMetric
		expectedHistogramMetrics    []statsdDefaultMetric
		expectedDistributionMetrics []statsdDefaultMetric
	}{
		{
			metricsEnabled: true,
			inputMetrics: []RatedMetric{
				{
					Metric: Metric{
						Name:  "countMetricA",
						Typ:   Count,
						Value: 1,
						Tags:  []string{"tag:mytag"},
					},
					Rate: 0.4,
				},
			},
			expectedCountMetrics: []statsdCountMetric{
				{
					name:  "countMetricA",
					value: 1,
					rate:  0.4,
					tags:  []string{"tag:mytag"},
				},
			},
		},
		{
			metricsEnabled: true,
			inputMetrics: []RatedMetric{
				{
					Metric: Metric{
						Name:  "histogramMetricA",
						Typ:   Histogram,
						Value: 4.2,
						Tags:  []string{"tag:mytagC"},
					},
					Rate: 0.2,
				},
				{
					Metric: Metric{
						Name:  "histogramMetricB",
						Typ:   Histogram,
						Value: 5.2,
						Tags:  []string{"tag:mytagC"},
					},
					Rate: 1,
				},
				{
					Metric: Metric{
						Name:  "gaugeMetricA",
						Typ:   Gauge,
						Value: 7.1,
						Tags:  []string{"tag:mytagD"},
					},
					Rate: 0.8,
				},
				{
					Metric: Metric{
						Name:  "distributionMetricA",
						Typ:   Distribution,
						Value: 1.3,
						Tags:  []string{"tag:mytagE"},
					},
					Rate: 0.5,
				},
			},
			expectedHistogramMetrics: []statsdDefaultMetric{
				{
					name:  "histogramMetricA",
					value: 4.2,
					rate:  0.2,
					tags:  []string{"tag:mytagC"},
				},
				{
					name:  "histogramMetricB",
					value: 5.2,
					rate:  1,
					tags:  []string{"tag:mytagC"},
				},
			},
			expectedGaugeMetrics: []statsdDefaultMetric{
				{
					name:  "gaugeMetricA",
					value: 7.1,
					rate:  0.8,
					tags:  []string{"tag:mytagD"},
				},
			},
			expectedDistributionMetrics: []statsdDefaultMetric{
				{
					name:  "distributionMetricA",
					value: 1.3,
					rate:  0.5,
					tags:  []string{"tag:mytagE"},
				},
			},
		},
		{
			metricsEnabled: false,
			inputMetrics: []RatedMetric{
				{
					Metric: Metric{
						Name:  "countMetricA",
						Typ:   Count,
						Value: 1,
						Tags:  []string{"tag:mytag"},
					},
					Rate: 0.4,
				},
				{
					Metric: Metric{
						Name:  "histogramMetricA",
						Typ:   Histogram,
						Value: 4.2,
						Tags:  []string{"tag:mytagC"},
					},
					Rate: 0.2,
				},
				{
					Metric: Metric{
						Name:  "distributionMetricA",
						Typ:   Distribution,
						Value: 1.3,
						Tags:  []string{"tag:mytagE"},
					},
					Rate: 0.5,
				},
			},
			expectedCountMetrics:        []statsdCountMetric{},
			expectedHistogramMetrics:    []statsdDefaultMetric{},
			expectedGaugeMetrics:        []statsdDefaultMetric{},
			expectedDistributionMetrics: []statsdDefaultMetric{},
		},
	}

	for _, tc := range testCases {
		mockStatsdClient := &mockStatsdClient{
			expectedCountMetrics:        tc.expectedCountMetrics,
			expectedGaugeMetrics:        tc.expectedGaugeMetrics,
			expectedHistogramMetrics:    tc.expectedHistogramMetrics,
			expectedDistributionMetrics: tc.expectedDistributionMetrics,
		}

		ddClient := &ddClient{
			enabled:   tc.metricsEnabled,
			statsdCli: mockStatsdClient,
		}

		for _, im := range tc.inputMetrics {
			ddClient.PushWithRate(im)
		}

		if err := mockStatsdClient.Verify(); err != nil {
			t.Fatalf("Error pushing metrics: %v", err)
		}
	}
}
