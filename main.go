package main

import (
	"flag"
	"fmt"
	"strings"
	"strconv"
	"time"
)

var (
	board_to_solve Board
	starting_placements []Placement

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
	flag.IntVar(&flag_board_width, "width", 6, "Board width. Height will be (width-1)")
	flag.StringVar(&board_flag, "board", "",
			"Board String (R)ed, (B)lue, (G)reen, (L)ight, (D)ark, (H)eart, (P)oison), (M)ortal Poison, (J)ammer, B(o)mb.")
	flag.IntVar(&flag_combo_minimum, "combo", 7, "Minimum number of combos to stop matching at.")
	flag.BoolVar(&flag_allow_diagonals, "diagonals", false, "Whether to allow diagonals.")
	flag.IntVar(&flag_combo_weight, "combo_weight", 20, "How much combos are scored relative to move count. Higher value calculates faster, lower prioritizes fewer moves.")
	flag.IntVar(&flag_move_weight, "move_weight", 1, "How much moves cost in value for heuristic.")
	flag.IntVar(&flag_max_moves, "max_moves", 50, "Maximum number of allowable moves.")
	flag.IntVar(&flag_timeout_ms, "timeout_ms", -1, "How long to keep calculating (ms) before giving up. Negative is indefinite.")
	flag.StringVar(&flag_starting_positions, "starting_positions", "", "Allowable starting positions in 0-indexed Y-X separated format. e.g. \"0,0|2,1|4,5\".")
	flag.Parse()

	if flag_board_width < 5 || flag_board_width > 7 {
		panic("Board width should be in the range [5,7]")
	}
	starting_placements = make([]Placement, 0)
	if flag_starting_positions != "" {
		for _, coordinate := range strings.Split(flag_starting_positions, "|") {
			coordinate_vals := strings.Split(coordinate, ",")
			if len(coordinate_vals) != 2 {
				panic("Coordinate string invalid: \"" + coordinate + "\"")
			}
			y, err := strconv.Atoi(coordinate_vals[0])
			if err != nil {
				panic(err)
			}
			x, err := strconv.Atoi(coordinate_vals[1])
			if err != nil {
				panic(err)
			}
			if y < 0 || y >= flag_board_width - 1 {
				panic(fmt.Sprintf("Y value %d outside of range [0,%d]", y, flag_board_width - 1))
			}
			if x < 0 || x >= flag_board_width {
				panic(fmt.Sprintf("X value %d outside of range [0,%d]", x, flag_board_width))
			}
			starting_placements = append(starting_placements, Placement{uint8(y), uint8(x)})
		}
	}

	if board_flag == "" {
		board_to_solve = CreateRandomBoard(uint8(flag_board_width))
	} else {
		if len(board_flag) != (flag_board_width * (flag_board_width - 1)) {
			err := fmt.Sprintf("Board size expected to be %d, got %d", flag_board_width * (flag_board_width - 1), len(board_flag))
			panic(err)
		}
		board_to_solve = CreateBoard(board_flag, flag_board_width)
	}
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
	timed_out := false
	if flag_timeout_ms > 0 {
		go func() {
			timeout_timer := time.NewTimer(time.Duration(flag_timeout_ms) * time.Millisecond)
			<-timeout_timer.C
			timed_out = true
		}()
	}
	acceptance_fn := func(state AStarState) bool {
		if timed_out {
			fmt.Printf("Timed out after %dms. Returning best value.\n", flag_timeout_ms)
			return true
		}
		return len(state.board.GetAllCombos()) >= flag_combo_minimum
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
		move_cost := 3
		total_cost := 3
		last_move := state.moves[0]
		for _, move := range state.moves[1:] {
			if move == last_move {
				if move_cost > 1 {
					move_cost--
				}
			} else {
				move_cost = 3
			}
			last_move = move
			total_cost += move_cost
		}
		return len(state.board.GetAllCombos()) * flag_combo_weight - (flag_move_weight * total_cost)
	}

	requirement := SolveRequirement{
		flag_allow_diagonals,
		// Determines if a state meets the goal.
		acceptance_fn,
		// Determines if a state should be ignored.
		rejection_fn,
		// Determines and updates a state's score.
		scoring_fn,
		// Allowable starting positions. If empty, search all.
		starting_placements,
	}

	fmt.Printf("Solving board:\n%s", board_to_solve)

	moves := AStarSolve(board_to_solve, requirement)

	fmt.Println(moves)
	fmt.Println(ToDawnglare(board_to_solve, moves))
}
