package main

import (
	"fmt"

	pb "github.com/sascha-andres/lockutil/internal/lockserver"

	"log"
	"net"

	"github.com/sascha-andres/lockutil/server"
	"google.golang.org/grpc"
)

func main() {
	// Set up a listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the lock service
	pb.RegisterLockServiceServer(grpcServer, server.NewLockServer())

	fmt.Println("gRPC server running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
