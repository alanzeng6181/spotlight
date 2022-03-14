package algorithm

import "sync"

const MAX_NODES = 100
type Float float64

type QuadTree[T Locationer] struct{
	Root *QuadNode[T]
}

func MakeQuadTree[T Locationer](x0 Float, y0 Float, x1 Float, y1 Float) QuadTree[T]{
	return QuadTree[T]{Root:NewQuadNode[T](x0,y0,x1,y1)}
}

func (tree QuadTree[T]) Add(t T){
	tree.Root.Add(t)
}

type QuadNode[T Locationer] struct{
	mu sync.Mutex
	Children   [4]*QuadNode[T]
	Identifier string
	Values     []T
	LowerLeft  []Float
	UpperRight []Float
	Size       int
}

func NewQuadNode[T Locationer](x0 Float, y0 Float, x1 Float, y1 Float) *QuadNode[T]{
	return &QuadNode[T]{
		Children:   [4]*QuadNode[T]{},
		Identifier: "id",
		Values:     make([]T,0),
		LowerLeft:  []Float{x0,y0},
		UpperRight: []Float{x1,y1},
	}
}

type Locationer interface {
	X() Float
	Y() Float
}

func (node *QuadNode[T]) Add(v T) bool{
	if !node.IsOverlap(v){
		return false
	}
	node.mu.Lock()
	node.Size=node.Size+1
	if node.Size<=MAX_NODES{
		node.Values = append(node.Values,v)
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
		node.mu.Unlock()
		for _, v := range node.Values{
			for _, c := range children{
				if c.IsOverlap(v){
					c.Add(v)
					break
				}
			}
		}
		node.Values=make([]T,0)
	} else{
		node.mu.Unlock()
	}

	for _, child := range node.Children{
		if child.IsOverlap(v){
			return child.Add(v)
		}
	}
	return true
}

func (node *QuadNode[T]) IsOverlap(v T) bool{
	if v.X()<node.LowerLeft[0]||v.X()>node.UpperRight[0]||v.Y()<node.LowerLeft[1]||v.Y()>node.UpperRight[1]{
		return false
	}
	return true
}