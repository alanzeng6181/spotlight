package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/alanzeng6181/proximity_server/algorithm"
	"image/color"
)

type Dot struct{
	*canvas.Circle
	color color.Color
	radius float32
	position fyne.Position
}

func NewDot(color color.Color, radius float32, position fyne.Position) *Dot{
	return &Dot{
		Circle:&canvas.Circle{
			Position1:   position.Add(fyne.NewDelta(-1*radius, radius)),
			Position2:    position.Add(fyne.NewDelta(radius, -1*radius)),
			Hidden:      false,
			FillColor:   color,
			StrokeColor: color,
			StrokeWidth: 0,
		},
		color:color,
		radius:radius,
		position:position,
	}
}

func (dot Dot) X() algorithm.Float{
	return algorithm.Float(dot.position.X)
}

func (dot Dot) Y() algorithm.Float{
	return algorithm.Float(dot.position.Y)
}

func (dot Dot) CreateRenderer() fyne.WidgetRenderer{
	return NewDotRenderer(&dot)
}

type dotRenderer struct{
	dot *Dot
}

func NewDotRenderer(dot *Dot) *dotRenderer{
	return &dotRenderer{dot:dot}
}

func (r dotRenderer) Layout(size fyne.Size){
	r.dot.Circle.Resize(size)
}

func (r dotRenderer) MinSize() fyne.Size{
	return r.dot.Circle.MinSize()
}

func (r dotRenderer) Refresh(){
	r.dot.Circle.Refresh()
}

func (r dotRenderer) Objects() []fyne.CanvasObject{
	return []fyne.CanvasObject{r.dot.Circle}
}

func (r dotRenderer) Destroy(){
}
