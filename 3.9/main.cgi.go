/*
Exercise 3.9: Write a webserver that renders fractals and writes the image
data to the client. Allow the client to specify the x, y, and zoom values
as parameters to the HTTP request.
*/

// CGI-Version

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"math/rand"
	"net/http"
	"net/http/cgi"
	"os"
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
var rand_pal = rand_palette(ps)
var pal = rand_pal
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
	//http.HandleFunc("/", http_handler)
	var req *http.Request
	var err error
	req, err = cgi.Request()
	if err != nil {
		errorResponse(500, "cannot get cgi request"+err.Error())
		return
	}
	cgi_handler(os.Stdout, req)
}

func errorResponse(code int, msg string) {
	fmt.Printf("Status:%d %s\r\n", code, msg)
	fmt.Printf("Content-Type: text/plain\r\n")
	fmt.Printf("\r\n")
	fmt.Printf("%s\r\n", msg)
}

func genfract(out io.Writer, p param) {
	var xmin, ymin, xmax, ymax float64

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
			img.Set(px, height-py, mandelbrot(z))
		}
	}
	png.Encode(out, img)
}

/*
*
* HTTP Request-Handler
*
 */
func cgi_handler(w io.Writer, r *http.Request) {
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
	fmt.Fprintf(w, "Content-Type: image/png\r\n\r\n")
	pal = rand_pal
	genfract(w, p)
	log.Printf("%+v DONE", p)
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
