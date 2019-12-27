// Package that converts between Celsius, Fahrenheit and Kelvin values
package tempconv

import (
  "fmt"
)

type Celsius float64
type Fahrenheit float64
type Kelvin float64

const(
  AbsZeroC Celsius = -273.15
  FreezingC Celsius = 0
  BoilingC Celsius = 100
  AbsZeroK Kelvin = 0
  FreezingK Kelvin = 273.15
  BoilinigK Kelvin = 373.15
)

func (c Celsius) String() string {return fmt.Sprintf("%.2f°C", c)}
func (f Fahrenheit) String() string {return fmt.Sprintf("%.2f°F", f)}
func (k Kelvin) String() string {return fmt.Sprintf("%.2f°K", k)}
