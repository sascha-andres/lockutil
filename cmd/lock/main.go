// client/main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sascha-andres/reuse/flag"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sascha-andres/lockutil/internal/lockserver" // Adjust the import path based on your project structure

	"google.golang.org/grpc"
)

const (

	// defaultPort defines the default port used for the server connection.
	defaultPort = ":50051"

	// defaultHost specifies the default hostname for the server connection.
	defaultHost = "localhost"

	// applicationName specifies the name of the application used for logging and configuration purposes.
	applicationName = "lock"

	// defaultLockJame specifies the default name used for locking mechanisms
	defaultLockJame = "default"
)

// operationType represents different types of operations within the system.
type operationType int

const (

	// opNone represents an operation that performs no action or is uninitialized.
	opNone operationType = iota

	// opAcquire represents an operation that acquires resources or locks within the system.
	opAcquire

	// opRelease represents an operation that releases resources or locks within the system.
	opRelease
)

var (
	lockName string
	port     string
	host     string
	help     bool
	verbose  bool
)

// init initializes the logger settings, environment, and command-line flags for the application.
func init() {
	log.SetPrefix(fmt.Sprintf("[%s] ", strings.ToUpper(applicationName)))
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	flag.SetEnvPrefix(strings.ToUpper(applicationName))
	flag.SetSeparated()
	flag.StringVar(&port, "port", defaultPort, "The port to connect to")
	flag.StringVar(&host, "host", defaultHost, "The host to connect to")
	flag.StringVar(&lockName, "lock", defaultLockJame, "The name of the lock to acquire")
	flag.BoolVar(&help, "help", false, "Prints this help message")
	flag.BoolVar(&verbose, "verbose", false, "Enables verbose logging")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	ot := operationType(0)
	if len(flag.GetSeparated()) == 0 {
		ot = opAcquire
	} else {
		if flag.GetSeparated()[0] == "release" {
			ot = opRelease
		}
	}

	if err := run(ot); err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}

// run executes the operation specified by the operationType.
// It either acquires or releases a lock by communicating with a gRPC LockServiceClient.
func run(ot operationType) error {
	if ot == opNone {
		log.Println("Please specify no operation to lock or 'release' to release a lock")
		return errors.New("no supported operation")
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()

	client := pb.NewLockServiceClient(conn)
	if ot == opAcquire {
		return acquire(client)
	}

	if ot == opRelease {
		return release(client)
	}

	return errors.New("no supported operation")
}

// release attempts to release a lock held by the current process using the provided LockServiceClient.
func release(client pb.LockServiceClient) error {
	lockName, _, pid := getLockParameters()
	releaseResp, err := client.ReleaseLock(context.Background(), &pb.ReleaseRequest{LockName: lockName, Pid: pid})
	if err != nil {
		return err
	}
	if !releaseResp.Success {
		fmt.Printf("Failed to release lock: %s - %s\n", lockName, releaseResp.Message)
	}
	return nil
}

// acquire attempts to obtain a lock by sending a request to the LockServiceClient.
// It uses predefined lock parameters from getLockParameters() for lock name, timeout, and process ID.
// If the lock is acquired successfully, the function will return nil. If not, an error or a failure message is printed.
func acquire(client pb.LockServiceClient) error {
	lockName, timeoutSeconds, pid := getLockParameters()
	req := &pb.LockRequest{
		LockName:       lockName,
		TimeoutSeconds: timeoutSeconds,
		Pid:            pid,
	}
	resp, err := client.RequestLock(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Success {
		fmt.Printf("Failed to acquire lock: %s - %s\n", lockName, resp.Message)
	}
	return nil
}

// getLockParameters returns the default lock parameters including lock name, timeout in seconds, and process ID.
func getLockParameters() (string, int32, int32) {
	// Define lock request parameters
	lockName := lockName
	timeoutSeconds := int32(5) // TODO timeout default and flag
	pid := int32(os.Getppid())
	return lockName, timeoutSeconds, pid
}
