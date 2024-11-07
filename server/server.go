package server

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc/peer"

	"github.com/sascha-andres/lockutil/internal/lockmanager"

	pb "github.com/sascha-andres/lockutil/internal/lockserver" // Import the generated proto package
)

type LockServer struct {
	pb.UnimplementedLockServiceServer
	manager     *lockmanager.LockManager
	verbose     bool
	secretToken string
}

// NewLockServer initializes a new LockServer
func NewLockServer(token string, verbose bool) *LockServer {
	return &LockServer{
		verbose:     verbose,
		manager:     lockmanager.NewLockManager(verbose),
		secretToken: token,
	}
}

// RequestLock handles lock requests from clients
func (s *LockServer) RequestLock(ctx context.Context, req *pb.LockRequest) (*pb.LockResponse, error) {
	addr := extractRemote(ctx)
	if s.verbose {
		log.Printf("RequestLock request for %s from %d with timeout %d", req.GetLockName(), req.GetPid(), req.GetTimeoutSeconds())
	}
	err := s.manager.RequestLock(req.LockName, req.Pid, addr, req.TimeoutSeconds)
	if err != nil {
		log.Printf("RequestLock failed for %s from %d: %s", req.GetLockName(), req.GetPid(), err.Error())
		return &pb.LockResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.LockResponse{Success: true, Message: "Lock acquired"}, nil
}

// extractRemote extracts the remote address from a context containing peer information and returns it as a string.
func extractRemote(ctx context.Context) string {
	p, _ := peer.FromContext(ctx)
	addr := p.Addr.String()
	lio := strings.LastIndex(addr, ":")
	if lio > 0 {
		addr = addr[:lio]
	}
	return addr
}

// ReleaseLock handles lock release requests from clients
func (s *LockServer) ReleaseLock(ctx context.Context, req *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {
	addr := extractRemote(ctx)
	if req.GetForceToken() != "" && s.secretToken == "" {
		return &pb.ReleaseResponse{Success: false, Message: "Forceful release deactivated"}, nil
	}
	if req.GetForceToken() != "" && req.GetForceToken() != s.secretToken {
		return &pb.ReleaseResponse{Success: false, Message: "Invalid secret token"}, nil
	}
	if s.verbose {
		log.Printf("ReleaseLock request for %s from %d, is forced: %t", req.GetLockName(), req.GetPid(), req.GetForceToken() != "")
	}
	var err error
	if req.GetForceToken() == "" {
		err = s.manager.ReleaseLock(req.LockName, req.Pid, addr)
	} else {
		err = s.manager.ReleaseLockByName(req.LockName)
	}
	if err != nil {
		return &pb.ReleaseResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.ReleaseResponse{Success: true, Message: "Lock released"}, nil
}

// List all locks
func (s *LockServer) List(ctx context.Context, _ *pb.ListRequest) (*pb.ListResponse, error) {
	addr := extractRemote(ctx)
	if s.verbose {
		log.Printf("Listrequest from %s", addr)
	}
	resp := &pb.ListResponse{Locks: make([]*pb.Lock, 0)}
	for _, lock := range s.manager.GetLocks() {
		resp.Locks = append(resp.Locks, &pb.Lock{
			Name:   lock.Name,
			Addr:   lock.Addr,
			Pid:    lock.Pid,
			Locked: lock.IsLocked,
		})
	}
	return resp, nil
}
