package failover

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockLock struct {
	ctx      context.Context
	lockChan chan bool

	aquired atomic.Value
	mu      sync.Mutex

	cancel context.CancelFunc
	name   string
}

func NewMockLock(ctx context.Context, lockChan chan bool, name string) *mockLock {
	ctx, cancel := context.WithCancel(ctx)

	ml := &mockLock{
		ctx:      ctx,
		lockChan: lockChan,
		cancel:   cancel,
		name:     name,
	}
	//Make a single atomic channel.  Make the channel so that it blocks in select statements.  This channel will be closed
	//  When the lock has been aquired.
	aquired := ml.aquired.Load()
	aquired = make(chan struct{})
	ml.aquired.Store(aquired)

	go ml.loop()

	return ml
}

func (m *mockLock) Close() {
	fmt.Println("Closing mocklock", m.name)
	m.cancel()
}

func (m *mockLock) loop() {
	for {
		select {
		case <-m.ctx.Done():
			fmt.Println("Done on", m.name)
			if m.HaveLock() {
				m.Release()
			}
			return

		case <-m.lockChan:
			m.mu.Lock()
			m.aquired.Store(closedChan) //We close the lock here so that select statements will unblock
			m.mu.Unlock()
		}
	}
}

func (m *mockLock) HaveLock() bool {
	select {
	case <-m.LockAquired():
		//this will fire if the 'aquired' channel is closed
		return true
	default:
	}

	return false
}

func (m *mockLock) LockAquired() <-chan struct{} {
	a := m.aquired.Load()
	return a.(chan struct{})
}

func (m *mockLock) Release() error {
	fmt.Println("Releasing mocklock", m.name)
	if !m.HaveLock() {
		return ErrDontHaveLock
	}
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
	defer cancel()

	ctx2 := context.Background()
	lockChan := make(chan bool)
	defer close(lockChan)

	ml := NewMockLock(ctx1, lockChan, "ml")
	lockChan <- true //Populate the lock channel so someone can aquire the lock

	select {
	case <-ml.LockAquired():
	case <-time.After(time.Millisecond):
		require.Fail(t, "expected lock to be aquired")
	}
	require.True(t, ml.HaveLock())

	//Make sure we only have 1 lock
	newMl := NewMockLock(ctx2, lockChan, "newMl")
	defer newMl.Close()
	select {
	case <-newMl.LockAquired():
		require.Fail(t, "expected lock to be unaquirable")
	case <-time.After(time.Millisecond):
	}
	require.False(t, newMl.HaveLock())

	//Cancel the initial context, which should release the lock
	ml.Close()
	err := ml.Release()
	require.NoError(t, err)

	//Make sure the Previous lock is cleaned up
	select {
	case <-newMl.LockAquired():
	case <-time.After(time.Second):
		require.Fail(t, "expected lock to be aquired")
	}
	require.True(t, newMl.HaveLock())
	fmt.Println("Shutting down")

}
