package servant

import (
	"context"

	pb "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service"
)

type RegistryServant struct {
	client pb.RegistryServiceClient
}

func NewRegistryServant(client pb.RegistryServiceClient) *RegistryServant {
	return &RegistryServant{client: client}
}

func (s *RegistryServant) IsAlive(ctx context.Context, req *pb.IsAliveRequest) (*pb.IsAliveResponse, error) {
	// Implement the health check logic
	return s.client.IsAlive(ctx, req)
}
