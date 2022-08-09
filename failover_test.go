package failover

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGracefultakeoverHappyPath(t *testing.T) {
	tt := []struct {
		name        string
		numChildren int
	}{
		{name: "oneChild", numChildren: 1},
		{name: "threeChildren", numChildren: 3},
	}

	for _, myTest := range tt {
		require.Greater(t, 0, myTest.numChildren, fmt.Sprintf("%s: failed due to illegal number of children", myTest.name))

		//NewCluster := FailableCluster()
		parent := NewFauxNode(strconv.Itoa(myTest.numChildren), nil, nil)

		for i := 0; i < myTest.numChildren; i++ {
			newChild := NewFauxNode(strconv.Itoa(i), nil, nil)
			parent.children = append(parent.children, newChild)
		}

	}
	require.Fail(t, "unimplemented")
}
