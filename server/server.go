package server

import (
	"context"
	"log"

	"github.com/sascha-andres/lockutil/internal/lockmanager"

	pb "github.com/sascha-andres/lockutil/internal/lockserver" // Import the generated proto package
)

type LockServer struct {
	pb.UnimplementedLockServiceServer
	manager *lockmanager.LockManager
	verbose bool
}

// NewLockServer initializes a new LockServer
func NewLockServer(verbose bool) *LockServer {
	return &LockServer{
		verbose: verbose,
		manager: lockmanager.NewLockManager(verbose),
	}
}

// RequestLock handles lock requests from clients
func (s *LockServer) RequestLock(_ context.Context, req *pb.LockRequest) (*pb.LockResponse, error) {
	if s.verbose {
		log.Printf("RequestLock request for %s from %d with timeout %d", req.GetLockName(), req.GetPid(), req.GetTimeoutSeconds())
	}
	err := s.manager.RequestLock(req.LockName, req.Pid, req.TimeoutSeconds)
	if err != nil {
		log.Printf("RequestLock failed for %s from %d: %s", req.GetLockName(), req.GetPid(), err.Error())
		return &pb.LockResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.LockResponse{Success: true, Message: "Lock acquired"}, nil
}

// ReleaseLock handles lock release requests from clients
func (s *LockServer) ReleaseLock(_ context.Context, req *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {
	if s.verbose {
		log.Printf("ReleaseLock request for %s from %d", req.GetLockName(), req.GetPid())
	}
	err := s.manager.ReleaseLock(req.LockName, req.Pid)
	if err != nil {
		return &pb.ReleaseResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.ReleaseResponse{Success: true, Message: "Lock released"}, nil
}
