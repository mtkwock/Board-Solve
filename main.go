package main

import (
	"flag"
	"fmt"
	"strings"
	"strconv"
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
		board_to_solve = CreateRandomBoard(uint8(board_width))
		return
	}

	if len(board_flag) != (board_width * (board_width - 1)) {
		err := fmt.Sprintf("Board size expected to be %d, got %d", board_width * board_width - 1, len(board_flag))
		panic(err)
	}
	board_to_solve = CreateBoard(board_flag, board_width)
}

// Get a string that links to Dawnglare for the given board and moves.
func ToDawnglare(board Board, moves Moves) string {
	dawnglare_pattern := "https://pad.dawnglare.com/?height=%d&width=%d&patt=%s&replay=%s"
	height := int(board.Height)
	width := int(board.Width)

	// Convert the movements into positions for Dawnglare to use.
	positions := make([]string, len(moves.Directions) + 1)
	current_placement := moves.StartingPosition
	positions[0] = strconv.Itoa(int(current_placement.ToPos(board)))
	for i, direction := range moves.Directions {
		current_placement = current_placement.Swap(direction)
		positions[i + 1] = strconv.Itoa(int(current_placement.ToPos(board)))
	}
	move_string := strings.Join(positions, "|")

	return fmt.Sprintf(dawnglare_pattern, height, width, board.SimpleString(), move_string)
}

func main() {
	requirement := SolveRequirement{combo_flag, diagonals_flag}

	fmt.Printf("Solving board:\n%s", board_to_solve)

	moves := BfsFourDirectionSolver{}.Solve(board_to_solve, requirement)

	fmt.Println(moves)
	fmt.Println(ToDawnglare(board_to_solve, moves))
}
