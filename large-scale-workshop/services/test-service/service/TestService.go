package TestService

import (
	"context"
	"fmt"
	"log"

	services "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services"
	. "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/common"
	TestServiceServant "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/servant"
	. "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type testServiceImplementation struct {
	UnimplementedTestServiceServer
}

func (obj *testServiceImplementation) HelloWorld(ctxt context.Context, _ *emptypb.Empty) (res *wrapperspb.StringValue, err error) {
	return wrapperspb.String(TestServiceServant.HelloWorld()), nil
}

func (obj *testServiceImplementation) HelloToUser(ctx context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	Logger.Printf("HelloToUser called")
	response := fmt.Sprintf("Hello %s", req.Value)
	return wrapperspb.String(response), nil
}
func (obj *testServiceImplementation) Store(ctx context.Context, req *StoreKeyValue) (*emptypb.Empty, error) {
	log.Printf("Store called with key: %s, value: %s", req.Key, req.Value)
	TestServiceServant.Store(req.Key, req.Value)
	return &emptypb.Empty{}, nil
}

func (obj *testServiceImplementation) Get(ctx context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	log.Println("Get called with key:", req.GetValue())
	value, exists := TestServiceServant.Get(req.Value)
	if !exists {
		return nil, fmt.Errorf("key not found: %s", req.Value)
	}
	return wrapperspb.String(value), nil
}
func Start(configData []byte) error {
	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterTestServiceServer(s, &testServiceImplementation{})
	}
	services.Start("TestService", 50051, bindgRPCToService)
	return nil
}
