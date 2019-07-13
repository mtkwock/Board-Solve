package main

import (
	"fmt"
)

type SolveRequirement struct {
	ComboMinimum int
	Diagonals bool
	// TODO: Add different types of requirements
	// Max number of moves
	// Minimum color combos
	// Minimum combo types
}

type Moves struct {
	StartingPosition Placement
	Directions []Direction
}

var DirectionToLetter map[Direction]string = map[Direction]string {
	RIGHT: "R",
	DOWN_RIGHT: "DR",
	DOWN: "D",
	DOWN_LEFT: "DL",
	LEFT: "L",
	UP_LEFT: "UL",
	UP: "U",
	UP_RIGHT: "UR",
}

func DirectionsToString(directions []Direction) string {
	result := ""
	for _, direction := range directions {
		result += fmt.Sprintf("%s, ", DirectionToLetter[direction])
	}
	return result
}

func (self Moves) String() string {
	result := fmt.Sprintf("[%d, %d]: ", self.StartingPosition.Y, self.StartingPosition.X)
	result += DirectionsToString(self.Directions)
	return result
}

type BoardSolver interface {
	Solve(Board, SolveRequirement) Moves
}

type BfsFourDirectionSolver struct {}

type BfsState struct {
	current_board Board
	starting_pos Placement
	current_pos Placement
	moves []Direction
}

func (self BfsState) Clone() BfsState {
	copied := BfsState{moves: make([]Direction, len(self.moves))}
	copied.starting_pos = self.starting_pos
	copied.current_pos = self.current_pos
	copied.current_board = self.current_board.Clone()
	copy(copied.moves, self.moves)
	// Note: Parent is not copied
	// fmt.Println(self.moves)
	return copied
}

type BestState struct {
	State BfsState
	Combos int
}

var DirectionReverse map[Direction]Direction = map[Direction]Direction {
	RIGHT: LEFT,
	DOWN_RIGHT: UP_LEFT,
	DOWN: UP,
	DOWN_LEFT: UP_RIGHT,
	LEFT: RIGHT,
	UP_LEFT: DOWN_RIGHT,
	UP: DOWN,
	UP_RIGHT: DOWN_LEFT,
}

func (self BfsState) NextStates(requirements SolveRequirement) []BfsState {
	next_states := make([]BfsState, 0)
	reverse_move := DirectionReverse[self.moves[len(self.moves) - 1]]
	moves := []Direction{RIGHT, DOWN, LEFT, UP}
	if requirements.Diagonals {
		moves = []Direction{RIGHT, DOWN_RIGHT, DOWN, DOWN_LEFT, LEFT, UP_LEFT, UP, UP_RIGHT}
	}
	for _, direction := range moves {
		// fmt.Printf("BfsState: %s", self)
		if direction == reverse_move {
			continue
		}
		next_placement := self.current_pos.Swap(direction)
		if next_placement.Y >= self.current_board.Height ||
		   next_placement.X >= self.current_board.Width {
			continue
		}
		next_state := self.Clone()
		new_board, err := next_state.current_board.Swap(self.current_pos, direction)
		if err != nil {
			continue
		}
		next_state.current_board = new_board
		next_state.moves = append(next_state.moves, direction)
		next_state.current_pos = next_placement
		// next_state.parent = &self
		next_states = append(next_states, next_state)
	}
	return next_states
}

func Validate(board Board, requirement SolveRequirement) bool {
	combos := board.GetAllCombos()
	return len(combos) >= requirement.ComboMinimum
}

// CircularQueue code copied from https://stackoverflow.com/a/11757161
// CircularQueue is a basic FIFO CircularQueue based on a circular list that resizes as needed.
type CircularQueue struct {
	nodes	[]*BfsState
	head	int
	tail	int
	count	int
}

// Push adds a node to the CircularQueue.
func (q *CircularQueue) Push(n *BfsState) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*BfsState, len(q.nodes)*2)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	// fmt.Printf("Pushing state %s at position: %d\n", DirectionsToString((*n).moves), q.tail)
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the CircularQueue in first to last order.
func (q *CircularQueue) Pop() *BfsState {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	// fmt.Printf("Popping state %s at position: %d\n", DirectionsToString((*node).moves), q.head)
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

func (s BfsFourDirectionSolver) Solve(board Board, requirements SolveRequirement) Moves {
	// Initialize States
	queue := CircularQueue{nodes: make([]*BfsState, 1 << 10)}
  for y := uint8(0); y < board.Height; y++ {
		for x := uint8(0); x < board.Width; x++ {
			starting_pos := Placement{y, x}
			if x < board.Width - 1 {
				board, err := board.Swap(starting_pos, RIGHT)
				if err == nil {
					// new_state :=
					queue.Push(&BfsState {
						board,
						starting_pos,
						starting_pos.Swap(RIGHT),
						[]Direction{RIGHT},
					})
				}
				if requirements.Diagonals && y < board.Height - 1 {
					board_diagonal, err := board.Swap(starting_pos, DOWN_RIGHT)
					if err == nil {
						// new_diagonal :=
						queue.Push(&BfsState {
							board_diagonal,
							starting_pos,
							starting_pos.Swap(DOWN_RIGHT),
							[]Direction{DOWN_RIGHT},
						})
					}
				}
			}
			if y < board.Height - 1 {
				board, err := board.Swap(starting_pos, DOWN)
				if err == nil {
					// new_state :=
					queue.Push(&BfsState {
						board,
						starting_pos,
						starting_pos.Swap(DOWN),
						[]Direction{DOWN},
					})
				}
				if requirements.Diagonals && x > 0 {
					board_diagonal, err := board.Swap(starting_pos, DOWN_LEFT)
					if err == nil {
						// new_diagonal :=
						queue.Push(&BfsState {
							board_diagonal,
							starting_pos,
							starting_pos.Swap(DOWN_LEFT),
							[]Direction{DOWN_LEFT},
						})
					}
				}			}
			if x > 0 {
				board, err := board.Swap(starting_pos, LEFT)
				if err == nil {
					// new_state :=
					queue.Push(&BfsState {
						board,
						starting_pos,
						starting_pos.Swap(LEFT),
						[]Direction{LEFT},
					})
					if requirements.Diagonals && y < 0 {
						board_diagonal, err := board.Swap(starting_pos, UP_LEFT)
						if err == nil {
							queue.Push(&BfsState {
								board_diagonal,
								starting_pos,
								starting_pos.Swap(UP_LEFT),
								[]Direction{UP_LEFT},
							})
						}
					}
				}
			}
			if y > 0 {
				board, err := board.Swap(starting_pos, UP)
				if err == nil {
					// new_state :=
					queue.Push(&BfsState {
						board,
						starting_pos,
						starting_pos.Swap(UP),
						[]Direction{UP},
					})
					if requirements.Diagonals && x < board.Width - 1 {
						board_diagonal, err := board.Swap(starting_pos, UP_RIGHT)
						if err == nil {
							// new_diagonal :=
							queue.Push(&BfsState {
								board_diagonal,
								starting_pos,
								starting_pos.Swap(UP_RIGHT),
								[]Direction{UP_RIGHT},
							})
						}
					}
				}
			}
		}
	}
	// fmt.Println(queue.count)

	best_state := BestState{BfsState{}, 0}

	checked := 0
	skipped := 0

	known_boards := make(map[string]bool)
	for state_ptr := queue.Pop();
			!Validate(best_state.State.current_board, requirements);
			state_ptr = queue.Pop() {
		if state_ptr == nil {
			// fmt.Println(known_boards)
			fmt.Printf("Known Board size: %d\n", len(known_boards))
			// fmt.Println(queue.count)
			break
		}
		current_state := *state_ptr
		// fmt.Printf("%s: %s\n", current_state.starting_pos, DirectionsToString(current_state.moves))
		// fmt.Println(current_state)
		// fmt.Printf("%d|%d, ", queue.count, len(current_state.moves))
		board_string := current_state.current_pos.String() + current_state.current_board.SimpleString()
		// if board_string == "(1,1)RDHRDBBBDBHLBHGLDBHRRBBDRLLRRL" {
			// fmt.Println("SPECIAL CASE")
			// parent := *current_state.parent
			// fmt.Printf("Parent: %s: %s\n", parent.starting_pos, DirectionsToString(parent.moves))
		// }
		if _, exists := known_boards[board_string]; exists {
			skipped++
			// if len(current_state.moves) < 20 {
			// 	fmt.Printf("Conflicts: %s, %s | %s, %s\n", board_string, val, current_state.starting_pos, DirectionsToString(current_state.moves))
			// 	if current_state.parent != nil {
			// 		parent := *current_state.parent
			// 		fmt.Printf("Parent: %s: %s\n", parent.starting_pos, DirectionsToString(parent.moves))
			// 	} else {
			// 		fmt.Println("No parent!")
			// 	}
			// }
			// if skipped > 1000 {
			// 	fmt.Println(known_boards)
			// 	break
			// }
			// break
			continue
		}
		known_boards[board_string] = true // current_state.starting_pos.String() + DirectionsToString(current_state.moves)

		checked++
		// fmt.Printf("%s\n", current_state.current_board.GetAllCombos())
		if len(current_state.current_board.GetAllCombos()) > best_state.Combos {
			best_state.State = current_state
			best_state.Combos = len(current_state.current_board.GetAllCombos())
			current_moves := Moves{best_state.State.starting_pos, best_state.State.moves}
			fmt.Printf("Current best (%dc): %s\n", best_state.Combos, current_moves)
		}
		next_states := current_state.NextStates(requirements)
		// 		for _, next_state := range current_state.NextStates() {
		for i := 0; i < len(next_states); i++ {
			// fmt.Printf("%d\n", len(next_states))
			next_state := next_states[i]
			// fmt.Printf("Adding state - %s: %s\n", next_state.starting_pos, DirectionsToString(next_state.moves))
			queue.Push(&next_state)
		}
		// Check every 1,000,000 iterations
		if checked % 1000000 == 0 {
			fmt.Printf("Checked %d\n", checked)
		}
	}
	fmt.Println(best_state.State.current_board)
	combos := best_state.State.current_board.GetAllCombos()
	for _, combo := range combos {
		combo.Print(board.Width)
	}
	return Moves{best_state.State.starting_pos, best_state.State.moves}
}
