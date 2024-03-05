package main

import (
	"context"
	"log"
	pb "master_agent/helloworld"
	"os"
	"time"

	"github.com/shirou/gopsutil/mem" // Import package for fetching RAM usage
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051" // Address of the master
	defaultName = "agent-1"
)

type agent struct {
	name string
	id   string
}

func (a *agent) SendRAMUsage(ctx context.Context, client pb.MasterAgentClient) {
	for {
		// Fetch RAM usage information
		vm, err := mem.VirtualMemory()
		if err != nil {
			log.Fatalf("Failed to fetch RAM usage: %v", err)
		}
		ramUsage := float32(vm.UsedPercent)

		// Send RAM usage to master
		_, err = client.SendRAMUsage(ctx, &pb.SystemInfo{AgentId: a.id, RamUsage: ramUsage})
		if err != nil {
			log.Fatalf("Could not send RAM usage to master: %v", err)
		}

		// Wait for 2 seconds before sending the next RAM usage
		time.Sleep(2 * time.Second)
	}
}

func main() {
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	defer conn.Close()

	client := pb.NewMasterAgentClient(conn)

	agent := &agent{
		name: name,
		id:   "unique-agent-id", // You need to generate a unique ID for each agent
	}

	ctx := context.Background()

	// Register agent with master
	_, err = client.RegisterAgent(ctx, &pb.AgentInfo{Name: agent.name, Id: agent.id})
	if err != nil {
		log.Fatalf("Could not register agent: %v", err)
	}

	// Send RAM usage to master every 2 seconds
	agent.SendRAMUsage(ctx, client)
}
