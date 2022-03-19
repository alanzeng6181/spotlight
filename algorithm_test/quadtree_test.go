package algorithm_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"
	"github.com/alanzeng6181/proximity_server/algorithm"
)

type Quad algorithm.Quad[algorithm.Location]

func TestQuadtreeConcurrentPopulating(t *testing.T) {
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
				tree.Add(algorithm.Location{algorithm.Float(rand.Int31n(width*100))/100.0, algorithm.Float((rand.Int31n(height*100))/100.0)})
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if ok, err := tree.Root.Verify(count * threadCount); !ok {
		t.Errorf("expected tree to have %d nodes, but it did not => %s", count*threadCount, err.Error())
	}
}

func TestQuadtreeQueryingAndPopulating(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	const width = 1800
	const height = 900
	tree := algorithm.MakeQuadTree[algorithm.Location](0.0, 0.0, algorithm.Float(width), algorithm.Float(height))
	const count = 500000
	locations:=make([]algorithm.Location,0)
	const dps = 3000
	var wg sync.WaitGroup
	const threads = 50
	wg.Add(2*threads)
	const eachThread = count/threads;
	for j:=0; j<threads; j++{
		go func() {
			for i := 0; i < eachThread; i++ {
				time.Sleep(time.Duration(1000000/(threads*dps)) * time.Microsecond)
				location := algorithm.Location{algorithm.Float(rand.Int31n(width*100))/100.0, algorithm.Float((rand.Int31n(height*100))/100.0)}
				tree.Add(location)
				locations = append(locations, location)
			}
			wg.Done()
		}()
	}
	for j:=0; j<threads; j++{
		go func() {
			for i := 0; i < 500; i++ {
				time.Sleep(10 * time.Millisecond)
				n := len(locations)
				effectRadius := float32(100.0)
				if n > 0 {
					picked := locations[rand.Int31n(int32(n))]
					tree.FindNearby(float64(picked.X()),
						float64(picked.Y()),
						float64(effectRadius))
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	if ok, err := tree.Root.Verify(eachThread*threads); !ok {
		t.Errorf("expected tree to have %d nodes, but it did not => %s", count, err.Error())
	}
}