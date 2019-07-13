package main

import (
	"flag"
	"fmt"
)

var board_flag string
var board_width int
var board_to_solve Board
var combo_flag int
var diagonals_flag bool

func init() {
	flag.StringVar(&board_flag, "board", "",
			"Board String (R)ed, (B)lue, (G)reen, (L)ight, (D)ark, (H)eart, (P)oison), (M)ortal Poison, (J)ammer, B(o)mb.")
	flag.IntVar(&board_width, "width", 6, "Board width. Defaults to 6x5")
	flag.IntVar(&combo_flag, "combo", 7, "Minimum number of combos to stop matching at.")
	flag.BoolVar(&diagonals_flag, "diagonals", false, "Whether to allow diagonals.")
	flag.Parse()

	if board_width < 5 || board_width > 7 {
		panic("Board width should be in the range [5,7]")
	}
	if board_flag == "" {
		board_to_solve = CreateBoard("BLRBRLHHLLBRHRBBHBRBLBRGLRBDRG", 6)
		return
	}

	if len(board_flag) != (board_width * (board_width - 1)) {
		err := fmt.Sprintf("Board size expected to be %d, got %d", board_width * board_width - 1, len(board_flag))
		panic(err)
	}
	board_to_solve = CreateBoard(board_flag, board_width)
}

func main() {
	// swng_board := CreateBoard("RRHHHDGDGGGDGDDRRRRDRGGDGRRDDG", 6)
	// combos := swng_board.GetAllCombos()
	// fmt.Println(combos)
	requirement := SolveRequirement{combo_flag, diagonals_flag}

	fmt.Printf("Solving board:\n%s", board_to_solve)

	moves := BfsFourDirectionSolver{}.Solve(board_to_solve, requirement)

	fmt.Println(moves)
}
