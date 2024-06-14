package test

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	client "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/client"
	service "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/service"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	client.RegisterRegistryServiceServer(s, service.NewRegistryServer())
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
	client := client.NewRegistryServiceClient(conn)

	t.Run("Register", func(t *testing.T) {
		_, err := client.Register(ctx, &client.RegisterRequest{ServiceName: "testService", NodeAddress: "127.0.0.1:5000"})
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
	})

	t.Run("Discover", func(t *testing.T) {
		res, err := client.Discover(ctx, &client.DiscoverRequest{ServiceName: "testService"})
		if err != nil {
			t.Fatalf("Discover failed: %v", err)
		}
		if len(res.NodeAddresses) == 0 {
			t.Fatalf("No nodes discovered")
		}
	})

	t.Run("Unregister", func(t *testing.T) {
		_, err := client.Unregister(ctx, &client.UnregisterRequest{ServiceName: "testService", NodeAddress: "127.0.0.1:5000"})
		if err != nil {
			t.Fatalf("Unregister failed: %v", err)
		}
	})

	t.Run("IsAlive", func(t *testing.T) {
		res, err := client.IsAlive(ctx, &client.IsAliveRequest{NodeAddress: "127.0.0.1:5000"})
		if err != nil {
			t.Fatalf("IsAlive failed: %v", err)
		}
		if !res.Alive {
			t.Fatalf("Node is not alive")
		}
	})
}
