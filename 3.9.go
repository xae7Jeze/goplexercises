/*
Exercise 3.9: Write a webserver that renders fractals and writes the image
data to the client. Allow the client to specify the x, y, and zoom values
as parameters to the HTTP request.
*/
package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type param struct {
	x     float64
	y     float64
	zoom  float64
	ftype rune
}

const (
	x, y  = 0.0, 0.0
	zoom  = 1.0
	ftype = 'm'
)

const (
	width, height  = 1024, 1024
	limit          = 2
	def_iterations = 200
)

var ws_pal = palette.WebSafe
var ps = uint32(len(ws_pal))
var pal = ws_pal
var iterations = ps - 1

var (
	yellow = color.RGBA{0xFF, 0xFF, 0x00, 0xff}
	red    = color.RGBA{0xFF, 0x00, 0x00, 0xff}
	green  = color.RGBA{0x00, 0xFF, 0x00, 0xff}
	blue   = color.RGBA{0x00, 0x00, 0xFF, 0xff}
	black  = color.RGBA{0x00, 0x00, 0x00, 0xff}
)

var cm = map[complex128]color.Color{
	complex(0, 1):  red,
	complex(1, 0):  green,
	complex(0, -1): yellow,
	complex(-1, 0): blue,
}

func main() {
	http.HandleFunc("/", http_handler)
	log.Print("Opening Listening port 12345")
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func genfract(out io.Writer, p param) {
	var xmin, ymin, xmax, ymax float64
	var mode rune = 'm'
	if p.ftype == 'm' {
		mode = 'm'
	} else if p.ftype == 'n' {
		mode = 'n'
	}

	if p.zoom <= 0 {
		return
	}
	zrange := (2 * limit) / p.zoom
	xmin = p.x - (zrange / 2)
	if xmin < (-1 * limit) {
		xmin = -1 * limit
	}
	xmax = p.x + (zrange / 2)
	if xmax > limit {
		xmax = limit
	}
	ymin = p.y - (zrange / 2)
	if ymin < -1*limit {
		ymin = -1 * limit
	}
	ymax = p.y + (zrange / 2)
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
	png.Encode(out, img)
}

/*
*
* HTTP Request-Handler
*
 */
func http_handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print("oops: parsing url failed")
		return
	}
	if r.Method != "GET" {
		log.Print("answering only to get")
		return
	}
	p := param{
		x:     x,
		y:     y,
		zoom:  zoom,
		ftype: ftype,
	}
	if r.URL.Path == "/n" {
		p.ftype = 'n'
	} else if r.URL.Path == "/m" {
		p.ftype = 'm'
	}
	for k, v := range r.Form {
		switch k {
		case "x":
			if f, err := strconv.ParseFloat(v[len(v)-1], 64); err == nil {
				p.x = f
			}
		case "y":
			if f, err := strconv.ParseFloat(v[len(v)-1], 64); err == nil {
				p.y = f
			}
		case "zoom":
			if f, err := strconv.ParseFloat(v[len(v)-1], 64); err == nil {
				p.zoom = f
			}
		default:
			log.Printf("Ignoring unknown parameter: '%s'", k)
		}
	}
	log.Printf("%+v", p)
	w.Header().Set("Content-Type", "image/png")
	//pal = rand_palette(ps)
	pal = ws_pal
	genfract(w, p)
	log.Printf("%+v DONE", p)
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

/*
* Generates mandelbrot fractals
 */

func mandelbrot(z complex128) color.Color {
	var v complex128
	for n := uint32(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > limit {
			return pal[n%ps]
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

/*
* Generates random rgba color palette of "size"
*
 */
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

/*
*
* Shades color "c" by subtracting "reduce" from each
* rgb color element
*
 */
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
