package metrics

// ClientPool represetns a metris client pool.
type clientPool struct {
	clients []Client
}

// clientConstructor represents a metrics client
// constructor function.
type clientConstructor func() (Client, error)

var (
	supportedClients = map[string]clientConstructor{
		"DataDog": clientConstructor(newDDClientFromEnv),
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
func newClientPool() (Client, error) {
	var clients []Client

	for _, cconstructor := range supportedClients {
		c, err := cconstructor()
		if err != nil {
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
