package common

import (
	"fmt"
	"net"

	RegistryServiceClient "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/client"
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
func registerAddress(serviceName string, registryAddresses []string,
	listeningAddress string) (unregister func()) {
	registryClient := RegistryServiceClient.NewRegistryServiceClient(registryAddresses)
	err := registryClient.Register(serviceName, listeningAddress)
	if err != nil {
		Logger.Fatalf("Failed to register to registry service: %v", err)
	}
	return func() {
		registryClient.Unregister(serviceName, listeningAddress)
	}
}
