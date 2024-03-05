package main

import (
	"context"
	"log"
	pb "master_agent/helloworld"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type masterAgentServer struct {
	pb.UnimplementedMasterAgentServer
}

func (s *masterAgentServer) RegisterAgent(ctx context.Context, in *pb.AgentInfo) (*pb.AgentInfo, error) {
	log.Printf("Received registration from agent %s with ID %s", in.Name, in.Id)
	return &pb.AgentInfo{Name: in.Name, Id: in.Id}, nil
}

func (s *masterAgentServer) SendRAMUsage(ctx context.Context, in *pb.SystemInfo) (*pb.SystemInfo, error) {
	log.Printf("Received RAM usage from agent %s with ID %s: %f", in.AgentId, in.AgentId, in.RamUsage)
	// You can further process or store the received RAM usage information here
	return in, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMasterAgentServer(s, &masterAgentServer{})
	log.Printf("Master listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
