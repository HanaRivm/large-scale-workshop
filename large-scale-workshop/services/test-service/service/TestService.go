package TestService

import (
	"context"
	"fmt"
	"log"
	"net"

	services "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/common"
	. "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/common"
	TestServiceServant "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/servant"
	. "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v2"
)

type testServiceImplementation struct {
	UnimplementedTestServiceServer
}

func (obj *testServiceImplementation) HelloWorld(ctx context.Context, _ *emptypb.Empty) (res *wrapperspb.StringValue, err error) {
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

func (obj *testServiceImplementation) WaitAndRand(seconds *wrapperspb.Int32Value, streamRet TestService_WaitAndRandServer) error {
	Logger.Printf("WaitAndRand called")
	streamClient := func(x int32) error {
		return streamRet.Send(wrapperspb.Int32(x))
	}
	return TestServiceServant.WaitAndRand(seconds.Value, streamClient)
}

func (obj *testServiceImplementation) IsAlive(ctx context.Context, req *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	return &wrapperspb.BoolValue{Value: true}, nil
}

func (obj *testServiceImplementation) ExtractLinksFromURL(ctx context.Context, req *ExtractLinksFromURLParameters) (*ExtractLinksFromURLReturnedValue, error) {
	links, err := TestServiceServant.ExtractLinksFromURL(req.Url, req.Depth)
	if err != nil {
		return nil, err
	}
	return &ExtractLinksFromURLReturnedValue{Links: links}, nil
}

func Start(configData []byte) error {
	var config struct {
		RegistryAddresses []string `yaml:"registry_addresses"`
	}
	err := yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("error unmarshaling service config: %v", err)
	}

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	port := lis.Addr().(*net.TCPAddr).Port
	listeningAddress := fmt.Sprintf("127.0.0.1:%d", port)
	log.Printf("TestService listening on %s", listeningAddress)

	unregister := services.RegisterAddress("TestService", config.RegistryAddresses, listeningAddress)
	defer unregister()

	s := grpc.NewServer()
	RegisterTestServiceServer(s, &testServiceImplementation{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	return nil
}
