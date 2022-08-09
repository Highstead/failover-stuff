package failover

type MySQLNode struct {
}

func NewMySQLNode(parent failable, child failable) *MySQLNode {
	return &MySQLNode{}
}
