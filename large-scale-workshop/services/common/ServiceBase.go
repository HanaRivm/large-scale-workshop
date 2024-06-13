package common

import (
	"fmt"
	"net"

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
