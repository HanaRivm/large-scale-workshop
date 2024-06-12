package TestService

import (
	"context"
	"fmt"

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
	username := req.GetValue()
	message := fmt.Sprintf("Hello %s", username)
	return wrapperspb.String(message), nil
}

func Start(configData []byte) error {
	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterTestServiceServer(s, &testServiceImplementation{})
	}
	services.Start("TestService", 50051, bindgRPCToService)
	return nil
}
