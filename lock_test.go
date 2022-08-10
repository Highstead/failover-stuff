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
	m.cancel()
}

func (m *mockLock) loop() {
	for {
		select {
		case <-m.ctx.Done():
			m.Release()
			return

		case <-m.lockChan:
			m.mu.Lock()
			l := m.aquired.Load().(chan struct{})
			close(l)
			m.aquired.Store(l) //We close the lock here so that select statements will unblock
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
	if !m.HaveLock() {
		return ErrDontHaveLock
	}

	m.mu.Lock()
	select {
	//This will happen during a shutdown of a test
	case _, ok := <-m.lockChan:
		//channel is closed, we got here during cleanup
		if !ok {
			fmt.Println("channel was closed, not releasing lock to channel")
			return nil
		}
		if ok { //We found something on the lock channel and that should be impossible
			panic("there should have never been another event on the lock channel")
		}
	default:
		//this would panic without the previous case statement during a test cleanup.  This **SHOULD** never
		//  block.  If it does something is wrong, and someone else has polluted our lock channel.
		m.lockChan <- true
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
	case <-time.After(time.Millisecond * 5):
		require.Fail(t, "expected lock to be aquired")

	}
	require.True(t, newMl.HaveLock())

}
