/*
Exercise 3.8: Rendering fractals at high zoom levels demands great
arithmetic precision. Implement the same fractal using four different
represent ation s of numbers: complex64 , complex128, big.Float, and
big.Rat. (The latter two typ es are found in the math/big package.
Float uses arbit rary but bounded-precision floating-point; Rat uses
unbounded-precision rational numbers.) How do they compare in per
formance and memory usage? At what zoom levels do rendering artifacts
become visible?
*/

// BigRat-Solution

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
	"math"
	"math/big"
	"os"
	"strconv"
)

const (
	width, height = 1024, 1024
	limit         = 2
)

var p = palette.WebSafe

var ps = uint32(len(p))
var iterations uint32 = 50

var (
	yellow = color.RGBA{0xFF, 0xFF, 0x00, 0xff}
	red    = color.RGBA{0xFF, 0x00, 0x00, 0xff}
	green  = color.RGBA{0x00, 0xFF, 0x00, 0xff}
	blue   = color.RGBA{0x00, 0x00, 0xFF, 0xff}
	black  = color.RGBA{0x00, 0x00, 0x00, 0xff}
)

var cm = map[complex128]color.Color{
	complex(0, 1):  color.Black,
	complex(1, 0):  color.Black,
	complex(0, -1): color.Black,
	complex(-1, 0): color.Black,
}

type complexbf struct {
	r *big.Rat
	i *big.Rat
}

func NewComplexBf(r, i float64) *complexbf {
	z := new(complexbf)
	z.r = new(big.Rat).SetFloat64(r)
	z.i = new(big.Rat).SetFloat64(i)
	return z
}

func FreeComplexBf(z *complexbf) {
	z.r = nil
	z.i = nil
	z = nil
}

func CopyComplexBf(s *complexbf) *complexbf {
	d := NewComplexBf(0, 0)
	d.i.Set(s.i)
	d.r.Set(s.r)
	return d
}

func AddComplexBf(s1, s2 *complexbf) *complexbf {
	z := NewComplexBf(0, 0)
	z.i.Add(s1.i, s2.i)
	z.r.Add(s1.r, s2.r)
	return z
}

func AbsComplexBf(z *complexbf) float64 {
	a := big.NewRat(0, 1)
	b := big.NewRat(0, 1)
	a.Set(z.r)
	b.Set(z.i)
	a.Mul(a, z.r)
	b.Mul(b, z.i)
	a.Add(a, b)
	r, _ := a.Float64()
	return math.Sqrt(r)
}
func MulComplexBf(s1, s2 *complexbf) *complexbf {
	// REAL = s1rXs2r - s1iXs2i
	// IMG  = s1rXs2i + s2rXs1i
	z := NewComplexBf(0, 0)
	s1r_x_s2r := big.NewRat(0, 1)
	s1i_x_s2i := big.NewRat(0, 1)
	s1r_x_s2i := big.NewRat(0, 1)
	s2r_x_s1i := big.NewRat(0, 1)

	s1r_x_s2r.Mul(s1.r, s2.r)
	s1i_x_s2i.Mul(s1.i, s2.i)
	s1i_x_s2i.Neg(s1i_x_s2i)
	z.r.Add(s1r_x_s2r, s1i_x_s2i)

	s1r_x_s2i.Mul(s1.r, s2.i)
	s2r_x_s1i.Mul(s2.r, s1.i)
	z.i.Add(s1r_x_s2i, s2r_x_s1i)
	return z
}

func main() {
	var xmin, ymin, xmax, ymax float64
	var input [3]float64
	if len(os.Args) != 5 {
		usage()
		os.Exit(1)
	}
	if os.Args[1] != "m" {
		fmt.Fprintf(os.Stderr, "%s: Invalid  input: '%s'\n", os.Args[0], os.Args[1])
		os.Exit(2)
	}

	for i, arg := range os.Args[2:] {
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
	zoom := input[2]
	if zoom <= 0 {
		fmt.Fprintf(os.Stderr, "%s: Invalid  input: Zoom: '%s'\n", os.Args[0], zoom)
		os.Exit(2)
	}
	zrange := (2 * limit) / zoom
	xmin = input[0] - (zrange / 2)
	if xmin < -1*limit {
		xmin = -1 * limit
	}
	xmax = input[0] + (zrange / 2)
	if xmax > limit {
		xmax = limit
	}
	ymin = input[1] - (zrange / 2)
	if ymin < -1*limit {
		ymin = -1 * limit
	}
	ymax = input[1] + (zrange / 2)
	if ymax > limit {
		ymax = limit
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := int64(0); py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := int64(0); px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := NewComplexBf(x, y)
			img.Set(int(px), int(height-py), mandelbrot(z))
		}
	}
	png.Encode(os.Stdout, img)
}

func mandelbrot(z *complexbf) color.Color {
	V := NewComplexBf(0, 0)
	Z := NewComplexBf(0, 0)
	Z.r.Set(z.r)
	Z.i.Set(z.i)
	for n := uint32(0); n < iterations; n++ {
		TMP := CopyComplexBf(V)
		FreeComplexBf(V)
		V = NewComplexBf(0, 0)
		V = MulComplexBf(TMP, TMP)
		V = AddComplexBf(V, Z)
		if AbsComplexBf(V) > limit {
			return p[n%ps]
		}
	}
	return color.Black
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <type> X Y Zoom\n"+
		"<type only can be 'm' (mandelbrot)*\n"+
		"X,Y describing the center of the picture\n"+
		"\n* Not implemented yet",
		os.Args[0])
}
