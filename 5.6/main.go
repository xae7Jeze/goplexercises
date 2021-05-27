/*
Exercise 5.6: Modify the corner function in gopl.io/ch3/surface (§3.2)
to use named results and a bare return statement.
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

func corner(i, j int) (xCoord, yCoord float64) {
	xCoord = math.NaN()
	yCoord = math.NaN()
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	if math.IsInf(z, 0) || math.IsNaN(z) {
		return
	}

	xCoord = width/2 + (x-y)*cos30*xyscale
	yCoord = height/2 + (x+y)*sin30*xyscale - z*zscale
	return
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
	return r
}
