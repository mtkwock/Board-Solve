package main

import (
	"fmt"
)

func main() {
	swng_board := CreateBoard("RRHHHDGDGGGDGDDRRRRDRGGDGRRDDG", 6)
	combos := swng_board.GetAllCombos()
	fmt.Println(combos)
}
