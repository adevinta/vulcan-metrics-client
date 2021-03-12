/*
Copyright 2021 Adevinta
*/

package metrics

// clientPool represents a metrics client pool.
type clientPool struct {
	clients []Client
}

// clientConstructor represents a metrics client
// constructor function.
type clientConstructor func() (Client, error)

var (
	supportedClients = map[string]clientConstructor{
		"DataDog": clientConstructor(newDDClient),
	}
)

// newClientPool creates a new ClientPool
// parsing configuration from environment.
//
// Supported clients:
//	- DataDog:
//		DOGSTATSD_ENABLED
//		DOGSTATSD_HOST
//		DOGSTATSD_PORT
func newClientPool() (*clientPool, error) {
	var clients []Client

	for _, cConstructor := range supportedClients {
		c, err := cConstructor()
		if err == nil {
			clients = append(clients, c)
		}
	}

	return &clientPool{
		clients: clients,
	}, nil
}

// Push pushes the input metric for each
// client in the pool.
func (p *clientPool) Push(metric Metric) {
	for _, c := range p.clients {
		c.Push(metric)
	}
}

// PushWithRate pushes the input rated metric
// for each client in the pool.
func (p *clientPool) PushWithRate(ratedMetric RatedMetric) {
	for _, c := range p.clients {
		c.PushWithRate(ratedMetric)
	}
}
