package common

import (
	"fmt"
	"log"
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
func Start(serviceName string, grpcListenPort int, bindgRPCToService func(s grpc.ServiceRegistrar)) {
	listeningAddress, grpcServer, startListening := startgRPC(grpcListenPort)
	Logger.Printf("Starting %s gRPC server on %s", serviceName, listeningAddress)
	bindgRPCToService(grpcServer)
	startListening()
}

func RegisterAddress(serviceName string, registryAddresses []string, listeningAddress string) (unregister func()) {
	registryClient := RegistryServiceClient.NewRegistryServiceClient(registryAddresses)
	log.Printf("Registering service: %s at address: %s with registry: %v", serviceName, listeningAddress, registryAddresses)

	err := registryClient.Register(serviceName, listeningAddress)
	if err != nil {
		Logger.Fatalf("Failed to register to registry service: %v", err)
	}
	log.Printf("Successfully registered service: %s at address: %s", serviceName, listeningAddress)

	return func() {
		registryClient.Unregister(serviceName, listeningAddress)
	}
}
