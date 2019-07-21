package main

import (
	"container/heap"
	"fmt"
)

type SolveRequirement struct {
	AllowDiagonals bool
	// Determines if a state meets the goal.
	FinishedFn func(AStarState) bool
	// Determines if a state should be ignored.
	RejectionFn func(AStarState) bool
	// Determines and updates a state's score.
	ScoreState func(AStarState) int
	// Determine allowable starting positions. If empty slice, search all.
	StartingPositions []Pair
}

type Moves struct {
	StartingPosition Pair
	Directions []Direction
}

func (self Moves) String() string {
	result := fmt.Sprintf("[%d, %d]: ", self.StartingPosition.Y, self.StartingPosition.X)
	result += DirectionsToString(self.Directions)
	return result
}

type AStarState struct {
	board Board
	starting_pos Pair
	current_pos Pair
	moves []Direction
	// combos []BoardCombo // Should we store this?
	score int
}

func (self AStarState) Clone() AStarState {
	copied := AStarState{moves: make([]Direction, len(self.moves))}
	copied.starting_pos = self.starting_pos
	copied.current_pos = self.current_pos
	copied.board = self.board.Clone()
	copy(copied.moves, self.moves)
	// Note: Parent is not copied
	// fmt.Println(self.moves)
	return copied
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

func (self AStarState) NextStates(requirements SolveRequirement) []AStarState {
	next_states := make([]AStarState, 0)
	reverse_move := DirectionReverse[self.moves[len(self.moves) - 1]]
	moves := []Direction{RIGHT, DOWN, LEFT, UP}
	if requirements.AllowDiagonals {
		moves = []Direction{RIGHT, DOWN_RIGHT, DOWN, DOWN_LEFT, LEFT, UP_LEFT, UP, UP_RIGHT}
	}
	for _, direction := range moves {
		if direction == reverse_move {
			continue
		}
		next_placement := self.current_pos.Swap(direction)
		if next_placement.Y >= self.board.Height ||
		   next_placement.X >= self.board.Width {
			continue
		}
		next_state := self.Clone()
		new_board, err := next_state.board.Swap(self.current_pos, direction)
		if err != nil {
			continue
		}
		next_state.board = new_board
		next_state.moves = append(next_state.moves, direction)
		next_state.current_pos = next_placement
		// next_state.parent = &self
		next_states = append(next_states, next_state)
	}
	return next_states
}

// Implementation of Heap which keeps a circular queue as underlying data.
type StatePriorityQueue struct {
	data []*AStarState
	head int
	tail int
	count int
}

func (self StatePriorityQueue) Len() int {
	return self.count
}

func (self StatePriorityQueue) ToLocal(i int) int {
	return (i + self.head) % len(self.data)
}

func (self *StatePriorityQueue) Less(i, j int) bool {
	return self.data[self.ToLocal(i)].score > self.data[self.ToLocal(j)].score
}

func (self *StatePriorityQueue) Swap(i, j int) {
	// fmt.Println(i, j)
	local_i, local_j := self.ToLocal(i), self.ToLocal(j)
	self.data[local_i], self.data[local_j] = self.data[local_j], self.data[local_i]
}

func (self *StatePriorityQueue) Push(x interface{}) {
	if self.head == self.tail && self.count > 0 {
		new_data := make([]*AStarState, len(self.data) * 4)
		copy(new_data, self.data[self.head:])
		copy(new_data[len(self.data) - self.head:], self.data[:self.head])
		self.head = 0
		self.tail = len(self.data)
		self.data = new_data
	}
	self.data[self.tail] = x.(*AStarState)
	self.tail = (self.tail + 1) % len(self.data)
	self.count++
	// fmt.Println(self.count)
}

func (self *StatePriorityQueue) Pop() interface{} {
	if self.count == 0 {
		return nil
	}
	// fmt.Printf("Popping from %d\n", self.head)
	// fmt.Println(self.data[self.head] == self.data[self.head + 1])
	state := self.data[self.head]
	// fmt.Printf("Popping state %s at position: %d\n", DirectionsToString((*node).moves), q.head)
	self.head = (self.head + 1) % len(self.data)
	self.count--
	return state
}

func MakePriorityQueue(size int) StatePriorityQueue {
	data := make([]*AStarState, size)
	queue := StatePriorityQueue{data, 0, 0, 0}
	heap.Init(&queue)
	return queue
}

// CircularQueue code copied from https://stackoverflow.com/a/11757161
// CircularQueue is a basic FIFO CircularQueue based on a circular list that resizes as needed.
type CircularQueue struct {
	nodes	[]*AStarState
	head	int
	tail	int
	count	int
}

// Push adds a node to the CircularQueue.
func (q *CircularQueue) Push(n *AStarState) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*AStarState, len(q.nodes)*2)
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
func (q *CircularQueue) Pop() *AStarState {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	// fmt.Printf("Popping state %s at position: %d\n", DirectionsToString((*node).moves), q.head)
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

func AStarSolve(board Board, requirements SolveRequirement) Moves {
	// Initialize States
	// queue := CircularQueue{nodes: make([]*AStarState, 1 << 10)}
	queue := MakePriorityQueue(1 << 10)
	moves := []Direction{RIGHT, DOWN, LEFT, UP}
	if requirements.AllowDiagonals {
		moves = []Direction{RIGHT, DOWN_RIGHT, DOWN, DOWN_LEFT, LEFT, UP_LEFT, UP, UP_RIGHT}
	}
	var starting_positions []Pair = requirements.StartingPositions
	if len(starting_positions) == 0 {
		for y := uint8(0); y < board.Height; y++ {
			for x := uint8(0); x < board.Width; x++ {
				starting_positions = append(starting_positions, Pair{y, x})
			}
		}
	}
  for _, starting_pos := range starting_positions {
		for i := 0; i < len(moves); i++ {
			move := moves[i]
			board, err := board.Swap(starting_pos, move)
			if err != nil {
				continue
			}
			new_state := AStarState{
				board,
				starting_pos,
				starting_pos.Swap(move),
				[]Direction{move},
				0,
			}
			if !requirements.RejectionFn(new_state) {
				heap.Push(&queue, &new_state)
				// queue.Push(&new_state)
			}
		}
	}
	fmt.Printf("Queue initial size: %d\n", queue.count)

	best_state := AStarState{score: -100000}

	checked := 0
	skipped := 0

	var last_ptr *AStarState = nil
	for state_ptr := heap.Pop(&queue).(*AStarState);
			!requirements.FinishedFn(best_state);
			state_ptr = heap.Pop(&queue).(*AStarState) {
		if state_ptr == nil {
			fmt.Println("Ran out of boards to check, exiting.")
			break
		}
		if state_ptr == last_ptr {
			// fmt.Println((*state_ptr).moves)
			panic("This should not happen.")
		}
		last_ptr = state_ptr
		current_state := *state_ptr

		requirements.ScoreState(current_state)
		// fmt.Printf("%s\n", current_state.board.GetAllCombos())
		if current_state.score > best_state.score {
			best_state = current_state
			// best_state.Combos = len(current_state.board.GetAllCombos())
			current_moves := Moves{current_state.starting_pos, current_state.moves}
			// fmt.Println(best_state.board)
			fmt.Printf("Current best with score of %d\n%s\n", best_state.score, current_moves)
		}
		next_states := current_state.NextStates(requirements)
		// 		for _, next_state := range current_state.NextStates() {
		for i := 0; i < len(next_states); i++ {
			// fmt.Printf("%d\n", len(next_states))
			next_state := next_states[i]
			next_state.score = requirements.ScoreState(next_state)
			if requirements.RejectionFn(next_state) {
				skipped++
				break
			}
			// fmt.Printf("Adding state - %s: %s\n", next_state.starting_pos, DirectionsToString(next_state.moves))
			// queue.Push(&next_state)
			heap.Push(&queue, &next_state)
		}
		// Check every 1,000,000 iterations
		checked++
		if (checked + skipped) % 1000000 == 0 {
			fmt.Printf("Checked: %d, Skipped: %d\n", checked, skipped)
		}
	}
	fmt.Printf("Finished after %d checks with %d skipped.\n", checked, skipped)
	fmt.Println(best_state.board)
	combos := best_state.board.GetAllCombos()
	for _, combo := range combos {
		combo.Print(board.Width)
	}
	return Moves{best_state.starting_pos, best_state.moves}
}
