package failover

import (
	"context"
	"fmt"
)

var (
	ErrAlreadyChild = fmt.Errorf("node is already a child to parent")
)

type failable interface {
	Position() string
	ComparePosition(failable) int
	Writable() bool

	Children() []failable
	Parent() failable
	SetParent(failable) error

	PrepareForTakeover() error
	PrepareToTakeover() error
	CompleteTakeover(isWritable bool) error

	RevertTakeoverAttempt() error
	UID() string
}

func NewFailoverService(lock lockService, rules ruleService) *FailoverService {
	return &FailoverService{
		lock:  lock,
		rules: rules,
	}
}

type FailoverService struct {
	lock  lockService
	rules ruleService
}

func (fs *FailoverService) GracefulFailover(ctx context.Context, parent failable, child failable) error {
	//Prune candidates based on rules (region?)

	//Select most ahead candidate

	//Prepare for takeover
	_, err := fs.rules.SuitableCandidates(ctx, parent)
	if err != nil {
		return err
	}
	if err := parent.PrepareForTakeover(); err != nil {
		return err
	}
	if err := child.PrepareToTakeover(); err != nil {
		return err
	}

	return nil
}

func MigrateChildren(source failable, dest failable) error {
	return nil
}
