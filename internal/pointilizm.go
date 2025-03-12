package internal

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand/v2"
	"os"

	"github.com/disintegration/gift"
	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
)

const Pi2 = math.Pi * 2

// GetAngles get vector field angles
func GetAngles(im image.Image, blurSigma float32) [][]float64 {
	angles := MakeGrid(im.Bounds().Size())

	filterScharrX := gift.Convolution(
		[]float32{
			3, 0, -3,
			10, 1, -10,
			3, 0, -3,
		},
		false, false, false, 0,
	)

	filterScharrY := gift.Convolution(
		[]float32{
			3, 10, 3,
			0, 0, 0,
			-3, -10, -3,
		},
		false, false, false, 0,
	)

	// scharr filters
	ScharrXFilter := gift.New(filterScharrX, gift.GaussianBlur(blurSigma))
	ScharrYFilter := gift.New(filterScharrY, gift.GaussianBlur(blurSigma))

	// calculate new bounds
	GrayXImageBounds := ScharrXFilter.Bounds(im.Bounds())
	GrayYImageBounds := ScharrYFilter.Bounds(im.Bounds())

	// new gray images
	GrayXImage := image.NewGray16(GrayXImageBounds)
	GrayYImage := image.NewGray16(GrayYImageBounds)

	// apply Kernels
	ScharrXFilter.Draw(GrayXImage, im)
	ScharrYFilter.Draw(GrayYImage, im)

	for y := range im.Bounds().Dy() {
		for x := range im.Bounds().Dx() {
			// brush angle
			angles[y][x] = math.Atan2(
				float64(GrayYImage.Gray16At(x, y).Y),
				float64(GrayXImage.Gray16At(x, y).Y),
			)
		}
	}

	NormalizeAngles(angles)

	return angles
}

// NormalizeAngles
func NormalizeAngles(arr [][]float64) {
	minVal := arr[0][0]
	maxVal := arr[0][0]
	for y := range arr {
		for x := range arr[y] {
			minVal = min(minVal, arr[y][x])
			maxVal = max(maxVal, arr[y][x])
		}
	}
	for y := range arr {
		for x := range arr[y] {
			arr[y][x] = MapRange(arr[y][x], minVal, maxVal, 0, Pi2)
		}
	}
}

// MapRange map range to another range
func MapRange(v, v1, v2, min, max float64) float64 {
	return min + ((max-min)/(v2-v1))*(v-v1)
}

// RandRange random number with range
func RandRange(min float64, max float64) float64 {
	r := min + rand.Float64()*(max-min)
	return r
}

// NearestColor get nearest color in palette
func NearestColor(target color.Color, p []colorful.Color) colorful.Color {
	index := 0
	clr, _ := colorful.MakeColor(target)
	var minDist float64 = p[0].DistanceLab(clr)
	var distList = make([]float64, 0, len(p))
	for _, c := range p {
		dist := c.DistanceLab(clr)
		distList = append(distList, dist)
		if dist < minDist {
			minDist = dist
		}
	}
	for i, dist := range distList {
		if dist == minDist {
			index = i
		}
	}
	col := p[index]
	return col
}

// VaryColor Adds HSV values to color.
func VaryColor(c colorful.Color, h float64, s float64, v float64) colorful.Color {
	hue, sat, val := c.Hsv()
	return colorful.Hsv(hue+h, sat+s, val+v).Clamped()
}

// VaryPalette Adds HSV values to every color in palette.
func VaryPalette(p []colorful.Color, h float64, s float64, v float64) []colorful.Color {
	pal := make([]colorful.Color, len(p))
	for i := range pal {
		pal[i] = VaryColor(p[i], h, s, v)
	}
	return pal
}

// ToColorfulPalette convert []color.Color or []color.RGBA to []colorful.Color
func ToColorfulPalette(plt []color.RGBA) []colorful.Color {
	var colorfulPalette []colorful.Color
	for _, v := range plt {
		clr, _ := colorful.MakeColor(v)
		colorfulPalette = append(colorfulPalette, clr)
	}
	return colorfulPalette
}

// LoadImage
func LoadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("os.Open failed: %v", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("image.Decode failed: %v", err)
	}
	return img
}

// SaveImage
func SaveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("os.Create failed: %v", err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		log.Fatalf("png.Encode failed: %v", err)
	}
}

func MakeGrid(size image.Point) [][]float64 {
	var tm [][]float64
	for range size.Y {
		tm = append(tm, make([]float64, size.X))
	}
	return tm
}

// paletteToImage saves pallette as image
func PaletteToImage(fileName string, p []colorful.Color, scale int, row int) {
	w := len(p) / row * scale
	h := scale * row
	c := gg.NewContext(len(p)/row*scale, scale*row)
	paletteIndex := 0
	for y := 0; y < h; y += scale {
		for x := 0; x < w; x += scale {
			c.SetColor(p[paletteIndex])
			c.DrawRectangle(float64(x), float64(y), float64(scale), float64(scale))
			c.Fill()
			paletteIndex++
		}
	}
	c.SavePNG(fileName)
}

// // ColorToImage save color as image
// func ColorToImage(fileName string, c colorful.Color) {
// 	ctx := gg.NewContext(30, 30)
// 	ctx.SetColor(c)
// 	ctx.DrawRectangle(0, 0, 30, 30)
// 	ctx.Fill()
// 	ctx.SavePNG(fileName)
// }

// // RegulatePalette vary palette
// func RegulatePalette(p []colorful.Color) []colorful.Color {
// 	temp := p
// 	p = append(p, VaryPalette(temp, 0, -0.3, -0.2)...)
// 	p = append(p, VaryPalette(temp, 20, 0.3, 0)...)
// 	p = append(p, VaryPalette(temp, 0, 1, 1)...)
// 	return p
// }
