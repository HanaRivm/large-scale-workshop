package test

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/common"
	service "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/service"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterRegistryServiceServer(s, service.NewRegistryServer())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestRegistryService(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewRegistryServiceClient(conn)

	t.Run("Register", func(t *testing.T) {
		_, err := client.Register(ctx, &pb.RegisterRequest{ServiceName: "testService", NodeAddress: "127.0.0.1:5000"})
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
	})

	t.Run("Discover", func(t *testing.T) {
		res, err := client.Discover(ctx, &pb.DiscoverRequest{ServiceName: "testService"})
		if err != nil {
			t.Fatalf("Discover failed: %v", err)
		}
		if len(res.NodeAddresses) == 0 {
			t.Fatalf("No nodes discovered")
		}
	})

	t.Run("Unregister", func(t *testing.T) {
		_, err := client.Unregister(ctx, &pb.UnregisterRequest{ServiceName: "testService", NodeAddress: "127.0.0.1:5000"})
		if err != nil {
			t.Fatalf("Unregister failed: %v", err)
		}
	})
}

func TestServiceRegistration(t *testing.T) {
	registry := service.NewRegistryServer()
	// Register a test service
	registry.Register(context.Background(), &pb.RegisterRequest{ServiceName: "TestService", NodeAddress: "127.0.0.1:5000"})

	// Discover the test service
	res, err := registry.Discover(context.Background(), &pb.DiscoverRequest{ServiceName: "TestService"})
	if err != nil {
		t.Fatalf("Failed to discover service: %v", err)
	}

	if len(res.NodeAddresses) != 1 || res.NodeAddresses[0] != "127.0.0.1:5000" {
		t.Fatalf("Expected node address 127.0.0.1:5000, got %v", res.NodeAddresses)
	}
}
