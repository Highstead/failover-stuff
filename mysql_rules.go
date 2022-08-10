package failover

import (
	"context"
	"fmt"
	"time"
)

/*
type ruleService interface {
	IsNodeHealthy(failable) bool
	WaitForNodeToBeHealthy(failable) error

	SuitableCandidates(parent failable)
}*/

var (
	ErrNodeIsBehind         = fmt.Errorf("node is behind")
	ErrNoSuitableCandidates = fmt.Errorf("no suitable candidates")
	mysqlPollFrequency      = time.Second
)

type MysqlRules struct {
}

func NewMysqlRules() ruleService {
	return &MysqlRules{}
}

func (m *MysqlRules) IsNodeHealthy(node failable) bool {
	if node.Parent() == nil && node.Writable() {
		return true
	}

	return false
}

func (m *MysqlRules) WaitForNodeToBeHealthy(ctx context.Context, node failable) error {
	child := node
	parent := node.Parent()

	for {
		select {
		case <-ctx.Done():
			return ErrNodeIsBehind
		case <-time.After(mysqlPollFrequency):
			if child.ComparePosition(parent) == 0 {
				return nil
			}
		}

	}

}

func (m *MysqlRules) SuitableCandidates(ctx context.Context, parent failable) ([]failable, error) {
	return nil, ErrNoSuitableCandidates
}

func (m *MysqlRules) BestCandidate(ctx context.Context, parent failable) (failable, error) {
	candidates, err := m.SuitableCandidates(ctx, parent)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, ErrNoSuitableCandidates
	}

	var bestCandidate failable
	for _, candidate := range candidates {
		if bestCandidate == nil {
			bestCandidate = candidate
		}

		//TODO: Prefer IOthread lag over sqlthread i think? at some point
		if bestCandidate.ComparePosition(candidate) == 1 {
			//Candidate is ahead of best candidate replace best candidate
			bestCandidate = candidate
		}
	}
	return bestCandidate, nil
}
