package lockmanager

import (
	"testing"
)

func TestRequestLock(t *testing.T) {
	manager := NewLockManager()

	tests := []struct {
		name    string
		pid     int32
		timeout int32
		wantErr bool
	}{
		{
			name:    "lock that doesn't exist",
			pid:     123,
			timeout: 1,
			wantErr: false,
		},
		{
			name:    "lock that already exists",
			pid:     123,
			timeout: 1,
			wantErr: true,
		},
		{
			name:    "lock request timeout",
			pid:     123,
			timeout: 0,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := manager.RequestLock("default", test.pid, test.timeout)
			if (err != nil) != test.wantErr {
				t.Errorf("got error = %v, wantErr = %v", err != nil, test.wantErr)
			}

			if test.wantErr {
				return
			}

			existingLock, exists := manager.locks[test.name]
			if !exists {
				t.Errorf("lock does not exist, but should")
			}
			if existingLock.pid != test.pid {
				t.Errorf("got PID = %v, want PID = %v", existingLock.pid, test.pid)
			}
			if !existingLock.isLocked {
				t.Errorf("lock is not held, but should be")
			}

			// Release lock after testing for subsequent tests to work properly
			manager.ReleaseLock(test.name, test.pid)
		})
	}
}
