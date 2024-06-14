package servant

import (
	"context"

	"github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/client"
)

type RegistryServant struct {
	client client.RegistryServiceClient
}

func NewRegistryServant(client client.RegistryServiceClient) *RegistryServant {
	return &RegistryServant{client: client}
}

func (s *RegistryServant) IsAlive(ctx context.Context, req *client.IsAliveRequest) (*client.IsAliveResponse, error) {
	// Implement the health check logic
	return s.client.IsAlive(ctx, req)
}
