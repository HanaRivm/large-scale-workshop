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
