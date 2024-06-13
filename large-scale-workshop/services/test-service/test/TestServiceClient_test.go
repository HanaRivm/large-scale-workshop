package TestService

import (
	"testing"

	client "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/client"
)

func TestHelloWorld(t *testing.T) {
	c := client.NewTestServiceClient("localhost:50051")
	r, err := c.HelloWorld()
	if err != nil {
		t.Fatalf("could not call HelloWorld: %v", err)
		return
	}
	t.Logf("Response: %v", r)
}

func TestHelloToUser(t *testing.T) {
	username := "Alice"
	c := client.NewTestServiceClient("localhost:50051")
	r, err := c.HelloToUser(username)
	if err != nil {
		t.Fatalf("could not call HelloToUser: %v", err)
	}
	expected := "Hello Alice"
	if r != expected {
		t.Errorf("unexpected response: got %s, want %s", r, expected)
	}
	t.Logf("Response: %v", r)
}
func TestStoreAndGet(t *testing.T) {
	c := client.NewTestServiceClient("localhost:50051")

	// Test Store
	err := c.Store("key1", "value1")
	if err != nil {
		t.Fatalf("could not call Store: %v", err)
	}

	// Test Get
	r, err := c.Get("key1")
	if err != nil {
		t.Fatalf("could not call Get: %v", err)
	}
	expected := "value1"
	if r != expected {
		t.Errorf("unexpected response: got %s, want %s", r, expected)
	}
	t.Logf("Response: %v", r)
}

func TestWaitAndRand(t *testing.T) {
	c := client.NewTestServiceClient("localhost:50051")
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
}
func TestIsAlive(t *testing.T) {
	c := client.NewTestServiceClient("localhost:50051")
	res, err := c.IsAlive()
	if err != nil {
		t.Fatalf("Calling IsAlive failed: %v", err)
		return
	}
	if !res {
		t.Fatalf("IsAlive returned false, expected true")
	}
	t.Logf("IsAlive returned: %v\n", res)
}
