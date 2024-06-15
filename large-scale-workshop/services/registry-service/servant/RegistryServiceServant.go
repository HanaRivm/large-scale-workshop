package RegistryServiceServant

import (
	"context"

	client "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/common"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type RegistryServant struct {
	client client.RegistryServiceClient
}

func NewRegistryServant(client client.RegistryServiceClient) *RegistryServant {
	return &RegistryServant{client: client}
}

func (s *RegistryServant) IsAlive(ctx context.Context, req *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	// Implement the health check logic
	return s.client.IsAlive(ctx, req)
}
