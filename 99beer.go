package main

func main() {
	var i int = 99
	var bof string = " bottles of "
	var bofs string = " bottle of "
	var td string = "Take one down, pass it around, "
	var otw string = " on the wall"
	for ; i > 1; i-- {
		Print(i, bof, "beer", otw, ", ", i, bof, "beer", ".\n",
			td, i, bof, "beer", otw, ".\n\n")
	}
	Print(i, bofs, "beer", otw, ", ", i, bofs, "beer", ".\n",
		td, "no more", bofs, "beer", otw, ".\n")
}
