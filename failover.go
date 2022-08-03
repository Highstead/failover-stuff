package failover

type MysqlNode struct {
}

type failable interface {
	Position() string
	ComparePosition(string) int

	Children() []failable
	Parent() failable

	PrepareForTakeover() error

	GracefulTakeover() error
	HostileTakeover() error
}
