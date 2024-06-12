package TestServiceClient

import (
	context "context"
	"fmt"

	services "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/common"
	service "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/common"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TestServiceClient struct {
	services.ServiceClientBase[service.TestServiceClient]
}

func NewTestServiceClient(address string) *TestServiceClient {
	return &TestServiceClient{

		ServiceClientBase: services.ServiceClientBase[service.TestServiceClient]{
			Address:      address,
			CreateClient: service.NewTestServiceClient,
		},
	}
}
func (obj *TestServiceClient) HelloWorld() (string, error) {
	c, closeFunc, err := obj.Connect()
	defer closeFunc()
	// Call the HelloWorld RPC function
	r, err := c.HelloWorld(context.Background(), &emptypb.Empty{})
	if err != nil {
		return "", fmt.Errorf("could not call HelloWorld: %v", err)
	}
	return r.Value, nil
}
