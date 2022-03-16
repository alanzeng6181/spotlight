package ui
import (
	"image/color"
	"image/draw"
	"image/png"
	"github.com/alanzeng6181/proximity_server/algorithm"
	"image"
	"os"
	"math/rand"
	"strconv"
)
type QuadTree algorithm.QuadTree[algorithm.Location]
func (tree QuadTree) Draw(fileName string){
	myimage := image.NewRGBA(image.Rect(0, 0, 10000, 10000)) // x1,y1,  x2,y2 of background rectangle
	queue := make([]*algorithm.QuadNode[algorithm.Location], 0)
	queue = append(queue, tree.Root)
	for len(queue)>0{
		node := queue[0]
		queue = queue[1:]
		const thickness = 10
		//two horizontal edges of the rect
		for x:=(int)(node.LowerLeft[0]); x<(int)(node.UpperRight[0]); x++{
			for i:=-1*thickness; i<thickness; i++ {
				myimage.Set(x, (int)(node.LowerLeft[1])+i, color.Black)
			}

			for i:=-1*thickness; i<thickness; i++ {
				myimage.Set(x, (int)(node.UpperRight[1])+i, color.Black)
			}
		}

		//two vertical lines of the rect
		for y:=(int)(node.LowerLeft[1]); y<(int)(node.UpperRight[1]); y++{
			for i:=-1*thickness; i<thickness; i++ {
				myimage.Set((int)(node.LowerLeft[0])+i, y, color.Black)
			}

			for i:=-1*thickness; i<thickness; i++ {
				myimage.Set((int)(node.UpperRight[0])+i, y, color.Black)
			}
		}

		for _, location := range node.Values{
			const dotSize=20
			dot := image.Rect((int)(location.X())-dotSize,(int)(location.Y())-dotSize, (int)(location.X())+dotSize, (int)(location.Y())+dotSize)
			draw.Draw(myimage, dot, &image.Uniform{color.Black}, image.ZP, draw.Src)
		}
		for _, c := range node.Children{
			if c!=nil{
				queue=append(queue, c)
			}
		}
	}

	myfile, err := os.Create(fileName)     // ... now lets save imag
	if err != nil {
		panic(err)
	}
	defer myfile.Close()
	png.Encode(myfile, myimage)
}

func generateAndDrawQuadTree() {
	tree := algorithm.MakeQuadTree[algorithm.Location](0.0, 0.0, 10000.0, 10000.0)
	count := 600
	if len(os.Args) > 1 {
		if i, err := strconv.Atoi(os.Args[1]); err == nil {
			count = i
		}
	}
	for i := 0; i < count; i++ {
		tree.Add(&algorithm.Location{algorithm.Float(rand.Int31n(1000000)) / 100.0, algorithm.Float((rand.Int31n(1000000)) / 100.0)})
	}
	(QuadTree)(tree).Draw("./bin/tree.png")
}
