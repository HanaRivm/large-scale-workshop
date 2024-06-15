package RegistryService

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type registryServer struct {
	pb.UnimplementedRegistryServiceServer
	mu       sync.Mutex
	services map[string][]string
}

func NewRegistryServer() *registryServer {
	return &registryServer{
		services: make(map[string][]string),
	}
}

func (s *registryServer) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Printf("Attempting to register service: %s at address: %s", req.ServiceName, req.NodeAddress)
	s.services[req.ServiceName] = append(s.services[req.ServiceName], req.NodeAddress)
	log.Printf("Successfully registered service: %s at address: %s", req.ServiceName, req.NodeAddress)
	return &emptypb.Empty{}, nil
}

func (s *registryServer) Unregister(ctx context.Context, req *pb.UnregisterRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	nodes := s.services[req.ServiceName]
	for i, addr := range nodes {
		if addr == req.NodeAddress {
			s.services[req.ServiceName] = append(nodes[:i], nodes[i+1:]...)
			break
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *registryServer) Discover(ctx context.Context, req *pb.DiscoverRequest) (*pb.DiscoverResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	nodes := s.services[req.ServiceName]
	return &pb.DiscoverResponse{NodeAddresses: nodes}, nil
}

func (s *registryServer) startHealthCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for _, nodes := range s.services {
			for _, node := range nodes {
				conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithBlock())
				if err != nil {
					continue
				}
				client := pb.NewRegistryServiceClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				res, err := client.IsAlive(ctx, &emptypb.Empty{})
				if err != nil {
					log.Printf("error is not nil")
				}
				log.Printf(res.String())
				//|| !res.Value {
				//
				//	log.Printf("Node %s is not responding, removing from registry", node)
				//	s.services[serviceName] = removeNode(s.services[serviceName], node)
				//}
				conn.Close()
			}
		}
		s.mu.Unlock()
	}
}
func removeNode(nodes []string, node string) []string {
	for i, n := range nodes {
		if n == node {
			return append(nodes[:i], nodes[i+1:]...)
		}
	}
	return nodes
}
func RunServer() {
	lis, err := net.Listen("tcp", ":8502")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	server := NewRegistryServer()
	go server.startHealthCheck()
	pb.RegisterRegistryServiceServer(s, server)
	reflection.Register(s)

	log.Println("Registry service is running on port 8502")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
