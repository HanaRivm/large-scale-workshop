package RegistryService

import (
	"context"
	"sync"

	pb "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/common"
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
	s.services[req.ServiceName] = append(s.services[req.ServiceName], req.NodeAddress)
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

func (s *registryServer) IsAlive(ctx context.Context, req *pb.IsAliveRequest) (*pb.IsAliveResponse, error) {
	// This is a placeholder for the actual health check logic
	return &pb.IsAliveResponse{Alive: true}, nil
}
