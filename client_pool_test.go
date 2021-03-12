/*
Copyright 2021 Adevinta
*/

package metrics

import "testing"

type mockClient struct {
	pushCalls         int
	pushWithRateCalls int
}

func (m *mockClient) Push(metric Metric) {
	m.pushCalls++
}

func (m *mockClient) PushWithRate(ratedMetric RatedMetric) {
	m.pushWithRateCalls++
}

func TestPoolPush(t *testing.T) {
	testCases := []struct {
		name                string
		inputMetrics        []Metric
		nClients            int
		expectedClientCalls int
	}{
		{
			name: "Happy path",
			inputMetrics: []Metric{
				{
					Name:  "countMetricA",
					Typ:   Count,
					Value: 1,
					Tags:  []string{"tag:mytag"},
				},
			},
			nClients:            1,
			expectedClientCalls: 1,
		},
		{
			name: "Should push 3 metrics with 4 clients",
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
			nClients:            4,
			expectedClientCalls: 3,
		},
	}

	for _, tc := range testCases {
		// given
		var mockClients []*mockClient
		for i := 0; i < tc.nClients; i++ {
			mockClients = append(mockClients, &mockClient{})
		}

		clientPool := &clientPool{}
		for i := 0; i < tc.nClients; i++ {
			clientPool.clients = append(clientPool.clients, mockClients[i])
		}

		// when
		for _, m := range tc.inputMetrics {
			clientPool.Push(m)
		}

		// then
		for _, mc := range mockClients {
			if mc.pushCalls != tc.expectedClientCalls {
				t.Fatalf("Error, expected push calls for each client to be: %d\nBut got: %d",
					tc.expectedClientCalls, mc.pushCalls)
			}
		}
	}
}

func TestPoolPushWithRate(t *testing.T) {
	testCases := []struct {
		name                string
		inputMetrics        []RatedMetric
		nClients            int
		expectedClientCalls int
	}{
		{
			name: "Happy path",
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
			nClients:            1,
			expectedClientCalls: 1,
		},
		{
			name: "Should push 2 metrics with 5 clients",
			inputMetrics: []RatedMetric{
				{
					Metric: Metric{
						Name:  "histogramMetricA",
						Typ:   Histogram,
						Value: 4.2,
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
			},
			nClients:            5,
			expectedClientCalls: 2,
		},
	}

	for _, tc := range testCases {
		// given
		var mockClients []*mockClient
		for i := 0; i < tc.nClients; i++ {
			mockClients = append(mockClients, &mockClient{})
		}

		clientPool := &clientPool{}
		for i := 0; i < tc.nClients; i++ {
			clientPool.clients = append(clientPool.clients, mockClients[i])
		}

		// when
		for _, m := range tc.inputMetrics {
			clientPool.PushWithRate(m)
		}

		// then
		for _, mc := range mockClients {
			if mc.pushWithRateCalls != tc.expectedClientCalls {
				t.Fatalf("Error, expected push calls for each client to be: %d\nBut got: %d",
					tc.expectedClientCalls, mc.pushWithRateCalls)
			}
		}
	}
}
