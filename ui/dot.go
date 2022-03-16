package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/alanzeng6181/proximity_server/algorithm"
)

type Dot canvas.Circle

func NewDot(color color.Color, radius float32, position fyne.Position) *Dot {
	circle := canvas.NewCircle(color)
	circle.Hidden = false
	circle.StrokeColor = color
	circle.FillColor = color
	circle.Move(position)
	circle.Resize(fyne.NewSize(radius, radius))
	return (*Dot)(circle)
}

func (dot Dot) X() algorithm.Float {
	return algorithm.Float((dot.Position1.X + dot.Position2.X) / 2.0)
}

func (dot Dot) Y() algorithm.Float {
	return algorithm.Float((dot.Position1.Y + dot.Position2.Y) / 2.0)
}

func (dot *Dot) Glow(sizeMultiply float32, duration time.Duration, glowColor color.Color) {
	go func() {
		originalSize := (*canvas.Circle)(dot).Size()
		originalColor := dot.FillColor
		(*canvas.Circle)(dot).Resize(fyne.NewSize(originalSize.Width*sizeMultiply, originalSize.Height*sizeMultiply))
		dot.FillColor = glowColor
		(*canvas.Circle)(dot).Refresh()
		time.Sleep(duration)
		(*canvas.Circle)(dot).Resize(originalSize)
		dot.FillColor = originalColor
		(*canvas.Circle)(dot).Refresh()
	}()
}
