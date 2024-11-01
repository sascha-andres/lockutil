package types

import "errors"

var (
	// ErrLockExists is returned when an attempt is made to acquire a lock that already exists and is currently held.
	ErrLockExists = errors.New("Lock already exists")

	// ErrStrangersLock is returned when an attempt is made to release a lock that is either not held by the given PID or does not exist.
	ErrStrangersLock = errors.New("lock not held by given PID or does not exist")
)

// Locker interface defines methods for acquiring and releasing locks.
type Locker interface {

	// Lock attempts to acquire a lock identified by the given name and associated with the provided process ID (pid).
	Lock(name string, pid int32) error

	// Unlock releases the lock identified by the given name and associated with the provided process ID (pid). Returns an error if the unlock operation fails.
	Unlock(name string, pid int32) error
}
