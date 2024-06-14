package common

import (
	"context"
	"fmt"
	"net"

	pb "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/common"

	. "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/utils"
	"google.golang.org/grpc"
)

func startgRPC(listenPort int) (listeningAddress string, grpcServer *grpc.Server,
	startListening func()) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		Logger.Fatalf("failed to listen: %v", err)
	}
	listeningAddress = lis.Addr().String()
	grpcServer = grpc.NewServer()
	startListening = func() {
		if err := grpcServer.Serve(lis); err != nil {
			Logger.Fatalf("failed to serve: %v", err)
		}
	}
	return
}
func registerAddress(serviceName string, registryAddresses []string, listeningAddress string) (unregister func()) {
	registryClient := pb.NewRegistryServiceClient(&ServiceClientBase(registryAddresses))
	_, err := registryClient.Register(context.Background(), &pb.RegisterRequest{
		ServiceName: serviceName,
		NodeAddress: listeningAddress,
	})
	if err != nil {
		Logger.Fatalf("Failed to register to registry service: %v", err)
	}
	return func() {
		_, err := registryClient.Unregister(context.Background(), &pb.UnregisterRequest{
			ServiceName: serviceName,
			NodeAddress: listeningAddress,
		})
		if err != nil {
			Logger.Printf("Failed to unregister from registry service: %v", err)
		}
	}
}
func Start(serviceName string, grpcListenPort int, bindgRPCToService func(s grpc.ServiceRegistrar)) {
	listeningAddress, grpcServer, startListening := startgRPC(grpcListenPort)
	Logger.Printf("Starting %s gRPC server on %s", serviceName, listeningAddress)
	bindgRPCToService(grpcServer)
	startListening()
}
