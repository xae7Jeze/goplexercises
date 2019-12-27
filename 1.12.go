/*
Exercise 1.12: Modify the Lissajous server to read parameter values
from the URL. For example, you might arrange it so that a URL like
http://localhost:8000/?cycles=20 sets the number of cycles to 20
instead of the default 5. Use the strconv.Atoi func tion to convert
the string parameter into an integer. You can see its documentation
with go doc strconv.Atoi .
*/
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

type param struct {
	cycles  int
	res     float64
	size    int
	nframes int
	delay   int
}

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
	p := param{
		cycles: 5, res: 0.001, size: 100, nframes: 64,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("oops: parsing url failed")
			return
		}
		if r.Method != "GET" {
			log.Print("answering only to get")
			return
		}
		log.Printf("%+v", p)
		for k, v := range r.Form {
			switch k {
			case "cycles":
				if f, err := strconv.Atoi(v[len(v)-1]); err == nil {
					p.cycles = f
				}
			case "res":
				if f, err := strconv.ParseFloat(v[len(v)-1], 64); err == nil {
					p.res = f
				}
			case "size":
				if f, err := strconv.Atoi(v[len(v)-1]); err == nil {
					p.size = f
				}
			case "nframes":
				if f, err := strconv.Atoi(v[len(v)-1]); err == nil {
					p.nframes = f
				}
			default:
				log.Printf("Ignoring unknown parameter: '%s'", k)
			}
		}
		lissajous(w, p)
		log.Printf("%+v DONE", p)
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func lissajous(out io.Writer, p param) {
	freq := rand.Float64() * 3.0          // relative frequency of y-oscillation
	anim := gif.GIF{LoopCount: p.nframes} // gif image object
	phase := 0.0                          // phase shift between x and y oscillations
	for i := 0; i < p.nframes; i++ {
		rect := image.Rect(0, 0, 2*p.size+1, 2*p.size+1) // create rectangle
		img := image.NewPaletted(rect, palette)          // create image based on rect and palette
		color_index := uint8(1 + (rand.Int() % (len(palette) - 1)))
		for t := 0.0; t < float64(p.cycles)*2*math.Pi; t += p.res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			size := float64(p.size)
			img.SetColorIndex(int(size+(x*size+0.5)), int(size+(y*size+0.5)), color_index)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, p.delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
