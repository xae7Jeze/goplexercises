/*
Exercise 3.7: Another simple fractal uses Newton’s method to find complex
solutions to a function such as z^4 − 1 = 0. Shade each starting point
by the number of iterat ions required to get close to one of the four
roots. Color each point by the root it approaches.
*/
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
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

func main() {
	var xmin, ymin, xmax, ymax float64
	var input [3]float64
	// map covering the 4 possible solutions
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
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			if mode == 'n' {
				img.Set(px, height-py, x4_newton(z))
			} else {
				img.Set(px, height-py, mandelbrot(z))
			}
		}
	}
	png.Encode(os.Stdout, img)
}

/*
 * find complex solutions for Z^4 - 1 = 0 using newton's method,
 * iterating Zn+1 = Zn + (f(Zn)/f´(Zn)) until the absolute value
 * of Zn - Zn+1 falls below limit
 * f(X)  = X^4 - 1
 * f´(X) = 4 * X^3
 * Xn+1 = (Xn - (Xn / 4) + (1 / 4 * Xn^3))
 *
 */
func x4_newton(z complex128) color.Color {
	var z1 complex128
	var i uint32

	limit := 1e-100
	for i = 0; i < iterations; z = z1 {
		z1 = z - (z / 4) + (1 / (4 * z * z * z))
		i++
		if cmplx.Abs(z-z1) < limit {
			break
		}
	}
	if i >= iterations {
		return color.Black
	}
	for k, _ := range cm {
		if cmplx.Abs(k-z1) < limit {
			return shade(cm[k].(color.RGBA), i)
		}
	}
	return color.Black
}

func mandelbrot(z complex128) color.Color {
	var v complex128
	for n := uint32(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > limit {
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

func shade(c color.RGBA, reduce uint32) color.RGBA {
	var rc color.RGBA
	rr := uint8(reduce % 256)
	rc.A = 255

	if c.R >= rr {
		rc.R = c.R - rr
	} else {
		rc.R = 0
	}
	if c.G >= rr {
		rc.G = c.G - rr
	} else {
		rc.G = 0
	}
	if c.B >= rr {
		rc.B = c.B - rr
	} else {
		rc.B = 0
	}
	return rc
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <type> X Y Zoom\n"+
		"<type can be 'm' (mandelbrot) or 'n' (newton method)\n"+
		"X,Y describe the bottom left corner\n",
		os.Args[0])
}
