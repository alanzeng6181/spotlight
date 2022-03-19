package main

import (
	"flag"
	"image/color"
	"log"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/alanzeng6181/proximity_server/algorithm"
	"github.com/alanzeng6181/proximity_server/ui"
)

func main() {
	width := flag.Float64("w", 1900, "width of screen, between 100 and 5000")
	height := flag.Float64("h", 1000, "height of screen, between 100 and 5000")
	count := flag.Int("c", 10000, "number of dots to populate, between 1 and 500000")
	glowCenter := flag.Float64("gc", 5.0, "center dot enlargement, between 1.0 and 10.0")
	glowSurrounding := flag.Float64("gs", 2.0, "surrounding dots enlargement, between 1.0 and 10.0")
	durationMS := flag.Int("gd", 2000, "glow duration in miliseconds, bewteen 1000 and 1000")
	effectiveRadius := flag.Float64("r", 100, "effective search radius, between 10 and 400")
	dps := flag.Int("dps", 500, "dots per second, between 1 and 10000")
	flag.Parse()

	if *width < 100 || *width > 5000 {
		*width = 1500
	}

	if *height < 100 || *height > 5000 {
		*height = 1500
	}

	if *count < 0 || *count > 500000 {
		*count = 10000
	}

	if *glowCenter < 1.0 || *glowCenter > 10.0 {
		*glowCenter = 5.0
	}

	if *glowSurrounding < 1.0 || *glowSurrounding > 10.0 {
		*glowSurrounding = 2.0
	}

	if *durationMS < 1000 || *durationMS > 10000 {
		*durationMS = 2000
	}

	if *effectiveRadius < 10 || *effectiveRadius > 400 {
		*effectiveRadius = 100
	}

	if *dps < 1 || *dps > 10000 {
		*dps = 500
	}

	myApp := app.New()
	window := myApp.NewWindow("Spotlight")

	ct := fyne.NewContainerWithoutLayout()

	screenSize := fyne.NewSize(float32(*width), float32(*height))
	ct.Resize(screenSize)
	tree := algorithm.MakeQuadTree[*ui.Dot](0.0, 0.0, algorithm.Float(*width), algorithm.Float(*height))

	dots := make([]*ui.Dot, 0)

	//goroutine for populating dots.
	go func() {
		for i := 0; i < *count; i++ {
			time.Sleep(time.Duration(1000000/(*dps)) * time.Microsecond)
			dot := ui.NewDot(color.White, 5.0, fyne.NewPos(float32(rand.Int31n(int32(*width)*99))/100.0,
				float32(rand.Int31n(int32(*height)*99))/100.0))
			if tree.Add(dot) {
				ct.Add(&dot.Circle)
				dots = append(dots, dot)
				if i%(*dps/4) == 0 {
					ct.Refresh()
				}
			} else {
				log.Fatalf("Couldn't add dot at x:%f, y:%f", dot.X(), dot.Y())
			}
		}
		ct.Refresh()
	}()

	//gorouting for selecting a random dot, highlight it and it's surroudings within a distance to show a spotlight effect.
	go func() {
		duration := time.Duration(*durationMS) * time.Millisecond
		for i := 0; i < 70; i++ {
			time.Sleep(3 * time.Second)
			n := len(dots)

			if n > 0 {
				red := color.RGBA{255, 0, 0, 1}
				picked := dots[rand.Int31n(int32(n))]
				picked.Glow(*glowCenter, duration, red)
				surroundings := tree.FindNearby(float64((picked.Position1.X+picked.Position2.X)/2),
					float64((picked.Position1.Y+picked.Position2.Y)/2),
					*effectiveRadius)

				for _, dot := range surroundings {
					if dot == picked {
						continue
					}
					dot.Glow(*glowSurrounding, duration, red)
				}
			}
		}
	}()

	window.SetContent(ct)
	window.Resize(screenSize)
	window.ShowAndRun()
}
