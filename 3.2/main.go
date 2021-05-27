/*
Exercise 3.2: Experiment with visualizations of other functions from the
math package. Can you produce an egg box, moguls, or a saddle?
*/
package main

import (
	"fmt"
	"math"
)

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)
var colors = [...]string{"red", "green", "blue", "yellow"}

func main() {
	fmt.Printf(
		"<svg xmlns='http://www.w3.org/2000/svg' "+
			"style='stroke: grey; fill: white; stroke-width: 0.7' "+
			"width='%d' height='%d'>\n", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			valid := true
			for _, v := range [...]float64{ax, ay, bx, by, cx, cy, dx, dy} {
				if math.IsNaN(v) {
					valid = false
					break
				}
			}
			if valid {
				fmt.Printf("<polygon points = '%g,%g,%g,%g,%g,%g,%g,%g'/>\n",
					ax, ay, bx, by, cx, cy, dx, dy)
			}
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int) (float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	if math.IsInf(z, 0) || math.IsNaN(z) {
		return math.NaN(), math.NaN()
	}

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

// Taken from https://www.mathcurve.com/surfaces.gb/boiteaoeufs/boiteaoeufs.shtml
func f(x, y float64) float64 {
	var a float64 = 0.1
	var b float64 = 3
	return a * ((math.Sin(x) / b) + (math.Sin(y) / b))
}
