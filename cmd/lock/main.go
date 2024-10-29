// client/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/sascha-andres/lockutil/internal/lockserver" // Adjust the import path based on your project structure

	"google.golang.org/grpc"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewLockServiceClient(conn)

	// Define lock request parameters
	lockName := "my_lock"
	timeoutSeconds := int32(5) // Wait up to 5 seconds to acquire the lock
	pid := int32(os.Getpid())  // Use the current process PID

	// Create a context with timeout to handle request deadlines
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds+2)*time.Second)
	defer cancel()

	// Request a lock
	req := &pb.LockRequest{
		LockName:       lockName,
		TimeoutSeconds: timeoutSeconds,
		Pid:            pid,
	}

	resp, err := client.RequestLock(ctx, req)
	if err != nil {
		log.Fatalf("Error while requesting lock: %v", err)
	}

	// Check if lock acquisition was successful
	if resp.Success {
		fmt.Printf("Successfully acquired lock: %s\n", lockName)
	} else {
		fmt.Printf("Failed to acquire lock: %s - %s\n", lockName, resp.Message)
	}

	// Optionally, release the lock after some work (simulate by sleeping)
	time.Sleep(2 * time.Second) // Simulate doing some work while holding the lock

	releaseResp, err := client.ReleaseLock(ctx, &pb.ReleaseRequest{LockName: lockName, Pid: pid})
	if err != nil {
		log.Fatalf("Error while releasing lock: %v", err)
	}
	if releaseResp.Success {
		fmt.Printf("Successfully released lock: %s\n", lockName)
	} else {
		fmt.Printf("Failed to release lock: %s - %s\n", lockName, releaseResp.Message)
	}
}
