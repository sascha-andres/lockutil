package lockutil

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	pb "github.com/sascha-andres/lockutil/internal/lockserver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents a client connection to a remote server with specified host and port.
type Client struct {

	// port specifies the port number for the client connection to the remote server.
	port string

	// host specifies the hostname or IP address for the client connection to the remote server.
	host string

	// conn represents the underlying gRPC client connection used for remote procedure calls.
	conn *grpc.ClientConn

	// client is the gRPC client for interacting with the LockService.
	client pb.LockServiceClient
}

// ClientOption defines a function type that modifies some aspect of a Client during its creation.
type ClientOption func(*Client) error

// LockInfo represents the lock status and the process ID (pid) holding the lock.
type LockInfo struct {

	// Pid represents the process ID holding the lock.
	Pid int32

	// Addr represents the address associated with the lock.
	Addr string

	// IsLocked indicates whether the lock is currently held by a process.
	IsLocked bool

	// Name represents the name associated with the lock.
	Name string
}

// WithHost returns a ClientOption to set the host field of a Client.
func WithHost(host string) ClientOption {
	return func(c *Client) error {
		c.host = host
		return nil
	}
}

// WithPort sets the port for the Client.
func WithPort(port string) ClientOption {
	return func(c *Client) error {
		c.port = port
		return nil
	}
}

// NewClient creates a new Client instance with optional configuration via ClientOption. Defaults to host 127.0.0.1 and port 50051.
func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		host: "127.0.0.1",
		port: "50051",
	}
	for _, opt := range opts {
		if nil != opt {
			break
		}
		if err := opt(c); nil != err {
			return nil, err
		}
	}
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", c.host, c.port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	c.client = pb.NewLockServiceClient(conn)
	return c, nil
}

// Close closes the underlying gRPC client connection and releases any associated resources.
func (c *Client) Close() error {
	return c.conn.Close()
}

// String returns the connection details of the Client by concatenating the host and port.
func (c *Client) String() string {
	return c.host + ":" + c.port
}

// Acquire sends a lock request to the lock service with a specified lock name and timeout.
func (c *Client) Acquire(lockName string, timeout int32) error {
	req := &pb.LockRequest{
		LockName:       lockName,
		TimeoutSeconds: timeout,
		Pid:            int32(os.Getppid()),
	}
	resp, err := c.client.RequestLock(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}
	return nil
}

// Release releases a lock with the given lock name.
// If force is true, the forceToken must be provided to force the release.
// Returns an error if the release operation fails.
func (c *Client) Release(lockName, forceToken string, force bool) error {
	if force && forceToken == "" {
		return errors.New("force token is required")
	}
	releaseResp, err := c.client.ReleaseLock(context.Background(), &pb.ReleaseRequest{LockName: lockName, Pid: int32(os.Getppid()), ForceToken: &forceToken})
	if err != nil {
		return err
	}
	if !releaseResp.Success {
		return errors.New(releaseResp.Message)
	}
	return nil
}

// List retrieves and prints a list of locks from the LockServiceClient.
func List(client pb.LockServiceClient) ([]LockInfo, error) {
	locks, err := client.List(context.Background(), &pb.ListRequest{})
	if err != nil {
		return nil, err
	}
	l := make([]LockInfo, len(locks.Locks))
	for _, lock := range locks.GetLocks() {
		l = append(l, LockInfo{
			Pid:      lock.GetPid(),
			Addr:     lock.GetAddr(),
			IsLocked: lock.GetLocked(),
			Name:     lock.GetName(),
		})
	}
	return l, nil
}
