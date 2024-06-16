package RegistryServiceClient

import (
	"context"
	"log"
	"math/rand"
	"time"

	service "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/common"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type RegistryServiceClient struct {
	client service.RegistryServiceClient
	conn   *grpc.ClientConn
}

func pickNode(addresses []string) string {
	rand.Seed(time.Now().UnixNano())
	return addresses[rand.Intn(len(addresses))]
}
func NewRegistryServiceClient(addresses []string) *RegistryServiceClient {
	address := pickNode(addresses)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil
	}
	client := service.NewRegistryServiceClient(conn)
	return &RegistryServiceClient{client: client, conn: conn}
}

func (c *RegistryServiceClient) Register(serviceName, nodeAddress string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := c.client.Register(ctx, &service.RegisterRequest{ServiceName: serviceName, NodeAddress: nodeAddress})
	return err
}

func (c *RegistryServiceClient) Unregister(serviceName, nodeAddress string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := c.client.Unregister(ctx, &service.UnregisterRequest{ServiceName: serviceName, NodeAddress: nodeAddress})
	return err
}

func (c *RegistryServiceClient) Discover(serviceName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	res, err := c.client.Discover(ctx, &service.DiscoverRequest{ServiceName: serviceName})
	if err != nil {
		return nil, err
	}
	return res.NodeAddresses, nil
}

func (c *RegistryServiceClient) Close() {
	if err := c.conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}

func (s *RegistryServiceClient) IsAlive(ctx context.Context, req *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	return &wrapperspb.BoolValue{Value: true}, nil
}
