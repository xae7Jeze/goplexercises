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

// BigFloat-Solution

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
	"math/big"
	"math/cmplx"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	width, height = 1024, 1024
	limit         = 2
)

var p = palette.WebSafe

//var p = rand_palette(iterations)
var ps = uint32(len(p))
var iterations uint32 = ps

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
	r *big.Float
	i *big.Float
}

func NewComplexBf() *complexbf {
	z := new(complexbf)
	z.r = big.NewFloat(0)
	z.i = big.NewFloat(0)
	return z
}

func AddComplexBf(s1, s2 *complexbf) *complexbf {
	z := NewComplexBf()
	z.i.Add(s1.i, s2.i)
	z.r.Add(s1.r, s2.r)
	return z
}

func AbsComplexBf(z *complexbf) *big.Float {
	r := big.NewFloat(0)
	a := big.NewFloat(0)
	b := big.NewFloat(0)
	a.Copy(z.r)
	b.Copy(z.i)
	a.Mul(a, z.r)
	b.Mul(b, z.i)
	a.Add(a, b)
	r.Sqrt(a)
	return r
}
func MulComplexBf(s1, s2 *complexbf) *complexbf {
	// REAL = s1rXs2r - s1iXs2i
	// IMG  = s1rXs2i + s2rXs1i
	/*
	     0-2i * 0-2i
	     REAL = 0 * 0 - (-2 * -2) = 0 - 4
	   	IMG  = 0
	*/
	z := NewComplexBf()
	s1rXs2r := big.NewFloat(0)
	s1iXs2i := big.NewFloat(0)
	s1rXs2i := big.NewFloat(0)
	s2rXs1i := big.NewFloat(0)

	s1rXs2r.Mul(s1.r, s2.r)
	s1iXs2i.Mul(s1.i, s2.i)
	s1iXs2i.Neg(s1iXs2i)
	z.r.Add(s1rXs2r, s1iXs2i)

	s1rXs2i.Mul(s1.r, s2.i)
	s2rXs1i.Mul(s2.r, s1.i)
	z.i.Add(s1rXs2i, s2rXs1i)
	return z
}

func main() {
	var xmin, ymin, xmax, ymax float64
	var input [3]float64
	// find 4 random solutions for
	// map covering the 4 possible solutions with
	rand.Seed(time.Now().UnixNano())
	for k, _ := range cm {
		cm[k] = p[uint32(rand.Intn(int(ps)))]
	}
	if len(os.Args) != 5 {
		usage()
		os.Exit(1)
	}
	var mode rune
	if os.Args[1] == "m" {
		mode = 'm'
	} else if os.Args[1] == "n" {
		mode = 'n'
	} else {
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
	z := NewComplexBf()
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z.r.SetFloat64(x)
			z.i.SetFloat64(y)
			if mode == 'm' {
				img.Set(px, height-py, mandelbrot(z))
			} else {
				fmt.Fprintf(os.Stderr, "%s: Sorry, but actually only mode 'm' is implemented\n", os.Args[0])
				os.Exit(3)
			}
		}
	}
	png.Encode(os.Stdout, img)
}

func mandelbrot(z *complexbf) color.Color {
	v := NewComplexBf()
	lim := big.NewFloat(limit)
	for n := uint32(0); n < iterations; n++ {
		v = MulComplexBf(v, v)
		v = AddComplexBf(v, z)
		if AbsComplexBf(v).Cmp(lim) > 0 {
			return p[n%ps]
		}
	}
	return color.Black
}

func rand_color() color.RGBA {
	var r, g, b uint8
	var v *uint8
	rand.Seed(time.Now().UnixNano())
	s := uint32(16)
	ri := uint32(rand.Intn(0x1000000))
	mask := uint32(0xff0000)
	for _, v = range [...]*uint8{&r, &g, &b} {
		*v = uint8((ri & mask) >> s)
		s -= 8
		mask >>= 8
	}
	return color.RGBA{r, g, b, 0xff}
}

func rand_palette(size uint32) []color.RGBA {
	var r, g, b uint8
	var v *uint8
	var p = make([]color.RGBA, size)
	for i, _ := range p {
		rand.Seed(time.Now().UnixNano())
		s := uint32(16)
		ri := uint32(rand.Intn(0x1000000))
		mask := uint32(0xff0000)
		for _, v = range [...]*uint8{&r, &g, &b} {
			*v = uint8((ri & mask) >> s)
			s -= 8
			mask >>= 8
		}
		p[i] = color.RGBA{r, g, b, 0xff}
	}
	return p
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <type> X Y Zoom\n"+
		"<type can be 'm' (mandelbrot) or 'n' (newton method)*\n"+
		"X,Y describe the bottom left corner\n"+
		"\n* Not implemented yet",
		os.Args[0])
}
