package TestService

import (
	client "[module name]/services/test-service/client"
	"testing"
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
