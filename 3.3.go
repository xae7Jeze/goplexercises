/*
Exercise 3.3: Color each polygon based on its height, so that the peaks
are colored red ( #ff0000 ) and the valleys blue ( #0000ff ).
*/
package main

import (
	"fmt"
	"math"
	//  "os"
)

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.2
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
			ax, ay, ac := corner(i+1, j)
			bx, by, bc := corner(i, j)
			cx, cy, cc := corner(i, j+1)
			dx, dy, dc := corner(i+1, j+1)
			color := "#FF00FF"
			if ac == bc && bc == cc && cc == dc {
				color = ac
			}
			valid := true
			for _, v := range [...]float64{ax, ay, bx, by, cx, cy, dx, dy} {
				if math.IsNaN(v) {
					valid = false
					break
				}
			}
			if valid {
				fmt.Printf("<polygon points = '%g,%g,%g,%g,%g,%g,%g,%g' "+
					"style='fill: %s'/>\n",
					ax, ay, bx, by, cx, cy, dx, dy, color)
			}
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int) (float64, float64, string) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	color := "#0000ff"
	if z >= 0 {
		color = "#ff0000"
	}
	//fmt.Fprintf(os.OpenFile(os.DevNull, O_WRONLY, 0666), "%g\n",z);
	if math.IsInf(z, 0) || math.IsNaN(z) {
		return math.NaN(), math.NaN(), ""
	}

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, color
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}
