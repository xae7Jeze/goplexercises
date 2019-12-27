package tempconv

func CtoF(c Celsius) Fahrenheit {return Fahrenheit(c * 9 / 5 + 32)}
func FtoC(f Fahrenheit) Celsius {return Celsius((f - 32) * 5 / 9)}
func KtoC(k Kelvin) Celsius {return Celsius(k) + AbsZeroC}
func CtoK(c Celsius) Kelvin {return Kelvin(c - AbsZeroC)}
func KtoF(k Kelvin) Fahrenheit {return CtoF(Celsius((k)) + AbsZeroC)}
func FtoK(f Fahrenheit) Kelvin {return Kelvin(FtoC(f) - AbsZeroC)}
