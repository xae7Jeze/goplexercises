/*
Exercise 1.6: Modify the Lissajous program to produce images in multiple
colors by adding more values to palette and then displaying them by
changing the third argument of SetColorIndex in some interesting way.
*/
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
	// "fmt"
)

const (
	blackIndex = 0
	redIndex   = 1
	greenIndex = 2
	blueIndex  = 3
)

var (
	black = color.RGBA{0x00, 0x00, 0x00, 0xff}
	red   = color.RGBA{0xFF, 0x00, 0x00, 0xff}
	green = color.RGBA{0x00, 0xFF, 0x00, 0xff}
	blue  = color.RGBA{0x00, 0x00, 0xFF, 0xff}
)

var palette = []color.Color{black, red, green, blue}

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // x-oscillator cycles
		res     = 0.001 // angular resolution
		size    = 100   // size of canvas
		nframes = 64    // number of frames
		delay   = 8     // Delay between frames, unit = 10ms
	)

	freq := rand.Float64() * 2.0        // relative frequency of y-oscillation
	anim := gif.GIF{LoopCount: nframes} // gif image object
	phase := 0.0                        // phase shift between x and y oscillations
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1) // create rectangle
		img := image.NewPaletted(rect, palette)      // create image based on rect and palette
		color_index := uint8(1 + (rand.Int() % (len(palette) - 1)))
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), color_index)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
