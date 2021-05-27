/*
Exercise 3.4: Following the approach of the Lissajous example in Section 1.7,
cons truct a webserver that computes surfaces and writes SVG data to the client.
The server must set the Con-tent-Type header like this:
w.Header().Set("Content-Type", "image/svg+xml") (This step was not required in
the Lissajous example because the ser ver uses standard heuristics to recognize
common formats like PNG from the first 512 bytes of the response and generates
the proper header.) Allow the client to specify values like height, width, and
color as HTTP request parameters.
*/
package main

import (
	"fmt"
	"golang.org/x/image/colornames"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
	color         = "white"
	stroke        = "grey"
)

type param struct {
	cells   int
	width   int
	height  int
	color   string
	stroke  string
	xyrange float64
}

var sin30, cos30 = math.Sin(angle), math.Cos(angle)
var colors = [...]string{"red", "green", "blue", "yellow"}

func main() {
	m := colornames.Map
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("oops: parsing url failed")
			return
		}
		if r.Method != "GET" {
			log.Print("answering only to get")
			return
		}
		p := param{
			cells: cells, width: width, height: height, color: color,
			xyrange: xyrange, stroke: stroke,
		}
		for k, v := range r.Form {
			switch k {
			case "width":
				if f, err := strconv.Atoi(v[len(v)-1]); err == nil {
					p.width = f
				}
			case "height":
				if f, err := strconv.Atoi(v[len(v)-1]); err == nil {
					p.height = f
				}
			case "cells":
				if f, err := strconv.Atoi(v[len(v)-1]); err == nil {
					p.cells = f
				}
			case "color":
				if _, ok := m[v[len(v)-1]]; ok == true {
					p.color = v[len(v)-1]
				}
			case "stroke":
				if _, ok := m[v[len(v)-1]]; ok == true {
					p.stroke = v[len(v)-1]
				}
			case "xyrange":
				if f, err := strconv.ParseFloat(v[len(v)-1], 64); err == nil {
					p.xyrange = f
				}

			default:
				log.Printf("Ignoring unknown parameter: '%s'", k)
			}
		}
		log.Printf("%+v", p)
		w.Header().Set("Content-Type", "image/svg+xml")
		gensvg(w, p)
		log.Printf("%+v DONE", p)
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func gensvg(out io.Writer, p param) {
	fmt.Fprintf(out,
		"<svg xmlns='http://www.w3.org/2000/svg' "+
			"style='stroke: %s; fill: %s; stroke-width: 0.7' "+
			"width='%d' height='%d'>\n", p.stroke, p.color, p.width, p.height)
	for i := 0; i < p.cells; i++ {
		for j := 0; j < p.cells; j++ {
			ax, ay := corner(i+1, j, p)
			bx, by := corner(i, j, p)
			cx, cy := corner(i, j+1, p)
			dx, dy := corner(i+1, j+1, p)
			valid := true
			for _, v := range [...]float64{ax, ay, bx, by, cx, cy, dx, dy} {
				if math.IsNaN(v) {
					valid = false
					break
				}
			}
			if valid {
				fmt.Fprintf(out, "<polygon points = '%g,%g,%g,%g,%g,%g,%g,%g'/>\n",
					ax, ay, bx, by, cx, cy, dx, dy)
			}
		}
	}
	fmt.Fprintf(out, "</svg>\n")
}

func corner(i, j int, p param) (float64, float64) {
	xyscale := float64(p.width) / 2 / p.xyrange
	zscale := float64(p.height) * 0.4
	x := p.xyrange * (float64(i)/cells - 0.5)
	y := p.xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	if math.IsInf(z, 0) || math.IsNaN(z) {
		return math.NaN(), math.NaN()
	}

	sx := float64(p.width)/2 + (x-y)*cos30*xyscale
	sy := float64(p.height)/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Cos(r) / r
}
