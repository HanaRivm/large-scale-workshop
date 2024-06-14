package service

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type registryServer struct {
	pb.UnimplementedRegistryServiceServer
	services map[string][]string
	mu       sync.Mutex
}

func NewRegistryServer() *registryServer {
	return &registryServer{
		services: make(map[string][]string),
	}
}

func (s *registryServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.services[req.ServiceName] = append(s.services[req.ServiceName], req.NodeAddress)
	return &pb.RegisterResponse{Success: true}, nil
}

func (s *registryServer) Unregister(ctx context.Context, req *pb.UnregisterRequest) (*pb.UnregisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	nodes := s.services[req.ServiceName]
	for i, addr := range nodes {
		if addr == req.NodeAddress {
			s.services[req.ServiceName] = append(nodes[:i], nodes[i+1:]...)
			break
		}
	}
	return &pb.UnregisterResponse{Success: true}, nil
}

func (s *registryServer) Discover(ctx context.Context, req *pb.DiscoverRequest) (*pb.DiscoverResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	nodes := s.services[req.ServiceName]
	return &pb.DiscoverResponse{NodeAddresses: nodes}, nil
}

func (s *registryServer) IsAlive(ctx context.Context, req *pb.IsAliveRequest) (*pb.IsAliveResponse, error) {
	// This is a placeholder for the actual health check logic
	return &pb.IsAliveResponse{Alive: true}, nil
}

func (s *registryServer) startHealthCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for serviceName, nodes := range s.services {
			for _, node := range nodes {
				conn, err := grpc.Dial(node, grpc.WithInsecure())
				if err != nil {
					continue
				}
				client := pb.NewRegistryServiceClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				res, err := client.IsAlive(ctx, &pb.IsAliveRequest{NodeAddress: node})
				if err != nil || !res.Alive {
					log.Printf("Node %s is not responding, removing from registry", node)
					s.services[serviceName] = removeNode(s.services[serviceName], node)
				}
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
