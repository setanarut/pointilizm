package main

import (
	"fmt"
	"image"
	"math"

	"github.com/cenkalti/dominantcolor"
	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/setanarut/pointilizm/v2/internal"
)

func main() {
	// median := gift.New(gift.Median(7, false))
	img := internal.LoadImage("/Users/haz/Documents/GitHub/Pointilizm/assets/lake.jpg")

	// brush directions
	dir := internal.GetAngles(img)
	dir.MapToRange(0.0, math.Pi*2)

	palette := internal.ToColorfulPalette(dominantcolor.FindN(img, 20))
	temp := palette

	medianImage := image.NewRGBA(img.Bounds())
	// median.Draw(medianImage, img)
	ctx := gg.NewContextForImage(medianImage)

	// plot(dir, img, pos, ctx, palette, 3, 3)

	palette = internal.VaryPalette(temp, -20, -0.2, 0.4)
	plot(dir, img, ctx, palette, 3, 30)

	// palette = internal.VaryPalette(temp, 0, 0, 0.1)
	// plot(dir, img, ctx, palette, 3, 20)

	// palette = internal.VaryPalette(temp, 10, 0.2, 0)
	// plot(dir, img, ctx, palette, 2, 10)

	// palette = internal.VaryPalette(temp, 0, -0.3, 0.1)
	// plot(dir, img, ctx, palette, 2, 10)

	// palette = internal.VaryPalette(temp, 0, 0.20, 0)
	// plot(dir, img, ctx, palette, 3, 20)

	// palette = internal.VaryPalette(temp, -20, 0.5, 0.1)
	// plot(dir, img, ctx, palette, 3, 10)

	// utils.PaletteToImage("assets/palette.png", palette, 30, 4)
	ctx.SavePNG("/Users/haz/Documents/GitHub/Pointilizm/lakeP.png")
}

// g = grid resolution, s = scale
func plot(dir internal.Mat, img image.Image, ctx *gg.Context, palette []colorful.Color, s float64, g int) {
	fmt.Println("drawing pass")
	posX, posY := 0.0, 0.0
	for y := 0; y < img.Bounds().Max.Y; y += g {
		for x := 0; x < img.Bounds().Max.X; x += g {
			posX, posY = float64(x), float64(y)
			posX += internal.RandRange(-10, 10)
			posY += internal.RandRange(-10, 10)
			clr := internal.NearestColor(img.At(int(posX), int(posY)), palette)
			rclr := internal.VaryColor(clr, internal.RandRange(-10.0, 10.0), internal.RandRange(-0.5, 0.5), 0)
			ctx.SetColor(rclr)
			ctx.Push()
			ctx.RotateAbout(dir.At(x, y), posX, posY)
			// brushSize := dir.At(x,y)
			ctx.DrawEllipse(posX, posY, 3, 5)
			ctx.Pop()
			ctx.Fill()
		}
	}
}
