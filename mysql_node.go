package failover

/*type failable interface {
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
}*/

type MySQLNode struct {
}

func NewMySQLNode(parent failable, child failable) *MySQLNode {
	return &MySQLNode{}
}
