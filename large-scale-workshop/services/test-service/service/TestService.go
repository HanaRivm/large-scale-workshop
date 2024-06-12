package TestService

import (
	services "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/common"
	. "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/common"

	"google.golang.org/grpc"
)

type testServiceImplementation struct {
	UnimplementedTestServiceServer
}

func Start(configData []byte) error {
	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterTestServiceServer(s, &testServiceImplementation{})
	}
	services.Start("TestService", 50051, bindgRPCToService)
}
