package registryservice

import (
	"sync"
	"time"
)

type RegistryService struct {
	mu       sync.Mutex
	services map[string]map[string]bool
}

func NewRegistryService() *RegistryService {
	return &RegistryService{
		services: make(map[string]map[string]bool),
	}
}

func (r *RegistryService) Register(serviceName, nodeAddress string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.services[serviceName]; !exists {
		r.services[serviceName] = make(map[string]bool)
	}
	r.services[serviceName][nodeAddress] = true
}

func (r *RegistryService) Unregister(serviceName, nodeAddress string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if nodes, exists := r.services[serviceName]; exists {
		delete(nodes, nodeAddress)
		if len(nodes) == 0 {
			delete(r.services, serviceName)
		}
	}
}

func (r *RegistryService) Discover(serviceName string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodes := r.services[serviceName]
	addresses := make([]string, 0, len(nodes))
	for address := range nodes {
		addresses = append(addresses, address)
	}
	return addresses
}

func (r *RegistryService) IsAlive() {
	for {
		time.Sleep(10 * time.Second)
		r.mu.Lock()
		for serviceName, nodes := range r.services {
			for nodeAddress := range nodes {
				// Call IsAlive gRPC function on nodeAddress
				// If node is not alive twice in a row, unregister it
			}
		}
		r.mu.Unlock()
	}
}
