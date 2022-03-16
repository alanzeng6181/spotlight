package main

import (
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/alanzeng6181/proximity_server/algorithm"
	"github.com/alanzeng6181/proximity_server/ui"
)

func main() {
	dotsSimulation()
}

func dotsSimulation() {
	myApp := app.New()
	window := myApp.NewWindow("Dots")
	const width = 1500.0
	const height = 1500.0

	ct := fyne.NewContainerWithoutLayout()
	ct.Resize(fyne.NewSize(width, height))
	tree := algorithm.MakeQuadTree[ui.Dot](0.0, 0.0, width, height)
	count := 10000
	if len(os.Args) > 1 {
		if i, err := strconv.Atoi(os.Args[1]); err == nil {
			count = i
		}
	}
	dots := make([]*ui.Dot, 0)
	go func() {
		for i := 0; i < count; i++ {
			time.Sleep(5 * time.Millisecond)
			dot := ui.NewDot(color.White, 5.0, fyne.NewPos(float32(rand.Int31n(width*99))/100.0, float32(rand.Int31n(height*99))/100.0))
			if tree.Add(dot) {
				ct.Add((*canvas.Circle)(dot))
				dots = append(dots, dot)
				if i%30 == 0 {
					ct.Refresh()
				}
			} else {
				log.Fatalf("Couldn't add dot at x:%f, y:%f", dot.X(), dot.Y())
			}
		}
		ct.Refresh()
	}()

	go func() {
		for i := 0; i < 50; i++ {
			time.Sleep(3 * time.Second)
			n := len(dots)
			effectRadius := float32(200.0)
			if n > 0 {
				red := color.RGBA{255, 0, 0, 1}
				picked := dots[rand.Int31n(int32(n))]
				picked.Glow(5.0, 2*time.Second, red)
				surroundings := tree.FindNearby((picked.Position1.X+picked.Position2.X)/2,
					(picked.Position1.Y+picked.Position2.Y)/2,
					effectRadius)

				for _, dot := range surroundings {
					if dot == picked {
						continue
					}
					dot.Glow(1.5, 2*time.Second, red)
				}
			}
		}
		if ok, err := tree.Root.Verify(count); !ok {
			log.Fatalf("tree verification failed...%v", err)
		} else {
			log.Println("tree verification passed....")
		}
	}()

	window.SetContent(ct)
	window.Resize(fyne.NewSize(width, height))
	window.ShowAndRun()
}
