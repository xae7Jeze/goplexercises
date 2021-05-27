/*
Exercise 3.6: Supersamp ling is a technique to reduce the effect of
pixelation by computing the color value at several points within each
pixel and taking the average. The simplest method is to divide each
pixel into four "subpixels" Implement it.
*/
package main

import (
	"image"
	//"math"
	//"image/color"
	"fmt"
	"image/color/palette"
	"image/png"
	"math/cmplx"
	"os"
	"strconv"
)

const (
	width, height = 1024, 1024
	limit         = 2
)

var p = palette.WebSafe
var ps = uint32(len(p))
var iterations uint32 = 200
var xmin, ymin, xmax, ymax float64

func main() {
	var input [3]float64
	if len(os.Args) != 4 {
		usage()
		os.Exit(1)
	}
	for i, arg := range os.Args[1:] {
		if u, err := strconv.ParseFloat(arg, 64); err != nil {
			fmt.Fprintf(os.Stderr, "%s: Invalid  input: '%s'\n", os.Args[0], arg)
			os.Exit(2)
		} else {
			input[i] = u
		}
	}
	if input[2] <= 0 {
		usage()
		os.Exit(2)
	}
	xmin, xmax, ymin, ymax =
		input[0], input[0]+input[2], input[1], input[1]+input[2]
	for _, u := range [...]float64{xmin, xmax, ymin, ymax} {
		if u < limit*-1 || u > limit {
			fmt.Fprintf(os.Stderr, "%s: Invalid input\n", os.Args[0])
			os.Exit(2)
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			// fmt.Fprintf(os.Stderr, "y: %.2f x: %.2f -> py: %v px -> %v\n", y, x, py, px)
			img.Set(px, py, p[oversample(uint(px), uint(py), 2)%ps])
		}
	}
	png.Encode(os.Stdout, img)
}

// returns zero if part of mandelbrot set, iterations to leave it otherwise
func mandelbrot(z complex128) uint32 {
	var v complex128
	for n := uint32(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > limit {
			return n
		}
	}
	return 0
}

// oversamples by dividing pixel in fac * fac subpixels
func oversample(px, py, fac uint) uint32 {
	v := float64(0)
	for dx := fac * px; dx < (fac*px + fac); dx++ {
		for dy := fac * py; dy < (fac*py + fac); dy++ {
			y := -1 * ((float64(dy)/float64(height*fac))*(ymax-ymin) + ymin)
			x := (float64(dx)/float64(width*fac))*(xmax-xmin) + xmin
			z := complex(x, y)
			v += float64(mandelbrot(z))
		}
	}
	return uint32(v / float64(fac*fac))
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s X Y DIM\n"+
		"X,Y describe the bottom left corner, DIM is the square's sidelength\n",
		os.Args[0])
}
