/*
Exercise 3.5: Implement a full-color Mandelbrot set using the function
image.NewRGBA and the type color.RGBA or color.YCbCr .
*/
package main

import (
	"image"
	//"math"
	"fmt"
	"image/color"
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

func main() {
	var xmin, ymin, xmax, ymax float64
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
			// fmt.Fprintf(os.Stderr, "%s: Invalid  input:\n", os.Args[0])
			os.Exit(2)
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// fmt.Fprintf(os.Stderr, "y: %.2f x: %.2f -> py: %v px -> %v\n", y, x, py, px)
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(os.Stdout, img)
}

func mandelbrot(z complex128) color.Color {
	var v complex128
	for n := uint32(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > limit {
			return p[n%ps]
		}
	}
	return color.RGBA{0x00, 0x00, 0x00, 0xff}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s X Y DIM\n"+
		"X,Y describe the bottom left corner, DIM is the square's sidelength\n",
		os.Args[0])
}
