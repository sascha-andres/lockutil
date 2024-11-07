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
	defaultPort = "50051"

	// defaultHost specifies the default hostname for the server connection.
	defaultHost = "localhost"

	// applicationName specifies the name of the application used for logging and configuration purposes.
	applicationName = "lock"

	// defaultLockJame specifies the default name used for locking mechanisms
	defaultLockJame = "default"

	defaultTimeout = 0
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

	// opList represents aan operation to list all currently existing locks
	opList

	// opForceRelease indicates an operation that forcibly releases resources or locks, without checking the current state.
	opForceRelease
)

var (
	lockName   string
	port       string
	host       string
	forceToken string
	help       bool
	verbose    bool
	timeout    int
)

// init initializes the logger settings, environment, and command-line flags for the application.
func init() {
	log.SetPrefix(fmt.Sprintf("[%s] ", strings.ToUpper(applicationName)))
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	flag.SetEnvPrefix(strings.ToUpper(applicationName))
	flag.StringVar(&port, "port", defaultPort, "The port to connect to")
	flag.StringVar(&host, "host", defaultHost, "The host to connect to")
	flag.StringVar(&lockName, "lock", defaultLockJame, "The name of the lock to acquire")
	flag.StringVar(&forceToken, "force-token", "", "The force token to use for force release")
	flag.IntVar(&timeout, "timeout", defaultTimeout, "The timeout in seconds for the lock")
	flag.BoolVar(&help, "help", false, "Prints this help message")
	flag.BoolVar(&verbose, "verbose", false, "Enables verbose logging")
}

// main is the entry point of the application which parses command-line flags and determines the operation to execute.
// It either runs the acquire or release operation based on the provided flags and handles any errors encountered.
func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	ot := operationType(0)
	if verbose {
		log.Printf("Flags: %v", flag.GetVerbs())
	}
	if len(flag.GetVerbs()) == 0 {
		ot = opAcquire
	} else {
		if flag.GetVerbs()[0] == "release" {
			ot = opRelease
		}
		if flag.GetVerbs()[0] == "list" {
			ot = opList
		}
		if flag.GetVerbs()[0] == "force-release" {
			ot = opForceRelease
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

	if verbose {
		otString := "ERR"
		if ot == opRelease {
			otString = "release"
		}
		if ot == opAcquire {
			otString = "acquire"
		}
		if ot == opList {
			otString = "list"
		}
		log.Printf("Running operation: %s", otString)
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

	if ot == opRelease || ot == opForceRelease {
		return release(client, ot == opForceRelease)
	}

	if ot == opList {
		return list(client)
	}

	return errors.New("no supported operation")
}

// list retrieves and prints a list of locks from the LockServiceClient.
func list(client pb.LockServiceClient) error {
	locks, err := client.List(context.Background(), &pb.ListRequest{})
	if err != nil {
		return err
	}
	for _, lock := range locks.GetLocks() {
		fmt.Printf("%s: from pid %d on %s is locked: %t\n", lock.GetName(), lock.GetPid(), lock.GetAddr(), lock.GetLocked())
	}
	return nil
}

// release attempts to release a lock held by the current process using the provided LockServiceClient.
func release(client pb.LockServiceClient, force bool) error {
	if force && forceToken == "" {
		return errors.New("force token is required")
	}
	name, _, pid := getLockParameters()
	if verbose {
		log.Printf("Releasing lock: %s, pid: %d", name, pid)
	}
	releaseResp, err := client.ReleaseLock(context.Background(), &pb.ReleaseRequest{LockName: name, Pid: pid})
	if err != nil {
		return err
	}
	if !releaseResp.Success {
		fmt.Printf("Failed to release lock: %s - %s\n", name, releaseResp.Message)
	}
	return nil
}

// acquire attempts to obtain a lock by sending a request to the LockServiceClient.
// It uses predefined lock parameters from getLockParameters() for lock name, timeout, and process ID.
// If the lock is acquired successfully, the function will return nil. If not, an error or a failure message is printed.
func acquire(client pb.LockServiceClient) error {
	lockName, timeoutSeconds, pid := getLockParameters()
	if verbose {
		log.Printf("Acquiring lock: %s, timeout: %d, pid: %d", lockName, timeoutSeconds, pid)
	}
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
	timeoutSeconds := int32(timeout) // TODO timeout default and flag
	pid := int32(os.Getppid())
	return lockName, timeoutSeconds, pid
}
