package algorithm

import (
	"fmt"
	"log"
	"math"
	"sync"
)

const MAX_NODES = 100
type Float float64

type QuadTree[T Locationer] struct{
	Root *Quad[T]
}

func MakeQuadTree[T Locationer](x0 Float, y0 Float, x1 Float, y1 Float) QuadTree[T]{
	return QuadTree[T]{Root:NewQuad[T](x0,y0,x1,y1)}
}

//Add adds a node to the tree.
func (tree QuadTree[T]) Add(t T) bool{
	return tree.Root.Add(t)
}

//FindNearby returns an array of nodes d distance from coordinate(x,y) that are within the tree.
func (tree QuadTree[T]) FindNearby(x float64, y float64, d float64) []T{
	return tree.Root.FindNearby(x,y,d)
}

type Quad[T Locationer] struct{
	mu sync.RWMutex
	Children   [4]*Quad[T]
	Identifier string
	Nodes     []T
	LowerLeft  []Float
	UpperRight []Float
	Size       int
}

func NewQuad[T Locationer](x0 Float, y0 Float, x1 Float, y1 Float) *Quad[T]{
	return &Quad[T]{
		Children:   [4]*Quad[T]{},
		Identifier: "id",
		Nodes:     make([]T,0),
		LowerLeft:  []Float{x0,y0},
		UpperRight: []Float{x1,y1},
	}
}

type Locationer interface {
	X() Float
	Y() Float
}

//Add adds a node to a quad.
func (quad *Quad[T]) Add(node T) bool{
	if !quad.IsOverlap(node){
		return false
	}
	quad.mu.Lock()
	quad.Size=quad.Size+1

	if quad.Size<=MAX_NODES{
		quad.Nodes = append(quad.Nodes, node)
		quad.mu.Unlock()
		return true
	}
	
	if quad.Size==MAX_NODES+1{
		var children [4]*Quad[T]
		xMid := (quad.UpperRight[0]+quad.LowerLeft[0])/2
		yMid := (quad.UpperRight[1]+quad.LowerLeft[1])/2
		children[0] = &Quad[T]{Identifier:quad.Identifier+"ll", LowerLeft:quad.LowerLeft, UpperRight: []Float{xMid, yMid}}
		children[1] = &Quad[T]{Identifier:quad.Identifier+"ul", LowerLeft:[]Float{quad.LowerLeft[0],yMid}, UpperRight: []Float{xMid, quad.UpperRight[1]}}
		children[2] = &Quad[T]{Identifier:quad.Identifier+"lr", LowerLeft:[]Float{xMid, quad.LowerLeft[1]}, UpperRight: []Float{ quad.UpperRight[0], yMid}}
		children[3] = &Quad[T]{Identifier:quad.Identifier+"ur", LowerLeft:[]Float{xMid, yMid}, UpperRight: quad.UpperRight}
		quad.Children = children
		quad.mu.Unlock()
		for _, _node := range quad.Nodes{
			var isAdded = false
			for _, c := range children{
				if c.IsOverlap(_node){
					c.Add(_node)
					isAdded=true
					break
				}
			}
			if !isAdded{
				for _, c:= range children{
					log.Println("%v, %v", c.LowerLeft, c.UpperRight)
				}
				log.Fatalln("unable to add to children x:%f, y%f, %v, %v", node.X(),node.Y(), quad.LowerLeft, quad.UpperRight)
			}
		}
		quad.Nodes=make([]T,0)
	} else{
		quad.mu.Unlock()
	}

	for _, c := range quad.Children{
		if c.IsOverlap(node){
			c.Add(node)
			return true
		}
	}
	return true
}

//IsOverlap checks whether v is within Quad
func (quad *Quad[T]) IsOverlap(node T) bool{
	if node.X()<quad.LowerLeft[0]||node.X()>quad.UpperRight[0]||node.Y()<quad.LowerLeft[1]||node.Y()>quad.UpperRight[1]{
		return false
	}
	return true
}

//FindNearby finds all nodes that are within d distance from point (x,y)
func (quad *Quad[T]) FindNearby(x float64, y float64, d float64) []T{
	//whether the circle touches the square.
	touches := false

	//check if one of the corners of square is inside the circle.
	corners:= [][]float64{[]float64{float64(quad.LowerLeft[0]), float64(quad.LowerLeft[1])}, 
	[]float64{float64(quad.LowerLeft[0]), float64(quad.UpperRight[1])},
	[]float64{float64(quad.UpperRight[0]), float64(quad.LowerLeft[1])},
	[]float64{float64(quad.UpperRight[0]), float64(quad.UpperRight[1])}}

	for _, c := range corners{
		if float64(math.Sqrt(math.Pow(float64(c[0]-x),2.0)+math.Pow(float64(c[1]-y),2.0))) <=d{
			touches = true
			break
		}
	}

	if !touches{
		//check most-left, most-right, most-up, most-down points of circle, see if one of them 
		// is inside the square.
		outerMost := [][]float64{[]float64{x-d,y},[]float64{x+d,y},[]float64{x,y+d},[]float64{x,y-d}}
		for _, o := range outerMost{
			if o[0]>=float64(quad.LowerLeft[0]) && o[0]<=float64(quad.UpperRight[0]) && 
			o[1]>=float64(quad.LowerLeft[1]) && 
			o[1]<=float64(quad.UpperRight[1]){
				touches = true
				break
			}
		}
	}

	if !touches{
		return nil
	}

	list := make([]T,0)
	quad.mu.RLock()
	if quad.Size <=MAX_NODES{
		for _, node := range quad.Nodes{
			if math.Sqrt(math.Pow(float64(node.X())-float64(x), 2.0) + math.Pow(float64(node.Y())-float64(y), 2.0)) <= float64(d){
				list = append(list, node)
			}
		}
		quad.mu.RUnlock()
		return list
	}
	quad.mu.RUnlock()
	for _, c:= range quad.Children{
		if sublist:=c.FindNearby(x, y, d); sublist!=nil{
			list = append(list, sublist...)
		}
	}
	return list
}

//Verify checks Quad integrity
func (quad *Quad[T]) Verify(count int) (bool, error) {
	if quad.Size <= MAX_NODES && (quad.Children[0] != nil || quad.Children[1] != nil || quad.Children[2] != nil || quad.Children[3] != nil) {
		return false, fmt.Errorf("when there are less than %d locations, children should be nil", MAX_NODES)
	}

	if quad.Size > MAX_NODES && (quad.Children[0] == nil || quad.Children[1] == nil ||
		quad.Children[2] == nil || quad.Children[3] == nil || (quad.Nodes != nil && len(quad.Nodes) != 0)) {
		return false, fmt.Errorf("when there are greater than %d locations, children should not be nil and values should be empty",
			MAX_NODES)
	}

	if quad.Size != count || (len(quad.Nodes) > 0 && quad.Size != len(quad.Nodes)) {
		return false, fmt.Errorf("node size %d does not equal to count %d, or node values array does not match node size when it's not empty",
			quad.Size, count)
	}

	if quad.Size > MAX_NODES {
		sum := 0
		for _, c := range quad.Children {
			sum += c.Size
			if ok, err := (*Quad[T])(c).Verify(c.Size); !ok {
				return false, err
			}
		}
		if sum != quad.Size {
			return false, fmt.Errorf("node size is %d, but sum of children node size is %d, len(values):%d", quad.Size, sum, len(quad.Nodes))
		}
	}
	return true, nil
}