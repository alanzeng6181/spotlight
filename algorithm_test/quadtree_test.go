package algorithm_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
	"github.com/alanzeng6181/proximity_server/algorithm"
)

type QuadNode algorithm.QuadNode[algorithm.Location]

func TestQuadtree(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	const width = 10000000
	const height = 10000000
	tree := algorithm.MakeQuadTree[algorithm.Location](0.0, 0.0, algorithm.Float(width), algorithm.Float(height))
	const count = 5000
	const threadCount = 5000
	var wg sync.WaitGroup
	wg.Add(threadCount)
	for j := 0; j < threadCount; j++ {
		go func() {
			for i := 0; i < count; i++ {
				tree.Add(&algorithm.Location{algorithm.Float(rand.Int31n(width*100))/100.0, algorithm.Float((rand.Int31n(height*100))/100.0)})
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if ok, err := (*QuadNode)(tree.Root).Verify(count * threadCount); !ok {
		t.Errorf("expected %d values, but it's not => %s", count*threadCount, err.Error())
	}
}

func (node *QuadNode) Verify(count int) (bool, error) {
	if node.Size <= algorithm.MAX_NODES && (node.Children[0] != nil || node.Children[1] != nil || node.Children[2] != nil || node.Children[3] != nil) {
		return false, fmt.Errorf("when there are less than %d locations, children should be nil", algorithm.MAX_NODES)
	}

	if node.Size > algorithm.MAX_NODES && (node.Children[0] == nil || node.Children[1] == nil ||
		node.Children[2] == nil || node.Children[3] == nil || (node.Values != nil && len(node.Values) != 0)) {
		return false, fmt.Errorf("when there are greater than %d locations, children should not be nil and values should be empty",
			algorithm.MAX_NODES)
	}

	if node.Size != count || (len(node.Values) > 0 && node.Size != len(node.Values)) {
		return false, fmt.Errorf("node size %d does not equal to count %d, or node values array does not match node size when it's not empty",
			node.Size, count)
	}

	if node.Size > algorithm.MAX_NODES {
		sum := 0
		for _, c := range node.Children {
			sum += c.Size
			if ok, err := (*QuadNode)(c).Verify(c.Size); !ok {
				return false, err
			}
		}
		if sum != node.Size {
			return false, fmt.Errorf("node size is %d, but sum of children node size is %d", node.Size, sum)
		}
	}
	return true, nil
}