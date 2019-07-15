package main

import (
	"flag"
	"fmt"
	"strings"
	"strconv"
)

var (
	board_to_solve Board

	// User defined flags.
	board_flag string
	flag_board_width int
	flag_combo_minimum int
	flag_allow_diagonals bool
	flag_combo_weight int
	flag_move_weight int
	flag_max_moves int
	flag_timeout_ms int
	flag_starting_positions string
)


func init() {
	flag.StringVar(&board_flag, "board", "",
			"Board String (R)ed, (B)lue, (G)reen, (L)ight, (D)ark, (H)eart, (P)oison), (M)ortal Poison, (J)ammer, B(o)mb.")
	flag.IntVar(&flag_board_width, "width", 6, "Board width. Height will be (width-1)")
	flag.IntVar(&flag_combo_minimum, "combo", 7, "Minimum number of combos to stop matching at.")
	flag.BoolVar(&flag_allow_diagonals, "diagonals", false, "Whether to allow diagonals.")
	flag.IntVar(&flag_combo_weight, "combo_weight", 20, "How much combos are scored relative to move count. Higher value calculates faster, lower prioritizes fewer moves.")
	flag.IntVar(&flag_move_weight, "move_weight", 1, "How much moves cost in value for heuristic.")
	flag.IntVar(&flag_max_moves, "max_moves", 50, "Maximum number of allowable moves.")
	flag.IntVar(&flag_timeout_ms, "timeout", -1, "How long to keep calculating (ms) before giving up. Negative is indefinite. (Unimplemented)")
	flag.StringVar(&flag_starting_positions, "starting_positions", "", "Allowable starting positions in 0-indexed Y-X separated format. e.g. \"0-0,2-1,9-9\". (Unimplemented)")

	flag.Parse()

	if flag_board_width < 5 || flag_board_width > 7 {
		panic("Board width should be in the range [5,7]")
	}
	if board_flag == "" {
		board_to_solve = CreateRandomBoard(uint8(flag_board_width))
		return
	}

	if len(board_flag) != (flag_board_width * (flag_board_width - 1)) {
		err := fmt.Sprintf("Board size expected to be %d, got %d", flag_board_width * (flag_board_width - 1), len(board_flag))
		panic(err)
	}
	board_to_solve = CreateBoard(board_flag, flag_board_width)
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
	acceptance_fn := func(state AStarState) bool {
		return len(state.board.GetAllCombos()) >= flag_combo_minimum // tate.score >= (flag_combo_minimum * 20)
	}
	known_boards := map[string]int{}
	rejection_fn := func(state AStarState) bool {
		key := state.current_pos.String() + state.board.SimpleString()
		// value := state.starting_pos.String() + DirectionsToString(state.moves)
		if old_val, exists := known_boards[key]; exists && state.score <= old_val {
			return true
		}
		if len(state.moves) > flag_max_moves {
			return true
		}
		known_boards[key] = state.score
		return false
	}
	scoring_fn := func(state AStarState) int {
		return len(state.board.GetAllCombos()) * flag_combo_weight - (flag_move_weight * len(state.moves))
	}

	requirement := SolveRequirement{
		flag_allow_diagonals,
		// Determines if a state meets the goal.
		acceptance_fn,
		// Determines if a state should be ignored.
		rejection_fn,
		// Determines and updates a state's score.
		scoring_fn,
	}

	fmt.Printf("Solving board:\n%s", board_to_solve)

	moves := AStarSolve(board_to_solve, requirement)

	fmt.Println(moves)
	fmt.Println(ToDawnglare(board_to_solve, moves))
}
