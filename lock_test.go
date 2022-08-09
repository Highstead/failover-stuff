package failover

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	ErrLockTimeout     = fmt.Errorf("unable to release lock within aloted time")
	defaultLockTimeout = time.Millisecond * 10
)

type mockLock struct {
	ctx      context.Context
	hasLock  bool
	lockChan chan bool

	cancel
}

func NewMockLock(ctx context.Context, lockChan chan bool) *mockLock {
	ml := &mockLock{
		ctx:      ctx,
		lockChan: lockChan,
	}

	go ml.loop()

	return ml
}

func (m *mockLock) loop() {
	for {
		select {
		case <-m.ctx.Done():
			if m.hasLock {
				m.Release()
			}
			return

		case m.hasLock = <-m.lockChan:
		}
	}
}

func (m *mockLock) HaveLock() bool {
	return m.hasLock
}

func (m *mockLock) Release() error {
	select {
	//If the lock chan is closed we will panic here.  But this is a mock so this should never happen
	case m.lockChan <- true:

	//This should never happen
	case <-time.After(defaultLockTimeout):
		return ErrLockTimeout
	}
	return nil
}

func TestMockLock(t *testing.T) {
	ctx1, cancel := context.WithCancel(context.Background())
	ctx2 := context.Background()
	lockChan := make(chan bool)
	defer close(lockChan)

	ml := NewMockLock(ctx1, lockChan)
	require.True(t, ml.HaveLock())

	//Make sure we only have 1 lock
	newMl := NewMockLock(ctx2, lockChan)
	require.False(t, newMl.HaveLock())

	cancel()

	require.True(t, newMl.HaveLock())

}
