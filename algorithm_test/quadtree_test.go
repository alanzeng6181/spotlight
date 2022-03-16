package algorithm_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"
	"github.com/alanzeng6181/proximity_server/algorithm"
)

type QuadNode algorithm.QuadNode[algorithm.Location]

/*
func TestQuadtree(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	const width = 70000
	const height = 70000
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
	if ok, err := tree.Root.Verify(count * threadCount); !ok {
		t.Errorf("expected %d values, but it's not => %s", count*threadCount, err.Error())
	}
}*/

func TestQuadtreeQueryingAndPopulating(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	const width = 70000
	const height = 70000
	tree := algorithm.MakeQuadTree[algorithm.Location](0.0, 0.0, algorithm.Float(width), algorithm.Float(height))
	const count = 75000000
	locations:=make([]*algorithm.Location,0)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < count; i++ {
			location := &algorithm.Location{algorithm.Float(rand.Int31n(width*100))/100.0, algorithm.Float((rand.Int31n(height*100))/100.0)}
			tree.Add(location)
			locations = append(locations, location)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(10 * time.Millisecond)
			n := len(locations)
			effectRadius := float32(1400.0)
			if n > 0 {
				picked := locations[rand.Int31n(int32(n))]
				tree.FindNearby(float32(picked.X()),
					float32(picked.Y()),
					effectRadius)
			}
		}
		wg.Done()
	}()

	wg.Wait()
	if ok, err := tree.Root.Verify(count); !ok {
		t.Errorf("expected %d values, but it's not => %s", count, err.Error())
	}
}