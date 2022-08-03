package failover

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type FauxNode struct {
	position string
	parent   *FauxNode
	children []*FauxNode

	lock     *sync.Mutex
	writable bool
}

// NewFauxNode is a tool for testing.. it will panic if any of the parents or children are not FauxNodes
func NewFauxNode(position string, parent failable, children []failable) *FauxNode {
	return &FauxNode{
		position: position,
		parent:   parent.(*FauxNode),
	}

}

func (f *FauxNode) Position() string {
	return f.Position()
}

func (f *FauxNode) ComparePosition(pos string) int {
	return strings.Compare(f.Position(), pos)
}
func (f *FauxNode) Children() []failable {
	var failables []failable
	for _, child := range f.children {
		failables = append(failables, child)
	}
	return failables
}

func (f *FauxNode) Parent() failable {
	return f.parent
}

func (f *FauxNode) PrepareForTakeover() error {
	return nil
}
func (f *FauxNode) GracefulTakeover() error {
	return nil
}
func (f *FauxNode) HostileTakeover() error {
	return nil
}

func TestGracefultakeoverHappyPath(t *testing.T) {
	tt := []struct {
		name        string
		numChildren int
	}{
		{name: "oneChild", numChildren: 1},
		{name: "threeChildren", numChildren: 3},
	}

	for _, myTest := range tt {
		require.Greater(t, 1, myTest.numChildren, fmt.Sprintf("%s: failed due to illegal number of children", myTest.name))

		//NewCluster := FailableCluster()
		parent := NewFauxNode(strconv.Itoa(myTest.numChildren), nil, nil)

		for i := 0; i < myTest.numChildren; i++ {
			newChild := NewFauxNode(strconv.Itoa(i), nil, nil)
			parent.children = append(parent.children, newChild)
		}

	}
	require.Fail(t, "unimplemented")
}
