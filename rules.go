package failover

import "context"

type ruleService interface {
	IsNodeHealthy(failable) bool
	WaitForNodeToBeHealthy(context.Context, failable) error

	SuitableCandidates(ctx context.Context, parent failable) ([]failable, error)
	BestCandidate(ctx context.Context, parent failable) (failable, error)
}
