package main

import (
	"errors"
	"fmt"
	"strings"
	"math/rand" // Runs deterministically unless we set a seed.
)

type Pair struct {
	Y uint8
	X uint8
}

func (self Pair) ToPos(board Board) uint8 {
	return self.Y * board.Width + self.X
}

func (self Pair) String() string {
	return fmt.Sprintf("(%d,%d)", self.Y, self.X)
}

func (self Pair) Clone() Pair {
	return Pair{self.Y, self.X}
}

func (self Pair) Swap(direction Direction) Pair {
	new_placement := self.Clone()
	// TODO: Consider changing this to a hashmap of placement, direction: placement.
	// Will need to profile if faster or not.
	switch direction {
	case UP:
		new_placement.Y--
  case UP_RIGHT:
		new_placement.Y--
		new_placement.X++
	case RIGHT:
		new_placement.X++
	case DOWN_RIGHT:
		new_placement.Y++
		new_placement.X++
	case DOWN:
		new_placement.Y++
	case DOWN_LEFT:
		new_placement.Y++
		new_placement.X--
	case LEFT:
		new_placement.X--
	case UP_LEFT:
		new_placement.Y--
		new_placement.X--
	}
	return new_placement
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

func (self Pair) ManhattanDistance(other Pair) int {
	return abs(int(self.X) - int(other.X)) + abs(int(self.Y) - int(other.Y))
}

type BoardSpace struct {
	Orb Orb
	State BoardSpaceStateFlag
}

func (self BoardSpace) Clone() BoardSpace {
	return BoardSpace{self.Orb.Clone(), self.State}
}

func (self BoardSpace) String() string {
	if self.State & CLOUD != 0 {
		return "{%}"
	}
	if self.State & TAPE != 0 {
		return fmt.Sprintf("<%s>", string(self.Orb.String()[1]))
	}
	return self.Orb.String()
}

type Board struct {
	Slots []BoardSpace
	Height uint8
	Width uint8
	MinimumMatch int
}

func (self Board) Clone() Board {
	new_slots := make([]BoardSpace, len(self.Slots))
	for i, board_space := range self.Slots {
		new_slots[i] = board_space.Clone()
	}
	return Board{new_slots, self.Height, self.Width, self.MinimumMatch}
}

func (self Board) String() string {
	border := "++" + strings.Repeat("-", 3 * int(self.Width)) + "++\n"
	body := ""
	for y := uint8(0); y < self.Height; y++ {
		body += "||"
		for x := uint8(0); x < self.Width; x++ {
			body += self.Slots[Pair{y, x}.ToPos(self)].String()
		}
		body += "||\n"
	}
	return fmt.Sprintf("%s%s%s", border, body, border)
}

func (self Board) Print() {
	fmt.Println(self)
}

func (self Board) SimpleString() string {
	result := make([]byte, len(self.Slots))
	for i, slot := range self.Slots {
		result[i] = byte(AttributeToLetter[slot.Orb.Attribute][0])
	}
	return string(result)
}

func (self Board) Swap(placement Pair, direction Direction) (Board, error) {
	new_placement := placement.Swap(direction)

	new_board := self.Clone()

	if new_placement.X >= self.Width || new_placement.Y >= self.Height {
		return new_board, errors.New("Invalid move. Off board.")
	}
	new_pos := new_placement.ToPos(self)
	old_pos := placement.ToPos(self)

  if self.Slots[new_pos].State & TAPE != 0 {
		return new_board, errors.New("Invalid move. New Position is Immobilized.")
	}
	if self.Slots[old_pos].State & TAPE != 0 {
		return new_board, errors.New("Invalid move. Original Position is Immobilized.")
	}
	temp := new_board.Slots[new_pos].Orb.Clone()
	new_board.Slots[new_pos].Orb = new_board.Slots[old_pos].Orb.Clone()
	new_board.Slots[old_pos].Orb = temp
	new_board.Slots[new_pos].Orb.State &= ^BLIND
	new_board.Slots[old_pos].Orb.State &= ^BLIND
	return new_board, nil
}

func (self Board) GetCounts() map[OrbAttribute]int {
	result := map[OrbAttribute]int{}
	for _, slot := range self.Slots {
		attribute := slot.Orb.Attribute
		if _, exists := result[attribute]; exists {
			result[attribute]++
		} else {
			result[attribute] = 1
		}
	}
	return result
}

type BoardCombo struct {
	Attribute OrbAttribute
	Positions []Pair
	// IsEnhanced bool
}

func (self BoardCombo) String() string {
	return fmt.Sprintf("%s: [%d]", AttributeToName[self.Attribute], len(self.Positions))
}

func (self BoardCombo) Print(width uint8) {
	fmt.Println(self)
	board := CreateEmptyBoard(width)
	for _, placement := range self.Positions {
		board.Slots[placement.ToPos(board)].Orb.Attribute = self.Attribute
	}
	fmt.Println(board)
}

// type BoardRestriction uint16
//
// const (
// 	RESTRICT_4PLUS BoardRestriction = 1 << iota
// 	RESTRICT_5PLUS
// 	RESTRICT_FIRE
// 	RESTRICT_WATER
// 	RESTRICT_LIGHT
// 	RESTRICT_DARK
// 	RESTRICT_HEART
// 	RESTRICT_POISON
// )

func (self Board) GetOrbAt(placement Pair) Orb {
	// Guard Clause
	if placement.Y >= self.Height || placement.X >= self.Width {
		return Orb{EMPTY, 0}
	}
	return self.Slots[placement.ToPos(self)].Orb
}

func (self Board) markCombos(dst []bool) {
	horizontal_max := self.Width - 2
	vertical_max := self.Height - 2
  for y := uint8(0); y < self.Height; y++ {
		for x := uint8(0); x < self.Width; x++ {
			placement := Pair{y, x}
			position := placement.ToPos(self)
			// Already known to be matchable, ignore.
			if dst[position] {
				continue
			}

			// Unmatchable orb, ignore
			orb := self.GetOrbAt(placement)
			if orb.Attribute == EMPTY || orb.State & UNMATCHABLE != 0 {
				continue
			}

			// Determine matches to the right.
			if (x < horizontal_max) {
        orb_next := self.GetOrbAt(Pair{placement.Y, placement.X + 1})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Pair{placement.Y, placement.X + 2})
					if orb.Attribute == orb_next_next.Attribute {
						dst[position] = true
						dst[position + 1] = true
						dst[position + 2] = true
						continue
					}
				}
			}

			// Determine matches to the left.
			if (x > 1) {
				orb_next := self.GetOrbAt(Pair{placement.Y, placement.X - 1})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Pair{placement.Y, placement.X - 2})
					if orb.Attribute == orb_next_next.Attribute {
						// fmt.Println(y, x)
						dst[position] = true
						dst[position - 1] = true
						dst[position - 2] = true
						continue
					}
				}
			}

			// Determine matches below.
			if (y < vertical_max) {
				orb_next := self.GetOrbAt(Pair{placement.Y + 1, placement.X})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Pair{placement.Y + 2, placement.X})
					if orb.Attribute == orb_next_next.Attribute {
						dst[position] = true
						dst[position + self.Width] = true
						dst[position + 2 * self.Width] = true
						continue
					}
				}
			}

			// Determine matches above.
			if (y > 1) {
				orb_next := self.GetOrbAt(Pair{placement.Y - 1, placement.X})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Pair{placement.Y - 2, placement.X})
					if orb.Attribute == orb_next_next.Attribute {
						dst[position] = true
						dst[position - self.Width] = true
						dst[position - 2 * self.Width] = true
						continue
					}
				}
			}

			if (x > 0 && x < self.Width - 1) {
				orb_next := self.GetOrbAt(Pair{placement.Y, placement.X - 1})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Pair{placement.Y, placement.X + 1})
					if orb.Attribute == orb_next_next.Attribute {
						dst[position] = true
						continue
					}
				}
			}

			if (y > 0 && y < self.Height - 1) {
				orb_next := self.GetOrbAt(Pair{placement.Y - 1, placement.X})
				if orb.Attribute == orb_next.Attribute {
					orb_next_next := self.GetOrbAt(Pair{placement.Y + 1, placement.X})
					if orb.Attribute == orb_next_next.Attribute {
						dst[position] = true
						continue
					}
				}
			}
		}
	}
}

func (self Board) ToPair(pos uint8) Pair {
	return Pair{pos / self.Width, pos % self.Width}
}

// TODO: Add BoardRestriction capabilities.
func (self Board) GetCombos() ([]BoardCombo, Board) {
	// Determine which orbs will be comboed out. Do not group them yet.
  marked_combos := make([]bool, len(self.Slots))
	self.markCombos(marked_combos)

	// Determine if there are any unmatched bombs.  If so, clear orbs first then
	// get combos from that new board.
	unmatched_bombs := make([]Pair, 0)
	for i := uint8(0); i < uint8(len(self.Slots)); i++ {
		if self.Slots[i].Orb.Attribute == BOMB && !marked_combos[i] {
			unmatched_bombs = append(unmatched_bombs, Pair{uint8(i / self.Width), i % self.Width})
		}
	}
	if len(unmatched_bombs) != 0 {
		new_board := self.Clone()
		for _, placement := range unmatched_bombs {
			for x2 := uint8(0); x2 < self.Width; x2++ {
				pos := Pair{placement.Y, x2}.ToPos(new_board)
				if new_board.Slots[pos].Orb.Attribute != BOMB {
					new_board.Slots[pos].Orb.Attribute = EMPTY
				}
			}
			for y2 := uint8(0); y2 < self.Height; y2++ {
				pos := Pair{y2, placement.X}.ToPos(new_board)
				if new_board.Slots[pos].Orb.Attribute != BOMB {
					new_board.Slots[pos].Orb.Attribute = EMPTY
				}
			}
			new_board.Slots[placement.ToPos(new_board)].Orb.Attribute = EMPTY
		}
		return new_board.GetCombos()
	}

  // Group orbs to combo out.
	is_used := make([]bool, len(self.Slots))
	combos := make([]BoardCombo, 0)

  // For each orb, use a DFS to find all connected orbs.
	for i, is_comboed := range marked_combos {
		if !is_comboed || is_used[i] {
			continue
		}
		is_used[i] = true
		placement := Pair{uint8(i) / self.Width, uint8(i) % self.Width}
		attribute := self.GetOrbAt(placement).Attribute
		placements := []Pair{placement}
		for j := 0; j < len(placements); j++ {
			current := placements[j]
			// Check to the right.
			pos := current.ToPos(self) + 1
			right := Pair{current.Y, current.X + 1}
			if self.GetOrbAt(right).Attribute == attribute && marked_combos[pos] && !is_used[pos] {
				placements = append(placements, right)
				is_used[pos] = true
			}
			// Check to the left
			pos -= 2
			left := Pair{current.Y, current.X - 1}
			if self.GetOrbAt(left).Attribute == attribute && marked_combos[pos] && !is_used[pos]{
				placements = append(placements, left)
				is_used[pos] = true
			}
			// Check below
			pos += self.Width + 1
			down := Pair{current.Y + 1, current.X}
			if self.GetOrbAt(down).Attribute == attribute && marked_combos[pos] && !is_used[pos]{
				placements = append(placements, down)
				is_used[pos] = true
			}
			// Check above
			pos -= 2 * self.Width
			up := Pair{current.Y - 1, current.X}
			if self.GetOrbAt(up).Attribute == attribute && marked_combos[pos] && !is_used[pos]{
				placements = append(placements, up)
				is_used[pos] = true
			}
		}
		// For leads such as Khepri/Keela, orb matches must contain at least
		// MinimumMatch orbs (5/4 respectively).
		if len(placements) < self.MinimumMatch {
			continue
		}
		combo := BoardCombo{attribute, placements}
		// TODO: Mark the combos with special values.
		combos = append(combos, combo)
	}

	return combos, self
}

func (self Board) dropOrbs() {
	for y := self.Height - 1; y > 0; y-- {
		for x := uint8(0); x < self.Width; x++ {
			if self.GetOrbAt(Pair{y, x}).Attribute != EMPTY {
				continue
			}
			pos := Pair{y, x}.ToPos(self)
			for yo := y - 1; yo < self.Height; yo-- {
				new_pos := Pair{yo, x}.ToPos(self)
				moved_orb := self.Slots[new_pos].Orb
				if moved_orb.Attribute != EMPTY {
					self.Slots[pos].Orb = moved_orb
					self.Slots[new_pos].Orb = Orb{EMPTY, 0}
					break
				}
			}
		}
	}
}

func (self Board) GetAllCombos() []BoardCombo {
	all_combos := make([]BoardCombo, 0)
	current_board := self.Clone()
	for new_combos, current_board := current_board.GetCombos(); len(new_combos) > 0; new_combos, current_board = current_board.GetCombos() {
		all_combos = append(all_combos, new_combos...)
		for _, combo := range new_combos {
			for _, placement := range combo.Positions {
				current_board.Slots[placement.ToPos(current_board)].Orb.Attribute = EMPTY
			}
		}
		current_board.dropOrbs()
	}
	return all_combos
}

func CreateEmptyBoard(width uint8) Board {
	return Board{make([]BoardSpace, width * (width - 1)), width - 1, width, 3}
}

func CreateRandomBoard(width uint8) Board {
	size := width * (width - 1)
	slots := make([]BoardSpace, size)
	for i := uint8(0); i < size; i++ {
		slots[i].Orb.Attribute = OrbAttribute(uint8(rand.Intn(6) + 1))
	}
	return Board{slots, width - 1, width, 3}
}

func CreateBoard(s string, width int) Board {
  slots := make([]BoardSpace, len(s))
	for i, rune := range s {
		slots[i].Orb.Attribute = LetterToAttribute[string(rune)]
	}
	return Board{slots, uint8(len(s) / width), uint8(width), 3}
}
