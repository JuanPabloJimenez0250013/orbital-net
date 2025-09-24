package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/JuanPabloJimenez0250013/orbital-net/sim-service/handler/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.NewClient(
		"localhost:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	// Create client
	client := pb.NewSimServiceClient(conn)

	// Call RPC
	resp, err := client.GetNodes(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatal("Error calling GetNodes:", err)
	}

	for _, n := range resp.Nodes {
		fmt.Printf("Node %s at (%.2f, %.2f)\n", n.Name, n.XPos, n.YPos)
	}

}
