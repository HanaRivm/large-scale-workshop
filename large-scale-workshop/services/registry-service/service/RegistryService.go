package RegistryService

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/common"
	"github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/servant/dht"
	"gopkg.in/yaml.v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type registryServer struct {
	pb.UnimplementedRegistryServiceServer
	mu       sync.Mutex
	services map[string][]string
	chord    *dht.Chord
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

func (s *registryServer) IsAlive(ctx context.Context, req *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &wrapperspb.BoolValue{Value: true}, nil
}

func (s *registryServer) startHealthCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		// Copy the map to avoid concurrent access issues
		nodes := make(map[string][]string, len(s.services))
		for k, v := range s.services {
			nodes[k] = append([]string{}, v...) // Deep copy of the value slice
		}
		s.mu.Unlock()
		// Now safely iterate over the copied map
		for serviceName, nodeList := range nodes {
			for _, node := range nodeList {
				conn, err := grpc.Dial(node, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Printf("Error dialing node %s: %v", node, err)
					continue
				}
				defer conn.Close()
				client := pb.NewRegistryServiceClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				res, err := client.IsAlive(ctx, &emptypb.Empty{})
				if err != nil || !res.Value {
					log.Printf("Node %s is not responding, removing from registry", node)
					// Consider using a separate goroutine for removing unhealthy nodes
					// to avoid blocking the health check loop
					s.mu.Lock()
					s.services[serviceName] = removeNode(s.services[serviceName], node)
					s.mu.Unlock()
				}
			}
		}
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

func RunServer(configData []byte) error {
	var config struct {
		Port     int32  `yaml:"port"`
		Name     string `yaml:"name"`
		RootName string `yaml:"rootName"`
	}

	err := yaml.Unmarshal(configData, &config)
	if err != nil {
		return fmt.Errorf("error unmarshaling service config: %v", err)
	}

	port := config.Port
	var lis net.Listener
	for {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Printf("Port %d in use, trying next port", port)
			port++
		} else {
			break
		}
	}

	defer lis.Close()
	var rootName = config.RootName
	var chordPort int32
	chordPort = 1099
	server := NewRegistryServer()
	if rootName == "" {
		newchord, err := dht.NewChord(config.Name, chordPort)
		server.chord = newchord
		log.Println("created chord")
		if err != nil {
			return fmt.Errorf("failed to create Chord: %v", err)
		}

	} else {
		server.chord, err = dht.JoinChord(config.Name, rootName, chordPort)
		if err != nil {
			return fmt.Errorf("failed to join Chord: %v", err)
		}
	}

	//go server.startHealthCheck()

	s := grpc.NewServer()
	pb.RegisterRegistryServiceServer(s, server)
	reflection.Register(s)

	log.Printf("Registry service is running on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	return nil
}
