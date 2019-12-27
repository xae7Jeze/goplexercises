package popcount

var pc [256]byte

func init() {
	for i, _ := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

/* BookVersion */
func PopCount(x uint64) int {
	return int(
		pc[byte(x>>(0*8))] +
			pc[byte(x>>(1*8))] +
			pc[byte(x>>(2*8))] +
			pc[byte(x>>(3*8))] +
			pc[byte(x>>(4*8))] +
			pc[byte(x>>(5*8))] +
			pc[byte(x>>(6*8))] +
			pc[byte(x>>(7*8))])
}

/* Exercise 2.3: Lookup pc with Loop-Variable */
func PopCountL(x uint64) int {
	var cnt byte
	for i := uint(0); i < 8; i++ {
		cnt += pc[byte(x>>(i*8))]
	}
	return int(cnt)
}

/* Exercise 2.4: shift right and compare LSB */
func PopCountS(x uint64) int {
	var cnt byte
	for ; x != 0; x >>= 1 {
		cnt += byte(x & 1)
	}
	return int(cnt)
}

/*
Exercise 2.5: Use the fact, that the expression x&(x-1) clears
the rightmost non-zero bit of x
*/
func PopCountR(x uint64) int {
	var cnt byte
	for ; x != 0; x = x & (x - 1) {
		cnt++
	}
	return int(cnt)
}
