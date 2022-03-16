package algorithm

import (
	"math"
	"sync"
	"fmt"
)

const MAX_NODES = 100
type Float float64

type QuadTree[T Locationer] struct{
	Root *QuadNode[T]
}

func MakeQuadTree[T Locationer](x0 Float, y0 Float, x1 Float, y1 Float) QuadTree[T]{
	return QuadTree[T]{Root:NewQuadNode[T](x0,y0,x1,y1)}
}

func (tree QuadTree[T]) Add(t *T) bool{
	return tree.Root.Add(t)
}

func (tree QuadTree[T]) FindNearby(x float32, y float32, d float32) []*T{
	return tree.Root.FindNearby(x,y,d)
}

type QuadNode[T Locationer] struct{
	mu sync.RWMutex
	Children   [4]*QuadNode[T]
	Identifier string
	Values     []*T
	LowerLeft  []Float
	UpperRight []Float
	Size       int
}

func NewQuadNode[T Locationer](x0 Float, y0 Float, x1 Float, y1 Float) *QuadNode[T]{
	return &QuadNode[T]{
		Children:   [4]*QuadNode[T]{},
		Identifier: "id",
		Values:     make([]*T,0),
		LowerLeft:  []Float{x0,y0},
		UpperRight: []Float{x1,y1},
	}
}

type Locationer interface {
	X() Float
	Y() Float
}

func (node *QuadNode[T]) Add(v *T) bool{
	if !node.IsOverlap(v){
		return false
	}
	node.mu.Lock()
	node.Size=node.Size+1

	if node.Size<=MAX_NODES{
		node.Values = append(node.Values, v)
		node.mu.Unlock()
		return true
	}
	
	if node.Size==MAX_NODES+1{
		var children [4]*QuadNode[T]
		xMid := (node.UpperRight[0]+node.LowerLeft[0])/2
		yMid := (node.UpperRight[1]+node.LowerLeft[1])/2
		children[0] = &QuadNode[T]{Identifier:node.Identifier+"ll", LowerLeft:node.LowerLeft, UpperRight: []Float{xMid, yMid}}
		children[1] = &QuadNode[T]{Identifier:node.Identifier+"ul", LowerLeft:[]Float{node.LowerLeft[0],yMid}, UpperRight: []Float{xMid, node.UpperRight[1]}}
		children[2] = &QuadNode[T]{Identifier:node.Identifier+"lr", LowerLeft:[]Float{xMid, node.LowerLeft[1]}, UpperRight: []Float{ node.UpperRight[0], yMid}}
		children[3] = &QuadNode[T]{Identifier:node.Identifier+"ur", LowerLeft:[]Float{xMid, yMid}, UpperRight: node.UpperRight}
		node.Children = children
		copyValues := make([]*T, len(node.Values))
		copy(copyValues, node.Values)
		node.Values=make([]*T,0)
		node.mu.Unlock()
		for _, v := range copyValues{
			for _, c := range children{
				if c.IsOverlap(v){
					c.Add(v)
					break
				}
			}
		}
	} else{
		node.mu.Unlock()
	}

	for _, c := range node.Children{
		if c.IsOverlap(v){
			c.Add(v)
			return true
		}
	}
	return true
}

func (node *QuadNode[T]) IsOverlap(v *T) bool{
	if (*v).X()<node.LowerLeft[0]||(*v).X()>node.UpperRight[0]||(*v).Y()<node.LowerLeft[1]||(*v).Y()>node.UpperRight[1]{
		return false
	}
	return true
}

func (node *QuadNode[T]) FindNearby(x float32, y float32, d float32) []*T{
	//check whether the circle touches the square.
	touches := false

	//check if one of the corners of square is inside the circle.
	corners:= [][]float32{[]float32{float32(node.LowerLeft[0]), float32(node.LowerLeft[1])}, 
	[]float32{float32(node.LowerLeft[0]), float32(node.UpperRight[1])},
	[]float32{float32(node.UpperRight[0]), float32(node.LowerLeft[1])},
	[]float32{float32(node.UpperRight[0]), float32(node.UpperRight[1])}}

	for _, c := range corners{
		if float32(math.Sqrt(math.Pow(float64(c[0]-x),2.0)+math.Pow(float64(c[1]-y),2.0))) <=d{
			touches = true
			break
		}
	}

	if !touches{
		//check most-left, most-right, most-up, most-down points of circle, see if one of them 
		// is inside the square.
		outerMost := [][]float32{[]float32{x-d,y},[]float32{x+d,y},[]float32{x,y+d},[]float32{x,y-d}}
		for _, o := range outerMost{
			if o[0]>=float32(node.LowerLeft[0]) && o[0]<=float32(node.UpperRight[0]) && 
			o[1]>=float32(node.LowerLeft[1]) && 
			o[1]<=float32(node.UpperRight[1]){
				touches = true
				break
			}
		}
	}

	if !touches{
		return nil
	}

	list := make([]*T,0)
	node.mu.RLock()
	if node.Size <=MAX_NODES{
		for _, v := range node.Values{
			if math.Sqrt(math.Pow(float64((*v).X())-float64(x), 2.0) + math.Pow(float64((*v).Y())-float64(y), 2.0)) <= float64(d){
				list = append(list, v)
			}
		}
		node.mu.RUnlock()
		return list
	}
	node.mu.RUnlock()
	for _, c:= range node.Children{
		if sublist:=c.FindNearby(x, y, d); sublist!=nil{
			list = append(list, sublist...)
		}
	}
	return list
}

func (node *QuadNode[T]) Verify(count int) (bool, error) {
	if node.Size <= MAX_NODES && (node.Children[0] != nil || node.Children[1] != nil || node.Children[2] != nil || node.Children[3] != nil) {
		return false, fmt.Errorf("when there are less than %d locations, children should be nil", MAX_NODES)
	}

	if node.Size > MAX_NODES && (node.Children[0] == nil || node.Children[1] == nil ||
		node.Children[2] == nil || node.Children[3] == nil || (node.Values != nil && len(node.Values) != 0)) {
		return false, fmt.Errorf("when there are greater than %d locations, children should not be nil and values should be empty",
			MAX_NODES)
	}

	if node.Size != count || (len(node.Values) > 0 && node.Size != len(node.Values)) {
		return false, fmt.Errorf("node size %d does not equal to count %d, or node values array does not match node size when it's not empty",
			node.Size, count)
	}

	if node.Size > MAX_NODES {
		sum := 0
		for _, c := range node.Children {
			sum += c.Size
			if ok, err := (*QuadNode[T])(c).Verify(c.Size); !ok {
				return false, err
			}
		}
		if sum != node.Size {
			return false, fmt.Errorf("node size is %d, but sum of children node size is %d, len(values):%d", node.Size, sum, len(node.Values))
		}
	}
	return true, nil
}