package TestService

import (
	"context"

	services "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/common"
	. "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/common"
	. "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/servant"
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

func Start(configData []byte) error {
	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterTestServiceServer(s, &testServiceImplementation{})
	}
	services.Start("TestService", 50051, bindgRPCToService)
	return nil
}
