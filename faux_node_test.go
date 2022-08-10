package failover

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var count int
var lock sync.Mutex

//getNewFauxNamer is just here to make sure we have unique identifiable names for our test instances
func getFauxUID() string {
	lock.Lock()
	defer lock.Unlock()

	count = count + 1
	return strconv.Itoa(count)
}

type FauxNode struct {
	position string
	parent   *FauxNode
	children []*FauxNode

	lock     *sync.Mutex
	writable bool

	uid string
}

// NewFauxNode is a tool for testing.. it will panic if any of the parents or children are not FauxNodes
func NewFauxNode(position string, parent failable, children []failable) *FauxNode {
	fn := &FauxNode{
		position: position,
		uid:      getFauxUID(),
	}

	if parent != nil {
		fn.parent = parent.(*FauxNode)
	}

	return fn
}

//String is here to implement stringer and make debugging easier
func (f *FauxNode) String() string {
	if f.parent != nil {
		return fmt.Sprintf(" uid: %s, pos: %s, parent: %s", f.uid, f.position, f.parent.uid)
	}
	return fmt.Sprintf(" uid: %s, pos: %s, parent: nil", f.uid, f.position)
}

func (f *FauxNode) Writable() bool {
	return f.writable
}

func (f *FauxNode) ComparePosition(other failable) int {
	o := other.(*FauxNode)
	return strings.Compare(f.Position(), o.Position())
}

func (f *FauxNode) CompleteTakeover(isWritable bool) error {
	f.writable = isWritable
	return nil
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
	if f.writable {
		f.writable = false
	}
	return nil
}

//our fake binlog position
func (f *FauxNode) Position() string {
	return f.position
}

func (f *FauxNode) PrepareToTakeover() error {
	return nil
}

func (f *FauxNode) RevertTakeoverAttempt() error {
	return nil
}

//setPosition is here for testing purposes
func (f *FauxNode) setPosition(pos string) {
	f.position = pos
}

func (f *FauxNode) SetParent(parent failable) error {
	p := parent.(*FauxNode)

	for _, child := range p.children {
		// this node is already a replica of the parent do nothing
		if child.uid == f.uid {
			return ErrAlreadyChild
		}
	}
	f.parent = parent.(*FauxNode)
	p.children = append(p.children, f)
	return nil
}

//UID is effectively the node name and should be set by 'getFauxUID' on creation
func (f *FauxNode) UID() string {
	return f.uid
}

/// Tests start here
func TestInitFauxNodeUID(t *testing.T) {
	var testNodes []*FauxNode
	for i := 0; i < 10; i++ {
		testNodes = append(testNodes, NewFauxNode(strconv.Itoa(i), nil, nil))
	}
	for i := 0; i < 10; i++ {
		for j := i + 1; j < 10; j++ {
			//Test for unique UID
			require.NotEqual(t, testNodes[i].UID(), testNodes[j].UID())

			if i == 0 {
				//Test that setParent doesnt fail the first time
				require.NoError(t, testNodes[j].SetParent(testNodes[0]))
			}

			//Test that we have already set the parent
			err := testNodes[j].SetParent(testNodes[0])
			require.Error(t, err)
			require.Equal(t, ErrAlreadyChild, err)
		}
	}

}

func TestComparePosition(t *testing.T) {
	parent := NewFauxNode("1", nil, nil)
	child := NewFauxNode("0", nil, nil)
	child.SetParent(parent)

	require.Equal(t, -1, child.ComparePosition(parent))
	require.Equal(t, 1, parent.ComparePosition(child))

	child.setPosition(parent.position)
	require.Equal(t, 0, parent.ComparePosition(child))
}
