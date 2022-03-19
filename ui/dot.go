package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/alanzeng6181/proximity_server/algorithm"
)

var _ algorithm.Locationer = (*Dot)(nil)

type Dot struct {
	canvas.Circle
	Center fyne.Position
}

func NewDot(color color.Color, radius float32, position fyne.Position) *Dot {
	circle := canvas.NewCircle(color)
	circle.Hidden = false
	circle.StrokeColor = color
	circle.FillColor = color
	circle.Move(position)
	circle.Resize(fyne.NewSize(radius, radius))
	return &Dot{
		Circle: *circle,
		Center: position,
	}
}

func (dot *Dot) X() algorithm.Float {
	return algorithm.Float(dot.Center.X)
}

func (dot *Dot) Y() algorithm.Float {
	return algorithm.Float(dot.Center.Y)
}

func (dot *Dot) Glow(sizeMultiply float64, duration time.Duration, glowColor color.Color) {
	go func() {
		if !sameColor(dot.Circle.FillColor, glowColor) {
			circle := &dot.Circle
			originalSize := circle.Size()
			circle.Resize(fyne.NewSize(originalSize.Width*float32(sizeMultiply), originalSize.Height*float32(sizeMultiply)))
			originalColor := circle.FillColor
			circle.FillColor = glowColor
			circle.Refresh()
			time.Sleep(duration)
			circle.Resize(originalSize)
			circle.FillColor = originalColor
			circle.Refresh()
		}
	}()
}

func sameColor(colorA color.Color, colorB color.Color) bool {
	ra, ga, ba, aa := colorA.RGBA()
	rb, gb, bb, ab := colorB.RGBA()
	if ra == rb && ga == gb && ba == bb && aa == ab {
		return true
	}
	return false
}
