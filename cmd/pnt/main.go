package main

import (
	"fmt"
	"image"
	"os"

	"github.com/cenkalti/dominantcolor"
	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
	it "github.com/setanarut/pointilizm/v2/internal"
)

var ctx *gg.Context

func main() {
	args := os.Args

	opts := &Options{
		Colors:                 30,
		GridSize:               4,
		BlurSigma:              30,
		BrushRadiusX:           3,
		BrushRadiusY:           6,
		GridRandomizeOffsetMin: -30,
		GridRandomizeOffsetMax: 30,
	}
	inputImage := it.LoadImage(args[1])

	ctx = gg.NewContextForImage(inputImage)

	brushAngleGrid := it.GetAngles(inputImage, opts.BlurSigma)

	palette := it.ToColorfulPalette(dominantcolor.FindN(inputImage, opts.Colors))
	paint(
		inputImage,
		palette,
		brushAngleGrid,
		opts,
	)
	fmt.Println("Done!")

	// it.PaletteToImage("./palette.png", palette, 30, 4)
	ctx.SavePNG(args[2])
}

func paint(inputImage image.Image, palette []colorful.Color, brushAngleGrid [][]float64, o *Options) {
	fmt.Println("painting...")
	posX, posY := 0.0, 0.0
	for y := 0; y < ctx.Height(); y += o.GridSize {
		for x := 0; x < ctx.Width(); x += o.GridSize {
			posX, posY = float64(x), float64(y)
			posX += it.RandRange(o.GridRandomizeOffsetMin, o.GridRandomizeOffsetMax)
			posY += it.RandRange(o.GridRandomizeOffsetMin, o.GridRandomizeOffsetMax)
			clr := it.NearestColor(inputImage.At(int(posX), int(posY)), palette)
			rclr := it.VaryColor(clr, it.RandRange(-10, 10), it.RandRange(0, 0.5), 0)
			ctx.SetColor(rclr)
			ctx.Push()
			ctx.RotateAbout(brushAngleGrid[y][x], posX, posY)
			r := it.RandRange(-2, 2)
			// r := 0.0
			ctx.DrawEllipse(
				posX,
				posY,
				o.BrushRadiusX+r,
				o.BrushRadiusY+r,
			)
			ctx.Pop()
			ctx.Fill()
		}
	}
}

type Options struct {
	Colors                 int
	GridSize               int
	BlurSigma              float32
	BrushRadiusX           float64
	BrushRadiusY           float64
	GridRandomizeOffsetMin float64
	GridRandomizeOffsetMax float64
}
