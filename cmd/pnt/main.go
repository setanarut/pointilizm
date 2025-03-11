package main

import (
	"fmt"

	"github.com/cenkalti/dominantcolor"
	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
	it "github.com/setanarut/pointilizm/v2/internal"
)

var ctx *gg.Context

func main() {
	opts := &Options{
		Colors:            20,
		GridSize:          5,
		BrushRadiusX:      3,
		BrushRadiusY:      5,
		RandomizeStokeMin: -10,
		RandomizeStokeMax: 10,
	}
	inputImage := it.LoadImage("./lake.jpg")
	ctx = gg.NewContextForImage(inputImage)

	brushAngleGrid := it.GetAngles(inputImage, 30)

	palette := it.ToColorfulPalette(dominantcolor.FindN(inputImage, opts.Colors))

	paint(palette, brushAngleGrid, opts)
	paint(
		it.VaryPalette(palette, it.RandRange(-10, 10), it.RandRange(-0.2, 0.2), 0),
		brushAngleGrid,
		opts,
	)

	// it.PaletteToImage("./palette.png", palette, 30, 4)
	ctx.SavePNG("./out.png")
}

func paint(palette []colorful.Color, brushAngleGrid [][]float64, o *Options) {
	fmt.Println("drawing pass")
	posX, posY := 0.0, 0.0
	for y := 0; y < ctx.Height(); y += o.GridSize {
		for x := 0; x < ctx.Width(); x += o.GridSize {
			posX, posY = float64(x), float64(y)
			posX += it.RandRange(o.RandomizeStokeMin, o.RandomizeStokeMax)
			posY += it.RandRange(o.RandomizeStokeMin, o.RandomizeStokeMax)
			clr := it.NearestColor(ctx.Image().At(int(posX), int(posY)), palette)
			rclr := it.VaryColor(clr, it.RandRange(-10.0, 10.0), it.RandRange(-0.5, 0.5), 0)
			ctx.SetColor(rclr)
			ctx.Push()
			ctx.RotateAbout(brushAngleGrid[y][x], posX, posY)
			ctx.DrawEllipse(posX, posY, o.BrushRadiusX, o.BrushRadiusY)
			ctx.Pop()
			ctx.Fill()
		}
	}
}

type Options struct {
	Colors            int
	GridSize          int
	BrushRadiusX      float64
	BrushRadiusY      float64
	RandomizeStokeMin float64
	RandomizeStokeMax float64
}
