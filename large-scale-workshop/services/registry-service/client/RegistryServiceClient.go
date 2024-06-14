package client

import (
	"context"
	"log"
	"time"

	service "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service"
	"google.golang.org/grpc"
)

type RegistryClient struct {
	client service.RegistryServiceClient
	conn   *grpc.ClientConn
}

func NewRegistryClient(address string) (*RegistryClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*5))
	if err != nil {
		return nil, err
	}
	client := service.NewRegistryServiceClient(conn)
	return &RegistryClient{client: client, conn: conn}, nil
}

func (c *RegistryClient) Register(serviceName, nodeAddress string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.client.Register(ctx, &service.RegisterRequest{ServiceName: serviceName, NodeAddress: nodeAddress})
	return err
}

func (c *RegistryClient) Unregister(serviceName, nodeAddress string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.client.Unregister(ctx, &service.UnregisterRequest{ServiceName: serviceName, NodeAddress: nodeAddress})
	return err
}

func (c *RegistryClient) Discover(serviceName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.client.Discover(ctx, &service.DiscoverRequest{ServiceName: serviceName})
	if err != nil {
		return nil, err
	}
	return res.NodeAddresses, nil
}

func (c *RegistryClient) Close() {
	if err := c.conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}
