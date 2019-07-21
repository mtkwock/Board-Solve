package main

import (
	"fmt"
)

type SetupCombo struct {
	Attribute OrbAttribute
	// Unordered coordinates.
	Positions []Pair
}

func (self SetupCombo) Clone() SetupCombo {
	positions := make([]Pair, len(self.Positions))
	for i, position := range self.Positions {
		positions[i] = position.Clone()
	}
	return SetupCombo{self.Attribute, positions}
}

type BoardSetup struct {
	Combos []SetupCombo

	// These values are set when calling Init.
	PositionToAttribute map[uint8]OrbAttribute
	width uint8
}

func (self *BoardSetup) Init(width uint8) {
	if len(self.PositionToAttribute) > 0 {
		return
	}
	self.PositionToAttribute = make(map[uint8]OrbAttribute, 0)
	self.width = width
	for _, combo := range self.Combos {
		for _, pair := range combo.Positions {
			self.PositionToAttribute[pair.Y * width + pair.X] = combo.Attribute
		}
	}
}

// Assumes a 6x5 board.
func (self BoardSetup) String() string {
	if self.width == 0 {
		panic("BoardSetup not Initialized with Init()!")
	}
	slots := make([]BoardSpace, self.width * (self.width - 1))

	for pos, attribute := range self.PositionToAttribute {
		slots[pos].Orb.Attribute = attribute
	}

	board := Board{slots, 5, 6, 0}

	return fmt.Sprintf("Board Setup:\n%s\n", board.String())
}

var WidthToOrder map[uint8][]uint8 = map[uint8][]uint8 {
	// 1 2 3 2 1
	// 2 3 4 3 2
	// 2 3 4 3 2
	// 1 2 3 2 1
	5: []uint8 {
		0, 4, 15, 19,
		1, 3, 5, 9, 10, 14, 16, 18,
		2, 6, 8, 11, 13, 17,
		7, 12,
	},
	// 1 2 3 3 2 1
	// 2 3 4 4 3 2
	// 3 4 5 5 4 3
	// 2 3 4 4 3 2
	// 1 2 3 3 2 1
	6: []uint8 {
		0, 5, 24, 29, // 1
		1, 4, 6, 11, 18, 23, 25, 28, // 2
		2, 3, 7, 10, 12, 17, 19, 22, 26, 27, // 3
		8, 9, 13, 16, 20, 21, // 4
		14, 15, // 5
	},
	// 1 2 3 4 3 2 1
	// 2 3 4 5 4 3 2
	// 3 4 5 6 5 4 3
	// 3 4 5 6 5 4 3
	// 2 3 4 5 4 3 2
	// 1 2 3 4 3 2 1
	7: []uint8 {
		0, 6, 35, 41,
		1, 5, 7, 13, 28, 34, 36, 40,
		2, 4, 8, 12, 14, 20, 21, 27, 29, 33, 37, 39,
		3, 9, 11, 15, 19, 22, 26, 30, 32, 38,
		10, 16, 18, 23, 25, 31,
		17, 24,
	},
}

func (self BoardSetup) Clone() BoardSetup {
	combos := make([]SetupCombo, len(self.Combos))
	for i, combo := range self.Combos {
		combos[i] = combo.Clone()
	}
	return BoardSetup{combos, make(map[uint8]OrbAttribute, 0), 0}
}

// Determine the Manhattan distance of this board compared to others.
// Note that the orbs in BoardSetup are *required* to be present in board.
// Analysis must be done beforehand.
// TODO: Do this thoroughly instead of Greedily.
func (self BoardSetup) ManhattanDistanceGreedyEdges(board Board) int {
	self.Init(board.Width)

	board_slots := make(map[OrbAttribute][]Pair, 0)
	for i, slot := range board.Slots {
		if val, exists := board_slots[slot.Orb.Attribute]; exists {
			board_slots[slot.Orb.Attribute] = append(val, board.ToPair(uint8(i)))
		} else {
			board_slots[slot.Orb.Attribute] = []Pair{board.ToPair(uint8(i))}
		}
	}

	total_distance := 0

	for _, pos := range WidthToOrder[board.Width] {
		if attribute, exists := self.PositionToAttribute[uint8(pos)]; exists {
			del_idx := -1
			lowest := 1000
			goal_pos := board.ToPair(uint8(pos))
			for i := 0; i < len(board_slots[attribute]); i++ {
				distance := goal_pos.ManhattanDistance(board_slots[attribute][i])
				if distance < lowest {
					del_idx = i
					lowest = distance
					// fmt.Println(lowest)
				}
			}
			total_distance += lowest
			// Remove the greedily taken value.
			board_slots[attribute] = append(board_slots[attribute][:del_idx],
                                      board_slots[attribute][del_idx + 1:]...)
		}
	}

	return total_distance
}

func (self BoardSetup) ManhattanDistanceAverage(board Board) float32 {
	total_distance := self.ManhattanDistanceGreedyEdges(board)
	return float32(total_distance) / float32(len(self.PositionToAttribute))
}

func (self BoardSetup) MirrorHorizontal(width uint8) BoardSetup {
	base_idx := width - 1
  new_setup := self.Clone()
	for i := 0; i < len(new_setup.Combos); i++ {
		for j := 0; j < len(new_setup.Combos[i].Positions); j++ {
			new_setup.Combos[i].Positions[j].X = base_idx - new_setup.Combos[i].Positions[j].X
		}
	}
	new_setup.Init(width)
	return new_setup
}

func (self BoardSetup) MirrorVertical(height uint8) BoardSetup {
	base_idx := height - 1
	new_setup := self.Clone()
	for i := 0; i < len(new_setup.Combos); i++ {
		for j := 0; j < len(new_setup.Combos[i].Positions); j++ {
			new_setup.Combos[i].Positions[j].Y = base_idx - new_setup.Combos[i].Positions[j].Y
		}
	}
	new_setup.Init(height + 1)
	return new_setup
}

// Used to determine valid orbs to start with.
// If all of one attribute is used for the setup, ignores those positions.
// This assumes that sufficient orbs exist, and will panic if this is violated.
func (self BoardSetup) UnusedOrbIdxs(board Board) []Pair {
	self.Init(board.Width)
	result := make([]Pair, 0)

  orb_counts := board.GetCounts()
	for _, attr := range self.PositionToAttribute {
		orb_counts[attr]--
	}

	for attr, count := range orb_counts {
		if count < 0 {
			panic(fmt.Sprintf("Precondition violated, insufficient orbs of %s", attr))
		}
		if count == 0 {
			continue
		}
		for idx, slot := range board.Slots {
			if slot.Orb.Attribute == attr {
				result = append(result, board.ToPair(uint8(idx)))
			}
		}
	}
	return result
}

// Given a BoardSetup, solve using a combo by combo basis.
func StrategySolve(board Board, setup BoardSetup, requirements SolveRequirement) Moves {
	unused_idxs := setup.UnusedOrbIdxs(board)
	var starting_positions []Pair = make([]Pair, 0)
	if len(requirements.StartingPositions) == 0 {
		starting_positions = unused_idxs
	} else {
		for _, pos := range unused_idxs {
			available := false
			for _, manual_pos := range requirements.StartingPositions {
				if pos.Y == manual_pos.Y && pos.X == manual_pos.X {
					available = true
					break
				}
			}
			if available {
				starting_positions = append(starting_positions, pos)
			}
		}
	}

	current_board := board.Clone()
	moves := Moves{}

	for _, combo := range setup.Combos {
		sub_known_boards := map[string]int{}

		sub_requirements := SolveRequirement {
			AllowDiagonals: requirements.AllowDiagonals,
			FinishedFn: func(state AStarState) bool {
				for _, pair := range combo.Positions {
					if state.board.GetOrbAt(pair).Attribute != combo.Attribute {
						return false
					}
				}
				return true
			},
			RejectionFn: func(state AStarState) bool {
				key := state.current_pos.String() + state.board.SimpleString()
				if old_val, exists := sub_known_boards[key]; exists && state.score <= old_val {
					return true
				}
				// if len(state.moves) > flag_max_moves {
				// 	return true
				// }
				sub_known_boards[key] = state.score
				return false
			},
			ScoreState: func(state AStarState) int {
				// Ignore the currently held orb.
				temp := state.board.GetOrbAt(state.current_pos).Attribute
				state.board.Slots[state.current_pos.ToPos(state.board)].Orb.Attribute = EMPTY
				temp_board_setup := BoardSetup{Combos: []SetupCombo{combo}}
				distance_cost := 13 * temp_board_setup.ManhattanDistanceGreedyEdges(state.board)
				state.board.Slots[state.current_pos.ToPos(state.board)].Orb.Attribute = temp

				move_cost := 3
				total_move_cost := 3
				last_move := state.moves[0]
				for _, move := range state.moves[1:] {
					if move == last_move {
						if move_cost > 1 {
							move_cost = 1
						}
					} else {
						move_cost = 3
					}
					last_move = move
					total_move_cost += move_cost
				}

				return 10000 - (total_move_cost + distance_cost)
			},
			// Determine allowable starting positions. If empty slice, search all.
			StartingPositions: starting_positions,
		}

		// Create a board that ignores all values that aren't the given attribute.
		sub_board := current_board.Clone()
		for i := 0; i < len(sub_board.Slots); i++ {
			if sub_board.Slots[i].Orb.Attribute != combo.Attribute {
				sub_board.Slots[i].Orb.Attribute = EMPTY
			}
		}

		next_moves := AStarSolve(sub_board, sub_requirements)
		// Set up the first instantiation.
		if len(moves.Directions) == 0 {
			moves.Directions = next_moves.Directions
			moves.StartingPosition = next_moves.StartingPosition
		} else {
			moves.Directions = append(moves.Directions, next_moves.Directions...)
		}

		position := next_moves.StartingPosition
		for _, direction := range next_moves.Directions {
			new_board, err := current_board.Swap(position, direction)
			if err != nil {
				panic("Error occurred when expecting it not to.")
			}
			current_board = new_board
			position = position.Swap(direction)
		}
		// Force the next starting position.
		starting_positions = []Pair{position}
		for _, pair := range combo.Positions {
			current_board.Slots[pair.ToPos(current_board)].State |= TAPE
		}
	}

	last_requirement := SolveRequirement {
		AllowDiagonals: requirements.AllowDiagonals,
		FinishedFn: requirements.FinishedFn,
		RejectionFn: requirements.RejectionFn,
		ScoreState: requirements.ScoreState,
		StartingPositions: starting_positions,
	}
	last_moves := AStarSolve(current_board, last_requirement)
	moves.Directions = append(moves.Directions, last_moves.Directions...)

	return moves
}
