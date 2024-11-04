package main

import (
	"fmt"
	"log"
	"strings"

	pb "github.com/sascha-andres/lockutil/internal/lockserver"

	"net"

	"github.com/sascha-andres/lockutil/server"
	"github.com/sascha-andres/reuse/flag"
	"google.golang.org/grpc"
)

const (
	defaultPort     = ":50051"
	defaultHost     = "localhost"
	applicationName = "lockd"
)

var (
	port    string
	host    string
	help    bool
	verbose bool
)

// init initializes the logger settings, environment, and command-line flags for the application.
func init() {
	log.SetPrefix(fmt.Sprintf("[%s] ", strings.ToUpper(applicationName)))
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	flag.SetEnvPrefix(strings.ToUpper(applicationName))

	flag.StringVar(&port, "port", defaultPort, "The port to listen on")
	flag.StringVar(&host, "host", defaultHost, "The host to listen on")
	flag.BoolVar(&help, "help", false, "Prints this help message")
	flag.BoolVar(&verbose, "verbose", false, "Enables verbose logging")
}

// main is the entry point of the program, handling command-line flag parsing and executing the main functionality.
func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	if err := run(); err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}

// run starts a gRPC server on the default host and port, registers the LockService, and begins serving client requests.
func run() error {
	// Set up a listener on port 50051
	lis, err := net.Listen("tcp", fmt.Sprintf("%s%s", defaultHost, defaultPort))
	if err != nil {
		return err
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the lock service
	pb.RegisterLockServiceServer(grpcServer, server.NewLockServer(verbose))

	log.Printf("gRPC server running on port %q:%q...", host, port)
	return grpcServer.Serve(lis)
}
