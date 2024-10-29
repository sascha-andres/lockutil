package server

import (
	"context"
	"github.com/sascha-andres/lockutil/internal/lockmanager"

	pb "github.com/sascha-andres/lockutil/internal/lockserver" // Import the generated proto package
)

type lockServer struct {
	pb.UnimplementedLockServiceServer
	manager *lockmanager.LockManager
}

// NewLockServer initializes a new lockServer
func NewLockServer() *lockServer {
	return &lockServer{
		manager: lockmanager.NewLockManager(),
	}
}

// RequestLock handles lock requests from clients
func (s *lockServer) RequestLock(ctx context.Context, req *pb.LockRequest) (*pb.LockResponse, error) {
	err := s.manager.RequestLock(req.LockName, req.Pid, req.TimeoutSeconds)
	if err != nil {
		return &pb.LockResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.LockResponse{Success: true, Message: "Lock acquired"}, nil
}

// ReleaseLock handles lock release requests from clients
func (s *lockServer) ReleaseLock(ctx context.Context, req *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {
	err := s.manager.ReleaseLock(req.LockName, req.Pid)
	if err != nil {
		return &pb.ReleaseResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.ReleaseResponse{Success: true, Message: "Lock released"}, nil
}
