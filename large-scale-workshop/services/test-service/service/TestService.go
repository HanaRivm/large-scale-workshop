package TestService

import (
	"context"
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

func (obj *testServiceImplementation) HelloToUser(ctxt context.Context, req *wrapperspb.StringValue) (res *wrapperspb.StringValue, err error) {
	username := req.GetValue()
	response := TestServiceServant.HelloToUser(username)
	log.Println("HelloToUser response:", response)
	return wrapperspb.String(TestServiceServant.HelloToUser(username)), nil
}

func Start(configData []byte) error {
	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterTestServiceServer(s, &testServiceImplementation{})
	}
	services.Start("TestService", 50051, bindgRPCToService)
	return nil
}
