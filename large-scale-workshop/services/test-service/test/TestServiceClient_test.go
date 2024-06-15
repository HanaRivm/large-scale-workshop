package TestService

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/common"
	registryClient "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/common"
	client "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/client"
)

func getClient(registryAddress string) (*common.ServiceClientBase[registryClient.RegistryServiceClient], error) {
	clientBase := &common.ServiceClientBase[registryClient.RegistryServiceClient]{
		RegistryAddresses: []string{registryAddress},
		CreateClient:      registryClient.NewRegistryServiceClient,
	}
	_, _, err := clientBase.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to registry service: %v", err)
	}
	return clientBase, nil
}

func discoverNodes(clientBase *common.ServiceClientBase[registryClient.RegistryServiceClient], serviceName string) ([]string, error) {
	client_t, _, err := clientBase.Connect()
	res, err := client_t.Discover(context.Background(), &registryClient.DiscoverRequest{ServiceName: serviceName})
	if err != nil {
		return nil, fmt.Errorf("discover failed: %v", err)
	}
	return res.NodeAddresses, nil
}

// Simulate load balancing by randomly picking a node from the provided addresses.
func pickRandomNode(nodes []string) string {
	rand.Seed(time.Now().UnixNano())
	return nodes[rand.Intn(len(nodes))]
}

func TestTestServiceClient(t *testing.T) {
	registryAddresses := []string{
		"127.0.0.1:8502",
		"127.0.0.1:8503",
	}

	clientBase, err := getClient(registryAddresses[0])
	if err != nil {
		t.Fatalf("Failed to connect to registry service: %v", err)
	}

	nodes, err := discoverNodes(clientBase, "TestService")
	if err != nil {
		t.Fatalf("Failed to discover nodes: %v", err)
	}

	if len(nodes) < 2 {
		t.Fatalf("Expected at least 2 TestService nodes, got %d", len(nodes))
	}

	t.Run("HelloWorld", func(t *testing.T) {
		nodeAddress := pickRandomNode(nodes)
		c := client.NewTestServiceClient([]string{nodeAddress})
		r, err := c.HelloWorld()
		if err != nil {
			t.Fatalf("could not call HelloWorld: %v", err)
			return
		}
		t.Logf("Response: %v", r)
	})

	t.Run("HelloToUser", func(t *testing.T) {
		username := "Alice"
		nodeAddress := pickRandomNode(nodes)
		c := client.NewTestServiceClient([]string{nodeAddress})
		r, err := c.HelloToUser(username)
		if err != nil {
			t.Fatalf("could not call HelloToUser: %v", err)
		}
		expected := "Hello Alice"
		if r != expected {
			t.Errorf("unexpected response: got %s, want %s", r, expected)
		}
		t.Logf("Response: %v", r)
	})

	t.Run("StoreAndGet", func(t *testing.T) {
		nodeAddress := pickRandomNode(nodes)
		c := client.NewTestServiceClient([]string{nodeAddress})

		err := c.Store("key1", "value1")
		if err != nil {
			t.Fatalf("could not call Store: %v", err)
		}

		r, err := c.Get("key1")
		if err != nil {
			t.Fatalf("could not call Get: %v", err)
		}
		expected := "value1"
		if r != expected {
			t.Errorf("unexpected response: got %s, want %s", r, expected)
		}
		t.Logf("Response: %v", r)
	})

	t.Run("WaitAndRand", func(t *testing.T) {
		nodeAddress := pickRandomNode(nodes)
		c := client.NewTestServiceClient([]string{nodeAddress})
		resPromise, err := c.WaitAndRand(3)
		if err != nil {
			t.Fatalf("Calling WaitAndRand failed: %v", err)
			return
		}
		res, err := resPromise()
		if err != nil {
			t.Fatalf("WaitAndRand failed: %v", err)
			return
		}
		t.Logf("Returned random number: %v\n", res)
	})

	t.Run("IsAlive", func(t *testing.T) {
		nodeAddress := pickRandomNode(nodes)
		c := client.NewTestServiceClient([]string{nodeAddress})
		res, err := c.IsAlive()
		if err != nil {
			t.Fatalf("Calling IsAlive failed: %v", err)
			return
		}
		if !res {
			t.Fatalf("IsAlive returned false, expected true")
		}
		t.Logf("IsAlive returned: %v\n", res)
	})

	t.Run("ExtractLinksFromURL", func(t *testing.T) {
		nodeAddress := pickRandomNode(nodes)
		c := client.NewTestServiceClient([]string{nodeAddress})
		links, err := c.ExtractLinksFromURL("https://www.microsoft.com", 1)
		if err != nil {
			t.Fatalf("Calling ExtractLinksFromURL failed: %v", err)
			return
		}
		if len(links) == 0 {
			t.Fatalf("ExtractLinksFromURL returned no links")
		}
		t.Logf("Returned links: %v\n", links)
	})
}
