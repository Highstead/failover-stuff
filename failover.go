package failover

import "fmt"

var (
	ErrAlreadyChild = fmt.Errorf("node is already a child to parent")
)

type failable interface {
	Position() string
	ComparePosition(failable) int

	Children() []failable
	Parent() failable
	SetParent(failable) error

	PrepareForTakeover() error
	PrepareToTakeover() error
	CompleteTakeover(isWritable bool) error

	RevertTakeoverAttempt() error
	UID() string
}

type FailoverService struct {
	lock lockService
}

func GracefulFailover(parent failable, child failable) error {
	//Prune candidates based on rules (region?)

	//Select most ahead candidate

	//Prepare for takeover
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
