package TestServiceClient

import (
	context "context"
	"fmt"

	services "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/common"
	service "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/common"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func (obj *TestServiceClient) HelloToUser(username string) (string, error) {
	c, closeFunc, err := obj.Connect()
	defer closeFunc()
	if err != nil {
		return "", fmt.Errorf("failed to connect: %v", err)
	}
	req := &wrapperspb.StringValue{Value: username}
	res, err := c.HelloToUser(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("could not call HelloToUser: %v", err)
	}
	return res.Value, nil
}

func (obj *TestServiceClient) Store(key, value string) error {
	c, closeFunc, err := obj.Connect()
	defer closeFunc()
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	req := &service.StoreKeyValue{
		Key:   key,
		Value: value,
	}
	_, err = c.Store(context.Background(), req)
	if err != nil {
		return fmt.Errorf("could not call Store: %v", err)
	}
	return nil
}
func (obj *TestServiceClient) Get(key string) (string, error) {
	c, closeFunc, err := obj.Connect()
	defer closeFunc()
	if err != nil {
		return "", fmt.Errorf("failed to connect: %v", err)
	}
	req := &wrapperspb.StringValue{Value: key}
	res, err := c.Get(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("could not call Get: %v", err)
	}
	return res.Value, nil
}
