package types

import "errors"

var (
	// ErrLockExists is returned when an attempt is made to acquire a lock that already exists and is currently held.
	ErrLockExists = errors.New("Lock already exists")

	// ErrStrangersLock is returned when an attempt is made to release a lock that is either not held by the given PID or does not exist.
	ErrStrangersLock = errors.New("lock not held by given PID or does not exist")
)

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

// Locker interface defines methods for acquiring and releasing locks.
type Locker interface {

	// Lock attempts to acquire a lock identified by the given name and associated with the provided process ID (pid).
	Lock(name string, pid int32, addr string) error

	// Unlock releases the lock identified by the given name and associated with the provided process ID (pid). Returns an error if the unlock operation fails.
	Unlock(name string, pid int32, addr string) error

	// GetLocks returns a slice of LockInfo representing all the current locks and their statuses.
	GetLocks() []LockInfo
}
