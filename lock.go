package failover

import (
	"fmt"
	"time"
)

var (
	ErrLockTimeout     = fmt.Errorf("unable to release lock within aloted time")
	ErrDontHaveLock    = fmt.Errorf("this lock does not posess the lock")
	defaultLockTimeout = time.Millisecond * 10
	closedChan         = make(chan struct{})
)

func init() {
	close(closedChan)
}

type lockService interface {
	Release() error
	HaveLock() bool
	LockAquired() <-chan struct{}
}
